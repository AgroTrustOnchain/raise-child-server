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
	"raise-child/util/ai"
	"raise-child/util/cache"
	"raise-child/util/db"
	"raise-child/util/image/cloudinary"
	on_chain "raise-child/util/on_chain"
	walrus_pkg "raise-child/util/walrus_pkg"
	"slices"
	"strings"
	"time"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/sui"
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
	aiProvider              ai.IAiClientProvider
	walrusProvider          walrus_pkg.IWalrusProvider
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
	aiProvider ai.IAiClientProvider,
	walrusProvider walrus_pkg.IWalrusProvider,
	clients map[string]sui.ISuiAPI,
	errLogger *log.Logger,
) business.IRegistrationRequestService {
	return &registrationRequestService{
		registrationRequestRepo: registrationRequestRepo,
		profileRepo:             profileRepo,
		aiProvider:              aiProvider,
		walrusProvider:          walrusProvider,
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
		ai.InitializeAiProvider(errLogger),
		walrus_pkg.InitializeWalrusProvider(errLogger),
		_networkAliases,
		errLogger,
	), nil
}

// // ConfirmRegistrationRequest implements business.IRegistrationRequestService.
// func (r *registrationRequestService) ConfirmRegistrationRequest(id string, ctx context.Context) (response.BuildTransactionResponse, error) {
// 	req, err := r.registrationRequestRepo.GetRegistrationRequest(id, ctx)
// 	if err != nil {
// 		return response.BuildTransactionResponse{}, err
// 	}

// 	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
// 	if req == nil {
// 		return response.BuildTransactionResponse{}, genericErr
// 	}

// 	var sender string = ctx.Value("address").(string)
// 	if req.CreatedBy != sender {
// 		return response.BuildTransactionResponse{}, errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
// 	}

// 	// Pending process
// 	if req.ClosedAt.After(time.Now()) {
// 		return response.BuildTransactionResponse{}, errors.New(noti.STILL_PENDING_REQUEST_MESSAGE)
// 	} else { // Request closed
// 		var rate float32 = float32(len(req.Approvers)) / float32(len(req.Approvers)+len(req.Refusers))
// 		var isDenied bool = false

// 		if rate >= approve_rate_limit {
// 			req.Status = request_approved_status
// 			req.IsConfirmRegister = true
// 		} else {
// 			req.Status = request_refused_status
// 			isDenied = true
// 		}

// 		req.UpdatedAt = time.Now()
// 		if err := r.registrationRequestRepo.UpdateRegistrationRequest(*req, ctx); err != nil {
// 			return response.BuildTransactionResponse{}, err
// 		}

// 		if isDenied {
// 			return response.BuildTransactionResponse{}, nil
// 		}
// 	}

// 	// Wait for background server to mint cap object to register
// 	if !req.IsAvailableToConfirm {
// 		return response.BuildTransactionResponse{}, nil
// 	}

// 	var client = r.clients[constant.SuiTestnet]
// 	var manageModule = on_chain.InitializeModuleManage()
// 	var capType string

// 	switch req.RegisterRole {
// 	case admin_role:
// 		capType = manageModule.GetRegisterAdminCapStruct()
// 	case local_leader_role:
// 		capType = manageModule.GetRegisterLeaderCapStruct()
// 	case volunteer_role:
// 		capType = manageModule.GetRegisterVolunteerCapStruct()
// 	}

// 	r.errLogger.Println("Cap type:", capType)

// 	caps, err := on_chain.GetOnChainOwnedObjects[entities.Cap](on_chain.GetOnChainOwnedObjectsRequest{
// 		Client:       client,
// 		OwnerAddress: sender,
// 		StructType:   fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), manageModule.GetModule(), capType),
// 		ErrLogger:    r.errLogger,
// 	}, ctx)
// 	if err != nil {
// 		return response.BuildTransactionResponse{}, err
// 	}

// 	var args []interface{}
// 	var function string
// 	var staffModule = on_chain.InitializeModuleStaff()
// 	switch req.RegisterRole {
// 	case admin_role:
// 		function = staffModule.GetFunctionRegisterAdmin()
// 		args = staffModule.ToRegisterAdminArguments(on_chain.RegisterAdminArguments{
// 			CapID:              caps[0].ID.ID,
// 			IdentityCode:       req.IdentityCode,
// 			IdentityCardBlobID: req.IdentityCardBlobID,
// 			AvatarBlobID:       req.AvatarBlobID,
// 			FirstName:          req.FirstName,
// 			LastName:           req.LastName,
// 			Gender:             req.Gender,
// 			DateOfBirth:        req.DateOfBirth,
// 			PhoneNumber:        req.PhoneNumber,
// 			Email:              req.Email,
// 		})
// 	case local_leader_role:
// 		pool, err := on_chain.GetOnChainObject[entities.MainPool](on_chain.GetOnChainObjectRequest{
// 			Client:    client,
// 			ObjectId:  os.Getenv(env.POOL_ID),
// 			ErrLogger: r.errLogger,
// 		}, ctx)
// 		if err != nil {
// 			return response.BuildTransactionResponse{}, err
// 		}

// 		var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)
// 		if pool == nil {
// 			return response.BuildTransactionResponse{}, internalErr
// 		}

// 		localPools, err := on_chain.GetOnChainObjects[entities.LocalPool](on_chain.GetOnChainObjectsRequest{
// 			Client:    client,
// 			ObjectIds: pool.LocalPools,
// 			ErrLogger: r.errLogger,
// 		}, ctx)
// 		if err != nil {
// 			return response.BuildTransactionResponse{}, err
// 		}

// 		if localPools == nil || len(localPools) == 0 {
// 			return response.BuildTransactionResponse{}, internalErr
// 		}

// 		var localPoolId string
// 		for _, localPool := range localPools {
// 			if localPool.Region == req.Region {
// 				localPoolId = localPool.ID.ID
// 				break
// 			}
// 		}

// 		if localPoolId == "" {
// 			localPoolId = os.Getenv(env.SHARED_LOCAL_POOL_ID)
// 		}

// 		function = staffModule.GetFunctionRegisterLeader()
// 		args = staffModule.ToRegisterNormalStaffArguments(on_chain.RegisterNormalStaffArguments{
// 			LocalPoolID: localPoolId,
// 			Region:      req.Region,
// 			RegisterAdminArguments: on_chain.RegisterAdminArguments{
// 				CapID:              caps[0].ID.ID,
// 				IdentityCode:       req.IdentityCode,
// 				IdentityCardBlobID: req.IdentityCardBlobID,
// 				AvatarBlobID:       req.AvatarBlobID,
// 				FirstName:          req.FirstName,
// 				LastName:           req.LastName,
// 				Gender:             req.Gender,
// 				DateOfBirth:        req.DateOfBirth,
// 				PhoneNumber:        req.PhoneNumber,
// 				Email:              req.Email,
// 			},
// 		})
// 	case volunteer_role:
// 		function = staffModule.GetFunctionRegisterVolunteer()
// 		args = staffModule.ToRegisterNormalStaffArguments(on_chain.RegisterNormalStaffArguments{
// 			Region: req.Region,
// 			RegisterAdminArguments: on_chain.RegisterAdminArguments{
// 				CapID:              caps[0].ID.ID,
// 				IdentityCode:       req.IdentityCode,
// 				IdentityCardBlobID: req.IdentityCardBlobID,
// 				AvatarBlobID:       req.AvatarBlobID,
// 				FirstName:          req.FirstName,
// 				LastName:           req.LastName,
// 				Gender:             req.Gender,
// 				DateOfBirth:        req.DateOfBirth,
// 				PhoneNumber:        req.PhoneNumber,
// 				Email:              req.Email,
// 			},
// 		})
// 	}

// 	var module = on_chain.InitializeModuleStaff()
// 	txBytes, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
// 		Client:    r.clients[constant.SuiTestnet],
// 		Sender:    sender,
// 		Module:    module.GetModule(),
// 		Function:  function,
// 		ErrLogger: r.errLogger,
// 		Arguments: args,
// 	}, ctx)

// 	return response.BuildTransactionResponse{
// 		TxBytes:         txBytes,
// 		RegistrationReq: id,
// 	}, err
// }

// ConfirmRegistrationRequest implements business.IRegistrationRequestService.
func (r *registrationRequestService) ConfirmRegistrationRequest(id string, ctx context.Context) error {
	req, err := r.registrationRequestRepo.GetRegistrationRequest(id, ctx)
	if err != nil {
		return err
	}

	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if req == nil {
		return genericErr
	}

	var sender string = ctx.Value("address").(string)
	if req.CreatedBy != sender {
		return errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	}

	// Pending process
	var curTime time.Time = time.Now()
	if req.ClosedAt.After(curTime) {
		return errors.New(noti.STILL_PENDING_REQUEST_MESSAGE)
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

		req.UpdatedAt = curTime
		if err := r.registrationRequestRepo.UpdateRegistrationRequest(*req, ctx); err != nil {
			return err
		}

		if isDenied {
			return nil
		}
	}

	var args []interface{}
	var function string
	var staffModule = on_chain.InitializeModuleStaff()
	var client = r.clients[constant.SuiTestnet]
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)
	switch req.RegisterRole {
	case admin_role:
		function = staffModule.GetFunctionRegisterAdmin()
		args = staffModule.ToRegisterAdminArguments(on_chain.RegisterAdminArguments{
			IdentityCode:       req.IdentityCode,
			IdentityCardBlobID: req.IdentityCardBlobID,
			AvatarBlobID:       req.AvatarBlobID,
			FirstName:          req.FirstName,
			LastName:           req.LastName,
			Gender:             req.Gender,
			DateOfBirth:        req.DateOfBirth,
			PhoneNumber:        req.PhoneNumber,
			Email:              req.Email,
			Owner:              req.CreatedBy,
		})
	case local_leader_role:
		pool, err := on_chain.GetOnChainObject[entities.MainPool](on_chain.GetOnChainObjectRequest{
			Client:    client,
			ObjectId:  os.Getenv(env.POOL_ID),
			ErrLogger: r.errLogger,
		}, ctx)
		if err != nil {
			return err
		}

		if pool == nil {
			return internalErr
		}

		localPools, err := on_chain.GetOnChainObjects[entities.LocalPool](on_chain.GetOnChainObjectsRequest{
			Client:    client,
			ObjectIds: pool.LocalPools,
			ErrLogger: r.errLogger,
		}, ctx)
		if err != nil {
			return err
		}

		if len(localPools) == 0 {
			return internalErr
		}

		var localPoolId string
		for _, localPool := range localPools {
			if localPool.Region == req.Region {
				localPoolId = localPool.ID.ID
				break
			}
		}

		if localPoolId == "" {
			localPoolId = os.Getenv(env.SHARED_LOCAL_POOL_ID)
		}

		function = staffModule.GetFunctionRegisterLeader()
		args = staffModule.ToRegisterNormalStaffArguments(on_chain.RegisterNormalStaffArguments{
			LocalPoolID: localPoolId,
			Region:      req.Region,
			RegisterAdminArguments: on_chain.RegisterAdminArguments{
				IdentityCode:       req.IdentityCode,
				IdentityCardBlobID: req.IdentityCardBlobID,
				AvatarBlobID:       req.AvatarBlobID,
				FirstName:          req.FirstName,
				LastName:           req.LastName,
				Gender:             req.Gender,
				DateOfBirth:        req.DateOfBirth,
				PhoneNumber:        req.PhoneNumber,
				Email:              req.Email,
				Owner:              req.CreatedBy,
			},
		})
	case volunteer_role:
		function = staffModule.GetFunctionRegisterVolunteer()
		args = staffModule.ToRegisterNormalStaffArguments(on_chain.RegisterNormalStaffArguments{
			Region: req.Region,
			RegisterAdminArguments: on_chain.RegisterAdminArguments{
				IdentityCode:       req.IdentityCode,
				IdentityCardBlobID: req.IdentityCardBlobID,
				AvatarBlobID:       req.AvatarBlobID,
				FirstName:          req.FirstName,
				LastName:           req.LastName,
				Gender:             req.Gender,
				DateOfBirth:        req.DateOfBirth,
				PhoneNumber:        req.PhoneNumber,
				Email:              req.Email,
				Owner:              req.CreatedBy,
			},
		})
	}

	for i := 1; i <= 3; i++ {
		if _, err := on_chain.ExecuteTransactionV2(on_chain.ExecuteTransactionRequestV2{
			Client:    client,
			Module:    staffModule.GetModule(),
			Function:  function,
			Arguments: args,
			ErrLogger: r.errLogger,
		}, ctx); err == nil {
			return nil
		}
	}

	return internalErr
}

// CreateRegistrationRequest implements business.IRegistrationRequestService.
func (r *registrationRequestService) CreateRegistrationRequest(req request.CreateRegistrationRequest, ctx context.Context) (*entities.RegistrationRequest, error) {
	profile, err := r.profileRepo.GetProfile(ctx.Value("sub").(string), ctx)
	if err != nil {
		return nil, err
	}

	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if profile == nil {
		r.errLogger.Println("Error profile nil!!!")
		return nil, genericErr
	}

	if profile.IdentityCode == nil {
		return nil, errors.New(noti.PROFILE_EMPTY_MESSAGE)
	}

	var client = r.clients[constant.SuiTestnet]
	manage, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: r.errLogger,
	}, ctx)
	if err != nil {
		return nil, err
	}

	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)
	if manage == nil {
		return nil, internalErr
	}

	var sender string = ctx.Value("address").(string)
	var role string = strings.TrimSpace(req.RegisterRole)
	if req.RegisterRole != admin_role {
		var searchLen int
		if len(manage.VolunteerIds) > len(manage.LocalLeaderIds) {
			searchLen = len(manage.VolunteerIds)
		} else {
			searchLen = len(manage.LocalLeaderIds)
		}

		var staffIds []string
		for i := 0; i < searchLen; i++ {
			if i < len(manage.VolunteerIds) {
				if sender == manage.VolunteerIds[i] {
					staffIds = append(staffIds, manage.VolunteerNfts[i])

					if len(staffIds) == 2 {
						break
					}
				}
			}

			if i < len(manage.LocalLeaderIds) {
				if sender == manage.LocalLeaderIds[i] {
					staffIds = append(staffIds, manage.LocalLeaderNfts[i])

					if len(staffIds) == 2 {
						break
					}
				}
			}
		}

		if len(staffIds) > 0 {
			if len(staffIds) == 2 {
				return nil, errors.New(noti.ALREADY_STAFF_ROLE_MESSAGE)
			}

			nft, err := on_chain.GetOnChainObject[entities.StaffNft](on_chain.GetOnChainObjectRequest{
				Client:    client,
				ObjectId:  staffIds[0],
				ErrLogger: r.errLogger,
			}, ctx)
			if err != nil {
				return nil, err
			}

			if nft == nil {
				return nil, internalErr
			}

			if nft.Region != req.Region {
				return nil, errors.New(noti.ALREADY_ANOTHER_REGION_STAFF_MESSAGE)
			}

			if nft.Role == role {
				return nil, errors.New(noti.ALREADY_STAFF_ROLE_MESSAGE)
			}
		}
	} else {
		if slices.Contains(manage.AdminIds, sender) {
			return nil, errors.New(noti.ALREADY_STAFF_ROLE_MESSAGE)
		}
	}

	reqs, err := r.registrationRequestRepo.GetWalletRegistrationRequests(sender, ctx)
	if err != nil {
		return nil, err
	}

	var curTime time.Time = time.Now()
	if len(reqs) > 0 {
		for _, req := range reqs {
			if req.RegisterRole == role && (req.Status == request_pending_status || req.Status == request_approved_status) && curTime.Before(req.ClosedAt) {
				r.errLogger.Println("Error previous reqs!!!")
				return nil, genericErr
			}
		}
	}

	// Admin registration request not contain region
	if role == admin_role {
		if req.Region != "" {
			r.errLogger.Println("Error Admin Has Region!!!")
			return nil, genericErr
		}
	}

	if role == volunteer_role {
		if !slices.Contains(manage.LocalRegions, req.Region) {
			return nil, errors.New(noti.REGION_NOT_HAVE_LEADER_MESSAGE)
		}
	}

	var identityCode string = util.StandardizeString(*profile.IdentityCode)
	var firstName string = strings.TrimSpace(*profile.FirstName)
	var lastName string = strings.TrimSpace(*profile.LastName)
	var request = entities.RegistrationRequest{
		ID:                 util.GenerateId(),
		ProfileID:          ctx.Value("sub").(string),
		RegisterRole:       role,
		IdentityCode:       identityCode,
		IdentityCardBlobID: strings.TrimSpace(req.IdentityCardCloudinaryID),
		AvatarBlobID:       strings.TrimSpace(req.AvatarBlobID),
		Region:             req.Region,
		FirstName:          firstName,
		LastName:           lastName,
		Gender:             *profile.Gender,
		DateOfBirth:        *profile.DateOfBirth,
		PhoneNumber:        *profile.PhoneNumber,
		Email:              *profile.Email,
		Status:             request_pending_status,
		CreatedBy:          sender,
		CreatedAt:          curTime,
		UpdatedAt:          curTime,
		ClosedAt:           util.GetRequestDuration(),
	}

	return &request, r.registrationRequestRepo.CreateRegistrationRequest(request, ctx)
}

// GetRegistrationRequest implements business.IRegistrationRequestService.
func (r *registrationRequestService) GetRegistrationRequest(id string, ctx context.Context) (response.RegistrationRequestResponse, error) {
	if id == "" {
		return response.RegistrationRequestResponse{}, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	req, err := r.registrationRequestRepo.GetRegistrationRequest(id, ctx)
	if err != nil {
		return response.RegistrationRequestResponse{}, err
	}

	var res response.RegistrationRequestResponse = req.ToRegistrationRequestResponse()
	res.IdentityCardImgUrl = cloudinary.GetImageUrl(req.IdentityCardBlobID)

	return res, nil
}

// GetRegistrationRequests implements business.IRegistrationRequestService.
func (r *registrationRequestService) GetRegistrationRequests(req request.GetRegistrationRequests, ctx context.Context) (response.PaginationDataResponse, error) {
	req.SortOrder = util.StandardizeSortOrder(req.SortOrder)
	if req.Page < 1 {
		req.Page = 1
	}

	if req.PageSize < 1 {
		req.PageSize = default_page_size
	}

	var res response.PaginationDataResponse
	data, pages, err := r.registrationRequestRepo.GetRegistrationRequests(req, ctx)
	var amount int
	if len(data) == 0 {
		amount = 0
	} else {
		amount = len(data)
	}

	var resData []response.RegistrationRequestResponse
	if amount > 0 {
		for _, req := range data {
			var reqRes = req.ToRegistrationRequestResponse()
			reqRes.IdentityCardImgUrl = cloudinary.GetImageUrl(req.IdentityCardBlobID)
			resData = append(resData, reqRes)
		}
	}

	res = response.PaginationDataResponse{
		Data:       resData,
		Amount:     amount,
		Page:       req.Page,
		TotalPages: pages,
	}

	return res, err
}

// GetWalletRegistrationRequests implements business.IRegistrationRequestService.
func (r *registrationRequestService) GetWalletRegistrationRequests(id string, ctx context.Context) ([]response.RegistrationRequestResponse, error) {
	if !util.IsValidSuiAddressStrict(id) {
		return nil, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	data, err := r.registrationRequestRepo.GetWalletRegistrationRequests(id, ctx)
	if err != nil {
		return nil, err
	}

	var res []response.RegistrationRequestResponse
	if len(data) > 0 {
		for _, req := range data {
			var reqRes = req.ToRegistrationRequestResponse()
			reqRes.IdentityCardImgUrl = cloudinary.GetImageUrl(req.IdentityCardBlobID)
			res = append(res, reqRes)
		}
	}

	return res, nil
}

// VoteRegistrationRequest implements business.IRegistrationRequestService.
func (r *registrationRequestService) VoteRegistrationRequest(id string, req request.VoteRequest, ctx context.Context) error {
	request, err := r.registrationRequestRepo.GetRegistrationRequest(id, ctx)
	if err != nil {
		return err
	}

	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if request == nil {
		return genericErr
	}

	if request.ClosedAt.Before(time.Now()) {
		return errors.New(noti.REQUEST_CLOSED_MESSAGE)
	}

	var voter string = ctx.Value("address").(string)
	if voter == request.CreatedBy {
		return errors.New(noti.OWNER_VOTE_WARN_MSG)
	}

	if slices.Contains(request.Approvers, voter) || slices.Contains(request.Refusers, voter) {
		return errors.New(noti.ALREADY_VOTE_MESSAGE)
	}

	var client = r.clients[constant.SuiTestnet]
	manageObj, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: r.errLogger,
	}, ctx)
	if err != nil {
		return err
	}

	// Not admin
	if !slices.Contains(manageObj.AdminIds, voter) {
		var genericRightErr error = errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
		switch request.RegisterRole {
		case admin_role:
			return genericRightErr
		case local_leader_role:
			if !slices.Contains(manageObj.LocalLeaderIds, voter) {
				return genericRightErr
			}
		default:
			var foundIdx int = -1
			for i, leader := range manageObj.LocalLeaderIds {
				if leader == voter {
					foundIdx = i
					break
				}
			}

			if foundIdx == -1 {
				return genericRightErr
			}

			staff, err := on_chain.GetOnChainObject[entities.StaffNft](on_chain.GetOnChainObjectRequest{
				Client:    client,
				ObjectId:  manageObj.LocalLeaderNfts[foundIdx],
				ErrLogger: r.errLogger,
			}, ctx)
			if err != nil {
				return err
			}

			if staff.Region != request.Region {
				return genericRightErr
			}
		}
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

	return fmt.Sprintf("registration_req:role:%s:kw:%s:r:%s:g:%s:status:%s:closed:%s:o:%s:s:%d:p:%d",
		role, keyword, region, gender, status, isClosed, req.SortOrder, req.PageSize, req.Page)
}
