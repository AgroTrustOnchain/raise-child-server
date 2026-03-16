package business

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"raise-child/constants/noti"
	"raise-child/constants/shared"
	"raise-child/interfaces/business"
	i_repository "raise-child/interfaces/repository"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/repository"
	"raise-child/util"
	"raise-child/util/db"
	"strings"
	"time"

	"github.com/block-vision/sui-go-sdk/sui"
)

type profileService struct {
	profileRepo i_repository.IProfileRepository
	clients     map[string]sui.ISuiAPI
	errLogger   *log.Logger
}

func InitializeProfileService(db *sql.DB, errLogger *log.Logger) business.IProfileService {
	return &profileService{
		profileRepo: repository.InitializeProfileRepository(db, errLogger),
		clients:     _networkAliases,
		errLogger:   errLogger,
	}
}

func initializeProfileService(
	profileRepo i_repository.IProfileRepository,
	clients map[string]sui.ISuiAPI,
	errLogger *log.Logger,
) business.IProfileService {
	return &profileService{
		profileRepo: profileRepo,
		clients:     clients,
		errLogger:   errLogger,
	}
}

func GenerateProfileService() (business.IProfileService, error) {
	var errLogger = util.GetLogConfig(shared.ERROR_LEVEL)

	cnn, err := db.ConnectDB(errLogger, db.InitializePostgreSQL())
	if err != nil {
		return nil, err
	}

	return InitializeProfileService(cnn, errLogger), nil
}

// UploadProfile implements business.IProfileService.
func (p *profileService) UploadProfile(id string, req request.UploadProfileRequest, ctx context.Context) (response.PersonalProfileResponse, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)

	var sub string = ctx.Value("sub").(string)
	if sub != id { // other edits his/her profile
		return response.PersonalProfileResponse{}, errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	}

	profile, err := p.profileRepo.GetProfile(id, ctx)
	if err != nil {
		return response.PersonalProfileResponse{}, err
	}

	if profile == nil {
		return response.PersonalProfileResponse{}, genericErr
	}

	// Already upload profile
	if profile.IdentityCode != "" {
		return response.PersonalProfileResponse{}, genericErr
	}

	var gender string = util.StanderizeGender(util.StanderizeString(req.Gender))
	if gender == "" {
		return response.PersonalProfileResponse{}, errors.New(noti.UNDEFINED_GENDER_MESSAGE)
	}

	var dateOfBirth string = strings.TrimSpace(req.DateOfBirth)
	if dob := util.RawDateToTime(dateOfBirth); dob.IsZero() {
		return response.PersonalProfileResponse{}, errors.New(noti.INVALID_DATE_FORMAT_WARN_MSG)
	}

	var identityCode string = strings.TrimSpace(req.IdentityCode)
	var phoneNumber string = strings.TrimSpace(req.PhoneNumber)
	var email string = strings.TrimSpace(req.Email)
	if !util.IsValidEmail(email) {
		return response.PersonalProfileResponse{}, genericErr
	}

	isInfoExist, err := p.profileRepo.IsPersonalInfoExist(identityCode, phoneNumber, email, ctx)
	if err != nil {
		return response.PersonalProfileResponse{}, err
	}

	if isInfoExist {
		return response.PersonalProfileResponse{}, errors.New(noti.GENERIC_PERSONAL_INFO_REGISTERED_WANR_MSG)
	}

	profile.IdentityCode = identityCode
	profile.FirstName = strings.TrimSpace(req.FirstName)
	profile.LastName = strings.TrimSpace(req.LastName)
	profile.Gender = gender
	profile.DateOfBirth = dateOfBirth
	profile.PhoneNumber = phoneNumber
	profile.Email = email
	profile.UpdatedAt = time.Now()

	return (*profile).ToPersonalProfile(), p.profileRepo.UploadProfile(*profile, ctx)
}
