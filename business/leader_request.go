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
	panic("unimplemented")
}

// CreateRequest implements business.ILocalLeaderRequestService.
func (l *leaderRequestService) CreateRequest(req request.CreateRegistrationRequest, ctx context.Context) (*entities.LocalLeaderRegistrationRequest, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)

	var sender string = ctx.Value("address").(string)
	if !utils.IsValidSuiAddress(models.SuiAddress(sender)) {
		return nil, genericErr
	}

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
		// 	IdentityCode:         util.StanderizeString(profile.IdentityCode),
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
	if !utils.IsValidSuiAddress(models.SuiAddress(id)) {
		return nil, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	return l.leaderRequestRepo.GetWalletRegistrationRequests(id, ctx)
}

// VoteRequest implements business.ILocalLeaderRequestService.
func (l *leaderRequestService) VoteRequest(id string, req request.VoteRequest, ctx context.Context) error {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)

	var voter string = ctx.Value("address").(string)
	if !utils.IsValidSuiAddress(models.SuiAddress(voter)) {
		return genericErr
	}

	request, err := l.leaderRequestRepo.GetRequest(id, ctx)
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
