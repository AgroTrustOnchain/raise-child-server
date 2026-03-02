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
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/model/entities"
	"raise-child/util"
	"raise-child/util/cache"
	on_chain "raise-child/util/on_chain"
	"strings"
	"time"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/block-vision/sui-go-sdk/utils"
)

type donorService struct {
	redisCache cache.IRedisCache
	clients    map[string]sui.ISuiAPI
	errLogger  *log.Logger
}

func InitializeDonorService(clients map[string]sui.ISuiAPI, errLogger *log.Logger) business.IDonorService {
	return &donorService{
		redisCache: cache.InitializeRedisCache(),
		clients:    clients,
		errLogger:  errLogger,
	}
}

func GenerateDonorService() (business.IDonorService, error) {
	return InitializeDonorService(_networkAliases, util.GetLogConfig(shared.ERROR_LEVEL)), nil
}

const (
	donor_records_limit int = 10
)

// GetDonor implements business.IDonorService.
func (s *donorService) GetDonor(id string, ctx context.Context) (response.DonorResponse, error) {
	if !utils.IsValidSuiAddress(models.SuiAddress(id)) {
		return response.DonorResponse{}, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	var client = s.clients[constant.SuiTestnet]

	var donorModule = on_chain.InitializeModuleDonor()
	donors, err := on_chain.GetOnChainOwnedObjects[entities.Donor](on_chain.GetOnChainOwnedObjectsRequest{
		Client:       client,
		OwnerAddress: id,
		StructType:   fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), donorModule.GetModule(), donorModule.GetDonorNftStruct()),
		ErrLogger:    s.errLogger,
	}, ctx)
	if err != nil {
		return response.DonorResponse{}, err
	}

	var res = donors[0].ToDonorResponse()
	var recordModule = on_chain.InitializeModuleRecord()
	txs, _ := on_chain.GetOnChainOwnedObjects[entities.Transaction](on_chain.GetOnChainOwnedObjectsRequest{
		Client:       client,
		OwnerAddress: id,
		StructType:   fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), recordModule.GetModule(), recordModule.GetTransactionRecordStruct()),
		ErrLogger:    s.errLogger,
	}, ctx)

	if txs != nil && len(txs) > 0 {
		var contributions []response.TransactionResponse
		for _, tx := range txs {
			contributions = append(contributions, tx.ToTransactionResponse())
		}
		res.Contributions = contributions
	}

	// Set donor object ID to user wallet address
	res.ID = id

	return res, nil
}

// GetDonors implements business.IDonorService.
func (s *donorService) GetDonors(req request.GetDonorsRequest, ctx context.Context) (response.PaginationDataResponse, error) {
	req.Keyword = util.StanderizeString(req.Keyword)
	if req.Page < 1 {
		req.Page = 1
	}

	if req.PageSize < 1 {
		req.PageSize = default_page_size
	}

	var res response.PaginationDataResponse
	var redisKey string = s.getGetDonorsRedisKey(req)
	if s.redisCache.Get(redisKey, &res, ctx) {
		return res, nil
	}

	var client = s.clients[constant.SuiTestnet]
	manageObj, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: s.errLogger,
	}, ctx)
	if err != nil {
		return response.PaginationDataResponse{}, err
	}

	donors, err := on_chain.GetOnChainObjects[entities.Donor](on_chain.GetOnChainObjectsRequest{
		Client:    client,
		ObjectIds: manageObj.DonorNfts,
		ErrLogger: s.errLogger,
	}, ctx)
	if err != nil {
		return response.PaginationDataResponse{}, err
	}

	if donors == nil {
		return response.PaginationDataResponse{}, nil
	}

	var filteredDonors []entities.Donor
	for i := len(donors) - 1; i >= 0; i++ {
		var donor entities.Donor = donors[i]

		if req.Gender != "" {
			if req.Gender != donor.Gender {
				continue
			}
		}

		if req.Keyword != "" {
			var firstName string = util.StanderizeString(donor.FirstName)
			var lastName string = util.StanderizeString(donor.LastName)
			var phoneNumber string = util.StanderizeString(donor.PhoneNumber)
			var email string = util.StanderizeString(donor.Email)
			if !strings.Contains(firstName, req.Keyword) && !strings.Contains(lastName, req.Keyword) && !strings.Contains(phoneNumber, req.Keyword) && !strings.Contains(email, req.Keyword) {
				continue
			}
		}

		filteredDonors = append(filteredDonors, donor)
	}

	var skippedRecords int = (req.Page - 1) * req.PageSize
	if len(filteredDonors) <= skippedRecords {
		return response.PaginationDataResponse{}, err
	}

	var data []response.DonorResponse
	for i := skippedRecords; i < len(filteredDonors); i++ {
		data = append(data, filteredDonors[i].ToDonorResponse())
	}

	res = response.PaginationDataResponse{
		Data:       data,
		Amount:     len(data),
		Page:       req.Page,
		TotalPages: int(math.Ceil(float64(len(filteredDonors)) / float64(req.PageSize))),
	}

	s.redisCache.Set(redisKey, res, time.Minute*5, ctx)

	return res, nil
}

func (s *donorService) getGetDonorsRedisKey(req request.GetDonorsRequest) string {
	var keyword string = "empty"
	if req.Keyword != "" {
		keyword = req.Keyword
	}

	var gender string = "empty"
	if req.Gender != "" {
		gender = req.Gender
	}

	return fmt.Sprintf("donor:kw:%s:g:%s:s:%d:p:%d", keyword, gender, req.PageSize, req.Page)
}
