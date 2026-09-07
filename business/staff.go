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
	"raise-child/util/image/cloudinary"
	on_chain "raise-child/util/on_chain"
	"sort"
	"strings"
	"time"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/sui"
)

type staffService struct {
	registrationRequestRepo i_repository.IRegistrationRequestRepository
	redisCache              cache.IRedisCache
	clients                 map[string]sui.ISuiAPI
	errLogger               *log.Logger
}

func initializeStaffService(
	registrationRequestRepo i_repository.IRegistrationRequestRepository,
	clients map[string]sui.ISuiAPI,
	errLogger *log.Logger,
) business.IStaffService {
	return &staffService{
		registrationRequestRepo: registrationRequestRepo,
		redisCache:              cache.InitializeRedisCache(),
		clients:                 clients,
		errLogger:               errLogger,
	}
}

func GenerateStaffService() (business.IStaffService, error) {
	var errLogger = util.GetLogConfig(shared.ERROR_LEVEL)

	cnn, err := db.ConnectDB(errLogger, db.InitializePostgreSQL())
	if err != nil {
		return nil, err
	}

	return initializeStaffService(
		repository.InitializeRegistrationRequestRepo(cnn, errLogger),
		_networkAliases,
		errLogger,
	), nil
}

const (
	staff_records_limit int = 10
)

// GetStaffByOwnerWallet implements business.IStaffService.
func (s *staffService) GetStaffByOwnerWallet(id string, ctx context.Context) (response.StaffNftResponse, error) {
	if !util.IsValidSuiAddressStrict(id) {
		return response.StaffNftResponse{}, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	var client = s.clients[constant.SuiTestnet]
	manage, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: s.errLogger,
	}, ctx)
	if err != nil {
		return response.StaffNftResponse{}, err
	}

	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)
	if manage == nil {
		return response.StaffNftResponse{}, internalErr
	}

	var searchLen int
	if len(manage.VolunteerIds) > len(manage.LocalLeaderIds) {
		searchLen = len(manage.VolunteerIds)
	} else {
		searchLen = len(manage.LocalLeaderIds)
	}

	var staffId string
	for i := 0; i < searchLen; i++ {
		if i < len(manage.VolunteerIds) {
			if id == manage.VolunteerIds[i] {
				staffId = manage.VolunteerNfts[i]
				break
			}
		}

		if i < len(manage.LocalLeaderIds) {
			if id == manage.LocalLeaderIds[i] {
				staffId = manage.LocalLeaderNfts[i]
				break
			}
		}
	}

	if staffId == "" {
		return response.StaffNftResponse{}, nil
	}

	staff, err := getOnChainObject[entities.StaffNft](client, staffId, s.errLogger, ctx)
	if err != nil {
		return response.StaffNftResponse{}, err
	}

	if staff == nil {
		return response.StaffNftResponse{}, internalErr
	}

	var res response.StaffNftResponse = staff.ToStaffNftResponse()
	if req, _ := s.registrationRequestRepo.GetRegistrationRequestByOnchainId(staffId, ctx); req != nil {
		res.IdentityCardImgUrl = cloudinary.GetImageUrl(req.IdentityCardBlobID)
	} else {
		orgReq, _ := s.registrationRequestRepo.GetRegistrationRequestWithDetail(entities.GetRegistrationRequestWithDetail{
			RegisterRole:      staff.Role,
			IdentityCode:      staff.IdentityCode,
			AvatarBlobID:      staff.AvatarBlobID,
			Region:            staff.Region,
			FirstName:         staff.FirstName,
			LastName:          staff.LastName,
			Gender:            staff.Gender,
			DateOfBirth:       staff.DateOfBirth,
			PhoneNumber:       staff.PhoneNumber,
			Email:             staff.Email,
			Status:            request_approved_status,
			IsConfirmRegister: true,
			CreatedBy:         id,
		}, ctx)
		if orgReq != nil {
			orgReq.OnchainID = &staffId
			orgReq.UpdatedAt = time.Now()
			s.registrationRequestRepo.UpdateRegistrationRequest(*orgReq, ctx)
			res.IdentityCardImgUrl = cloudinary.GetImageUrl(orgReq.IdentityCardBlobID)
		}
	}

	return res, nil
}

// GetStaff implements business.IStaffService.
func (s *staffService) GetStaff(id string, ctx context.Context) (response.StaffNftResponse, error) {
	if !util.IsValidSuiAddressStrict(id) {
		return response.StaffNftResponse{}, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	var client = s.clients[constant.SuiTestnet]
	staff, err := getOnChainObject[entities.StaffNft](client, id, s.errLogger, ctx)
	if err != nil {
		return response.StaffNftResponse{}, err
	}

	var res response.StaffNftResponse = staff.ToStaffNftResponse()
	if req, _ := s.registrationRequestRepo.GetRegistrationRequestByOnchainId(id, ctx); req != nil {
		res.IdentityCardImgUrl = cloudinary.GetImageUrl(req.IdentityCardBlobID)
	} else {
		orgReq, _ := s.registrationRequestRepo.GetRegistrationRequestWithDetail(entities.GetRegistrationRequestWithDetail{
			RegisterRole:      staff.Role,
			IdentityCode:      staff.IdentityCode,
			AvatarBlobID:      staff.AvatarBlobID,
			Region:            staff.Region,
			FirstName:         staff.FirstName,
			LastName:          staff.LastName,
			Gender:            staff.Gender,
			DateOfBirth:       staff.DateOfBirth,
			PhoneNumber:       staff.PhoneNumber,
			Email:             staff.Email,
			Status:            request_approved_status,
			IsConfirmRegister: true,
			CreatedBy:         staff.Owner,
		}, ctx)
		if orgReq != nil {
			orgReq.OnchainID = &id
			orgReq.UpdatedAt = time.Now()
			s.registrationRequestRepo.UpdateRegistrationRequest(*orgReq, ctx)
			res.IdentityCardImgUrl = cloudinary.GetImageUrl(orgReq.IdentityCardBlobID)
		}
	}

	return res, nil
}

// GetStaffs implements business.IStaffService.
func (s *staffService) GetStaffs(req request.GetStaffsRequest, ctx context.Context) (response.PaginationDataResponse, error) {
	req.Keyword = util.StandardizeString(req.Keyword)
	req.Region = util.StandardizeString(req.Region)
	req.SortOrder = util.StandardizeSortOrder(req.SortOrder)
	if req.Page < 1 {
		req.Page = 1
	}

	if req.PageSize < 1 {
		req.PageSize = default_page_size
	}

	var res response.PaginationDataResponse
	var client = s.clients[constant.SuiTestnet]
	manageObj, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: s.errLogger,
	}, ctx)
	if err != nil {
		return response.PaginationDataResponse{}, err
	}

	var staffIds []string = manageObj.LocalLeaderNfts
	staffIds = append(staffIds, manageObj.VolunteerNfts...)
	staffs, err := on_chain.GetOnChainObjects[entities.StaffNft](on_chain.GetOnChainObjectsRequest{
		Client:    client,
		ObjectIds: staffIds,
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
			if util.StandardizeString(staff.Role) != util.StandardizeString(req.Role) { // Not matched
				continue
			}
		}

		if req.Region != "" {
			if util.StandardizeString(staff.Region) != req.Region { // Not matched
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
			var firstName string = util.StandardizeString(staff.FirstName)
			var lastName string = util.StandardizeString(staff.LastName)
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
		var staff = filteredStaffs[i].ToStaffNftResponse()
		if req, _ := s.registrationRequestRepo.GetRegistrationRequestByOnchainId(staff.ID, ctx); req != nil {
			staff.IdentityCardImgUrl = cloudinary.GetImageUrl(req.IdentityCardBlobID)
		} else {
			orgReq, _ := s.registrationRequestRepo.GetRegistrationRequestWithDetail(entities.GetRegistrationRequestWithDetail{
				RegisterRole:      staff.Role,
				IdentityCode:      staff.IdentityCode,
				AvatarBlobID:      staff.AvatarBlobID,
				Region:            staff.Region,
				FirstName:         staff.FirstName,
				LastName:          staff.LastName,
				Gender:            staff.Gender,
				DateOfBirth:       filteredStaffs[i].DateOfBirth,
				PhoneNumber:       staff.PhoneNumber,
				Email:             staff.Email,
				Status:            request_approved_status,
				IsConfirmRegister: true,
				CreatedBy:         staff.Owner,
			}, ctx)
			if orgReq != nil {
				orgReq.OnchainID = &staff.ID
				orgReq.UpdatedAt = time.Now()
				s.registrationRequestRepo.UpdateRegistrationRequest(*orgReq, ctx)
				staff.IdentityCardImgUrl = cloudinary.GetImageUrl(orgReq.IdentityCardBlobID)
			}
		}

		data = append(data, staff)
		if len(data) == req.PageSize {
			break
		}
	}

	res = response.PaginationDataResponse{
		Data:       data,
		Amount:     len(data),
		Page:       req.Page,
		TotalPages: int(math.Ceil(float64(len(filteredStaffs)) / float64(req.PageSize))),
	}

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
