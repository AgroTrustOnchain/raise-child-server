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

type uploadChildRequestService struct {
	uploadChildRequestRepo i_repository.IUploadChildRequestRepository
	clients                map[string]sui.ISuiAPI
	errLogger              *log.Logger
}

func InitializeUploadChildRequestService(db *sql.DB, errLogger *log.Logger) business.IUploadChildRequestService {
	return &uploadChildRequestService{
		uploadChildRequestRepo: repository.InitializeUploadChildRequestRepo(db, errLogger),
		clients:                _networkAliases,
		errLogger:              errLogger,
	}
}

func GenerateUploadChildRequestService() (business.IUploadChildRequestService, error) {
	var errLogger = util.GetLogConfig(shared.ERROR_LEVEL)

	cnn, err := db.ConnectDB(errLogger, db.InitializePostgreSQL())
	if err != nil {
		return nil, err
	}

	return InitializeUploadChildRequestService(cnn, errLogger), nil
}

// ConfirmUploadChildRequest implements business.IUploadChildRequestService.
func (u *uploadChildRequestService) ConfirmUploadChildRequest(id string, ctx context.Context) (response.BuildTransactionResponse, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)

	var sender string = ctx.Value("address").(string)
	if !utils.IsValidSuiAddress(models.SuiAddress(sender)) {
		return response.BuildTransactionResponse{}, genericErr
	}

	req, err := u.uploadChildRequestRepo.GetUploadChildRequest(id, ctx)
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
			req.IsConfirmUpload = true
		} else {
			req.Status = request_refused_status
			isDenied = true
		}

		req.UpdatedAt = time.Now()
		if err := u.uploadChildRequestRepo.UpdateUploadChildRequest(*req, ctx); err != nil {
			return response.BuildTransactionResponse{}, err
		}

		if isDenied {
			return response.BuildTransactionResponse{}, nil
		}
	}

	var res response.BuildTransactionResponse
	var errRes error
	if req.IsConfirmUpload {
		var module = on_chain.InitializeModuleChild()
		txBytes, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
			Client:    u.clients[constant.SuiTestnet],
			Sender:    sender,
			Module:    module.GetModule(),
			Function:  module.GetFunctionAddChild(),
			ErrLogger: u.errLogger,
			Arguments: module.ToAddChildArguments(on_chain.AddChildArguments{
				IdentityCode: req.IdentityCode,
				FirstName:    req.FirstName,
				LastName:     req.LastName,
				Gender:       req.Gender,
				DateOfBirth:  req.DateOfBirth,
				AvatarBlobId: req.AvatarBlobId,
			}),
		}, ctx)

		res.TxBytes = txBytes
		errRes = err
	}

	return res, errRes
}

// CreateUploadChildRequest implements business.IUploadChildRequestService.
func (u *uploadChildRequestService) CreateUploadChildRequest(req request.UploadChildRequest, ctx context.Context) (*entities.UploadChildRequest, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)

	var sender string = ctx.Value("address").(string)
	if !utils.IsValidSuiAddress(models.SuiAddress(sender)) {
		return nil, genericErr
	}

	var identityCode string = strings.TrimSpace(req.IdentityCode)
	isRequested, err := u.uploadChildRequestRepo.IsChildRequested(identityCode, ctx)
	if err != nil {
		return nil, err
	}

	if isRequested {
		return nil, errors.New(noti.CHILD_STILL_REQUESTED_MESSAGE)
	}

	manageObj, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    u.clients[constant.SuiTestnet],
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: u.errLogger,
	}, ctx)
	if err != nil {
		return nil, err
	}

	var region string = strings.TrimSpace(req.Region)
	var isRegionAvailable bool = false
	for i, addedRegion := range manageObj.LocalRegions {
		if addedRegion == region && manageObj.CenterConfirmStatuses[i] {
			isRegionAvailable = true
			break
		}
	}

	if !isRegionAvailable {
		return nil, errors.New(noti.REGION_NOT_ADDED_WARN_MSG)
	}

	var gender string = util.StanderizeGender(util.StanderizeString(req.Gender))
	if gender == "" {
		return nil, errors.New(noti.UNDEFINED_GENDER_MESSAGE)
	}

	var dateOfBirth string = strings.TrimSpace(req.DateOfBirth)
	if dob := util.RawDateToTime(dateOfBirth); dob.IsZero() {
		return nil, errors.New(noti.INVALID_DATE_FORMAT_WARN_MSG)
	} else {
		if !isChildAgeInSupport(dob.Year()) {
			return nil, errors.New(noti.CHILD_AGE_OUT_OF_SUPPORT_MESSAGE)
		}
	}

	var curTime time.Time = time.Now()
	var request = entities.UploadChildRequest{
		ID:           util.GenerateId(),
		IdentityCode: identityCode,
		AvatarBlobId: strings.TrimSpace(req.AvatarBlobId),
		Region:       region,
		FirstName:    strings.TrimSpace(req.FirstName),
		LastName:     strings.TrimSpace(req.LastName),
		Gender:       gender,
		DateOfBirth:  dateOfBirth,
		Status:       request_pending_status,
		CreatedBy:    sender,
		CreatedAt:    curTime,
		UpdatedAt:    curTime,
		ClosedAt:     util.GetRequestDuration(),
	}

	return &request, u.uploadChildRequestRepo.CreateUploadChildRequest(request, ctx)
}

// GetUploadChildRequest implements business.IUploadChildRequestService.
func (u *uploadChildRequestService) GetUploadChildRequest(id string, ctx context.Context) (*entities.UploadChildRequest, error) {
	if id == "" {
		return nil, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	return u.uploadChildRequestRepo.GetUploadChildRequest(id, ctx)
}

// GetUploadChildRequests implements business.IUploadChildRequestService.
func (u *uploadChildRequestService) GetUploadChildRequests(req request.GetUploadChildRequests, ctx context.Context) (response.PaginationDataResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}

	if req.PageSize < 1 {
		req.PageSize = default_page_size
	}

	data, pages, err := u.uploadChildRequestRepo.GetUploadChildRequests(req, ctx)

	return response.PaginationDataResponse{
		Data:       data,
		Page:       req.Page,
		TotalPages: pages,
	}, err
}

// GetWalletUploadChildRequests implements business.IUploadChildRequestService.
func (u *uploadChildRequestService) GetWalletUploadChildRequests(id string, page int, ctx context.Context) (response.PaginationDataResponse, error) {
	if !utils.IsValidSuiAddress(models.SuiAddress(id)) {
		return response.PaginationDataResponse{}, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	if page < 1 {
		page = 1
	}

	data, pages, err := u.uploadChildRequestRepo.GetWalletUploadChildRequests(id, page, ctx)

	return response.PaginationDataResponse{
		Data:       data,
		Page:       page,
		TotalPages: pages,
	}, err
}

// VoteUploadChildRequest implements business.IUploadChildRequestService.
func (u *uploadChildRequestService) VoteUploadChildRequest(id string, req request.VoteRequest, ctx context.Context) error {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)

	var voter string = ctx.Value("address").(string)
	if !utils.IsValidSuiAddress(models.SuiAddress(voter)) {
		return genericErr
	}

	request, err := u.uploadChildRequestRepo.GetUploadChildRequest(id, ctx)
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
		Client:    u.clients[constant.SuiTestnet],
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: u.errLogger,
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

	return u.uploadChildRequestRepo.UpdateUploadChildRequest(*request, ctx)
}
