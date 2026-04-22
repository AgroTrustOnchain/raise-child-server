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
	"raise-child/util/db"
	on_chain "raise-child/util/on_chain"
	"slices"
	"strings"
	"time"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/sui"
)

type leaderRequestService struct {
	leaderRequestRepo i_repository.ILocalLeaderRequestRepository
	profileRepo       i_repository.IProfileRepository
	clients           map[string]sui.ISuiAPI
	errLogger         *log.Logger
}

func InitializeLocalLeaderRequestService(db *sql.DB, errLogger *log.Logger) business.ILocalLeaderRequestService {
	return &leaderRequestService{
		leaderRequestRepo: repository.InitializeLocalLeaderRequestRepository(db, errLogger),
		profileRepo:       repository.InitializeProfileRepository(db, errLogger),
		clients:           _networkAliases,
		errLogger:         errLogger,
	}
}

// GenerateLocalLeaderRequestService generates the local leader request service.
func GenerateLocalLeaderRequestService() (business.ILocalLeaderRequestService, error) {
	var errLogger = util.GetLogConfig(shared.ERROR_LEVEL)

	cnn, err := db.ConnectDB(errLogger, db.InitializePostgreSQL())
	if err != nil {
		return nil, err
	}

	return InitializeLocalLeaderRequestService(cnn, errLogger), nil
}

// ConfirmRequest implements business.ILocalLeaderRequestService.
func (l *leaderRequestService) ConfirmRequest(id string, ctx context.Context) (response.BuildTransactionResponse, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)

	var sender string = ctx.Value("address").(string)
	if !util.IsValidSuiAddressStrict(sender) {
		return response.BuildTransactionResponse{}, genericErr
	}

	req, err := l.leaderRequestRepo.GetRequest(id, ctx)
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
		if err := l.leaderRequestRepo.UpdateRegistrationRequest(*req, ctx); err != nil {
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

	var client = l.clients[constant.SuiTestnet]
	var mangeModule = on_chain.InitializeModuleManage()
	caps, err := on_chain.GetOnChainOwnedObjects[entities.Cap](on_chain.GetOnChainOwnedObjectsRequest{
		Client:       client,
		OwnerAddress: sender,
		StructType:   fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), mangeModule.GetModule(), mangeModule.GetRegisterLeaderCapStruct()),
		ErrLogger:    l.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	var staffModule = on_chain.InitializeModuleStaff()
	txBytes, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
		Client:    client,
		Sender:    sender,
		Module:    staffModule.GetModule(),
		Function:  staffModule.GetFunctionRegisterLeader(),
		ErrLogger: l.errLogger,
		Arguments: staffModule.ToRegisterLeaderArguments(on_chain.RegisterLeaderArguments{
			CenterAddress:     req.CenterAddress,
			CenterPhoneNumber: req.CenterPhoneNumber,
			CenterImageBlobID: req.CenterImageBlobID,
			RegisterVolunteerArguments: on_chain.RegisterVolunteerArguments{
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
			},
		}),
	}, ctx)

	return response.BuildTransactionResponse{
		TxBytes: txBytes,
	}, err
}

// CreateRequest implements business.ILocalLeaderRequestService.
func (l *leaderRequestService) CreateRequest(req request.CreateRegistrationRequest, ctx context.Context) (*entities.LocalLeaderRegistrationRequest, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	var sender string = ctx.Value("address").(string)
	reqs, err := l.leaderRequestRepo.GetWalletRegistrationRequests(sender, ctx)
	if err != nil {
		return nil, err
	}

	if reqs != nil && len(reqs) > 0 {
		for _, req := range reqs {
			if req.Status == request_pending_status || req.Status == request_approved_status {
				return nil, genericErr
			}
		}
	}

	profile, err := l.profileRepo.GetProfile(ctx.Value("sub").(string), ctx)
	if err != nil {
		return nil, err
	}

	if profile == nil {
		return nil, genericErr
	}

	//var curTime time.Time = time.Now()
	var request = entities.LocalLeaderRegistrationRequest{
		CenterAddress: "",
		// AdminRegistrationRequest: entities.AdminRegistrationRequest{
		// 	ID:                   util.GenerateId(),
		// 	IdentityCode:         util.StandardizeString(profile.IdentityCode),
		// 	IdentityCardBlobID:   strings.TrimSpace(req.IdentityCardBlobID),
		// 	AvatarBlobID:         strings.TrimSpace(req.AvatarBlobID),
		// 	FirstName:            strings.TrimSpace(profile.FirstName),
		// 	LastName:             strings.TrimSpace(profile.LastName),
		// 	Gender:               profile.Gender,
		// 	DateOfBirth:          profile.DateOfBirth,
		// 	PhoneNumber:          profile.PhoneNumber,
		// 	Email:                profile.Email,
		// 	Status:               request_pending_status,
		// 	IsAvailableToConfirm: false,
		// 	IsConfirmRegister:    false,
		// 	CreatedBy:            sender,
		// 	CreatedAt:            curTime,
		// 	UpdatedAt:            curTime,
		// 	ClosedAt:             util.GetRequestDuration(),
		// },
	}

	return &request, l.leaderRequestRepo.CreateRegistrationRequest(request, ctx)
}

// GetRequest implements business.ILocalLeaderRequestService.
func (l *leaderRequestService) GetRequest(id string, ctx context.Context) (*entities.LocalLeaderRegistrationRequest, error) {
	if id == "" {
		return nil, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	return l.leaderRequestRepo.GetRequest(id, ctx)
}

// GetRequests implements business.ILocalLeaderRequestService.
func (l *leaderRequestService) GetRequests(req request.GetNormalStaffRegistrationRequests, ctx context.Context) (response.PaginationDataResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}

	data, pages, err := l.leaderRequestRepo.GetRegistrationRequests(req, ctx)

	return response.PaginationDataResponse{
		Data:       data,
		Page:       req.Page,
		TotalPages: pages,
	}, err
}

// GetWalletRequests implements business.ILocalLeaderRequestService.
func (l *leaderRequestService) GetWalletRequests(id string, ctx context.Context) ([]entities.LocalLeaderRegistrationRequest, error) {
	if !util.IsValidSuiAddressStrict(id) {
		return nil, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	return l.leaderRequestRepo.GetWalletRegistrationRequests(id, ctx)
}

// VoteRequest implements business.ILocalLeaderRequestService.
func (l *leaderRequestService) VoteRequest(id string, req request.VoteRequest, ctx context.Context) error {
	request, err := l.leaderRequestRepo.GetRequest(id, ctx)
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

	manageObj, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    l.clients[constant.SuiTestnet],
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: l.errLogger,
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

	return l.leaderRequestRepo.UpdateRegistrationRequest(*request, ctx)
}
