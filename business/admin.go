package business

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"math"
	"os"
	"raise-child/constants/env"
	"raise-child/constants/noti"
	"raise-child/constants/shared"
	"raise-child/interfaces/business"
	i_repository "raise-child/interfaces/repository"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/model/entities"
	"raise-child/repository"
	"raise-child/util"
	"raise-child/util/db"
	on_chain "raise-child/util/on_chain"
	"sort"
	"strings"
	"time"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/signer"
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

const (
	admin_records_limit int = 10
)

// GetAdmins implements business.IAdminService.
func (a *adminService) GetAdmins(req request.GetAdminsRequest, ctx context.Context) (response.PaginationDataResponse, error) {
	var client = a.clients[constant.SuiTestnet]
	manageObj, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: a.errLogger,
	}, ctx)
	if err != nil {
		return response.PaginationDataResponse{}, err
	}

	admins, err := on_chain.GetOnChainObjects[entities.AdminNft](on_chain.GetOnChainObjectsRequest{
		Client:    client,
		ObjectIds: manageObj.AdminNfts,
		ErrLogger: a.errLogger,
	}, ctx)
	if err != nil {
		return response.PaginationDataResponse{}, err
	}

	var keyword string = util.StanderizeString(req.Keyword)
	var filteredAdmins []entities.AdminNft
	for i := len(admins) - 1; i >= 0; i++ {
		var admin entities.AdminNft = admins[i]

		if keyword != "" {
			var firstName string = util.StanderizeString(admin.FirstName)
			var lastName string = util.StanderizeString(admin.LastName)
			if !strings.Contains(firstName, keyword) && !strings.Contains(lastName, keyword) && !strings.Contains(admin.IdentityCode, keyword) && !strings.Contains(admin.PhoneNumber, keyword) && !strings.Contains(admin.Email, keyword) { // Not matched
				continue
			}
		}

		if req.Gender != "" {
			if admin.Gender != req.Gender { // Not matched
				continue
			}
		}

		if req.YearOfBirth != nil {
			var dob time.Time = util.RawDateToTime(admin.DateOfBirth)
			if dob.Year() != *req.YearOfBirth { // Not matched
				continue
			}
		}

		filteredAdmins = append(filteredAdmins, admin)
	}

	if req.SortOrder != "" {
		sort.Slice(filteredAdmins, func(i, j int) bool {
			if req.SortCriteria == "date_of_birth" {
				var dob1 time.Time = util.RawDateToTime(filteredAdmins[i].DateOfBirth)
				var dob2 time.Time = util.RawDateToTime(filteredAdmins[j].DateOfBirth)
				if req.SortOrder == "desc" {
					return dob2.After(dob1)
				}

				return dob2.Before(dob1)
			}

			if req.SortOrder == "asc" {
				return false
			}

			return true
		})
	}

	var page int = req.Page
	if page < 1 {
		page = 1
	}

	var skippedRecords int = (page - 1) * admin_records_limit
	if len(filteredAdmins) <= skippedRecords {
		return response.PaginationDataResponse{}, nil
	}

	var data []response.AdminNftResponse
	for i := skippedRecords; i < len(filteredAdmins); i++ {
		data = append(data, filteredAdmins[i].ToAdminNftResponse())
	}

	return response.PaginationDataResponse{
		Data:       data,
		Page:       page,
		TotalPages: int(math.Ceil(float64(len(filteredAdmins)) / float64(admin_records_limit))),
	}, nil
}

// UpdatePublisherInfo implements business.IAdminService.
func (a *adminService) UpdatePublisherInfo(req request.UpdatePublisherInfoRequest, ctx context.Context) (response.BuildTransactionResponse, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)

	signer, _ := signer.NewSignerWithSecretKey("")
	a.clients[""].SignAndExecuteTransactionBlock(ctx, models.SignAndExecuteTransactionBlockRequest{
		PriKey: signer.PriKey,
	})

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
		Arguments: module.ToUpdatePublisherNftArguments(on_chain.UpdatePublisherNftArguments{
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
