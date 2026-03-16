package business

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
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
	"slices"
	"strings"
	"time"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/block-vision/sui-go-sdk/utils"
)

const (
	admin_role        string = "Admin"
	local_leader_role string = "Local Leader"
	volunteer_role    string = "Volunteer"
	donor_role        string = "Donor"
	user_role         string = "User"
)

const (
	request_pending_status  string = "Pending"
	request_refused_status  string = "Refused"
	request_approved_status string = "Approved"
)

const (
	approve_rate_limit float32 = 0.7 // 70%
)

type registrationRequestService struct {
	registrationRequestRepo i_repository.IRegistrationRequestRepository
	profileRepo             i_repository.IProfileRepository
	redisCache              cache.IRedisCache
	clients                 map[string]sui.ISuiAPI
	errLogger               *log.Logger
}

func InitializeRegistrationRequestService(db *sql.DB, errLogger *log.Logger) business.IRegistrationRequestService {
	return &registrationRequestService{
		registrationRequestRepo: repository.InitializeRegistrationRequestRepo(db, errLogger),
		profileRepo:             repository.InitializeProfileRepository(db, errLogger),
		clients:                 _networkAliases,
		errLogger:               errLogger,
	}
}

func initializeRegistrationRequestService(
	registrationRequestRepo i_repository.IRegistrationRequestRepository,
	profileRepo i_repository.IProfileRepository,
	clients map[string]sui.ISuiAPI,
	errLogger *log.Logger,
) business.IRegistrationRequestService {
	return &registrationRequestService{
		registrationRequestRepo: registrationRequestRepo,
		profileRepo:             profileRepo,
		redisCache:              cache.InitializeRedisCache(),
		clients:                 clients,
		errLogger:               errLogger,
	}
}

func GenerateRegistrationRequestService() (business.IRegistrationRequestService, error) {
	var errLogger = util.GetLogConfig(shared.ERROR_LEVEL)

	cnn, err := db.ConnectDB(errLogger, db.InitializePostgreSQL())
	if err != nil {
		return nil, err
	}

	//return InitializeRegistrationRequestService(cnn, errLogger), nil
	return initializeRegistrationRequestService(
		repository.InitializeRegistrationRequestRepo(cnn, errLogger),
		repository.InitializeProfileRepository(cnn, errLogger),
		_networkAliases,
		errLogger,
	), nil
}

// ConfirmRegistrationRequest implements business.IRegistrationRequestService.
func (r *registrationRequestService) ConfirmRegistrationRequest(id string, ctx context.Context) (response.BuildTransactionResponse, error) {
	req, err := r.registrationRequestRepo.GetRegistrationRequest(id, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if req == nil {
		return response.BuildTransactionResponse{}, genericErr
	}

	var sender string = ctx.Value("address").(string)
	if req.CreatedBy != sender {
		return response.BuildTransactionResponse{}, errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	}

	// Pending process
	if req.ClosedAt.After(time.Now()) {
		return response.BuildTransactionResponse{}, errors.New(noti.STILL_PENDING_REQUEST_MESSAGE)
	} else { // Request closed
		var rate float32 = float32(len(req.Approvers)) / float32(len(req.Approvers)+len(req.Refusers))
		var isDenied bool = false

		if rate >= approve_rate_limit {
			req.Status = request_approved_status
			req.IsConfirmRegister = true
		} else {
			req.Status = request_refused_status
			isDenied = true
		}

		req.UpdatedAt = time.Now()
		if err := r.registrationRequestRepo.UpdateRegistrationRequest(*req, ctx); err != nil {
			return response.BuildTransactionResponse{}, err
		}

		if isDenied {
			return response.BuildTransactionResponse{}, nil
		}
	}

	// Wait for background server to mint cap object to register
	if !req.IsAvailableToConfirm {
		return response.BuildTransactionResponse{}, nil
	}

	var client = r.clients[constant.SuiTestnet]
	var manageModule = on_chain.InitializeModuleManage()
	var capType string

	switch req.RegisterRole {
	case admin_role:
		capType = manageModule.GetRegisterAdminCapStruct()
	case local_leader_role:
		capType = manageModule.GetRegisterLeaderCapStruct()
	case volunteer_role:
		capType = manageModule.GetRegisterVolunteerCapStruct()
	}

	caps, err := on_chain.GetOnChainOwnedObjects[entities.Cap](on_chain.GetOnChainOwnedObjectsRequest{
		Client:       client,
		OwnerAddress: sender,
		StructType:   fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), manageModule.GetModule(), capType),
		ErrLogger:    r.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	var args []interface{}
	var function string
	var staffModule = on_chain.InitializeModuleStaff()
	switch req.RegisterRole {
	case admin_role:
		function = staffModule.GetFunctionRegisterAdmin()
		args = staffModule.ToRegisterAdminArguments(on_chain.RegisterAdminArguments{
			CapID:              caps[0].ID.ID,
			IdentityCode:       req.IdentityCode,
			IdentityCardBlobID: req.IdentityCardBlobID,
			AvatarBlobID:       req.AvatarBlobID,
			FirstName:          req.FirstName,
			LastName:           req.LastName,
			Gender:             req.Gender,
			DateOfBirth:        req.DateOfBirth,
			PhoneNumber:        req.PhoneNumber,
			Email:              req.Email,
		})
	case local_leader_role:
		function = staffModule.GetFunctionRegisterLeader()
		args = staffModule.ToRegisterNormalStaffArguments(on_chain.RegisterNormalStaffArguments{
			Region: req.Region,
			RegisterAdminArguments: on_chain.RegisterAdminArguments{
				CapID:              caps[0].ID.ID,
				IdentityCode:       req.IdentityCode,
				IdentityCardBlobID: req.IdentityCardBlobID,
				AvatarBlobID:       req.AvatarBlobID,
				FirstName:          req.FirstName,
				LastName:           req.LastName,
				Gender:             req.Gender,
				DateOfBirth:        req.DateOfBirth,
				PhoneNumber:        req.PhoneNumber,
				Email:              req.Email,
			},
		})
	case volunteer_role:
		function = staffModule.GetFunctionRegisterVolunteer()
		args = staffModule.ToRegisterNormalStaffArguments(on_chain.RegisterNormalStaffArguments{
			Region: req.Region,
			RegisterAdminArguments: on_chain.RegisterAdminArguments{
				CapID:              caps[0].ID.ID,
				IdentityCode:       req.IdentityCode,
				IdentityCardBlobID: req.IdentityCardBlobID,
				AvatarBlobID:       req.AvatarBlobID,
				FirstName:          req.FirstName,
				LastName:           req.LastName,
				Gender:             req.Gender,
				DateOfBirth:        req.DateOfBirth,
				PhoneNumber:        req.PhoneNumber,
				Email:              req.Email,
			},
		})
	}

	var module = on_chain.InitializeModuleStaff()
	txBytes, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
		Client:    r.clients[constant.SuiTestnet],
		Sender:    sender,
		Module:    module.GetModule(),
		Function:  function,
		ErrLogger: r.errLogger,
		Arguments: args,
	}, ctx)

	return response.BuildTransactionResponse{
		TxBytes:        txBytes,
		RegistraionReq: id,
	}, err
}

// CreateRegistrationRequest implements business.IRegistrationRequestService.
func (r *registrationRequestService) CreateRegistrationRequest(req request.CreateRegistrationRequest, ctx context.Context) (*entities.RegistrationRequest, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)

	var sender string = ctx.Value("address").(string)
	if !utils.IsValidSuiAddress(models.SuiAddress(sender)) {
		return nil, genericErr
	}

	reqs, err := r.registrationRequestRepo.GetWalletRegistrationRequests(sender, ctx)
	if err != nil {
		return nil, err
	}

	var role string = strings.TrimSpace(req.RegisterRole)
	if reqs != nil && len(reqs) > 0 {
		for _, req := range reqs {
			if req.RegisterRole == role && (req.Status == request_pending_status || req.Status == request_approved_status) {
				return nil, genericErr
			}
		}
	}

	profile, err := r.profileRepo.GetProfile(ctx.Value("sub").(string), ctx)
	if err != nil {
		return nil, err
	}

	if profile == nil {
		return nil, genericErr
	}

	// Admin registration request not contain region
	if role == admin_role {
		if req.Region != "" {
			return nil, genericErr
		}
	}

	if !isRegionExist(req.Region) {
		return nil, genericErr
	}

	var client = r.clients[constant.SuiTestnet]
	manageObj, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: r.errLogger,
	}, ctx)
	if err != nil {
		return nil, err
	}

	if role == volunteer_role {
		if !slices.Contains(manageObj.LocalRegions, req.Region) {
			return nil, errors.New(noti.REGION_NOT_ADDED_WARN_MSG)
		}
	}

	// todo: validate identity code
	var curTime time.Time = time.Now()
	var request = entities.RegistrationRequest{
		ID:                   util.GenerateId(),
		ProfileID:            ctx.Value("sub").(string),
		RegisterRole:         role,
		IdentityCode:         util.StanderizeString(profile.IdentityCode),
		IdentityCardBlobID:   strings.TrimSpace(req.IdentityCardBlobID),
		AvatarBlobID:         strings.TrimSpace(req.AvatarBlobID),
		Region:               req.Region,
		FirstName:            strings.TrimSpace(profile.FirstName),
		LastName:             strings.TrimSpace(profile.LastName),
		Gender:               profile.Gender,
		DateOfBirth:          profile.DateOfBirth,
		PhoneNumber:          profile.PhoneNumber,
		Email:                profile.Email,
		Status:               request_pending_status,
		IsAvailableToConfirm: false,
		IsConfirmRegister:    false,
		CreatedBy:            sender,
		CreatedAt:            curTime,
		UpdatedAt:            curTime,
		ClosedAt:             util.GetRequestDuration(),
	}

	return &request, r.registrationRequestRepo.CreateRegistrationRequest(request, ctx)
}

// GetRegistrationRequest implements business.IRegistrationRequestService.
func (r *registrationRequestService) GetRegistrationRequest(id string, ctx context.Context) (*entities.RegistrationRequest, error) {
	if id == "" {
		return nil, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	return r.registrationRequestRepo.GetRegistrationRequest(id, ctx)
}

// GetRegistrationRequests implements business.IRegistrationRequestService.
func (r *registrationRequestService) GetRegistrationRequests(req request.GetRegistrationRequests, ctx context.Context) (response.PaginationDataResponse, error) {
	req.SortOrder = util.StanderizeSortOrder(req.SortOrder)
	if req.Page < 1 {
		req.Page = 1
	}

	if req.PageSize < 1 {
		req.PageSize = default_page_size
	}

	var res response.PaginationDataResponse
	var redisKey string = r.getGetRegistrationRequestsRedisKey(req)
	if r.redisCache.Get(redisKey, &res, ctx) {
		return res, nil
	}

	data, pages, err := r.registrationRequestRepo.GetRegistrationRequests(req, ctx)
	var amount int
	if data == nil || len(data) == 0 {
		amount = 0
	} else {
		amount = len(data)
	}

	res = response.PaginationDataResponse{
		Data:       data,
		Amount:     amount,
		Page:       req.Page,
		TotalPages: pages,
	}

	r.redisCache.Set(redisKey, res, time.Minute*5, ctx)

	return res, err
}

// GetWalletRegistrationRequests implements business.IRegistrationRequestService.
func (r *registrationRequestService) GetWalletRegistrationRequests(id string, ctx context.Context) ([]entities.RegistrationRequest, error) {
	if !utils.IsValidSuiAddress(models.SuiAddress(id)) {
		return nil, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	return r.registrationRequestRepo.GetWalletRegistrationRequests(id, ctx)
}

// VoteRegistrationRequest implements business.IRegistrationRequestService.
func (r *registrationRequestService) VoteRegistrationRequest(id string, req request.VoteRequest, ctx context.Context) error {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)

	var voter string = ctx.Value("address").(string)
	if !utils.IsValidSuiAddress(models.SuiAddress(voter)) {
		return genericErr
	}

	request, err := r.registrationRequestRepo.GetRegistrationRequest(id, ctx)
	if err != nil {
		return err
	}

	if request == nil {
		return genericErr
	}

	if request.ClosedAt.Before(time.Now()) {
		return errors.New(noti.REQUEST_CLOSED_MESSAGE)
	}

	if voter == request.CreatedBy {
		return errors.New(noti.OWNER_VOTE_WARN_MSG)
	}

	if slices.Contains(request.Approvers, voter) || slices.Contains(request.Refusers, voter) {
		return errors.New(noti.ALREADY_VOTE_MESSAGE)
	}

	manageObj, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    r.clients[constant.SuiTestnet],
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: r.errLogger,
	}, ctx)
	if err != nil {
		return err
	}

	// Not admins or local leaders
	if !slices.Contains(manageObj.AdminIds, voter) && !slices.Contains(manageObj.LocalLeaderIds, voter) {
		return errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	}

	if req.IsVoteYes {
		request.Approvers = append(request.Approvers, voter)
	} else {
		request.Refusers = append(request.Refusers, voter)
		if req.RefuseReason == "" {
			return errors.New(noti.FIELD_EMPTY_WARN_MSG)
		}

		request.RefuseReasons = append(request.RefuseReasons, strings.TrimSpace(req.RefuseReason))
	}

	request.UpdatedAt = time.Now()

	return r.registrationRequestRepo.UpdateRegistrationRequest(*request, ctx)
}

func (r *registrationRequestService) getGetRegistrationRequestsRedisKey(req request.GetRegistrationRequests) string {
	var role string = "empty"
	if req.RegisterRole != "" {
		role = req.RegisterRole
	}

	var isAvailable string = "empty"
	if req.IsAvailableToConfirm != nil {
		isAvailable = fmt.Sprintf("%v", *req.IsAvailableToConfirm)
	}

	var keyword string = "empty"
	if req.Keyword != "" {
		keyword = req.Keyword
	}

	var region string = "empty"
	if req.Region != "" {
		region = req.Region
	}

	var gender string = "empty"
	if req.Gender != "" {
		gender = req.Gender
	}

	var status string = "empty"
	if req.Status != "" {
		status = req.Status
	}

	var isClosed string = "empty"
	if req.IsClosed != nil {
		isClosed = fmt.Sprintf("%v", *req.IsClosed)
	}

	return fmt.Sprintf("registration_req:role:%s:available:%s:kw:%s:r:%s:g:%s:status:%s:closed:%s:o:%s:s:%d:p:%d",
		role, isAvailable, keyword, region, gender, status, isClosed, req.SortOrder, req.PageSize, req.Page)
}
