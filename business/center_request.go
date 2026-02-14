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
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/block-vision/sui-go-sdk/utils"
)

type centerRequestService struct {
	centerRequestRepo i_repository.ICenterRequestRepository
	profileRepo       i_repository.IProfileRepository
	clients           map[string]sui.ISuiAPI
	errLogger         *log.Logger
}

const (
	default_page_size int = 10
)

var (
	min_region_staffs int = 5
)

func InitializeCenterRequestService(db *sql.DB, errLogger *log.Logger) business.ICenterRequestService {
	return &centerRequestService{
		centerRequestRepo: repository.InitializeCenterRequestRepository(db, errLogger),
		profileRepo:       repository.InitializeProfileRepository(db, errLogger),
		clients:           _networkAliases,
		errLogger:         errLogger,
	}
}

func GenerateCenterRequestService() (business.ICenterRequestService, error) {
	var errLogger = util.GetLogConfig(shared.ERROR_LEVEL)

	cnn, err := db.ConnectDB(errLogger, db.InitializePostgreSQL())
	if err != nil {
		return nil, err
	}

	return InitializeCenterRequestService(cnn, errLogger), nil
}

// ConfirmRequest implements business.ICenterRequestService.
func (c *centerRequestService) ConfirmRequest(id string, ctx context.Context) (response.BuildTransactionResponse, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)

	var sender string = ctx.Value("address").(string)
	if !utils.IsValidSuiAddress(models.SuiAddress(sender)) {
		return response.BuildTransactionResponse{}, genericErr
	}

	req, err := c.centerRequestRepo.GetRequest(id, ctx)
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
	}

	var rate float32 = float32(len(req.Approvers)) / float32(len(req.Approvers)+len(req.Refusers))
	var isDenied bool = false

	if rate >= approve_rate_limit {
		req.Status = request_approved_status
		req.IsConfirmRegister = true
		if req.IsAvailableToConfirm {
			req.IsConfirmRegister = true
		}
	} else {
		req.Status = request_refused_status
		isDenied = true
	}

	req.UpdatedAt = time.Now()
	if err := c.centerRequestRepo.UpdateRegistrationRequest(*req, ctx); err != nil {
		return response.BuildTransactionResponse{}, err
	}

	if isDenied {
		return response.BuildTransactionResponse{}, nil
	}

	// Wait for background server to mint cap object to register
	if !req.IsAvailableToConfirm {
		return response.BuildTransactionResponse{}, nil
	}

	var manageModule = on_chain.InitializeModuleManage()
	var client = c.clients[constant.SuiTestnet]
	caps, err := on_chain.GetOnChainOwnedObjects[entities.Cap](on_chain.GetOnChainOwnedObjectsRequest{
		Client:       client,
		OwnerAddress: sender,
		StructType:   fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), manageModule.GetModule(), manageModule.GetUploadCenterCapStruct()),
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	if caps == nil || len(caps) == 0 {
		return response.BuildTransactionResponse{}, genericErr
	}

	var childModule = on_chain.InitializeModuleChild()
	txBytes, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
		Client:    client,
		Sender:    sender,
		Module:    childModule.GetModule(),
		Function:  childModule.GetFunctionUploadCenter(),
		ErrLogger: c.errLogger,
		Arguments: childModule.ToCreateCenterArguments(on_chain.CreateCenterArguments{
			CapID:       caps[0].ID.ID,
			Region:      req.Region,
			Address:     req.Address,
			PhoneNumber: req.PhoneNumber,
			ImageBlobID: req.ImageBlobID,
		}),
	}, ctx)

	return response.BuildTransactionResponse{
		TxBytes: txBytes,
	}, err
}

// CreateRequest implements business.ICenterRequestService.
func (c *centerRequestService) CreateRequest(req request.CreateCenterRequest, ctx context.Context) (*entities.CenterRequest, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)

	var sender string = ctx.Value("address").(string)
	if !utils.IsValidSuiAddress(models.SuiAddress(sender)) {
		return nil, genericErr
	}

	var client = c.clients[constant.SuiTestnet]
	var module = on_chain.InitializeModuleStaff()
	staffNfts, err := on_chain.GetOnChainOwnedObjects[entities.StaffNft](on_chain.GetOnChainOwnedObjectsRequest{
		Client:       client,
		OwnerAddress: sender,
		StructType:   fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), module.GetModule(), module.GetStaffNftObjectStruct()),
		ErrLogger:    c.errLogger,
	}, ctx)
	if err != nil {
		return nil, err
	}

	var genericRightErr error = errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	if staffNfts == nil || len(staffNfts) == 0 {
		return nil, genericRightErr
	}

	var isRegionStaff bool = false
	for _, nft := range staffNfts {
		if nft.Region == req.Region {
			isRegionStaff = true
			break
		}
	}

	if !isRegionStaff {
		return nil, genericRightErr
	}

	manageObj, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return nil, err
	}

	var isRegionHaveStaffs bool = false
	for i := 0; i < len(manageObj.LocalRegions); i++ {
		if manageObj.LocalRegions[i] == req.Region && !manageObj.CenterConfirmStatuses[i] {
			isRegionHaveStaffs = true
			break
		}
	}

	if !isRegionHaveStaffs {
		return nil, genericErr
	}

	var nftIds []string = append(manageObj.LocalLeaderNfts, manageObj.VolunteerNfts...)
	nfts, err := on_chain.GetOnChainObjects[entities.StaffNft](on_chain.GetOnChainObjectsRequest{
		Client:    client,
		ObjectIds: nftIds,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return nil, err
	}

	var staffCount int = 0
	for _, nft := range nfts {
		if nft.Region == req.Region {
			staffCount++
			if staffCount == min_region_staffs {
				break
			}
		}
	}

	if staffCount < min_region_staffs {
		return nil, genericErr
	}

	// todo: AI validate address if matched with region
	var address string = strings.TrimSpace(req.Address)
	var phoneNumber string = strings.TrimSpace(req.PhoneNumber)
	var curTime time.Time = time.Now()
	var request = entities.CenterRequest{
		ID:          util.GenerateId(),
		Region:      req.Region,
		Address:     address,
		PhoneNumber: phoneNumber,
		ImageBlobID: req.ImageBlobID,
		Status:      request_pending_status,
		CreatedBy:   sender,
		CreatedAt:   curTime,
		UpdatedAt:   curTime,
		ClosedAt:    util.GetRequestDuration(),
	}

	return &request, c.centerRequestRepo.CreateRegistrationRequest(request, ctx)
}

// GetRequest implements business.ICenterRequestService.
func (c *centerRequestService) GetRequest(id string, ctx context.Context) (*entities.CenterRequest, error) {
	if id == "" {
		return nil, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	return c.centerRequestRepo.GetRequest(id, ctx)
}

// GetRequests implements business.ICenterRequestService.
func (c *centerRequestService) GetRequests(req request.GetCenterRequests, ctx context.Context) (response.PaginationDataResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}

	if req.PageSize < 1 {
		req.PageSize = default_page_size
	}

	data, pages, err := c.centerRequestRepo.GetRegistrationRequests(req, ctx)

	return response.PaginationDataResponse{
		Data:       data,
		Page:       req.Page,
		TotalPages: pages,
	}, err
}

// GetWalletRequests implements business.ICenterRequestService.
func (c *centerRequestService) GetWalletRequests(id string, ctx context.Context) ([]entities.CenterRequest, error) {
	if !utils.IsValidSuiAddress(models.SuiAddress(id)) {
		return nil, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	return c.centerRequestRepo.GetWalletRegistrationRequests(id, ctx)
}

// VoteRequest implements business.ICenterRequestService.
func (c *centerRequestService) VoteRequest(id string, req request.VoteRequest, ctx context.Context) error {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)

	var voter string = ctx.Value("address").(string)
	if !utils.IsValidSuiAddress(models.SuiAddress(voter)) {
		return genericErr
	}

	request, err := c.centerRequestRepo.GetRequest(id, ctx)
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

	var module = on_chain.InitializeModuleStaff()
	nfts, err := on_chain.GetOnChainOwnedObjects[entities.StaffNft](on_chain.GetOnChainOwnedObjectsRequest{
		Client:       c.clients[constant.SuiTestnet],
		OwnerAddress: voter,
		StructType:   fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), module.GetModule(), module.GetStaffNftObjectStruct()),
		ErrLogger:    c.errLogger,
	}, ctx)
	if err != nil {
		return err
	}

	var genericRightErr error = errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	if nfts == nil || len(nfts) == 0 {
		return genericRightErr
	}

	var isRegionStaff bool = false
	for _, nft := range nfts {
		if nft.Role == request.Region {
			isRegionStaff = true
			break
		}
	}

	if !isRegionStaff {
		return genericRightErr
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

	return c.centerRequestRepo.UpdateRegistrationRequest(*request, ctx)
}
