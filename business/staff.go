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
	"sort"
	"strings"
	"time"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/block-vision/sui-go-sdk/utils"
)

type staffService struct {
	redisCache cache.IRedisCache
	clients    map[string]sui.ISuiAPI
	errLogger  *log.Logger
}

func initializeStaffService(clients map[string]sui.ISuiAPI, errLogger *log.Logger) business.IStaffService {
	return &staffService{
		redisCache: cache.InitializeRedisCache(),
		clients:    clients,
		errLogger:  errLogger,
	}
}

func GenerateStaffService() (business.IStaffService, error) {
	return initializeStaffService(_networkAliases, util.GetLogConfig(shared.ERROR_LEVEL)), nil
}

const (
	staff_records_limit int = 10
)

// GetStaff implements business.IStaffService.
func (s *staffService) GetStaff(id string, ctx context.Context) (response.StaffResponse, error) {
	if !utils.IsValidSuiAddress(models.SuiAddress(id)) {
		return response.StaffResponse{}, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	var client = s.clients[constant.SuiTestnet]
	staff, err := getOnChainObject[entities.Staff](client, id, s.errLogger, ctx)
	if err != nil {
		return response.StaffResponse{}, err
	}

	var res response.StaffResponse = staff.ToStaffResponse()

	var module = on_chain.InitializeModuleStaff()
	if nfts, _ := on_chain.GetOnChainOwnedObjects[entities.StaffNft](on_chain.GetOnChainOwnedObjectsRequest{
		Client:       client,
		OwnerAddress: staff.User,
		StructType:   fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), module.GetModule(), module.GetStaffNftObjectStruct()),
		ErrLogger:    s.errLogger,
	}, ctx); nfts != nil {
		var nftsRes []response.StaffNftResponse
		for _, nft := range nfts {
			nftsRes = append(nftsRes, nft.ToStaffNftResponse())
		}

		res.Nfts = nftsRes
	}

	return res, nil
}

// GetStaffs implements business.IStaffService.
func (s *staffService) GetStaffs(req request.GetStaffsRequest, ctx context.Context) (response.PaginationDataResponse, error) {
	req.Keyword = util.StanderizeString(req.Keyword)
	req.Region = util.StanderizeString(req.Region)
	req.SortOrder = util.StanderizeSortOrder(req.SortOrder)
	if req.Page < 1 {
		req.Page = 1
	}

	if req.PageSize < 1 {
		req.PageSize = default_page_size
	}

	var res response.PaginationDataResponse
	var redisKey string = s.getGetStaffsRedisKey(req)
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

	staffs, err := on_chain.GetOnChainObjects[entities.StaffNft](on_chain.GetOnChainObjectsRequest{
		Client:    client,
		ObjectIds: manageObj.ChildIds,
		ErrLogger: s.errLogger,
	}, ctx)

	if err != nil {
		return response.PaginationDataResponse{}, err
	}

	if staffs == nil {
		return response.PaginationDataResponse{}, nil
	}

	var filteredStaffs []entities.StaffNft
	for i := len(staffs) - 1; i >= 0; i-- {
		var staff entities.StaffNft = staffs[i]

		if req.Role != "" {
			if util.StanderizeString(staff.Role) != util.StanderizeString(req.Role) { // Not matched
				continue
			}
		}

		if req.Region != "" {
			if util.StanderizeString(staff.Region) != req.Region { // Not matched
				continue
			}
		}

		if req.Gender != "" {
			if staff.Gender != req.Gender { // Not matched
				continue
			}
		}

		if req.YearOfBirth != nil {
			var dob time.Time = util.RawDateToTime(staff.DateOfBirth)
			if dob.Year() != *req.YearOfBirth { // Not matched
				continue
			}
		}

		if req.Keyword != "" {
			var firstName string = util.StanderizeString(staff.FirstName)
			var lastName string = util.StanderizeString(staff.LastName)
			if !strings.Contains(firstName, req.Keyword) && !strings.Contains(lastName, req.Keyword) && !strings.Contains(staff.IdentityCode, req.Keyword) && !strings.Contains(staff.PhoneNumber, req.Keyword) && !strings.Contains(staff.Email, req.Keyword) { // Not matched
				continue
			}
		}

		filteredStaffs = append(filteredStaffs, staff)
	}

	sort.Slice(filteredStaffs, func(i, j int) bool {
		var name1 string = filteredStaffs[i].LastName + " " + filteredStaffs[i].FirstName
		var name2 string = filteredStaffs[j].LastName + " " + filteredStaffs[j].FirstName

		if req.SortOrder == "ASC" {
			return name1 < name2
		}

		return name2 > name1
	})

	var skippedRecords int = (req.Page - 1) * req.PageSize
	if len(filteredStaffs) <= skippedRecords {
		return response.PaginationDataResponse{}, err
	}

	var data []response.StaffNftResponse
	for i := skippedRecords; i < len(filteredStaffs); i++ {
		data = append(data, filteredStaffs[i].ToStaffNftResponse())
	}

	res = response.PaginationDataResponse{
		Data:       data,
		Amount:     len(data),
		Page:       req.Page,
		TotalPages: int(math.Ceil(float64(len(filteredStaffs)) / float64(req.PageSize))),
	}

	s.redisCache.Set(redisKey, res, time.Minute*5, ctx)

	return res, nil
}

// GetStaffsV2 implements business.IStaffService.
func (s *staffService) GetStaffsV2(req request.GetStaffsRequest, ctx context.Context) (response.PaginationDataResponse, error) {
	panic("unimplemented")
}

func (s *staffService) getGetStaffsRedisKey(req request.GetStaffsRequest) string {
	var keyword string = "empty"
	if req.Keyword != "" {
		keyword = req.Keyword
	}

	var role string = "empty"
	if req.Role != "" {
		role = req.Role
	}

	var region string = "empty"
	if req.Region != "" {
		region = req.Region
	}

	var gender string = "empty"
	if req.Gender != "" {
		gender = req.Gender
	}

	var yob string = "empty"
	if req.YearOfBirth != nil {
		yob = fmt.Sprintf("%d", *req.YearOfBirth)
	}

	return fmt.Sprintf("staff:kw:%s:role:%s:region:%s:g:%s:y:%s:o:%s:s:%d:p:%d",
		keyword, role, region, gender, yob, req.SortOrder, req.PageSize, req.Page)
}
