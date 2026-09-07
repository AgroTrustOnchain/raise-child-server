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
	"raise-child/util/ai"
	"raise-child/util/cache"
	"raise-child/util/db"
	"raise-child/util/image/cloudinary"
	on_chain "raise-child/util/on_chain"
	walrus_pkg "raise-child/util/walrus_pkg"
	"strings"
	"time"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/sui"
)

type uploadChildRequestService struct {
	platformConfigRepo     i_repository.IPlatformConfigRepository
	uploadChildRequestRepo i_repository.IUploadChildRequestRepository
	aiProvider             ai.IAiClientProvider
	walrusProvider         walrus_pkg.IWalrusProvider
	redisCache             cache.IRedisCache
	clients                map[string]sui.ISuiAPI
	errLogger              *log.Logger
}

func initializeUploadChildRequestService(
	platformConfigRepo i_repository.IPlatformConfigRepository,
	uploadChildRequestRepo i_repository.IUploadChildRequestRepository,
	aiProvider ai.IAiClientProvider,
	walrusProvider walrus_pkg.IWalrusProvider,
	clients map[string]sui.ISuiAPI,
	errLogger *log.Logger,
) business.IUploadChildRequestService {
	return &uploadChildRequestService{
		platformConfigRepo:     platformConfigRepo,
		uploadChildRequestRepo: uploadChildRequestRepo,
		aiProvider:             aiProvider,
		walrusProvider:         walrusProvider,
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

	var uploadChildRequestService = initializeUploadChildRequestService(
		repository.InitializePlatformConfigRepository(cnn, errLogger),
		repository.InitializeUploadChildRequestRepo(cnn, errLogger),
		ai.InitializeAiProvider(errLogger),
		walrus_pkg.InitializeWalrusProvider(errLogger),
		_networkAliases,
		errLogger,
	)

	return uploadChildRequestService, nil
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

	if request.Status != request_pending_status || request.ReviewedBy != nil {
		return errors.New(noti.REQUEST_REVIEWED_MESSAGE)
	}

	var client = u.clients[constant.SuiTestnet]
	manageObj, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: u.errLogger,
	}, ctx)
	if err != nil {
		return err
	}

	var foundIdx int = -1
	var centerId string
	var sender string = ctx.Value("address").(string)
	for i, leader := range manageObj.LocalLeaderIds {
		if centerId == "" {
			if i < len(manageObj.LocalRegions) {
				if manageObj.LocalRegions[i] == request.Region {
					centerId = manageObj.ChildrenCenters[i]
				}
			}
		}

		u.errLogger.Println("Leader through loop:", leader)
		if leader == sender {
			foundIdx = i
			break
		}
	}

	var genericRightErr error = errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	if foundIdx == -1 {
		return genericRightErr
	}

	nft, err := on_chain.GetOnChainObject[entities.StaffNft](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  manageObj.LocalLeaderNfts[foundIdx],
		ErrLogger: u.errLogger,
	}, ctx)
	if err != nil {
		return err
	}

	if nft.Region != request.Region {
		return genericRightErr
	}

	var curTime time.Time = time.Now()
	request.ReviewedBy = &sender
	request.UpdatedAt = curTime

	if !req.IsVoteYes {
		request.Status = request_refused_status
		return u.uploadChildRequestRepo.UpdateUploadChildRequest(*request, ctx)
	}

	var secondGuardianFullName, secondGuardianPhone, secondGuardianRelation string
	if request.SecondGuardianProfile != nil {
		secondGuardianFullName = request.SecondGuardianProfile.FullName
		secondGuardianPhone = request.SecondGuardianProfile.PhoneNumber
		secondGuardianRelation = request.SecondGuardianProfile.Relation
	}

	var childModule = on_chain.InitializeModuleChild()
	res, err := on_chain.ExecuteTransactionV2(on_chain.ExecuteTransactionRequestV2{
		Client:   client,
		Module:   childModule.GetModule(),
		Function: childModule.GetFunctionAddChild(),
		Arguments: childModule.ToAddChildArguments(on_chain.AddChildArguments{
			Center:                 centerId,
			IdentityCode:           request.IdentityCode,
			FirstName:              request.FirstName,
			LastName:               request.LastName,
			Gender:                 request.Gender,
			DateOfBirth:            request.DateOfBirth,
			HomeAddress:            request.HomeAddress,
			Region:                 request.Region,
			AvatarBlobId:           request.AvatarBlobId,
			HomeBlobID:             request.HomeBlobID,
			FirstGuardianFullName:  request.FirstGuardianProfile.FullName,
			FirstGuardianPhone:     request.FirstGuardianProfile.PhoneNumber,
			FirstGuardianRelation:  request.FirstGuardianProfile.Relation,
			SecondGuardianFullName: secondGuardianFullName,
			SecondGuardianPhone:    secondGuardianPhone,
			SecondGuardianRelation: secondGuardianRelation,
			Sender:                 request.CreatedBy,
		}),
		ErrLogger: u.errLogger,
	}, ctx)
	if err != nil {
		return err
	}

	var onchainId string
	var module = on_chain.InitializeModuleChild()
	var eventType string = fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), module.GetModule(), module.GetChildEventEmittedStruct())
	for _, event := range res.Events {
		if event.Type == eventType {
			if eventId, ok := event.ParsedJson["id"].(string); ok {
				onchainId = eventId
				break
			}
		}
	}

	request.Status = request_approved_status
	request.OnchainID = &onchainId

	for i := 1; i <= 3; i++ {
		if u.uploadChildRequestRepo.UpdateUploadChildRequest(*request, ctx) == nil {
			return nil
		}
	}

	return errors.New(noti.INTERNALL_ERR_MSG)
}

// // ConfirmUploadChildRequest implements business.IUploadChildRequestService.
// func (u *uploadChildRequestService) ConfirmUploadChildRequest(id string, ctx context.Context) (response.BuildTransactionResponse, error) {
// 	req, err := u.uploadChildRequestRepo.GetUploadChildRequest(id, ctx)
// 	if err != nil {
// 		return response.BuildTransactionResponse{}, err
// 	}

// 	if req == nil {
// 		return response.BuildTransactionResponse{}, errors.New(noti.GENERIC_ERROR_WARN_MSG)
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
// 			req.IsConfirmUpload = true
// 		} else {
// 			req.Status = request_refused_status
// 			isDenied = true
// 		}

// 		req.UpdatedAt = time.Now()
// 		if err := u.uploadChildRequestRepo.UpdateUploadChildRequest(*req, ctx); err != nil {
// 			return response.BuildTransactionResponse{}, err
// 		}

// 		if isDenied {
// 			return response.BuildTransactionResponse{}, nil
// 		}
// 	}

// 	var res response.BuildTransactionResponse
// 	var errRes error
// 	if req.IsConfirmUpload {
// 		var secondGuardianFullName, secondGuardianPhone, secondGuardianRelation, secondGuardianIdentityCardBlobID string
// 		if req.SecondGuardianProfile != nil {
// 			secondGuardianFullName = req.SecondGuardianProfile.FullName
// 			secondGuardianPhone = req.SecondGuardianProfile.PhoneNumber
// 			secondGuardianRelation = req.SecondGuardianProfile.Relation
// 			secondGuardianIdentityCardBlobID = req.SecondGuardianProfile.IdentityCardBlobID
// 		}

// 		var module = on_chain.InitializeModuleChild()
// 		txBytes, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
// 			Client:    u.clients[constant.SuiTestnet],
// 			Sender:    sender,
// 			Module:    module.GetModule(),
// 			Function:  module.GetFunctionAddChild(),
// 			ErrLogger: u.errLogger,
// 			Arguments: module.ToAddChildArguments(on_chain.AddChildArguments{
// 				IdentityCode:                     req.IdentityCode,
// 				FirstName:                        req.FirstName,
// 				LastName:                         req.LastName,
// 				Gender:                           req.Gender,
// 				DateOfBirth:                      req.DateOfBirth,
// 				HomeAddress:                      req.HomeAddress,
// 				Region:                           req.Region,
// 				AvatarBlobId:                     req.AvatarBlobId,
// 				HomeBlobID:                       req.HomeBlobID,
// 				FirstGuardianFullName:            req.FirstGuardianProfile.FullName,
// 				FirstGuardianPhone:               req.FirstGuardianProfile.PhoneNumber,
// 				FirstGuardianRelation:            req.FirstGuardianProfile.Relation,
// 				FirstGuardianIdentityCardBlobID:  req.FirstGuardianProfile.IdentityCardBlobID,
// 				SecondGuardianFullName:           secondGuardianFullName,
// 				SecondGuardianPhone:              secondGuardianPhone,
// 				SecondGuardianRelation:           secondGuardianRelation,
// 				SecondGuardianIdentityCardBlobID: secondGuardianIdentityCardBlobID,
// 			}),
// 		}, ctx)

// 		res.TxBytes = txBytes
// 		errRes = err
// 	}

// 	return res, errRes
// }

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

	var gender string = util.StandardizeGender(util.StandardizeString(req.Gender))
	if gender == "" {
		return nil, errors.New(noti.UNDEFINED_GENDER_MESSAGE)
	}

	var dateOfBirth string = strings.TrimSpace(req.DateOfBirth)
	if dob := util.RawDateToTime(dateOfBirth); dob.IsZero() {
		return nil, errors.New(noti.INVALID_DATE_FORMAT_WARN_MSG)
	} else {
		var maxAgeAccepted, minAgeAccepted int64
		var tmp *entities.PlatformConfig
		var numericTmp *entities.NumericConfig

		if config, _ := u.platformConfigRepo.GetConfigByKey(numericTmp.GetTable(), tmp.GetChildAgeMaxLimitConfigKey(), ctx); config != nil {
			if v, ok := config.Value.(int64); ok {
				maxAgeAccepted = v
			} else {
				maxAgeAccepted = child_max_age_accepted
			}
		} else {
			maxAgeAccepted = child_max_age_accepted
		}

		if config, _ := u.platformConfigRepo.GetConfigByKey(numericTmp.GetTable(), tmp.GetChildAgeMinLimitConfigKey(), ctx); config != nil {
			if v, ok := config.Value.(int64); ok {
				minAgeAccepted = v
			} else {
				minAgeAccepted = child_min_age_accepted
			}
		} else {
			minAgeAccepted = child_min_age_accepted
		}

		if !isChildAgeInSupport(dob.Year(), maxAgeAccepted, minAgeAccepted) {
			return nil, errors.New(noti.CHILD_AGE_OUT_OF_SUPPORT_MESSAGE)
		}
	}

	var firstGuardianProfile = entities.ChildGuardianProfile{
		FullName:           strings.TrimSpace(req.FirstGuardian.FullName),
		PhoneNumber:        strings.TrimSpace(req.FirstGuardian.PhoneNumber),
		Relation:           util.StandardizeRelation(req.FirstGuardian.Relation),
		IdentityCardBlobID: strings.TrimSpace(req.FirstGuardian.IdentityCardBlobID),
	}

	if firstGuardianProfile.FullName == "" || firstGuardianProfile.PhoneNumber == "" || firstGuardianProfile.Relation == "" || firstGuardianProfile.IdentityCardBlobID == "" {
		return nil, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	var secondGuardianProfile *entities.ChildGuardianProfile
	if req.SecondGuardian != nil {
		var relation string = util.StandardizeRelation(req.SecondGuardian.Relation)
		if relation == "" {
			return nil, errors.New(noti.UNDEFINED_RELATIONSHIP_MESSAGE)
		}

		secondGuardianProfile = &entities.ChildGuardianProfile{
			FullName:           strings.TrimSpace(req.SecondGuardian.FullName),
			PhoneNumber:        strings.TrimSpace(req.SecondGuardian.PhoneNumber),
			Relation:           relation,
			IdentityCardBlobID: strings.TrimSpace(req.SecondGuardian.IdentityCardBlobID),
		}
	}

	var firstNasme string = strings.TrimSpace(req.FirstName)
	var lastName string = strings.TrimSpace(req.LastName)
	var homeAddr string = strings.TrimSpace(req.HomeAddress)
	var curTime time.Time = time.Now()
	var request = entities.UploadChildRequest{
		ID:                     util.GenerateId(),
		ProfileID:              ctx.Value("sub").(string),
		IdentityCode:           identityCode,
		BirthCertificateBlobID: req.BirthCertificateBlobID,
		AvatarBlobId:           req.AvatarBlobId,
		HomeBlobID:             req.HomeBlobID,
		Region:                 region,
		FirstName:              firstNasme,
		LastName:               lastName,
		Gender:                 gender,
		DateOfBirth:            dateOfBirth,
		HomeAddress:            homeAddr,
		FirstGuardianProfile:   firstGuardianProfile,
		SecondGuardianProfile:  secondGuardianProfile,
		Status:                 request_pending_status,
		CreatedBy:              ctx.Value("address").(string),
		CreatedAt:              curTime,
		UpdatedAt:              curTime,
	}

	return &request, u.uploadChildRequestRepo.CreateUploadChildRequest(request, ctx)
}

// GetUploadChildRequest implements business.IUploadChildRequestService.
func (u *uploadChildRequestService) GetUploadChildRequest(id string, ctx context.Context) (response.UploadChildRequestResponse, error) {
	req, err := u.uploadChildRequestRepo.GetUploadChildRequest(id, ctx)
	if err != nil {
		return response.UploadChildRequestResponse{}, nil
	}

	var res = req.ToUploadChildRequestResponse()
	res.BirthCertificateImgUrl = cloudinary.GetImageUrl(req.BirthCertificateBlobID)
	res.FirstGuardianProfile.IdentityCardImgUrl = cloudinary.GetImageUrl(req.FirstGuardianProfile.IdentityCardBlobID)
	if req.SecondGuardianProfile != nil {
		res.SecondGuardianProfile.IdentityCardImgUrl = cloudinary.GetImageUrl(req.SecondGuardianProfile.IdentityCardBlobID)
	}

	return res, nil
}

// GetUploadChildRequests implements business.IUploadChildRequestService.
func (u *uploadChildRequestService) GetUploadChildRequests(req request.GetUploadChildRequests, ctx context.Context) (response.PaginationDataResponse, error) {
	req.SortOrder = util.StandardizeSortOrder(req.SortOrder)
	if req.Page < 1 {
		req.Page = 1
	}

	if req.PageSize < 1 {
		req.PageSize = default_page_size
	}

	var res response.PaginationDataResponse
	data, pages, err := u.uploadChildRequestRepo.GetUploadChildRequests(req, ctx)
	if err != nil {
		return response.PaginationDataResponse{}, err
	}

	var amount int
	if len(data) == 0 {
		amount = 0
	} else {
		amount = len(data)
	}

	var resData []response.UploadChildRequestResponse
	if amount > 0 {
		for _, req := range data {
			var reqRes = req.ToUploadChildRequestResponse()
			reqRes.BirthCertificateImgUrl = cloudinary.GetImageUrl(req.BirthCertificateBlobID)
			reqRes.FirstGuardianProfile.IdentityCardImgUrl = cloudinary.GetImageUrl(req.FirstGuardianProfile.IdentityCardBlobID)
			if req.SecondGuardianProfile != nil {
				reqRes.SecondGuardianProfile.IdentityCardImgUrl = cloudinary.GetImageUrl(req.SecondGuardianProfile.IdentityCardBlobID)
			}

			resData = append(resData, reqRes)
		}
	}

	res = response.PaginationDataResponse{
		Data:       resData,
		Amount:     amount,
		Page:       req.Page,
		TotalPages: pages,
	}

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
	data, pages, err := u.uploadChildRequestRepo.GetWalletUploadChildRequests(id, page, ctx)
	if err != nil {
		return response.PaginationDataResponse{}, err
	}

	var amount int
	if len(data) == 0 {
		amount = 0
	} else {
		amount = len(data)
	}

	var resData []response.UploadChildRequestResponse
	if amount > 0 {
		for _, req := range data {
			var reqRes = req.ToUploadChildRequestResponse()
			reqRes.BirthCertificateImgUrl = cloudinary.GetImageUrl(req.BirthCertificateBlobID)
			reqRes.FirstGuardianProfile.IdentityCardImgUrl = cloudinary.GetImageUrl(req.FirstGuardianProfile.IdentityCardBlobID)
			if req.SecondGuardianProfile != nil {
				reqRes.SecondGuardianProfile.IdentityCardImgUrl = cloudinary.GetImageUrl(req.SecondGuardianProfile.IdentityCardBlobID)
			}

			resData = append(resData, reqRes)
		}
	}

	res = response.PaginationDataResponse{
		Data:       resData,
		Amount:     amount,
		Page:       page,
		TotalPages: pages,
	}

	return res, nil
}

// // VoteUploadChildRequest implements business.IUploadChildRequestService.
// func (u *uploadChildRequestService) VoteUploadChildRequest(id string, req request.VoteRequest, ctx context.Context) error {
// 	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)

// 	var voter string = ctx.Value("address").(string)
// 	if !utils.IsValidSuiAddress(models.SuiAddress(voter)) {
// 		return genericErr
// 	}

// 	request, err := u.uploadChildRequestRepo.GetUploadChildRequest(id, ctx)
// 	if err != nil {
// 		return err
// 	}

// 	if request == nil {
// 		return genericErr
// 	}

// 	if request.ClosedAt.Before(time.Now()) {
// 		return errors.New(noti.REQUEST_CLOSED_MESSAGE)
// 	}

// 	if voter == request.CreatedBy {
// 		return errors.New(noti.OWNER_VOTE_WARN_MSG)
// 	}

// 	if slices.Contains(request.Approvers, voter) || slices.Contains(request.Refusers, voter) {
// 		return errors.New(noti.ALREADY_VOTE_MESSAGE)
// 	}

// 	var manageObj entities.Manage
// 	if !u.redisCache.Get(manageObj.GetRedisKey(), &manageObj, ctx) {
// 		res, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
// 			Client:    u.clients[constant.SuiTestnet],
// 			ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
// 			ErrLogger: u.errLogger,
// 		}, ctx)
// 		if err != nil {
// 			return err
// 		}

// 		if res != nil {
// 			u.redisCache.Set(manageObj.GetRedisKey(), res, time.Minute, ctx)
// 			manageObj = *res
// 		}
// 	}

// 	// Not admins or local leaders
// 	if !slices.Contains(manageObj.AdminIds, voter) && !slices.Contains(manageObj.LocalLeaderIds, voter) {
// 		return errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
// 	}

// 	if req.IsVoteYes {
// 		request.Approvers = append(request.Approvers, voter)
// 	} else {
// 		request.Refusers = append(request.Refusers, voter)
// 		if req.RefuseReason == "" {
// 			return errors.New(noti.FIELD_EMPTY_WARN_MSG)
// 		}

// 		request.RefuseReasons = append(request.RefuseReasons, strings.TrimSpace(req.RefuseReason))
// 	}

// 	request.UpdatedAt = time.Now()

// 	return u.uploadChildRequestRepo.UpdateUploadChildRequest(*request, ctx)
// }

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
