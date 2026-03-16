package business

import (
	"context"
	"errors"
	"fmt"
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
	"raise-child/util/cache"
	"raise-child/util/db"
	on_chain "raise-child/util/on_chain"
	"sort"
	"strings"
	"time"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/sui"
)

type adminService struct {
	profileRepo i_repository.IProfileRepository
	redisCache  cache.IRedisCache
	clients     map[string]sui.ISuiAPI
	errLogger   *log.Logger
}

func initializeAdminService(profileRepo i_repository.IProfileRepository, clients map[string]sui.ISuiAPI, errLogger *log.Logger) business.IAdminService {
	return &adminService{
		profileRepo: profileRepo,
		redisCache:  cache.InitializeRedisCache(),
		clients:     clients,
		errLogger:   errLogger,
	}
}

func GenerateAdminService() (business.IAdminService, error) {
	var errLogger = util.GetLogConfig(shared.ERROR_LEVEL)

	cnn, err := db.ConnectDB(errLogger, db.InitializePostgreSQL())
	if err != nil {
		return nil, err
	}

	//return InitializeAdminService(cnn, errLogger), nil

	return initializeAdminService(
		repository.InitializeProfileRepository(cnn, errLogger),
		_networkAliases,
		errLogger,
	), nil
}

const (
	admin_records_limit int = 10
)

// GetAdmins implements business.IAdminService.
func (a *adminService) GetAdmins(req request.GetAdminsRequest, ctx context.Context) (response.PaginationDataResponse, error) {
	req.SortOrder = util.StanderizeSortOrder(req.SortOrder)
	req.Keyword = util.StanderizeString(req.Keyword)
	if req.Page < 1 {
		req.Page = 1
	}

	if req.PageSize < 1 {
		req.PageSize = default_page_size
	}

	var res response.PaginationDataResponse
	var redisKey string = a.getGetAdminsRedisKey(req)
	if a.redisCache.Get(redisKey, &res, ctx) {
		return res, nil
	}

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

	var filteredAdmins []entities.AdminNft
	for i := len(admins) - 1; i >= 0; i-- {
		var admin entities.AdminNft = admins[i]

		if req.Keyword != "" {
			var firstName string = util.StanderizeString(admin.FirstName)
			var lastName string = util.StanderizeString(admin.LastName)
			if !strings.Contains(firstName, req.Keyword) && !strings.Contains(lastName, req.Keyword) && !strings.Contains(admin.IdentityCode, req.Keyword) && !strings.Contains(admin.PhoneNumber, req.Keyword) && !strings.Contains(admin.Email, req.Keyword) { // Not matched
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

	sort.Slice(filteredAdmins, func(i, j int) bool {
		if req.SortCriteria == "date_of_birth" {
			var dob1 time.Time = util.RawDateToTime(filteredAdmins[i].DateOfBirth)
			var dob2 time.Time = util.RawDateToTime(filteredAdmins[j].DateOfBirth)
			if req.SortOrder == "DESC" {
				return dob2.After(dob1)
			}

			return dob2.Before(dob1)
		}

		if req.SortOrder == "ASC" {
			return false
		}

		return true
	})

	var skippedRecords int = (req.Page - 1) * req.PageSize
	if len(filteredAdmins) <= skippedRecords {
		return response.PaginationDataResponse{}, nil
	}

	var data []response.AdminNftResponse
	for i := skippedRecords; i < len(filteredAdmins); i++ {
		data = append(data, filteredAdmins[i].ToAdminNftResponse())
	}

	res = response.PaginationDataResponse{
		Data:       data,
		Amount:     len(data),
		Page:       req.Page,
		TotalPages: int(math.Ceil(float64(len(filteredAdmins)) / float64(req.PageSize))),
	}

	a.redisCache.Set(redisKey, res, time.Minute*5, ctx)

	return res, nil
}

// UpdatePublisherInfo implements business.IAdminService.
func (a *adminService) UpdatePublisherInfo(req request.UpdatePublisherInfoRequest, ctx context.Context) (response.BuildTransactionResponse, error) {
	var profileId string = ctx.Value("sub").(string)
	var profile *entities.Profile
	var errRes error
	var isFoundAdmin bool = false
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)

	for i := 1; i <= 3; i++ {
		profile, errRes = a.profileRepo.GetProfileOfFirsts(i, ctx)
		if errRes != nil {
			return response.BuildTransactionResponse{}, errRes
		}

		if profileId == profile.ID {
			if profile.IdentityCode != "" {
				return response.BuildTransactionResponse{}, genericErr
			}

			isFoundAdmin = true
			break
		}
	}

	if !isFoundAdmin {
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

	var client = a.clients[constant.SuiTestnet]
	var sender string = ctx.Value("address").(string)
	var manangeModule = on_chain.InitializeModuleManage()
	nfts, err := on_chain.GetOnChainOwnedObjects[entities.AdminNft](on_chain.GetOnChainOwnedObjectsRequest{
		Client:       client,
		OwnerAddress: sender,
		StructType:   fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), manangeModule.GetModule(), manangeModule.GetAdminNftStruct()),
		ErrLogger:    a.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	editProfileCaps, err := on_chain.GetOnChainOwnedObjects[entities.Cap](on_chain.GetOnChainOwnedObjectsRequest{
		Client:       client,
		OwnerAddress: sender,
		StructType:   fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), manangeModule.GetModule(), manangeModule.GetUpdateAdminInfoCapStruct()),
		ErrLogger:    a.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	var identityCode string = strings.TrimSpace(req.IdentityCode)
	var firstName string = strings.TrimSpace(req.FirstName)
	var lastName string = strings.TrimSpace(req.LastName)
	var phoneNumber string = strings.TrimSpace(req.PhoneNumber)

	// todo: validate identity code
	var module = on_chain.InitializeModuleManage()
	txBytes, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
		Client:    client,
		Sender:    sender,
		Module:    module.GetModule(),
		Function:  module.GetFunctionUpdatePublisherNft(),
		ErrLogger: a.errLogger,
		Arguments: module.ToUpdatePublisherNftArguments(on_chain.UpdatePublisherNftArguments{
			AdminCap:           editProfileCaps[0].ID.ID,
			AdminNft:           nfts[0].ID.ID,
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

func (a *adminService) getGetAdminsRedisKey(req request.GetAdminsRequest) string {
	var keyword string = "empty"
	if req.Keyword != "" {
		keyword = req.Keyword
	}

	var gender string = "empty"
	if req.Gender != "" {
		gender = req.Gender
	}

	var yob string = "empty"
	if req.YearOfBirth != nil {
		yob = fmt.Sprintf("%d", *req.YearOfBirth)
	}

	var sortCriteria string = "empty"
	if req.SortCriteria != "" {
		sortCriteria = req.SortCriteria
	}

	return fmt.Sprintf("admin:kw:%s:g:%s:y:%s:sc:%s:o:%s:s:%d:p:%d",
		keyword, gender, yob, sortCriteria, req.SortOrder, req.PageSize, req.Page,
	)
}
