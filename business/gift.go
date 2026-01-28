package business

import (
	"context"
	"database/sql"
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
	"raise-child/util/db"
	on_chain "raise-child/util/on_chain"
	"strings"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/block-vision/sui-go-sdk/utils"
)

type giftService struct {
	profileRepo i_repository.IProfileRepository
	clients     map[string]sui.ISuiAPI
	errLogger   *log.Logger
}

func InitializeGiftService(db *sql.DB, errLogger *log.Logger) business.IGiftService {
	return &giftService{
		profileRepo: repository.InitializeProfileRepository(db, errLogger),
		clients:     _networkAliases,
		errLogger:   errLogger,
	}
}

func GenerateGiftService() (business.IGiftService, error) {
	var errLogger = util.GetLogConfig(shared.ERROR_LEVEL)

	cnn, err := db.ConnectDB(errLogger, db.InitializePostgreSQL())
	if err != nil {
		return nil, err
	}

	return InitializeGiftService(cnn, errLogger), nil
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

	child, err := on_chain.GetOnChainObject[entities.Child](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  gift.ChildID,
		ErrLogger: g.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	if child == nil {
		return response.BuildTransactionResponse{}, genericErr
	}

	var staffId string
	for _, role := range roles {
		if role.Region == child.Region {
			staffId = role.ID.ID
			break
		}
	}

	if staffId == "" {
		return response.BuildTransactionResponse{}, genericRightErr
	}

	var giftModule = on_chain.InitializeModuleGift()
	txBytes, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
		Client:    client,
		Sender:    sender,
		Module:    giftModule.GetModule(),
		Function:  giftModule.GetFunctionConfirmReceiveGift(),
		ErrLogger: g.errLogger,
		Arguments: giftModule.ToConfirmReceiveGiftArguments(on_chain.ConfirmReceiveGiftArguments{
			GiftID:      id,
			ChildID:     child.ID.ID,
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

	if !utils.IsValidSuiAddress(models.SuiAddress(req.ChildID)) {
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
	var sponsorModule = on_chain.InitializeModuleSponsor()
	nfts, err := on_chain.GetOnChainOwnedObjects[entities.Sponsor](on_chain.GetOnChainOwnedObjectsRequest{
		Client:       client,
		OwnerAddress: sender,
		StructType:   fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), sponsorModule.GetModule(), sponsorModule.GetSponsorNftStruct()),
		ErrLogger:    g.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	var nftId string
	if nfts == nil || len(nfts) == 0 {
		nftId = os.Getenv(env.PUB_SPONSOR_NFT_ID)
	} else {
		nftId = nfts[0].ID.ID
	}

	var giftModule = on_chain.InitializeModuleGift()
	txBytes, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
		Client:    client,
		Sender:    sender,
		Module:    giftModule.GetModule(),
		Function:  giftModule.GetFunctionCreateGift(),
		ErrLogger: g.errLogger,
		Arguments: giftModule.ToCreateGiftArguments(on_chain.CreateGiftArguments{
			SponsorID:       nftId,
			ChildID:         req.ChildID,
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
		}),
	}, ctx)

	return response.BuildTransactionResponse{
		TxBytes: txBytes,
	}, err
}

// GetGift implements business.IGiftService.
func (g *giftService) GetGift(id string, ctx context.Context) (response.GiftResponse, error) {
	if !utils.IsValidSuiAddress(models.SuiAddress(id)) {
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
	if !utils.IsValidSuiAddress(models.SuiAddress(id)) {
		return response.PaginationDataResponse{}, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	var client = g.clients[constant.SuiTestnet]
	child, err := on_chain.GetOnChainObject[entities.Child](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  id,
		ErrLogger: g.errLogger,
	}, ctx)
	if err != nil {
		return response.PaginationDataResponse{}, err
	}

	if child.Gifts == nil || len(child.Gifts) == 0 {
		return response.PaginationDataResponse{}, nil
	}

	gifts, err := on_chain.GetOnChainObjects[entities.Gift](on_chain.GetOnChainObjectsRequest{
		Client:    client,
		ObjectIds: child.Gifts,
		ErrLogger: g.errLogger,
	}, ctx)
	if err != nil {
		return response.PaginationDataResponse{}, err
	}

	var filteredGifts []entities.Gift
	var keyword string = util.StanderizeString(req.Keyword)
	if req.SortOrder == "desc" || req.SortOrder == "" {
		for i := len(gifts) - 1; i >= 0; i-- {
			var gift = gifts[i]
			if isGiftMatchedFilter(gift, keyword, req.Status, req.Category) {
				filteredGifts = append(filteredGifts, gift)
			}
		}
	} else {
		for _, gift := range gifts {
			if isGiftMatchedFilter(gift, keyword, req.Status, req.Category) {
				filteredGifts = append(filteredGifts, gift)
			}
		}
	}

	var page int = req.Page
	if page < 1 {
		page = 1
	}

	var skippedRecords int = (page - 1) * gift_limit_record
	if len(filteredGifts) <= skippedRecords {
		return response.PaginationDataResponse{}, nil
	}

	var data []response.GiftResponse
	for i := skippedRecords; i < len(filteredGifts); i++ {
		data = append(data, filteredGifts[i].ToGiftResponse())
	}

	return response.PaginationDataResponse{
		Data:       data,
		Page:       page,
		TotalPages: int(math.Ceil(float64(len(filteredGifts)) / float64(gift_limit_record))),
	}, nil
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
