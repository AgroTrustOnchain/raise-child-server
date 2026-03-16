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
	on_chain "raise-child/util/on_chain"
	"strings"
	"time"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/block-vision/sui-go-sdk/utils"
)

type giftService struct {
	profileRepo i_repository.IProfileRepository
	redisCache  cache.IRedisCache
	clients     map[string]sui.ISuiAPI
	errLogger   *log.Logger
}

func initializeGiftService(profileRepo i_repository.IProfileRepository, clients map[string]sui.ISuiAPI, errLogger *log.Logger) business.IGiftService {
	return &giftService{
		profileRepo: profileRepo,
		redisCache:  cache.InitializeRedisCache(),
		clients:     clients,
		errLogger:   errLogger,
	}
}

func GenerateGiftService() (business.IGiftService, error) {
	var errLogger = util.GetLogConfig(shared.ERROR_LEVEL)

	cnn, err := db.ConnectDB(errLogger, db.InitializePostgreSQL())
	if err != nil {
		return nil, err
	}

	// return InitializeGiftService(cnn, errLogger), nil
	return initializeGiftService(repository.InitializeProfileRepository(cnn, errLogger), _networkAliases, errLogger), nil
}

const (
	gift_limit_record int    = 10
	food_category     string = "Food Items"
	non_food_category string = "Non-Food Items"
)

// CancelGift implements business.IGiftService.
func (g *giftService) CancelGift(id string, req request.CancelGiftRequest, ctx context.Context) (response.BuildTransactionResponse, error) {
	var sender string = ctx.Value("address").(string)
	var genericRightErr error = errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	if !utils.IsValidSuiAddress(models.SuiAddress(sender)) {
		return response.BuildTransactionResponse{}, genericRightErr
	}

	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if !utils.IsValidSuiAddress(models.SuiAddress(id)) {
		return response.BuildTransactionResponse{}, genericErr
	}

	var client = g.clients[constant.SuiTestnet]
	gift, err := on_chain.GetOnChainObject[entities.Gift](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  id,
		ErrLogger: g.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	if gift.Status == "Delivered" || gift.Status == "Canceled" {
		return response.BuildTransactionResponse{}, genericErr
	}

	if gift.Sender != sender {
		return response.BuildTransactionResponse{}, genericRightErr
	}

	var giftModule = on_chain.InitializeModuleGift()
	txBytes, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
		Client:    g.clients[constant.SuiTestnet],
		Sender:    sender,
		Module:    giftModule.GetModule(),
		Function:  giftModule.GetFunctionCancelGift(),
		ErrLogger: g.errLogger,
		Arguments: giftModule.ToCancelGiftArguments(on_chain.CancelGiftArguments{
			GiftID:       id,
			CancelReason: req.CancelReason,
		}),
	}, ctx)

	return response.BuildTransactionResponse{
		TxBytes: txBytes,
	}, err
}

// ConfirmReceiveGift implements business.IGiftService.
func (g *giftService) ConfirmReceiveGift(id string, req request.ConfirmReceiveGiftRequest, ctx context.Context) (response.BuildTransactionResponse, error) {
	var sender string = ctx.Value("address").(string)
	var genericRightErr error = errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	if !utils.IsValidSuiAddress(models.SuiAddress(sender)) {
		return response.BuildTransactionResponse{}, genericRightErr
	}

	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if !utils.IsValidSuiAddress(models.SuiAddress(id)) {
		return response.BuildTransactionResponse{}, genericErr
	}

	var client = g.clients[constant.SuiTestnet]
	gift, err := on_chain.GetOnChainObject[entities.Gift](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  id,
		ErrLogger: g.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	if gift.Status == "Delivered" || gift.Status == "Canceled" {
		return response.BuildTransactionResponse{}, genericErr
	}

	if gift.Sender == sender {
		return response.BuildTransactionResponse{}, genericRightErr
	}

	var staffModule = on_chain.InitializeModuleStaff()
	roles, err := on_chain.GetOnChainOwnedObjects[entities.StaffNft](on_chain.GetOnChainOwnedObjectsRequest{
		Client:       client,
		OwnerAddress: sender,
		StructType:   fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), staffModule.GetModule(), staffModule.GetStaffNftObjectStruct()),
		ErrLogger:    g.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	if roles == nil || len(roles) == 0 {
		return response.BuildTransactionResponse{}, genericRightErr
	}

	var function string
	var region string
	var recipient string
	var giftModule = on_chain.InitializeModuleGift()
	var getObjReq = on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  gift.Recipient,
		ErrLogger: g.errLogger,
	}

	if gift.IsForChild {
		child, err := on_chain.GetOnChainObject[entities.Child](getObjReq, ctx)
		if err != nil {
			return response.BuildTransactionResponse{}, err
		}

		if child == nil {
			return response.BuildTransactionResponse{}, genericErr
		}

		function = giftModule.GetFunctionConfirmReceiveChildGift()
		region = child.Region
		recipient = child.ID.ID
	} else {
		center, err := on_chain.GetOnChainObject[entities.Center](getObjReq, ctx)
		if err != nil {
			return response.BuildTransactionResponse{}, err
		}

		if center == nil {
			return response.BuildTransactionResponse{}, genericErr
		}

		function = giftModule.GetFunctionConfirmReceiveCenterGift()
		region = center.Region
		recipient = center.ID.ID
	}

	var staffId string
	for _, role := range roles {
		if role.Region == region {
			staffId = role.ID.ID
			break
		}
	}

	if staffId == "" {
		return response.BuildTransactionResponse{}, genericRightErr
	}

	txBytes, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
		Client:    client,
		Sender:    sender,
		Module:    giftModule.GetModule(),
		Function:  function,
		ErrLogger: g.errLogger,
		Arguments: giftModule.ToConfirmReceiveGiftArguments(on_chain.ConfirmReceiveGiftArguments{
			GiftID:      id,
			Recipient:   recipient,
			StaffID:     staffId,
			ImageBlobID: req.DeliveredImageBlobID,
		}),
	}, ctx)

	return response.BuildTransactionResponse{
		TxBytes: txBytes,
	}, err
}

// CreateGift implements business.IGiftService.
func (g *giftService) CreateGift(req request.CreateGiftRequest, ctx context.Context) (response.BuildTransactionResponse, error) {
	var genericRightErr error = errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)

	var sender string = ctx.Value("address").(string)
	if !utils.IsValidSuiAddress(models.SuiAddress(sender)) {
		return response.BuildTransactionResponse{}, genericRightErr
	}

	if !utils.IsValidSuiAddress(models.SuiAddress(req.Recipient)) {
		return response.BuildTransactionResponse{}, genericErr
	}

	profile, err := g.profileRepo.GetProfile(ctx.Value("sub").(string), ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	if profile.FirstName == "" {
		return response.BuildTransactionResponse{}, errors.New(noti.PROFILE_EMPTY_MESSAGE)
	}

	if req.Category == "" {
		req.Category = non_food_category
	}

	var giftText string = "Gift"
	var description string = strings.TrimSpace(req.Description)
	var msg string = strings.TrimSpace(req.Message)
	if description == "" {
		description = giftText
	}

	if msg == "" {
		msg = giftText
	}

	// todo: AI validate other fields
	var client = g.clients[constant.SuiTestnet]
	var donorModule = on_chain.InitializeModuleDonor()
	nfts, err := on_chain.GetOnChainOwnedObjects[entities.Donor](on_chain.GetOnChainOwnedObjectsRequest{
		Client:       client,
		OwnerAddress: sender,
		StructType:   fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), donorModule.GetModule(), donorModule.GetDonorNftStruct()),
		ErrLogger:    g.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	var nftId string
	if nfts == nil || len(nfts) == 0 {
		nftId = os.Getenv(env.PUB_DONOR_NFT_ID)
	} else {
		nftId = nfts[0].ID.ID
	}

	var giftModule = on_chain.InitializeModuleGift()
	var function string
	var args []interface{}
	var getRecipientObjReq = on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  req.Recipient,
		ErrLogger: g.errLogger,
	}

	var child *entities.Child
	var errRes error
	child, errRes = on_chain.GetOnChainObject[entities.Child](getRecipientObjReq, ctx)
	if errRes != nil {
		return response.BuildTransactionResponse{}, errRes
	}

	// Gift for child
	if child != nil {
		manageObj, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
			Client:    client,
			ObjectId:  os.Getenv(env.PACKAGE_ID),
			ErrLogger: g.errLogger,
		}, ctx)
		if err != nil {
			return response.BuildTransactionResponse{}, err
		}

		var centerId string
		for i, region := range manageObj.LocalRegions {
			if region == child.Region {
				centerId = manageObj.ChildrenCenters[i]
				break
			}
		}

		function = giftModule.GetFunctionCreateGiftForChild()
		args = giftModule.ToCreateGiftForChildArguments(on_chain.CreateGiftForChildArguments{
			ChildID: req.Recipient,
			CreateGiftForCenterArguments: on_chain.CreateGiftForCenterArguments{
				DonorID:         nftId,
				CenterID:        centerId,
				TrackingCode:    req.TrackingCode,
				Carrier:         req.Carrier,
				GiftImageBlobID: req.GiftImageBlobID,
				Category:        req.Category,
				Amount:          req.GiftValue,
				FirstName:       profile.FirstName,
				LastName:        profile.LastName,
				Gender:          profile.Gender,
				PhoneNumber:     profile.PhoneNumber,
				Email:           profile.Email,
				Message:         msg,
				Description:     description,
			},
		})
	} else {
		center, _ := on_chain.GetOnChainObject[entities.Center](getRecipientObjReq, ctx)
		if center == nil {
			return response.BuildTransactionResponse{}, genericErr
		}

		function = giftModule.GetFunctionCreateGiftForCenter()
		args = giftModule.ToCreateGiftForCenterArguments(on_chain.CreateGiftForCenterArguments{
			DonorID:         nftId,
			CenterID:        req.Recipient,
			TrackingCode:    req.TrackingCode,
			Carrier:         req.Carrier,
			GiftImageBlobID: req.GiftImageBlobID,
			Category:        req.Category,
			Amount:          req.GiftValue,
			FirstName:       profile.FirstName,
			LastName:        profile.LastName,
			Gender:          profile.Gender,
			PhoneNumber:     profile.PhoneNumber,
			Email:           profile.Email,
			Message:         msg,
			Description:     description,
		})
	}

	// if child, _ := on_chain.GetOnChainObject[entities.Child](getRecipientObjReq, ctx); child != nil {
	// 	function = giftModule.GetFunctionCreateGiftForChild()
	// 	args = giftModule.ToCreateGiftForChildArguments(on_chain.CreateGiftForChildArguments{
	// 		ChildID: req.Recipient,
	// 		CreateGiftForCenterArguments: on_chain.CreateGiftForCenterArguments{
	// 			DonorID:         nftId,
	// 			TrackingCode:    req.TrackingCode,
	// 			Carrier:         req.Carrier,
	// 			GiftImageBlobID: req.GiftImageBlobID,
	// 			Category:        req.Category,
	// 			Amount:          req.GiftValue,
	// 			FirstName:       profile.FirstName,
	// 			LastName:        profile.LastName,
	// 			Gender:          profile.Gender,
	// 			PhoneNumber:     profile.PhoneNumber,
	// 			Email:           profile.Email,
	// 			Message:         msg,
	// 			Description:     description,
	// 		},
	// 	})
	// } else {
	// 	if center, _ := on_chain.GetOnChainObject[entities.Center](getRecipientObjReq, ctx); center != nil {
	// 		function = giftModule.GetFunctionCreateGiftForCenter()
	// 		args = giftModule.ToCreateGiftForCenterArguments(on_chain.CreateGiftForCenterArguments{
	// 			DonorID:         nftId,
	// 			TrackingCode:    req.TrackingCode,
	// 			Carrier:         req.Carrier,
	// 			GiftImageBlobID: req.GiftImageBlobID,
	// 			Category:        req.Category,
	// 			Amount:          req.GiftValue,
	// 			FirstName:       profile.FirstName,
	// 			LastName:        profile.LastName,
	// 			Gender:          profile.Gender,
	// 			PhoneNumber:     profile.PhoneNumber,
	// 			Email:           profile.Email,
	// 			Message:         msg,
	// 			Description:     description,
	// 		})
	// 	}
	// }

	if function == "" {
		return response.BuildTransactionResponse{}, genericErr
	}

	txBytes, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
		Client:    client,
		Sender:    sender,
		Module:    giftModule.GetModule(),
		Function:  function,
		ErrLogger: g.errLogger,
		Arguments: args,
	}, ctx)

	return response.BuildTransactionResponse{
		TxBytes: txBytes,
	}, err
}

// GetGift implements business.IGiftService.
func (g *giftService) GetGift(id string, ctx context.Context) (response.GiftResponse, error) {
	if !util.IsValidSuiAddressStrict(id) {
		return response.GiftResponse{}, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	res, err := on_chain.GetOnChainObject[entities.Gift](on_chain.GetOnChainObjectRequest{
		Client:    g.clients[constant.SuiTestnet],
		ObjectId:  id,
		ErrLogger: g.errLogger,
	}, ctx)

	return res.ToGiftResponse(), err
}

// GetGiftsOfChild implements business.IGiftService.
func (g *giftService) GetGiftsOfChild(id string, req request.GetGiftsRequest, ctx context.Context) (response.PaginationDataResponse, error) {
	if !util.IsValidSuiAddressStrict(id) {
		return response.PaginationDataResponse{}, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	if req.Page < 1 {
		req.Page = 1
	}

	if req.PageSize < 1 {
		req.PageSize = default_page_size
	}

	req.Keyword = util.StanderizeString(req.Keyword)
	req.SortOrder = util.StanderizeSortOrder(req.SortOrder)
	req.Category = strings.TrimSpace(req.Category)
	var res response.PaginationDataResponse
	var redisKey string = g.getGetGiftsRedisKey(id, req)
	if g.redisCache.Get(redisKey, &res, ctx) {
		return res, nil
	}

	var client = g.clients[constant.SuiTestnet]
	var giftIds []string
	var getOnchainObjReq = on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  id,
		ErrLogger: g.errLogger,
	}

	if child, _ := on_chain.GetOnChainObject[entities.Child](getOnchainObjReq, ctx); child != nil {
		giftIds = child.Gifts
	} else {
		if center, _ := on_chain.GetOnChainObject[entities.Center](getOnchainObjReq, ctx); center != nil {
			giftIds = center.Gifts
		}
	}

	if len(giftIds) == 0 {
		return response.PaginationDataResponse{}, nil
	}

	gifts, err := on_chain.GetOnChainObjects[entities.Gift](on_chain.GetOnChainObjectsRequest{
		Client:    client,
		ObjectIds: giftIds,
		ErrLogger: g.errLogger,
	}, ctx)
	if err != nil {
		return response.PaginationDataResponse{}, err
	}

	var filteredGifts []entities.Gift
	if req.SortOrder == "DESC" || req.SortOrder == "" {
		for i := len(gifts) - 1; i >= 0; i-- {
			var gift = gifts[i]
			if isGiftMatchedFilter(gift, req.Keyword, req.Status, req.Category) {
				filteredGifts = append(filteredGifts, gift)
			}
		}
	} else {
		for _, gift := range gifts {
			if isGiftMatchedFilter(gift, req.Keyword, req.Status, req.Category) {
				filteredGifts = append(filteredGifts, gift)
			}
		}
	}

	var skippedRecords int = (req.Page - 1) * req.PageSize
	if len(filteredGifts) <= skippedRecords {
		return response.PaginationDataResponse{}, nil
	}

	var data []response.GiftResponse
	for i := skippedRecords; i < len(filteredGifts); i++ {
		data = append(data, filteredGifts[i].ToGiftResponse())
	}

	res = response.PaginationDataResponse{
		Data:       data,
		Amount:     len(data),
		Page:       req.Page,
		TotalPages: int(math.Ceil(float64(len(filteredGifts)) / float64(req.PageSize))),
	}

	g.redisCache.Set(redisKey, res, time.Minute*5, ctx)

	return res, nil
}

// GetGiftsOfRegion implements business.IGiftService.
func (g *giftService) GetGiftsOfRegion(region string, req request.GetGiftsRequest, ctx context.Context) (response.PaginationDataResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}

	if req.PageSize < 1 {
		req.PageSize = default_page_size
	}

	req.Keyword = util.StanderizeString(req.Keyword)
	req.SortOrder = util.StanderizeSortOrder(req.SortOrder)
	req.Category = strings.TrimSpace(req.Category)
	var res response.PaginationDataResponse
	var redisKey string = g.getGetGiftsOfegionRedisKey(region, req)
	if g.redisCache.Get(redisKey, &res, ctx) {
		return res, nil
	}

	var client = g.clients[constant.SuiTestnet]
	manageObj, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.PACKAGE_ID),
		ErrLogger: g.errLogger,
	}, ctx)
	if err != nil {
		return response.PaginationDataResponse{}, err
	}

	var centerId string
	for i, localReion := range manageObj.LocalRegions {
		if localReion == region {
			centerId = manageObj.ChildrenCenters[i]
			break
		}
	}

	if centerId == "" {
		return response.PaginationDataResponse{}, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	center, err := on_chain.GetOnChainObject[entities.Center](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  centerId,
		ErrLogger: g.errLogger,
	}, ctx)
	if err != nil {
		return response.PaginationDataResponse{}, err
	}

	gifts, err := on_chain.GetOnChainObjects[entities.Gift](on_chain.GetOnChainObjectsRequest{
		Client:    client,
		ObjectIds: center.AllGifts,
		ErrLogger: g.errLogger,
	}, ctx)
	if err != nil {
		return response.PaginationDataResponse{}, err
	}

	var filteredGifts []entities.Gift
	if req.SortOrder == "DESC" || req.SortOrder == "" {
		for i := len(gifts) - 1; i >= 0; i-- {
			var gift = gifts[i]
			if isGiftMatchedFilter(gift, req.Keyword, req.Status, req.Category) {
				filteredGifts = append(filteredGifts, gift)
			}
		}
	} else {
		for _, gift := range gifts {
			if isGiftMatchedFilter(gift, req.Keyword, req.Status, req.Category) {
				filteredGifts = append(filteredGifts, gift)
			}
		}
	}

	var skippedRecords int = (req.Page - 1) * req.PageSize
	if len(filteredGifts) <= skippedRecords {
		return response.PaginationDataResponse{}, nil
	}

	var data []response.GiftResponse
	for i := skippedRecords; i < len(filteredGifts); i++ {
		data = append(data, filteredGifts[i].ToGiftResponse())
	}

	res = response.PaginationDataResponse{
		Data:       data,
		Amount:     len(data),
		Page:       req.Page,
		TotalPages: int(math.Ceil(float64(len(filteredGifts)) / float64(req.PageSize))),
	}

	g.redisCache.Set(redisKey, res, time.Minute*5, ctx)

	return res, nil
}

func isGiftMatchedFilter(gift entities.Gift, keyword, status, category string) bool {
	if keyword != "" {
		if !strings.Contains(gift.Sender, keyword) && !strings.Contains(gift.TrackingCode, keyword) && !strings.Contains(gift.Carrier, keyword) && !strings.Contains(gift.Description, keyword) && !strings.Contains(gift.Message, keyword) {
			return false
		}
	}

	if status != "" {
		if gift.Status != status {
			return false
		}
	}

	if category != "" {
		if gift.Category != category {
			return false
		}
	}

	return true
}

func (g *giftService) getGetGiftsRedisKey(id string, req request.GetGiftsRequest) string {
	var keyword string = "empty"
	if req.Keyword != "" {
		keyword = req.Keyword
	}

	var status string = "empty"
	if req.Status != "" {
		status = req.Status
	}

	var category string = "empty"
	if req.Category != "" {
		category = req.Category
	}

	return fmt.Sprintf("gift:of:%s:kw:%s:status:%s:c:%s:o:%s:s:%d:p:%d",
		id, keyword, status, category, req.SortOrder, req.PageSize, req.Page)
}

func (g *giftService) getGetGiftsOfegionRedisKey(region string, req request.GetGiftsRequest) string {
	var keyword string = "empty"
	if req.Keyword != "" {
		keyword = req.Keyword
	}

	var status string = "empty"
	if req.Status != "" {
		status = req.Status
	}

	var category string = "empty"
	if req.Category != "" {
		category = req.Category
	}

	return fmt.Sprintf("gift:region:%s:kw:%s:status:%s:c:%s:o:%s:s:%d:p:%d",
		region, keyword, status, category, req.SortOrder, req.PageSize, req.Page)
}
