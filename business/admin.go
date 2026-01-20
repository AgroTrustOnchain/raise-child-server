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
	on_chain "raise-child/util/on_chain"
	"strings"
	"time"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/block-vision/sui-go-sdk/utils"
)

type adminService struct {
	profileRepo i_repository.IProfileRepository
	clients     map[string]sui.ISuiAPI
	errLogger   *log.Logger
}

func InitializeAdminService(db *sql.DB, errLogger *log.Logger) business.IAdminService {
	return &adminService{
		profileRepo: repository.InitializeProfileRepository(db, errLogger),
		clients:     _networkAliases,
		errLogger:   errLogger,
	}
}

func GenerateAdminService() (business.IAdminService, error) {
	var errLogger = util.GetLogConfig(shared.ERROR_LEVEL)

	cnn, err := db.ConnectDB(errLogger, db.InitializePostgreSQL())
	if err != nil {
		return nil, err
	}

	return InitializeAdminService(cnn, errLogger), nil
}

// UpdatePublisherInfo implements business.IAdminService.
func (a *adminService) UpdatePublisherInfo(req request.UpdatePublisherInfoRequest, ctx context.Context) (response.BuildTransactionResponse, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)

	var sender string = ctx.Value("address").(string)
	if !utils.IsValidSuiAddress(models.SuiAddress(sender)) {
		return response.BuildTransactionResponse{}, genericErr
	}

	var profileId string = ctx.Value("sub").(string)
	profile, err := a.profileRepo.GetFirstProfile(ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	if profile == nil || profileId != profile.ID || profile.IdentityCode != "" {
		return response.BuildTransactionResponse{}, genericErr
	}

	var dateOfBirth string = strings.TrimSpace(req.DateOfBirth)
	if dob := util.RawDateToTime(dateOfBirth); dob.IsZero() {
		return response.BuildTransactionResponse{}, genericErr
	}

	var gender string = util.StanderizeGender(util.StanderizeString(req.Gender))
	if gender == "" {
		return response.BuildTransactionResponse{}, genericErr
	}

	var email string = strings.TrimSpace(req.Email)
	if !util.IsValidEmail(email) {
		return response.BuildTransactionResponse{}, genericErr
	}

	var identityCode string = strings.TrimSpace(req.IdentityCode)
	var firstName string = strings.TrimSpace(req.FirstName)
	var lastName string = strings.TrimSpace(req.LastName)
	var phoneNumber string = strings.TrimSpace(req.PhoneNumber)

	// todo: validate identity code
	var module = on_chain.InitializeModuleManage()
	txBytes, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
		Client:    a.clients[constant.SuiTestnet],
		Sender:    sender,
		Module:    module.GetModule(),
		Function:  module.GetFunctionUpdatePublisherNft(),
		ErrLogger: a.errLogger,
		Arguments: module.ToUpdatePublisherNftArguements(on_chain.UpdatePublisherNftArguements{
			IdentityCode:       identityCode,
			IdentityCardBlobID: strings.TrimSpace(req.IdentityCardBlobID),
			AvatarBlobID:       strings.TrimSpace(req.AvatarBlobID),
			FirstName:          firstName,
			LastName:           lastName,
			Gender:             gender,
			DateOfBirth:        dateOfBirth,
			PhoneNumber:        phoneNumber,
			Email:              email,
		}),
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	profile.IdentityCode = identityCode
	profile.FirstName = firstName
	profile.LastName = lastName
	profile.Gender = gender
	profile.DateOfBirth = dateOfBirth
	profile.PhoneNumber = phoneNumber
	profile.Email = email
	profile.UpdatedAt = time.Now()

	return response.BuildTransactionResponse{
		TxBytes: txBytes,
	}, a.profileRepo.UploadProfile(*profile, ctx)
}
