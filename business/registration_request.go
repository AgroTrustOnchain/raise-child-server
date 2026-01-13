package business

import (
	"context"
	"database/sql"
	"errors"
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
	local_leader_role string = "Local Leader"
	volunteer_role    string = "Volunteer"
)

const (
	request_pending_status  string = "Pending"
	request_refused_status  string = "Refused"
	request_approved_status string = "Aprroved"
)

const (
	approve_rate_limit float32 = 0.7 // 70%
)

type registrationRequestService struct {
	registrationRequestRepo i_repository.IRegistrationRequestRepository
	profileRepo             i_repository.IProfileRepository
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

func GenerateRegistrationRequestService() (business.IRegistrationRequestService, error) {
	var errLogger = util.GetLogConfig(shared.ERROR_LEVEL)

	cnn, err := db.ConnectDB(errLogger, db.InitializePostgreSQL())
	if err != nil {
		return nil, err
	}

	return InitializeRegistrationRequestService(cnn, errLogger), nil
}

// ConfirmRegistrationRequest implements business.IRegistrationRequestService.
func (r *registrationRequestService) ConfirmRegistrationRequest(id string, ctx context.Context) (response.BuildTransactionResponse, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)

	var sender string = ctx.Value("address").(string)
	if !utils.IsValidSuiAddress(models.SuiAddress(sender)) {
		return response.BuildTransactionResponse{}, genericErr
	}

	req, err := r.registrationRequestRepo.GetRegistrationRequest(id, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	if req == nil {
		return response.BuildTransactionResponse{}, genericErr
	}

	if req.CreatedBy != sender {
		return response.BuildTransactionResponse{}, errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	}

	profile, err := r.profileRepo.GetProfile(ctx.Value("sub").(string), ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	if profile == nil {
		return response.BuildTransactionResponse{}, genericErr
	}

	// Pending preocess
	var isHaveToUpdateRequest bool = false
	if req.ClosedAt.Before(time.Now()) {
		return response.BuildTransactionResponse{}, errors.New(noti.STILL_PENDING_REQUEST_MESSAGE)
	} else { // Request closed
		var rate float32 = float32(len(req.Aprrovers)) / float32(len(req.Aprrovers)+len(req.Refusers))
		if rate >= approve_rate_limit {
			req.Status = request_approved_status
			req.IsConfirmRegister = true
		} else {
			req.Status = request_refused_status
		}

		isHaveToUpdateRequest = true
		req.UpdatedAt = time.Now()
	}

	if isHaveToUpdateRequest {
		if err := r.registrationRequestRepo.UpdateRegistrationRequest(*req, ctx); err != nil {
			return response.BuildTransactionResponse{}, err
		}
	}

	var res response.BuildTransactionResponse
	var errRes error
	if req.IsConfirmRegister {
		var module = on_chain.InitializeModuleStaff()
		txBytes, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
			Client:    r.clients[constant.SuiTestnet],
			Sender:    sender,
			Module:    module.GetModule(),
			Function:  module.GetFunctionRegisterStaff(),
			ErrLogger: r.errLogger,
			Arguments: module.ToRegisterStaffArguements(on_chain.RegisterStaffArguements{
				IdentityCode:       profile.IdentityCode,
				IdentityCardBlobID: req.IdentityCardBlobID,
				Role:               req.RegisterRole,
				Region:             profile.Region,
				FirstName:          profile.FirstName,
				LastName:           profile.LastName,
				Gender:             profile.Gender,
				PhoneNumber:        profile.PhoneNumber,
				Email:              profile.Email,
			}),
		}, ctx)

		res.TxBytes = txBytes
		errRes = err
	}

	return res, errRes
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

	if role == volunteer_role {
		var client = r.clients[constant.SuiTestnet]
		manageObj, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
			Client:    client,
			ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
			ErrLogger: r.errLogger,
		}, ctx)
		if err != nil {
			return nil, err
		}

		leaders, err := on_chain.GetOnChainObjects[entities.Staff](on_chain.GetOnChainObjectsRequest{
			Client:    client,
			ObjectIds: manageObj.LocalLeaderIds,
			ErrLogger: r.errLogger,
		}, ctx)
		if err != nil {
			return nil, err
		}

		var isRegionAvailable bool = false
		for _, leader := range leaders {
			if leader.Region == profile.Region {
				isRegionAvailable = true
				break
			}
		}

		if !isRegionAvailable {
			return nil, errors.New(noti.REGION_NOT_ADDED_WARN_MSG)
		}
	}

	// todo: validate identity code
	var curTime time.Time = time.Now()
	var request = entities.RegistrationRequest{
		ID:                 util.GenerateId(),
		RegisterRole:       role,
		IdentityCode:       util.StanderizeString(profile.IdentityCode),
		IdentityCardBlobID: strings.TrimSpace(req.IdentityCardBlobID),
		AvatarBlobID:       strings.TrimSpace(req.AvatarBlobID),
		Region:             profile.Region,
		FirstName:          strings.TrimSpace(profile.FirstName),
		LastName:           strings.TrimSpace(profile.LastName),
		Gender:             profile.Gender,
		DateOfBirth:        profile.DateOfBirth,
		PhoneNumber:        profile.PhoneNumber,
		Email:              profile.Email,
		Status:             request_pending_status,
		IsConfirmRegister:  false,
		CreatedBy:          sender,
		CreatedAt:          curTime,
		UpdatedAt:          curTime,
		ClosedAt:           util.GetRequestDuration(),
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
	if req.Page <= 0 {
		req.Page = 1
	}

	data, pages, err := r.registrationRequestRepo.GetRegistrationRequests(req, ctx)

	return response.PaginationDataResponse{
		Data:       data,
		Page:       req.Page,
		TotalPages: pages,
	}, err
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

	if slices.Contains(request.Aprrovers, voter) || slices.Contains(request.Refusers, voter) {
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
		request.Aprrovers = append(request.Aprrovers, voter)
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
