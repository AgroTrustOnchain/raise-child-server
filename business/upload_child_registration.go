package business

import (
	"context"
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

type uploadChildRequestService struct {
	uploadChildRequestRepo i_repository.IUploadChildRequestRepository
	redisCache             cache.IRedisCache
	clients                map[string]sui.ISuiAPI
	errLogger              *log.Logger
}

func initializeUploadChildRequestService(
	uploadChildRequestRepo i_repository.IUploadChildRequestRepository,
	clients map[string]sui.ISuiAPI,
	errLogger *log.Logger,
) business.IUploadChildRequestService {
	return &uploadChildRequestService{
		uploadChildRequestRepo: uploadChildRequestRepo,
		redisCache:             cache.InitializeRedisCache(),
		clients:                clients,
		errLogger:              errLogger,
	}
}

func GenerateUploadChildRequestService() (business.IUploadChildRequestService, error) {
	var errLogger = util.GetLogConfig(shared.ERROR_LEVEL)

	cnn, err := db.ConnectDB(errLogger, db.InitializePostgreSQL())
	if err != nil {
		return nil, err
	}

	return initializeUploadChildRequestService(repository.InitializeUploadChildRequestRepo(cnn, errLogger), _networkAliases, errLogger), nil
}

// ReviewUploadChildRequest implements business.IUploadChildRequestService.
func (u *uploadChildRequestService) ReviewUploadChildRequest(id string, req request.VoteRequest, ctx context.Context) error {
	request, err := u.uploadChildRequestRepo.GetUploadChildRequest(id, ctx)
	if err != nil {
		return err
	}

	if request == nil {
		return errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	if request.ReviewStatus != request_pending_status || request.ReviewedBy != nil {
		return errors.New(noti.REQUEST_REVIEWED_MESSAGE)
	}

	var sender string = ctx.Value("address").(string)
	manageObj, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    u.clients[constant.SuiTestnet],
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: u.errLogger,
	}, ctx)
	if err != nil {
		return err
	}

	if !slices.Contains(manageObj.AdminIds, sender) {
		return errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	}

	var status string
	var closedAt *time.Time
	if req.IsVoteYes {
		status = request_approved_status
		var tmp = util.GetRequestDuration()
		closedAt = &tmp
	} else {
		status = request_refused_status
	}

	return u.uploadChildRequestRepo.SetReviewStatus(id, status, sender, closedAt, ctx)
}

// ConfirmUploadChildRequest implements business.IUploadChildRequestService.
func (u *uploadChildRequestService) ConfirmUploadChildRequest(id string, ctx context.Context) (response.BuildTransactionResponse, error) {
	req, err := u.uploadChildRequestRepo.GetUploadChildRequest(id, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	if req == nil {
		return response.BuildTransactionResponse{}, errors.New(noti.GENERIC_ERROR_WARN_MSG)
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
		var secondGuardianFullName, secondGuardianPhone, secondGuardianRelation, secondGuardianIdentityCardBlobID string
		if req.SecondGuardianProfile != nil {
			secondGuardianFullName = req.SecondGuardianProfile.FullName
			secondGuardianPhone = req.SecondGuardianProfile.PhoneNumber
			secondGuardianRelation = req.SecondGuardianProfile.Relation
			secondGuardianIdentityCardBlobID = req.SecondGuardianProfile.IdentityCardBlobID
		}

		var module = on_chain.InitializeModuleChild()
		txBytes, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
			Client:    u.clients[constant.SuiTestnet],
			Sender:    sender,
			Module:    module.GetModule(),
			Function:  module.GetFunctionAddChild(),
			ErrLogger: u.errLogger,
			Arguments: module.ToAddChildArguments(on_chain.AddChildArguments{
				IdentityCode:                     req.IdentityCode,
				FirstName:                        req.FirstName,
				LastName:                         req.LastName,
				Gender:                           req.Gender,
				DateOfBirth:                      req.DateOfBirth,
				HomeAddress:                      req.HomeAddress,
				Region:                           req.Region,
				AvatarBlobId:                     req.AvatarBlobId,
				HomeBlobID:                       req.HomeBlobID,
				FirstGuardianFullName:            req.FirstGuardianProfile.FullName,
				FirstGuardianPhone:               req.FirstGuardianProfile.PhoneNumber,
				FirstGuardianRelation:            req.FirstGuardianProfile.Relation,
				FirstGuardianIdentityCardBlobID:  req.FirstGuardianProfile.IdentityCardBlobID,
				SecondGuardianFullName:           secondGuardianFullName,
				SecondGuardianPhone:              secondGuardianPhone,
				SecondGuardianRelation:           secondGuardianRelation,
				SecondGuardianIdentityCardBlobID: secondGuardianIdentityCardBlobID,
			}),
		}, ctx)

		res.TxBytes = txBytes
		errRes = err
	}

	return res, errRes
}

// CreateUploadChildRequest implements business.IUploadChildRequestService.
func (u *uploadChildRequestService) CreateUploadChildRequest(req request.UploadChildRequest, ctx context.Context) (*entities.UploadChildRequest, error) {
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

	var firstGuardianProfile = entities.ChildGuardianProfile{
		FullName:           strings.TrimSpace(req.FirstGuardian.FullName),
		PhoneNumber:        strings.TrimSpace(req.FirstGuardian.PhoneNumber),
		Relation:           req.FirstGuardian.Relation,
		IdentityCardBlobID: strings.TrimSpace(req.FirstGuardian.IdentityCardBlobID),
	}

	if firstGuardianProfile.FullName == "" || firstGuardianProfile.PhoneNumber == "" || firstGuardianProfile.Relation == "" || firstGuardianProfile.IdentityCardBlobID == "" {
		return nil, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	var secondGuardianProfile *entities.ChildGuardianProfile
	if req.SecondGuardian != nil {
		secondGuardianProfile = &entities.ChildGuardianProfile{
			FullName:           strings.TrimSpace(req.SecondGuardian.FullName),
			PhoneNumber:        strings.TrimSpace(req.SecondGuardian.PhoneNumber),
			Relation:           req.SecondGuardian.Relation,
			IdentityCardBlobID: strings.TrimSpace(req.SecondGuardian.IdentityCardBlobID),
		}
	}

	// todo: AI validation
	var curTime time.Time = time.Now()
	var request = entities.UploadChildRequest{
		ID:                    util.GenerateId(),
		ProfileID:             ctx.Value("sub").(string),
		IdentityCode:          identityCode,
		AvatarBlobId:          strings.TrimSpace(req.AvatarBlobId),
		HomeBlobID:            req.HomeBlobID,
		Region:                region,
		FirstName:             strings.TrimSpace(req.FirstName),
		LastName:              strings.TrimSpace(req.LastName),
		Gender:                gender,
		DateOfBirth:           dateOfBirth,
		HomeAddress:           strings.TrimSpace(req.HomeAddress),
		FirstGuardianProfile:  firstGuardianProfile,
		SecondGuardianProfile: secondGuardianProfile,
		AIEvaluation:          "",
		CreatedBy:             ctx.Value("address").(string),
		CreatedAt:             curTime,
		UpdatedAt:             curTime,
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
	req.SortOrder = util.StanderizeSortOrder(req.SortOrder)
	if req.Page < 1 {
		req.Page = 1
	}

	if req.PageSize < 1 {
		req.PageSize = default_page_size
	}

	var res response.PaginationDataResponse
	var redisKey string = u.getGetUploadChildRequestsRedisKey(req)
	if u.redisCache.Get(redisKey, &res, ctx) {
		return res, nil
	}

	data, pages, err := u.uploadChildRequestRepo.GetUploadChildRequests(req, ctx)
	if err != nil {
		return response.PaginationDataResponse{}, err
	}

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

	u.redisCache.Set(redisKey, res, time.Minute*5, ctx)

	return res, nil
}

// GetWalletUploadChildRequests implements business.IUploadChildRequestService.
func (u *uploadChildRequestService) GetWalletUploadChildRequests(id string, page int, ctx context.Context) (response.PaginationDataResponse, error) {
	if !util.IsValidSuiAddressStrict(id) {
		return response.PaginationDataResponse{}, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	if page < 1 {
		page = 1
	}

	var res response.PaginationDataResponse
	var redisKey string = u.getGetUploadChildRequestsOfWalletRedisKey(id, page)
	if u.redisCache.Get(redisKey, &res, ctx) {
		return res, nil
	}

	data, pages, err := u.uploadChildRequestRepo.GetWalletUploadChildRequests(id, page, ctx)
	if err != nil {
		return response.PaginationDataResponse{}, err
	}

	var amount int
	if data == nil || len(data) == 0 {
		amount = 0
	} else {
		amount = len(data)
	}

	res = response.PaginationDataResponse{
		Data:       data,
		Amount:     amount,
		Page:       page,
		TotalPages: pages,
	}

	u.redisCache.Set(redisKey, res, time.Minute*5, ctx)

	return res, nil
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

func (u *uploadChildRequestService) getGetUploadChildRequestsRedisKey(req request.GetUploadChildRequests) string {
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

	return fmt.Sprintf("upload_child_req:kw:%s:r:%s:g:%s:status:%s:closed:%s:o:%s:s:%d:p:%d",
		keyword, region, gender, status, isClosed, req.SortOrder, req.PageSize, req.Page)
}

func (u *uploadChildRequestService) getGetUploadChildRequestsOfWalletRedisKey(id string, page int) string {
	return fmt.Sprintf("upload_child_req:of:%s:s:%d:p;%d", id, default_page_size, page)
}
