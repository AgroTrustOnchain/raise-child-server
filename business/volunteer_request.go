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

type volunteerRequestService struct {
	volunteerRequestRepo i_repository.IVolunteerRequestRepository
	profileRepo          i_repository.IProfileRepository
	clients              map[string]sui.ISuiAPI
	errLogger            *log.Logger
}

func InitializeVolunteerRequestService(db *sql.DB, errLogger *log.Logger) business.IVolunteerRequestService {
	return &volunteerRequestService{
		volunteerRequestRepo: repository.InitializeVolunteerRequestRepository(db, errLogger),
		profileRepo:          repository.InitializeProfileRepository(db, errLogger),
		clients:              _networkAliases,
		errLogger:            errLogger,
	}
}

func GenerateVolunteerRequestService() (business.IVolunteerRequestService, error) {
	var errLogger = util.GetLogConfig(shared.ERROR_LEVEL)

	cnn, err := db.ConnectDB(errLogger, db.InitializePostgreSQL())
	if err != nil {
		return nil, err
	}

	return InitializeVolunteerRequestService(cnn, errLogger), nil
}

// ConfirmRequest implements business.IVolunteerRequestService.
func (v *volunteerRequestService) ConfirmRequest(id string, ctx context.Context) (response.BuildTransactionResponse, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)

	var sender string = ctx.Value("address").(string)
	if !utils.IsValidSuiAddress(models.SuiAddress(sender)) {
		return response.BuildTransactionResponse{}, genericErr
	}

	req, err := v.volunteerRequestRepo.GetRequest(id, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	if req == nil {
		return response.BuildTransactionResponse{}, genericErr
	}

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
		if err := v.volunteerRequestRepo.UpdateRegistrationRequest(*req, ctx); err != nil {
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

	var client = v.clients[constant.SuiTestnet]
	var staffModule = on_chain.InitializeModuleStaff()
	txBytes, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
		Client:    client,
		Sender:    sender,
		Module:    staffModule.GetModule(),
		Function:  staffModule.GetFunctionRegisterVolunteer(),
		ErrLogger: v.errLogger,
		Arguments: staffModule.ToRegisterVolunteerArguments(on_chain.RegisterVolunteerArguments{
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
			},
		}),
	}, ctx)

	return response.BuildTransactionResponse{
		TxBytes: txBytes,
	}, err
}

// CreateRequest implements business.IVolunteerRequestService.
func (v *volunteerRequestService) CreateRequest(req request.VolunteerRegistrationRequest, ctx context.Context) (*entities.VolunteerRegistrationRequest, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)

	var sender string = ctx.Value("address").(string)
	if !utils.IsValidSuiAddress(models.SuiAddress(sender)) {
		return nil, genericErr
	}

	reqs, err := v.volunteerRequestRepo.GetWalletRegistrationRequests(sender, ctx)
	if err != nil {
		return nil, err
	}

	if len(reqs) > 0 {
		for _, req := range reqs {
			if req.Status == request_pending_status || req.Status == request_approved_status {
				return nil, genericErr
			}
		}
	}

	profile, err := v.profileRepo.GetProfile(ctx.Value("sub").(string), ctx)
	if err != nil {
		return nil, err
	}

	if profile == nil {
		return nil, genericErr
	}

	var curTime time.Time = time.Now()
	var request = entities.VolunteerRegistrationRequest{
		Region: req.Region,
		AdminRegistrationRequest: entities.AdminRegistrationRequest{
			ID:                   util.GenerateId(),
			IdentityCode:         util.StandardizeString(*profile.IdentityCode),
			IdentityCardBlobID:   strings.TrimSpace(req.IdentityCardBlobID),
			AvatarBlobID:         strings.TrimSpace(req.AvatarBlobID),
			FirstName:            strings.TrimSpace(*profile.FirstName),
			LastName:             strings.TrimSpace(*profile.LastName),
			Gender:               *profile.Gender,
			DateOfBirth:          *profile.DateOfBirth,
			PhoneNumber:          *profile.PhoneNumber,
			Email:                *profile.Email,
			Status:               request_pending_status,
			IsAvailableToConfirm: false,
			IsConfirmRegister:    false,
			CreatedBy:            sender,
			CreatedAt:            curTime,
			UpdatedAt:            curTime,
			ClosedAt:             util.GetRequestDuration(),
		},
	}

	return &request, v.volunteerRequestRepo.CreateRegistrationRequest(request, ctx)
}

// GetRequest implements business.IVolunteerRequestService.
func (v *volunteerRequestService) GetRequest(id string, ctx context.Context) (*entities.VolunteerRegistrationRequest, error) {
	if id == "" {
		return nil, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	return v.volunteerRequestRepo.GetRequest(id, ctx)
}

// GetRequests implements business.IVolunteerRequestService.
func (v *volunteerRequestService) GetRequests(req request.GetNormalStaffRegistrationRequests, ctx context.Context) (response.PaginationDataResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}

	data, pages, err := v.volunteerRequestRepo.GetRegistrationRequests(req, ctx)

	return response.PaginationDataResponse{
		Data:       data,
		Page:       req.Page,
		TotalPages: pages,
	}, err
}

// GetWalletRequests implements business.IVolunteerRequestService.
func (v *volunteerRequestService) GetWalletRequests(id string, ctx context.Context) ([]entities.VolunteerRegistrationRequest, error) {
	if !utils.IsValidSuiAddress(models.SuiAddress(id)) {
		return nil, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	return v.volunteerRequestRepo.GetWalletRegistrationRequests(id, ctx)
}

// VoteRequest implements business.IVolunteerRequestService.
func (v *volunteerRequestService) VoteRequest(id string, req request.VoteRequest, ctx context.Context) error {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)

	var voter string = ctx.Value("address").(string)
	if !utils.IsValidSuiAddress(models.SuiAddress(voter)) {
		return genericErr
	}

	request, err := v.volunteerRequestRepo.GetRequest(id, ctx)
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
		Client:    v.clients[constant.SuiTestnet],
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: v.errLogger,
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

	return v.volunteerRequestRepo.UpdateRegistrationRequest(*request, ctx)
}
