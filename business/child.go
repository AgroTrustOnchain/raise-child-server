package business

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"raise-child/constants/env"
	"raise-child/constants/env/payment"
	"raise-child/constants/noti"
	internal_sui "raise-child/constants/on-chain/sui"
	"raise-child/repository"
	"sort"
	"strconv"
	"strings"
	"time"

	"raise-child/constants/shared"
	"raise-child/interfaces/business"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/model/entities"
	"raise-child/util"
	"raise-child/util/ai"
	"raise-child/util/cache"
	"raise-child/util/db"
	on_chain "raise-child/util/on_chain"
	walrus_pkg "raise-child/util/walrus_pkg"
	"slices"

	i_repository "raise-child/interfaces/repository"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/payOSHQ/payos-lib-golang"
)

type childService struct {
	pendingChildSpecialNeedProposalRepo i_repository.IPendingChildSpecialNeedProposalRepository
	pendingWithdrawProposalRepo         i_repository.IPendingWithdrawProposalRepository
	withdrawRepo                        i_repository.IOffChainWithdrawProposalRepository
	donationRepo                        i_repository.IOffChainDonationRepository
	mealDurationRepo                    i_repository.IMealSupportDurationRepository
	paymentRepo                         i_repository.IPaymentRepository
	profileRepo                         i_repository.IProfileRepository
	bankRepo                            i_repository.IBankProfileRepository
	leaderNotiRepo                      i_repository.ILeaderNotiRepository
	aiProvider                          ai.IAiClientProvider
	walrusProvider                      walrus_pkg.IWalrusProvider
	redisCache                          cache.IRedisCache
	clients                             map[string]sui.ISuiAPI
	errLogger                           *log.Logger
}

func initializeChildService(
	pendingChildSpecialNeedProposalRepo i_repository.IPendingChildSpecialNeedProposalRepository,
	pendingWithdrawProposalRepo i_repository.IPendingWithdrawProposalRepository,
	withdrawRepo i_repository.IOffChainWithdrawProposalRepository,
	donationRepo i_repository.IOffChainDonationRepository,
	mealDurationRepo i_repository.IMealSupportDurationRepository,
	paymentRepo i_repository.IPaymentRepository,
	profileRepo i_repository.IProfileRepository,
	bankRepo i_repository.IBankProfileRepository,
	leaderNotiRepo i_repository.ILeaderNotiRepository,
	aiProvider ai.IAiClientProvider,
	walrusProvider walrus_pkg.IWalrusProvider,
	clients map[string]sui.ISuiAPI,
	errLogger *log.Logger,
) business.IChildService {
	return &childService{
		pendingChildSpecialNeedProposalRepo: pendingChildSpecialNeedProposalRepo,
		pendingWithdrawProposalRepo:         pendingWithdrawProposalRepo,
		withdrawRepo:                        withdrawRepo,
		donationRepo:                        donationRepo,
		mealDurationRepo:                    mealDurationRepo,
		paymentRepo:                         paymentRepo,
		profileRepo:                         profileRepo,
		bankRepo:                            bankRepo,
		leaderNotiRepo:                      leaderNotiRepo,
		aiProvider:                          aiProvider,
		walrusProvider:                      walrusProvider,
		redisCache:                          cache.InitializeRedisCache(),
		clients:                             clients,
		errLogger:                           errLogger,
	}
}

func GenerateChildService() (business.IChildService, error) {
	var errLogger = util.GetLogConfig(shared.ERROR_LEVEL)

	cnn, err := db.ConnectDB(errLogger, db.InitializePostgreSQL())
	if err != nil {
		return nil, err
	}

	return initializeChildService(
		repository.InitializePendingChildSpecialNeedProposalRepo(cnn, errLogger),
		repository.InitializePendingWithdrawProposalRepo(cnn, errLogger),
		repository.InitializeOffChainWithdrawProposalRepository(cnn, errLogger),
		repository.InitializeOffChainDonationRepository(cnn, errLogger),
		repository.InitializeMealSupportDurationRepository(cnn, errLogger),
		repository.InitializePaymentRepository(cnn, errLogger),
		repository.InitializeProfileRepository(cnn, errLogger),
		repository.InitializeBankProfileRepository(cnn, errLogger),
		repository.InitializeLeaderNotiRepository(cnn, errLogger),
		ai.InitializeAiProvider(nil, errLogger),
		walrus_pkg.InitializeWalrusProvider(errLogger),
		_networkAliases,
		errLogger,
	), nil
}

// GetChild implements business.IChildService.
func (c *childService) GetChild(id string, ctx context.Context) (response.ChildResponse, error) {
	if !util.IsValidSuiAddressStrict(id) {
		return response.ChildResponse{}, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	var res response.ChildResponse
	var redisKey string = c.getGetChildRediskey(id)
	if c.redisCache.Get(redisKey, &res, ctx) {
		return res, nil
	}

	var client = c.clients[constant.SuiTestnet]
	child, err := on_chain.GetOnChainObject[entities.Child](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  id,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return response.ChildResponse{}, err
	}

	res = child.ToChildResponse()
	if len(res.DynamicFields) > 0 {
		// Has dynamic fields
		if dynamicValues, _ := on_chain.GetDynamicFields(id, client, c.errLogger, ctx); dynamicValues != nil {
			res.DynamicValues = dynamicValues
		}
	}

	c.redisCache.Set(redisKey, res, time.Minute*5, ctx)

	return res, err
}

// GetChilds implements business.IChildService.
func (c *childService) GetChildren(req request.GetChildrenRequest, ctx context.Context) (response.PaginationDataResponse, error) {
	req.SortOrder = util.StanderizeSortOrder(req.SortOrder)
	req.Keyword = util.StanderizeString(req.Keyword)
	if req.Page < 1 {
		req.Page = 1
	}

	if req.PageSize < 1 {
		req.PageSize = default_page_size
	}

	var res response.PaginationDataResponse
	var redisKey string = c.getGetChildrenRediskey(req)
	if c.redisCache.Get(redisKey, &res, ctx) {
		return res, nil
	}

	var client = c.clients[constant.SuiTestnet]
	manageObj, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return response.PaginationDataResponse{}, err
	}

	children, err := on_chain.GetOnChainObjects[entities.Child](on_chain.GetOnChainObjectsRequest{
		Client:    client,
		ObjectIds: manageObj.ChildIds,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return response.PaginationDataResponse{}, err
	}

	if children == nil {
		return response.PaginationDataResponse{}, nil
	}

	if req.Page < 1 {
		req.Page = 1
	}

	if req.PageSize < 1 {
		req.PageSize = default_page_size
	}

	var keyword string = util.StanderizeString(req.Keyword)
	var region string = util.StanderizeString(req.Region)
	var filteredChildren []entities.Child
	for i := len(children) - 1; i >= 0; i-- {
		var child entities.Child = children[i]

		if region != "" {
			if util.StanderizeString(child.Region) != region { // Not matched
				continue
			}
		}

		if req.Gender != "" {
			if child.Gender != req.Gender { // Not matched
				continue
			}
		}

		if req.YearOfBirth != nil {
			var dob time.Time = util.RawDateToTime(child.DateOfBirth)
			if dob.Year() != *req.YearOfBirth { // Not matched
				continue
			}
		}

		if keyword != "" {
			var firstName string = util.StanderizeString(child.FirstName)
			var lastName string = util.StanderizeString(child.LastName)
			if !strings.Contains(firstName, keyword) && !strings.Contains(lastName, keyword) && !strings.Contains(child.IdentityCode, keyword) { // Not matched
				continue
			}
		}

		filteredChildren = append(filteredChildren, child)
	}

	sort.Slice(filteredChildren, func(i, j int) bool {
		var name1 string = filteredChildren[i].LastName + " " + filteredChildren[i].FirstName
		var name2 string = filteredChildren[j].LastName + " " + filteredChildren[j].FirstName

		if req.SortOrder == "ASC" {
			return name1 < name2
		}

		return name2 > name1
	})

	var skippedRecords int = (req.Page - 1) * req.PageSize
	if len(filteredChildren) <= skippedRecords {
		return response.PaginationDataResponse{}, nil
	}

	var data []response.ChildResponse
	for i := skippedRecords; i < len(filteredChildren); i++ {
		data = append(data, filteredChildren[i].ToMinimumChildResponse())
		if len(data) == req.PageSize {
			break
		}
	}

	res = response.PaginationDataResponse{
		Data:       data,
		Amount:     len(data),
		Page:       req.Page,
		TotalPages: int(math.Ceil(float64(len(filteredChildren)) / float64(req.PageSize))),
	}

	c.redisCache.Set(redisKey, res, time.Minute*5, ctx)
	return res, nil
}

// UploadChild implements business.IChildService.
func (c *childService) UploadChild(req request.UploadChildRequest, ctx context.Context) (response.BuildTransactionResponse, error) {
	// todo: validate if this child is existed or not
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)

	var rawDate string = strings.TrimSpace(req.DateOfBirth)
	if dob := util.RawDateToTime(rawDate); dob.IsZero() { // Invalid date
		return response.BuildTransactionResponse{}, genericErr
	}

	var gender string = util.StanderizeGender(req.Gender)
	if gender == "" {
		return response.BuildTransactionResponse{}, genericErr
	}

	var module = on_chain.InitializeModuleChild()
	res, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
		Client:    c.clients[constant.SuiTestnet],
		Sender:    ctx.Value("address").(string),
		Module:    module.GetModule(),
		Function:  module.GetFunctionAddChild(),
		ErrLogger: c.errLogger,
		Arguments: module.ToAddChildArguments(on_chain.AddChildArguments{
			IdentityCode: util.StanderizeString(req.IdentityCode),
			FirstName:    util.StanderizeString(req.FirstName),
			LastName:     util.StanderizeString(req.LastName),
			Gender:       gender,
			DateOfBirth:  rawDate,
			AvatarBlobId: util.StanderizeString(req.AvatarBlobId),
		}),
	}, ctx)

	return response.BuildTransactionResponse{
		TxBytes: res,
	}, err

}

// AddNumberMetada implements business.IChildService.
func (c *childService) AddNumberMetada(id string, req request.AddChildNumberMetadaRequest, ctx context.Context) (response.BuildTransactionResponse, error) {
	if !util.IsValidSuiAddressStrict(id) {
		return response.BuildTransactionResponse{}, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	var client = c.clients[constant.SuiTestnet]
	child, err := getOnChainObject[entities.Child](client, id, c.errLogger, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	var key string = util.StanderizeString(req.Key)
	if existed := slices.Contains(child.DynamicFields, key); existed { // Field existed
		return response.BuildTransactionResponse{}, errors.New(noti.METADATA_EXISTED_MESSAGE)
	}

	var module = on_chain.InitializeModuleChild()
	res, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
		Client:    client,
		Sender:    ctx.Value("address").(string),
		Module:    module.GetModule(),
		Function:  module.GetFunctionAddNumberMetadata(),
		ErrLogger: c.errLogger,
		Arguments: []interface{}{
			id,
			key,
			req.Value,
			internal_sui.CLOCK_OBJECT_ID,
		},
	}, ctx)

	return response.BuildTransactionResponse{
		TxBytes: res,
	}, err
}

// AddStringMetada implements business.IChildService.
func (c *childService) AddStringMetada(id string, req request.AddChildStringMetadaRequest, ctx context.Context) (response.BuildTransactionResponse, error) {
	if !util.IsValidSuiAddressStrict(id) {
		return response.BuildTransactionResponse{}, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	var client = c.clients[constant.SuiTestnet]
	child, err := getOnChainObject[entities.Child](client, id, c.errLogger, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	var key string = util.StanderizeString(req.Key)
	if existed := slices.Contains(child.DynamicFields, key); existed { // Field existed
		return response.BuildTransactionResponse{}, errors.New(noti.METADATA_EXISTED_MESSAGE)
	}

	var module = on_chain.InitializeModuleChild()
	res, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
		Client:    client,
		Sender:    ctx.Value("address").(string),
		Module:    module.GetModule(),
		Function:  module.GetFunctionAddStringMetadata(),
		ErrLogger: c.errLogger,
		Arguments: []interface{}{
			id,
			key,
			util.StanderizeString(req.Value),
			internal_sui.CLOCK_OBJECT_ID,
		},
	}, ctx)

	return response.BuildTransactionResponse{
		TxBytes: res,
	}, err
}

// CreateBooksNeedWithdrawProposal implements business.IChildService.
func (c *childService) CreateBooksNeedWithdrawProposal(req request.CreateNormalNeedWithdrawProposalRequest, ctx context.Context) (response.BuildTransactionResponse, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if !util.IsValidSuiAddressStrict(req.NeedID) {
		return response.BuildTransactionResponse{}, genericErr
	}

	var client = c.clients[constant.SuiTestnet]
	need, err := on_chain.GetOnChainObject[entities.BooksNeed](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  req.NeedID,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	// Already withdraw all
	if len(need.Donations) == len(need.WithdrawsForNeed) {
		return response.BuildTransactionResponse{}, errors.New("")
	}

	child, err := on_chain.GetOnChainObject[entities.Child](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  need.ChildID,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	var staffModule = on_chain.InitializeModuleStaff()
	var sender string = ctx.Value("address").(string)
	staffNfts, err := on_chain.GetOnChainOwnedObjects[entities.StaffNft](on_chain.GetOnChainOwnedObjectsRequest{
		Client:       client,
		OwnerAddress: sender,
		StructType:   fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), staffModule.GetModule(), staffModule.GetStaffNftObjectStruct()),
		ErrLogger:    c.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	var genericRightErr error = errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	if staffNfts == nil || len(staffNfts) == 0 {
		return response.BuildTransactionResponse{}, genericRightErr
	}

	var isLeaderOfRegion bool = false
	for _, nft := range staffNfts {
		if nft.Role == local_leader_role && nft.Region == child.Region {
			isLeaderOfRegion = true
			break
		}
	}

	if !isLeaderOfRegion {
		return response.BuildTransactionResponse{}, genericRightErr
	}

	pool, err := on_chain.GetOnChainObject[entities.MainPool](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.POOL_ID),
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	localPools, err := on_chain.GetOnChainObjects[entities.LocalPool](on_chain.GetOnChainObjectsRequest{
		Client:    client,
		ObjectIds: pool.LocalPools,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	var localPoolId string
	for _, localPool := range localPools {
		if localPool.Region == child.Region {
			localPoolId = localPool.ID.ID
			break
		}
	}

	var childModule = on_chain.InitializeModuleChild()
	txBytes, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
		Client:    client,
		Sender:    sender,
		Module:    childModule.GetModule(),
		Function:  childModule.GetFunctionCreateChildBooksNeedWithdrawProposal(),
		ErrLogger: c.errLogger,
		Arguments: childModule.ToCreateChildNormalNeedWithdrawProposalArguments(on_chain.CreateChildNormalNeedWithdrawProposalArguments{
			NeedID:      req.NeedID,
			ChildID:     need.ChildID,
			LocalPool:   localPoolId,
			Description: fmt.Sprintf("Withdraw Books Need Semester %s - %s for child %s %s", need.Semster, need.Year, child.LastName, child.FirstName),
			ClosedAt:    util.ToMilliseconds(util.GetRequestDuration()),
		}),
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	var proposalId string = util.GenerateId()
	return response.BuildTransactionResponse{
			TxBytes:    txBytes,
			ProposalId: proposalId,
		}, c.withdrawRepo.CreateOffChainWithdrawProposal(entities.OffChainWithdrawProposal{
			ID:        proposalId,
			Purpose:   string(entities.BOOKS_NEED_PURPOSE),
			Target:    req.NeedID,
			CreatedAt: time.Now(),
		}, ctx)
}

// CreateMealNeedWithdrawProposal implements business.IChildService.
func (c *childService) CreateMealNeedWithdrawProposal(req request.CreateNormalNeedWithdrawProposalRequest, ctx context.Context) (response.BuildTransactionResponse, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if !util.IsValidSuiAddressStrict(req.NeedID) {
		return response.BuildTransactionResponse{}, genericErr
	}

	var client = c.clients[constant.SuiTestnet]
	need, err := on_chain.GetOnChainObject[entities.MealNeed](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  req.NeedID,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	totalSupportedMonths, _ := strconv.Atoi(need.TotalSupportedMonths)
	var expectedDuration int = totalSupportedMonths - len(need.WithdrawsForNeed)
	// Already withdraw all
	if expectedDuration == 0 {
		return response.BuildTransactionResponse{}, errors.New("")
	}

	var previousDuration int = 0
	var expectedDay, expectedMonth int
	var startDate, endDate time.Time
	var curTime time.Time = time.Now()
	for i := len(need.Durations) - 1; i >= 0; i-- {
		var duration = need.Durations[0]
		var startPeriod time.Time = util.RawDateToTime(duration.Fields.StartPeriod)
		var endPeriod time.Time = util.RawDateToTime(duration.Fields.EndPeriod)
		var startMonth int = int(startPeriod.Month())
		var endMonth int = int(endPeriod.Month())
		if endMonth == 1 { // To next year
			endMonth = 13
		}

		var currentDuration int = endMonth - startMonth
		var totalDuration int = currentDuration + previousDuration
		var months int = totalDuration - expectedDuration
		if months >= 0 {
			startDate = startPeriod.AddDate(0, months, 0)
			endDate = startDate.AddDate(0, 1, 0)

			var expectedDate = startDate.AddDate(0, -3, 0)
			expectedDay = expectedDate.Day()
			expectedMonth = int(expectedDate.Month())
			break
		}

		previousDuration = totalDuration
	}

	// Still not date to withdraw
	if int(curTime.Month()) != expectedMonth || curTime.Day() != expectedDay {
		return response.BuildTransactionResponse{}, errors.New("")
	}

	child, err := on_chain.GetOnChainObject[entities.Child](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  need.ChildID,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	var sender string = ctx.Value("address").(string)
	var staffModule = on_chain.InitializeModuleStaff()
	staffNfts, err := on_chain.GetOnChainOwnedObjects[entities.StaffNft](on_chain.GetOnChainOwnedObjectsRequest{
		Client:       client,
		OwnerAddress: sender,
		StructType:   fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), staffModule.GetModule(), staffModule.GetStaffNftObjectStruct()),
		ErrLogger:    c.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	var genericRightErr error = errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	if staffNfts == nil || len(staffNfts) == 0 {
		return response.BuildTransactionResponse{}, genericRightErr
	}

	var isLeaderOfRegion bool = false
	for _, nft := range staffNfts {
		if nft.Role == local_leader_role && nft.Region == child.Region {
			isLeaderOfRegion = true
			break
		}
	}

	if !isLeaderOfRegion {
		return response.BuildTransactionResponse{}, genericRightErr
	}

	pool, err := on_chain.GetOnChainObject[entities.MainPool](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.POOL_ID),
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	localPools, err := on_chain.GetOnChainObjects[entities.LocalPool](on_chain.GetOnChainObjectsRequest{
		Client:    client,
		ObjectIds: pool.LocalPools,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	var localPoolId string
	for _, localPool := range localPools {
		if localPool.Region == child.Region {
			localPoolId = localPool.ID.ID
			break
		}
	}

	var childModule = on_chain.InitializeModuleChild()
	txBytes, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
		Client:    client,
		Sender:    sender,
		Module:    childModule.GetModule(),
		Function:  childModule.GetFunctionCreateChildMealNeedWithdrawProposal(),
		ErrLogger: c.errLogger,
		Arguments: childModule.ToCreateChildNormalNeedWithdrawProposalArguments(on_chain.CreateChildNormalNeedWithdrawProposalArguments{
			NeedID:      req.NeedID,
			ChildID:     need.ChildID,
			LocalPool:   localPoolId,
			Description: fmt.Sprintf("Withdraw Meal Need %s - %s for child %s %s", util.TimeToRawDate(startDate), util.TimeToRawDate(endDate), child.LastName, child.FirstName),
			ClosedAt:    util.ToMilliseconds(util.GetRequestDuration()),
		}),
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	var proposalId string = util.GenerateId()
	return response.BuildTransactionResponse{
			TxBytes:    txBytes,
			ProposalId: proposalId,
		}, c.withdrawRepo.CreateOffChainWithdrawProposal(entities.OffChainWithdrawProposal{
			ID:        proposalId,
			Purpose:   string(entities.MEAL_NEED_PURPOSE),
			Target:    req.NeedID,
			CreatedAt: curTime,
		}, ctx)
}

// CreateBooksNeedWithdrawProposalV2 implements business.IChildService.
func (c *childService) CreateBooksNeedWithdrawProposalV2(req request.CreateNormalNeedWithdrawProposalRequest, ctx context.Context) (*entities.PendingWithdrawProposal, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if !util.IsValidSuiAddressStrict(req.NeedID) {
		return nil, genericErr
	}

	var client = c.clients[constant.SuiTestnet]
	need, err := on_chain.GetOnChainObject[entities.BooksNeed](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  req.NeedID,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return nil, err
	}

	if need == nil {
		return nil, genericErr
	}

	// Already withdraw all
	if len(need.Donations) == len(need.WithdrawsForNeed) {
		return nil, errors.New(noti.NEED_WITHDRAWN_MESSAGE)
	}

	child, err := on_chain.GetOnChainObject[entities.Child](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  need.ChildID,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return nil, err
	}

	if child == nil {
		return nil, genericErr
	}

	var staffModule = on_chain.InitializeModuleStaff()
	var sender string = ctx.Value("address").(string)
	staffNfts, err := on_chain.GetOnChainOwnedObjects[entities.StaffNft](on_chain.GetOnChainOwnedObjectsRequest{
		Client:       client,
		OwnerAddress: sender,
		StructType:   fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), staffModule.GetModule(), staffModule.GetStaffNftObjectStruct()),
		ErrLogger:    c.errLogger,
	}, ctx)
	if err != nil {
		return nil, err
	}

	var genericRightErr error = errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	if staffNfts == nil || len(staffNfts) == 0 {
		return nil, genericRightErr
	}

	var isLeaderOfRegion bool = false
	for _, nft := range staffNfts {
		if nft.Role == local_leader_role && nft.Region == child.Region {
			isLeaderOfRegion = true
			break
		}
	}

	if !isLeaderOfRegion {
		return nil, genericRightErr
	}

	leaderNoti, err := c.leaderNotiRepo.GetNotiByNeed(req.NeedID, ctx)
	if err != nil {
		return nil, err
	}

	var curTime time.Time = time.Now()
	var expectedStartDate time.Time
	var index int = -1
	for i := len(leaderNoti.ExpectedWithdrawPeriods) - 1; i >= 0; i-- {
		var rawExpectedDate string = leaderNoti.ExpectedWithdrawPeriods[i]
		var expectedDate time.Time = util.ToStartOfDate(util.RawDateToTime(rawExpectedDate))

		if !curTime.Before(expectedDate) {
			expectedStartDate = expectedDate
			index = i
			break
		}
	}

	var notWithdrawDateErr error = errors.New(noti.NOT_WITHDRAW_EXPECTED_DATE_MESSAGE)
	if index == -1 {
		return nil, notWithdrawDateErr
	}

	var expectedEndDate time.Time = util.ToEndOfDate(expectedStartDate.AddDate(0, 0, 7))
	if curTime.Before(expectedStartDate) || curTime.After(expectedEndDate) {
		c.errLogger.Println("Die at below")
		return nil, notWithdrawDateErr
	}

	var description string = leaderNoti.Contents[index]
	var purpose string = string(entities.BOOKS_NEED_PURPOSE)
	withdrawAmount, _ := strconv.ParseInt(need.Value, 10, 64)
	isProposed, err := c.pendingWithdrawProposalRepo.IsPendingWithdrawProposalProposedWithSpecificInfo(purpose, req.NeedID, description, withdrawAmount, ctx)
	if err != nil {
		return nil, err
	}

	if isProposed {
		return nil, errors.New(noti.STILL_PENDING_REQUEST_MESSAGE)
	}

	pool, err := on_chain.GetOnChainObject[entities.MainPool](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.POOL_ID),
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return nil, err
	}

	localPools, err := on_chain.GetOnChainObjects[entities.LocalPool](on_chain.GetOnChainObjectsRequest{
		Client:    client,
		ObjectIds: pool.LocalPools,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return nil, err
	}

	var localPoolId, localPoolName string
	for _, localPool := range localPools {
		if localPool.Region == child.Region {
			localPoolId = localPool.ID.ID
			localPoolName = localPool.Region
			break
		}
	}

	var aiEvaluation string
	if req.ProofBlobID != nil {
		proofBytes, _ := c.walrusProvider.FetchBytesImage(*req.ProofBlobID)
		if proofBytes != nil {
			aiEvaluation = c.aiProvider.ValidateWithdrawProposal(ai.ValidateWithdrawProposal{
				Purpose:         purpose,
				WithdrawAmount:  withdrawAmount,
				Description:     description,
				ProofBytesImage: proofBytes,
			}, ctx)
		}
	}

	// todo: AI validation
	var res = entities.PendingWithdrawProposal{
		ID:             util.GenerateId(),
		ProfileID:      ctx.Value("sub").(string),
		Creator:        sender,
		PoolID:         localPoolId,
		PoolName:       localPoolName,
		Purpose:        purpose,
		Target:         req.NeedID,
		WithdrawAmount: withdrawAmount,
		ProofBlobID:    req.ProofBlobID,
		Description:    description,
		Status:         request_pending_status,
		AIEvaluation:   aiEvaluation,
		CreatedAt:      curTime,
		UpdatedAt:      curTime,
	}

	return &res, c.pendingWithdrawProposalRepo.CreatePendingWithdrawProposal(res, ctx)
}

// CreateHealthInsuranceNeedWithdrawProposalV2 implements business.IChildService.
func (c *childService) CreateHealthInsuranceNeedWithdrawProposalV2(req request.CreateNormalNeedWithdrawProposalRequest, ctx context.Context) (*entities.PendingWithdrawProposal, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if !util.IsValidSuiAddressStrict(req.NeedID) {
		return nil, genericErr
	}

	var client = c.clients[constant.SuiTestnet]
	need, err := on_chain.GetOnChainObject[entities.HealthInsuranceNeed](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  req.NeedID,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return nil, err
	}

	if need == nil {
		return nil, genericErr
	}

	// Already withdraw all
	if len(need.Donations) == len(need.WithdrawsForNeed) {
		return nil, errors.New(noti.NEED_WITHDRAWN_MESSAGE)
	}

	child, err := on_chain.GetOnChainObject[entities.Child](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  need.ChildID,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return nil, err
	}

	if child == nil {
		return nil, genericErr
	}

	var staffModule = on_chain.InitializeModuleStaff()
	var sender string = ctx.Value("address").(string)
	staffNfts, err := on_chain.GetOnChainOwnedObjects[entities.StaffNft](on_chain.GetOnChainOwnedObjectsRequest{
		Client:       client,
		OwnerAddress: sender,
		StructType:   fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), staffModule.GetModule(), staffModule.GetStaffNftObjectStruct()),
		ErrLogger:    c.errLogger,
	}, ctx)
	if err != nil {
		return nil, err
	}

	var genericRightErr error = errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	if staffNfts == nil || len(staffNfts) == 0 {
		return nil, genericRightErr
	}

	var isLeaderOfRegion bool = false
	for _, nft := range staffNfts {
		if nft.Role == local_leader_role && nft.Region == child.Region {
			isLeaderOfRegion = true
			break
		}
	}

	if !isLeaderOfRegion {
		return nil, genericRightErr
	}

	leaderNoti, err := c.leaderNotiRepo.GetNotiByNeed(req.NeedID, ctx)
	if err != nil {
		return nil, err
	}

	var curTime time.Time = time.Now()
	var index int = len(leaderNoti.ExpectedWithdrawPeriods) - 1
	var expectedStartDate time.Time = util.ToStartOfDate(util.RawDateToTime(leaderNoti.ExpectedWithdrawPeriods[index]))
	var expectedEndDate time.Time = util.ToEndOfDate(expectedStartDate.AddDate(0, 0, 7))
	if curTime.Before(expectedStartDate) || curTime.After(expectedEndDate) {
		return nil, errors.New(noti.NOT_WITHDRAW_EXPECTED_DATE_MESSAGE)
	}

	var description string = leaderNoti.Contents[index]
	var purpose string = string(entities.HEALTH_INSURANCE_NEED_PURPOSE)
	withdrawAmount, _ := strconv.ParseInt(need.Value, 10, 64)
	isProposed, err := c.pendingWithdrawProposalRepo.IsPendingWithdrawProposalProposedWithSpecificInfo(purpose, req.NeedID, description, withdrawAmount, ctx)
	if err != nil {
		return nil, err
	}

	if isProposed {
		return nil, errors.New(noti.STILL_PENDING_REQUEST_MESSAGE)
	}

	pool, err := on_chain.GetOnChainObject[entities.MainPool](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.POOL_ID),
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return nil, err
	}

	localPools, err := on_chain.GetOnChainObjects[entities.LocalPool](on_chain.GetOnChainObjectsRequest{
		Client:    client,
		ObjectIds: pool.LocalPools,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return nil, err
	}

	var localPoolId, localPoolName string
	for _, localPool := range localPools {
		if localPool.Region == child.Region {
			localPoolId = localPool.ID.ID
			localPoolName = localPool.Region
			break
		}
	}

	var aiEvaluation string
	if req.ProofBlobID != nil {
		proofBytes, _ := c.walrusProvider.FetchBytesImage(*req.ProofBlobID)
		if proofBytes != nil {
			aiEvaluation = c.aiProvider.ValidateWithdrawProposal(ai.ValidateWithdrawProposal{
				Purpose:         purpose,
				WithdrawAmount:  withdrawAmount,
				Description:     description,
				ProofBytesImage: proofBytes,
			}, ctx)
		}
	}

	// todo: AI validation
	var res = entities.PendingWithdrawProposal{
		ID:             util.GenerateId(),
		ProfileID:      ctx.Value("sub").(string),
		Creator:        sender,
		PoolID:         localPoolId,
		PoolName:       localPoolName,
		Purpose:        purpose,
		Target:         req.NeedID,
		WithdrawAmount: withdrawAmount,
		ProofBlobID:    req.ProofBlobID,
		Description:    description,
		Status:         request_pending_status,
		AIEvaluation:   aiEvaluation,
		CreatedAt:      curTime,
		UpdatedAt:      curTime,
	}

	return &res, c.pendingWithdrawProposalRepo.CreatePendingWithdrawProposal(res, ctx)
}

// SupportHealthInsuranceNeed implements business.IChildService.
func (c *childService) SupportHealthInsuranceNeed(id string, ctx context.Context) (response.UrlAPIResponse, error) {
	profile, err := c.profileRepo.GetProfile(ctx.Value("sub").(string), ctx)
	if err != nil {
		return response.UrlAPIResponse{}, err
	}

	if profile == nil || profile.IdentityCode == nil {
		return response.UrlAPIResponse{}, errors.New(noti.PROFILE_EMPTY_MESSAGE)
	}

	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if !util.IsValidSuiAddressStrict(id) {
		return response.UrlAPIResponse{}, genericErr
	}

	var client = c.clients[constant.SuiTestnet]
	need, err := on_chain.GetOnChainObject[entities.HealthInsuranceNeed](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  id,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return response.UrlAPIResponse{}, err
	}

	if need == nil {
		return response.UrlAPIResponse{}, genericErr
	}

	if slices.Contains(need.SupportedYears, need.Year) {
		return response.UrlAPIResponse{}, errors.New(noti.NEED_SUPPORTED_MESSAGE)
	}

	var paymentId string = util.GenerateId()
	var orderCode int = util.GenerateNumber()
	var callbackUrl string = os.Getenv(payment.PAYMENT_CALLBACK_URL) + paymentId
	var paymentDescription string = entities.HEALTH_INSRUANCE_PAYMENT_DESCRIPTION.GenerateSupportPaymentDescription()
	amount, _ := strconv.ParseInt(need.Value, 10, 64)
	data, err := payos.CreatePaymentLink(payos.CheckoutRequestType{
		OrderCode:   int64(orderCode),
		Amount:      int(amount),
		Description: paymentDescription,
		ReturnUrl:   callbackUrl,
		CancelUrl:   callbackUrl,
	})

	if err != nil {
		c.errLogger.Println("Err: ", err.Error())
		return response.UrlAPIResponse{}, errors.New(noti.INTERNALL_ERR_MSG)
	}

	var donationId string = util.GenerateId()
	var curTime time.Time = time.Now()
	if err := c.donationRepo.CreateDonation(entities.OffChainDonation{
		ID:        donationId,
		Purpose:   string(entities.HEALTH_INSURANCE_NEED_PURPOSE),
		Target:    id,
		CreatedAt: curTime,
	}, ctx); err != nil {
		return response.UrlAPIResponse{}, err
	}

	var description string = fmt.Sprintf("Support Health Insurance Need %s for child %s", need.Year, util.FormatAddress(need.ChildID))
	return response.UrlAPIResponse{
			Url: data.CheckoutUrl,
		}, c.paymentRepo.CreatePayment(entities.Payment{
			ID:            paymentId,
			Actor:         ctx.Value("address").(string),
			ProfileID:     profile.ID,
			DonationID:    &donationId,
			IsDonateTx:    true,
			TransactionId: fmt.Sprint(orderCode),
			Amount:        amount,
			Currency:      shared.VIETNAMDONG_CURRENCY,
			Status:        payment_pending_status,
			Method:        shared.PAYMENT_PAYOS_METHOD,
			Message:       description,
			ExpiredAt:     time.Unix(int64(*data.ExpiredAt), 0),
			CreatedAt:     curTime,
			UpdatedAt:     curTime,
		}, ctx)
}

// CreateMealNeedWithdrawProposalV2 implements business.IChildService.
func (c *childService) CreateMealNeedWithdrawProposalV2(req request.CreateNormalNeedWithdrawProposalRequest, ctx context.Context) (*entities.PendingWithdrawProposal, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if !util.IsValidSuiAddressStrict(req.NeedID) {
		return nil, genericErr
	}

	var client = c.clients[constant.SuiTestnet]
	need, err := on_chain.GetOnChainObject[entities.MealNeed](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  req.NeedID,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return nil, err
	}

	if need == nil {
		return nil, genericErr
	}

	totalSupportedMonths, _ := strconv.Atoi(need.TotalSupportedMonths)
	var expectedDuration int = totalSupportedMonths - len(need.WithdrawsForNeed)
	// Already withdraw all
	if expectedDuration == 0 {
		return nil, errors.New(noti.NEED_WITHDRAWN_MESSAGE)
	}

	var previousDuration int = 0
	var expectedDate time.Time
	var curTime time.Time = time.Now()
	for i := len(need.Durations) - 1; i >= 0; i-- {
		var duration = need.Durations[0]
		var startPeriod time.Time = util.ToStartOfDate(util.RawDateToTime(duration.Fields.StartPeriod))
		var endPeriod time.Time = util.ToEndOfDate(util.RawDateToTime(duration.Fields.EndPeriod))
		var startMonth int = int(startPeriod.Month())
		var endMonth int = int(endPeriod.Month())
		if endMonth == 1 { // To next year
			endMonth = 13
		}

		var currentDuration int = endMonth - startMonth
		var totalDuration int = currentDuration + previousDuration
		var months int = totalDuration - expectedDuration
		if months >= 0 {
			var startDate = startPeriod.AddDate(0, months, 0)
			expectedDate = startDate.AddDate(0, 0, -3)
			break
		}

		previousDuration = totalDuration
	}

	var expectedLimDate time.Time = expectedDate.AddDate(0, 0, 7)

	// Still not date to withdraw
	if curTime.Before(expectedDate) || curTime.After(expectedLimDate) {
		return nil, errors.New(noti.NOT_WITHDRAW_EXPECTED_DATE_MESSAGE)
	}

	child, err := on_chain.GetOnChainObject[entities.Child](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  need.ChildID,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return nil, err
	}

	if child == nil {
		return nil, genericErr
	}

	var staffModule = on_chain.InitializeModuleStaff()
	var sender string = ctx.Value("address").(string)
	staffNfts, err := on_chain.GetOnChainOwnedObjects[entities.StaffNft](on_chain.GetOnChainOwnedObjectsRequest{
		Client:       client,
		OwnerAddress: sender,
		StructType:   fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), staffModule.GetModule(), staffModule.GetStaffNftObjectStruct()),
		ErrLogger:    c.errLogger,
	}, ctx)
	if err != nil {
		return nil, err
	}

	var genericRightErr error = errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	if staffNfts == nil || len(staffNfts) == 0 {
		return nil, genericRightErr
	}

	var isLeaderOfRegion bool = false
	for _, nft := range staffNfts {
		if nft.Role == local_leader_role && nft.Region == child.Region {
			isLeaderOfRegion = true
			break
		}
	}

	if !isLeaderOfRegion {
		return nil, genericRightErr
	}

	leaderNoti, err := c.leaderNotiRepo.GetNotiByNeed(req.NeedID, ctx)
	if err != nil {
		return nil, err
	}

	//var rawExpectedDate string = util.TimeToRawDate(expectedDate)
	var index int
	for i := len(leaderNoti.ExpectedWithdrawPeriods) - 1; i >= 0; i-- {
		// var rawDate string = leaderNoti.ExpectedWithdrawPeriods[i]
		// if rawDate == rawExpectedDate {
		// 	index = i
		// 	break
		// }

		var rawDate string = leaderNoti.ExpectedWithdrawPeriods[i]
		if !expectedDate.Before(util.ToStartOfDate(util.RawDateToTime(rawDate))) {
			index = i
			break
		}
	}

	var description string = leaderNoti.Contents[index]
	var purpose string = string(entities.MEAL_NEED_PURPOSE)
	withdrawAmount, _ := strconv.ParseInt(need.Value, 10, 64)
	isProposed, err := c.pendingWithdrawProposalRepo.IsPendingWithdrawProposalProposedWithSpecificInfo(purpose, req.NeedID, description, withdrawAmount, ctx)
	if err != nil {
		return nil, err
	}

	if isProposed {
		return nil, errors.New(noti.STILL_PENDING_REQUEST_MESSAGE)
	}

	pool, err := on_chain.GetOnChainObject[entities.MainPool](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.POOL_ID),
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return nil, err
	}

	localPools, err := on_chain.GetOnChainObjects[entities.LocalPool](on_chain.GetOnChainObjectsRequest{
		Client:    client,
		ObjectIds: pool.LocalPools,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return nil, err
	}

	var localPoolId, localPoolName string
	for _, localPool := range localPools {
		if localPool.Region == child.Region {
			localPoolId = localPool.ID.ID
			localPoolName = localPool.Region
			break
		}
	}

	var aiEvaluation string
	if req.ProofBlobID != nil {
		proofBytes, _ := c.walrusProvider.FetchBytesImage(*req.ProofBlobID)
		if proofBytes != nil {
			aiEvaluation = c.aiProvider.ValidateWithdrawProposal(ai.ValidateWithdrawProposal{
				Purpose:         purpose,
				WithdrawAmount:  withdrawAmount,
				Description:     description,
				ProofBytesImage: proofBytes,
			}, ctx)
		}
	}

	// todo: AI validation
	var res = entities.PendingWithdrawProposal{
		ID:             util.GenerateId(),
		ProfileID:      ctx.Value("sub").(string),
		Creator:        sender,
		PoolID:         localPoolId,
		PoolName:       localPoolName,
		Purpose:        purpose,
		Target:         req.NeedID,
		WithdrawAmount: withdrawAmount,
		ProofBlobID:    req.ProofBlobID,
		Description:    description,
		Status:         request_pending_status,
		AIEvaluation:   aiEvaluation,
		CreatedAt:      curTime,
		UpdatedAt:      curTime,
	}

	return &res, c.pendingWithdrawProposalRepo.CreatePendingWithdrawProposal(res, ctx)
}

// CreateSpecialNeedWithdrawProposalV2 implements business.IChildService.
func (c *childService) CreateSpecialNeedWithdrawProposalV2(req request.CreateSpecialNeedWithdrawProposalRequest, ctx context.Context) (*entities.PendingWithdrawProposal, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if !util.IsValidSuiAddressStrict(req.CampaignID) {
		return nil, genericErr
	}

	var client = c.clients[constant.SuiTestnet]
	campaign, err := on_chain.GetOnChainObject[entities.SpecialNeedCampaign](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  req.CampaignID,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return nil, err
	}

	if campaign == nil {
		return nil, genericErr
	}

	var sender string = ctx.Value("address").(string)
	if campaign.Creator != sender {
		return nil, errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	}

	totalWithdrawAmount, _ := strconv.ParseInt(campaign.WithdrawAmount, 10, 64)
	totalDonation, _ := strconv.ParseInt(campaign.TotalDonated, 10, 64)
	if req.Amount > totalDonation-totalWithdrawAmount {
		return nil, errors.New(noti.CURRENT_BUDGET_NOT_ENOUGH_MESSAGE)
	}

	pool, err := on_chain.GetOnChainObject[entities.MainPool](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.POOL_ID),
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return nil, err
	}

	localPools, err := on_chain.GetOnChainObjects[entities.LocalPool](on_chain.GetOnChainObjectsRequest{
		Client:    client,
		ObjectIds: pool.LocalPools,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return nil, err
	}

	child, err := on_chain.GetOnChainObject[entities.Child](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  campaign.ChildID,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return nil, err
	}

	var localPoolId, localPoolName string
	for _, localPool := range localPools {
		if localPool.Region == child.Region {
			localPoolId = localPool.ID.ID
			localPoolName = localPool.Region
			break
		}
	}

	var description string = strings.TrimSpace(req.Description)
	var purpose string = string(entities.SPECIAL_NEED_PURPOSE)
	var aiEvaluation string
	if req.ProofBlobID != nil {
		proofBytes, _ := c.walrusProvider.FetchBytesImage(*req.ProofBlobID)
		if proofBytes != nil {
			aiEvaluation = c.aiProvider.ValidateWithdrawProposal(ai.ValidateWithdrawProposal{
				Purpose:         purpose,
				WithdrawAmount:  req.Amount,
				Description:     description,
				ProofBytesImage: proofBytes,
			}, ctx)
		}
	}

	var curTime time.Time = time.Now()
	var res = entities.PendingWithdrawProposal{
		ID:             util.GenerateId(),
		ProfileID:      ctx.Value("sub").(string),
		Creator:        sender,
		PoolID:         localPoolId,
		PoolName:       localPoolName,
		Purpose:        purpose,
		Target:         req.CampaignID,
		WithdrawAmount: req.Amount,
		ProofBlobID:    req.ProofBlobID,
		Description:    description,
		Status:         request_pending_status,
		AIEvaluation:   aiEvaluation,
		CreatedAt:      curTime,
		UpdatedAt:      curTime,
	}

	return &res, c.pendingWithdrawProposalRepo.CreatePendingWithdrawProposal(res, ctx)
}

// ConfirmSpecialNeedProposal implements business.IChildService.
func (c *childService) ConfirmSpecialNeedProposal(id string, ctx context.Context) (response.BuildTransactionResponse, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if !util.IsValidSuiAddressStrict(id) {
		return response.BuildTransactionResponse{}, genericErr
	}

	var client = c.clients[constant.SuiTestnet]
	proposal, err := on_chain.GetOnChainObject[entities.SpecialNeedProposal](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  id,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	if proposal == nil {
		return response.BuildTransactionResponse{}, genericErr
	}

	var sender string = ctx.Value("address").(string)
	if proposal.Creator != sender {
		return response.BuildTransactionResponse{}, errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	}

	if proposal.IsConfirm {
		return response.BuildTransactionResponse{}, errors.New(noti.SPECIAL_NEED_PROPOSAL_CONFIRMED_MESSAGE)
	}

	closedAt, _ := strconv.ParseInt(proposal.ClosedAt, 10, 64)
	if util.MilliSecToTime(closedAt).After(time.Now()) {
		return response.BuildTransactionResponse{}, errors.New(noti.STILL_PENDING_REQUEST_MESSAGE)
	}

	dao, err := on_chain.GetOnChainObject[entities.DaoStruct](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.SPECIAL_NEED_DAO_ID),
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	if !isProposalRateAvailableToConfirm(*dao, len(proposal.Approvers), len(proposal.Refusers), proposal.ApproveWeight, proposal.RefuseWeight) {
		return response.BuildTransactionResponse{}, errors.New(noti.PROPOSAL_FAIL_CONDITION_TO_CONFIRM_MESSAGE)
	}

	var childModule = on_chain.InitializeModuleChild()
	txBytes, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
		Client:    client,
		Sender:    sender,
		Module:    childModule.GetModule(),
		Function:  childModule.GetFunctionConfirmChildSpecialNeedProposal(),
		ErrLogger: c.errLogger,
		Arguments: childModule.ToConfirmChildSpecialNeedProposalArguments(on_chain.ConfirmChildSpecialNeedProposalArguments{
			ProposalID: id,
			ChildID:    proposal.ChildID,
		}),
	}, ctx)

	return response.BuildTransactionResponse{
		TxBytes: txBytes,
	}, err
}

// CreateSpecialNeedProposal implements business.IChildService.
func (c *childService) CreateSpecialNeedProposal(req request.CreateSpecialNeedProposalRequest, ctx context.Context) (response.BuildTransactionResponse, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if !util.IsValidSuiAddressStrict(req.ChildID) {
		return response.BuildTransactionResponse{}, genericErr
	}

	var client = c.clients[constant.SuiTestnet]
	child, err := on_chain.GetOnChainObject[entities.Child](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  req.ChildID,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	if child == nil {
		return response.BuildTransactionResponse{}, genericErr
	}

	pool, err := on_chain.GetOnChainObject[entities.MainPool](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.POOL_ID),
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	localPools, err := on_chain.GetOnChainObjects[entities.LocalPool](on_chain.GetOnChainObjectsRequest{
		Client:    client,
		ObjectIds: pool.LocalPools,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	var sender string = ctx.Value("address").(string)
	var localPoolId string
	var isLeaderOfRegion bool = false
	for _, localPool := range localPools {
		if localPool.Region == child.Region {
			localPoolId = localPool.ID.ID
			if slices.Contains(localPool.Mods, sender) {
				isLeaderOfRegion = true
			}
			break
		}
	}

	if !isLeaderOfRegion {
		return response.BuildTransactionResponse{}, errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	}

	// Not leader of region
	// if !isLeaderOfRegion {
	// 	var manageModule = on_chain.InitializeModuleManage()
	// 	nfts, err := on_chain.GetOnChainOwnedObjects[entities.AdminNft](on_chain.GetOnChainOwnedObjectsRequest{
	// 		Client:       client,
	// 		OwnerAddress: sender,
	// 		StructType:   fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), manageModule.GetModule(), manageModule.GetAdminNftStruct()),
	// 		ErrLogger:    c.errLogger,
	// 	}, ctx)
	// 	if err != nil {
	// 		return response.BuildTransactionResponse{}, err
	// 	}

	// 	if nfts == nil || len(nfts) == 0 {
	// 		return response.BuildTransactionResponse{}, errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	// 	}
	// }

	bankProfile, err := c.bankRepo.GetBankProfileByOwner(sender, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	if bankProfile == nil {
		return response.BuildTransactionResponse{}, errors.New(noti.LEADER_NOT_UPLOAD_BANK_PROFILE_MESSAGE)
	}

	var childModule = on_chain.InitializeModuleChild()
	txBytes, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
		Client:    client,
		Sender:    sender,
		Module:    childModule.GetModule(),
		Function:  childModule.GetFunctionCreateChildSpecialNeedProposal(),
		ErrLogger: c.errLogger,
		Arguments: childModule.ToCreateChildSpecialNeedProposalArguments(on_chain.CreateChildSpecialNeedProposalArguments{
			ChildID:     req.ChildID,
			LocalPool:   localPoolId,
			Target:      req.Target,
			Description: req.Description,
			ClosedAt:    util.ToMilliseconds(util.GetRequestDuration()),
		}),
	}, ctx)

	return response.BuildTransactionResponse{
		TxBytes: txBytes,
	}, err
}

// CreateSpecialNeedProposalV2 implements business.IChildService.
func (c *childService) CreateSpecialNeedProposalV2(req request.CreateSpecialNeedProposalRequest, ctx context.Context) (*entities.PendingChildSpecialNeedProposal, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if !util.IsValidSuiAddressStrict(req.ChildID) {
		return nil, genericErr
	}

	var client = c.clients[constant.SuiTestnet]
	child, err := on_chain.GetOnChainObject[entities.Child](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  req.ChildID,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return nil, err
	}

	if child == nil {
		return nil, genericErr
	}

	var sender string = ctx.Value("address").(string)
	var staffModule = on_chain.InitializeModuleStaff()
	staffNfts, err := on_chain.GetOnChainOwnedObjects[entities.StaffNft](on_chain.GetOnChainOwnedObjectsRequest{
		Client:       client,
		OwnerAddress: sender,
		StructType:   fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), staffModule.GetModule(), staffModule.GetStaffNftObjectStruct()),
		ErrLogger:    c.errLogger,
	}, ctx)
	if err != nil {
		return nil, err
	}

	if staffNfts == nil || len(staffNfts) == 0 {
		return nil, errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	}

	var isLeaderOfRegion bool = false
	for _, nft := range staffNfts {
		if nft.Region == child.Region && nft.Role == local_leader_role {
			isLeaderOfRegion = true
			break
		}
	}

	if !isLeaderOfRegion {
		return nil, errors.New(noti.LEADER_NOT_OF_REGION_MESSAGE)
	}

	bankProfile, err := c.bankRepo.GetBankProfileByOwner(sender, ctx)
	if err != nil {
		return nil, err
	}

	if bankProfile == nil {
		return nil, errors.New(noti.LEADER_NOT_UPLOAD_BANK_PROFILE_MESSAGE)
	}

	var description string = strings.TrimSpace(req.Description)
	var aiEvaluation string
	if req.ProofBlobID != nil {
		proofBytes, _ := c.walrusProvider.FetchBytesImage(*req.ProofBlobID)
		if proofBytes != nil {
			aiEvaluation = c.aiProvider.ValidateChildSpecialNeedProposal(ai.ValidateChildSpecialNeedProposal{
				CamapaignTarget: req.Target,
				Description:     description,
				ProofBytesImage: proofBytes,
			}, ctx)
		}
	}

	// todo: AI validation
	var curTime time.Time = time.Now()
	var proposal = entities.PendingChildSpecialNeedProposal{
		ID:             util.GenerateId(),
		ChildID:        req.ChildID,
		Region:         child.Region,
		ActorProfileID: ctx.Value("sub").(string),
		ActorAddress:   sender,
		Target:         req.Target,
		Description:    strings.TrimSpace(req.Description),
		ProofBlobID:    req.ProofBlobID,
		AIEvaluation:   aiEvaluation,
		CreatedAt:      curTime,
		UpdatedAt:      curTime,
	}

	return &proposal, c.pendingChildSpecialNeedProposalRepo.CreatePendingChildSpecialNeedProposal(proposal, ctx)
}

// CreateSpecialNeedWithdrawProposal implements business.IChildService.
func (c *childService) CreateSpecialNeedWithdrawProposal(req request.CreateSpecialNeedWithdrawProposalRequest, ctx context.Context) (response.BuildTransactionResponse, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if !util.IsValidSuiAddressStrict(req.CampaignID) {
		return response.BuildTransactionResponse{}, genericErr
	}

	var client = c.clients[constant.SuiTestnet]
	campaign, err := on_chain.GetOnChainObject[entities.SpecialNeedCampaign](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  req.CampaignID,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	if campaign == nil {
		return response.BuildTransactionResponse{}, genericErr
	}

	var sender string = ctx.Value("address").(string)
	if campaign.Creator != sender {
		return response.BuildTransactionResponse{}, errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	}

	totalWithdrawAmount, _ := strconv.ParseInt(campaign.WithdrawAmount, 10, 64)
	totalDonation, _ := strconv.ParseInt(campaign.TotalDonated, 10, 64)
	if req.Amount > totalDonation-totalWithdrawAmount {
		return response.BuildTransactionResponse{}, errors.New(noti.CURRENT_BUDGET_NOT_ENOUGH_MESSAGE)
	}

	pool, err := on_chain.GetOnChainObject[entities.MainPool](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.POOL_ID),
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	localPools, err := on_chain.GetOnChainObjects[entities.LocalPool](on_chain.GetOnChainObjectsRequest{
		Client:    client,
		ObjectIds: pool.LocalPools,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	child, err := on_chain.GetOnChainObject[entities.Child](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  campaign.ChildID,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	var localPoolId string
	for _, localPool := range localPools {
		if localPool.Region == child.Region {
			localPoolId = localPool.ID.ID
			break
		}
	}

	var childModule = on_chain.InitializeModuleChild()
	txBytes, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
		Client:    client,
		Sender:    sender,
		Module:    childModule.GetModule(),
		Function:  childModule.GetFunctionCreateChildBooksNeedWithdrawProposal(),
		ErrLogger: c.errLogger,
		Arguments: childModule.ToCreateChildSpecialNeedWithdrawProposalArguments(on_chain.CreateChildSpecialNeedWithdrawProposalArguments{
			CampaignID:     req.CampaignID,
			LocalPool:      localPoolId,
			ChildID:        campaign.ChildID,
			WithdrawAmount: req.Amount,
			Description:    req.Description,
			ClosedAt:       util.ToMilliseconds(util.GetRequestDuration()),
		}),
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	var proposalId string = util.GenerateId()

	return response.BuildTransactionResponse{
			TxBytes:    txBytes,
			ProposalId: proposalId,
		}, c.withdrawRepo.CreateOffChainWithdrawProposal(entities.OffChainWithdrawProposal{
			ID:        proposalId,
			Purpose:   string(entities.SPECIAL_NEED_PURPOSE),
			Target:    req.CampaignID,
			CreatedAt: time.Now(),
		}, ctx)
}

// EditSpecialNeedDao implements business.IChildService.
func (c *childService) EditSpecialNeedDao(req request.EditDaoRequest, ctx context.Context) (response.BuildTransactionResponse, error) {
	panic("unimplemented")
}

// ConfirmProvideMealForChild implements business.IChildService.
func (c *childService) ConfirmProvideMealForChild(id string, req request.ConfirmProvideMealForChildRequest, ctx context.Context) (response.BuildTransactionResponse, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if !util.IsValidSuiAddressStrict(id) {
		return response.BuildTransactionResponse{}, genericErr
	}

	var client = c.clients[constant.SuiTestnet]
	child, err := on_chain.GetOnChainObject[entities.Child](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  id,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	if child == nil {
		return response.BuildTransactionResponse{}, genericErr
	}

	var staffModule = on_chain.InitializeModuleStaff()
	var sender string = ctx.Value("address").(string)
	staffNfts, err := on_chain.GetOnChainOwnedObjects[entities.StaffNft](on_chain.GetOnChainOwnedObjectsRequest{
		Client:       client,
		OwnerAddress: sender,
		StructType:   fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), staffModule.GetModule(), staffModule.GetStaffNftObjectStruct()),
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	var genericRightErr error = errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	if staffNfts == nil || len(staffNfts) == 0 {
		return response.BuildTransactionResponse{}, genericRightErr
	}

	var isStaffOfRegion bool = false
	for _, nft := range staffNfts {
		if nft.Region == child.Region {
			isStaffOfRegion = true
			break
		}
	}

	if !isStaffOfRegion {
		return response.BuildTransactionResponse{}, genericRightErr
	}

	need, err := on_chain.GetOnChainObject[entities.MealNeed](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  child.MealNeed,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	var curTime time.Time = time.Now()
	var rawProvideDate string = util.TimeToRawDate(curTime)
	if slices.Contains(need.ProvideMealDates, rawProvideDate) {
		return response.BuildTransactionResponse{}, errors.New(noti.CHILD_PROVIDED_MEAL_MESSAGE)
	}

	var isProvideDateInDuration bool = false
	for i := len(need.Durations) - 1; i >= 0; i-- {
		var duration = need.Durations[i]
		var startPeriod = util.RawDateToTime(duration.Fields.StartPeriod)
		var endPeriod = util.RawDateToTime(duration.Fields.EndPeriod)
		if !startPeriod.After(curTime) && !curTime.After(endPeriod) {
			isProvideDateInDuration = true
			break
		}
	}

	if !isProvideDateInDuration {
		return response.BuildTransactionResponse{}, errors.New(noti.CHILD_NOT_IN_MEAL_SUPPORT)
	}

	// todo: implement AI to validate image
	var childModule = on_chain.InitializeModuleChild()
	txBytes, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
		Client:    client,
		Sender:    sender,
		Module:    childModule.GetModule(),
		Function:  childModule.GetFunctionConfirmProvideMealForChild(),
		ErrLogger: c.errLogger,
		Arguments: childModule.ToConfirmProvideMealForChildArguments(on_chain.ConfirmProvideMealForChildArguments{
			ChildID:     id,
			NeedID:      need.ID.ID,
			StaffNft:    staffNfts[0].ID.ID,
			ImageBlobID: req.ImageBlobID,
			ProvideDate: rawProvideDate,
		}),
	}, ctx)

	return response.BuildTransactionResponse{
		TxBytes: txBytes,
	}, err
}

// SupportBooksNeed implements business.IChildService.
func (c *childService) SupportBooksNeed(id string, ctx context.Context) (response.UrlAPIResponse, error) {
	profile, err := c.profileRepo.GetProfile(ctx.Value("sub").(string), ctx)
	if err != nil {
		return response.UrlAPIResponse{}, err
	}

	if profile == nil || profile.IdentityCode == nil {
		return response.UrlAPIResponse{}, errors.New(noti.PROFILE_EMPTY_MESSAGE)
	}

	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if !util.IsValidSuiAddressStrict(id) {
		return response.UrlAPIResponse{}, genericErr
	}

	var client = c.clients[constant.SuiTestnet]
	need, err := on_chain.GetOnChainObject[entities.BooksNeed](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  id,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return response.UrlAPIResponse{}, err
	}

	if need == nil {
		return response.UrlAPIResponse{}, genericErr
	}

	if slices.Contains(need.SupportedYears, need.Year) {
		return response.UrlAPIResponse{}, errors.New(noti.NEED_SUPPORTED_MESSAGE)
	}

	var paymentId string = util.GenerateId()
	var orderCode int = util.GenerateNumber()
	var callbackUrl string = os.Getenv(payment.PAYMENT_CALLBACK_URL) + paymentId
	amount, _ := strconv.ParseInt(need.Value, 10, 64)
	var paymentDescription string = entities.BOOKS_NEED_PAYMENT_DESCRIPTION.GenerateSupportPaymentDescription()
	data, err := payos.CreatePaymentLink(payos.CheckoutRequestType{
		OrderCode:   int64(orderCode),
		Amount:      int(amount),
		Description: paymentDescription,
		ReturnUrl:   callbackUrl,
		CancelUrl:   callbackUrl,
	})

	if err != nil {
		c.errLogger.Println("Err: ", err.Error())
		return response.UrlAPIResponse{}, errors.New(noti.INTERNALL_ERR_MSG)
	}

	// leaderNoti, err := c.leaderNotiRepo.GetNotiByMealNeed(id, ctx)
	// if err != nil {
	// 	return response.UrlAPIResponse{}, err
	// }

	// var curTime time.Time = time.Now()
	// if leaderNoti != nil {
	// 	leaderNoti.ExpectedWithdrawPeriods = append(leaderNoti.ExpectedWithdrawPeriods, "")
	// 	if err := c.leaderNotiRepo.UpdateNoti(*leaderNoti, ctx); err != nil {
	// 		return response.UrlAPIResponse{}, err
	// 	}
	// } else {
	// 	child, err := on_chain.GetOnChainObject[entities.Child](on_chain.GetOnChainObjectRequest{
	// 		Client:    client,
	// 		ObjectId:  need.ChildID,
	// 		ErrLogger: c.errLogger,
	// 	}, ctx)
	// 	if err != nil {
	// 		return response.UrlAPIResponse{}, err
	// 	}

	// 	pool, err := on_chain.GetOnChainObject[entities.MainPool](on_chain.GetOnChainObjectRequest{
	// 		Client:    client,
	// 		ObjectId:  os.Getenv(env.POOL_ID),
	// 		ErrLogger: c.errLogger,
	// 	}, ctx)
	// 	if err != nil {
	// 		return response.UrlAPIResponse{}, err
	// 	}

	// 	withdrawDates, err := on_chain.GetOnChainObject[entities.BooksNeedWithdrawDates](on_chain.GetOnChainObjectRequest{
	// 		Client:    client,
	// 		ObjectId:  os.Getenv(env.BOOKS_NEED_WITHDRAW_DATES_ID),
	// 		ErrLogger: c.errLogger,
	// 	}, ctx)
	// 	if err != nil {
	// 		return response.UrlAPIResponse{}, err
	// 	}

	// 	var withdrawDate string
	// 	if need.Semster == "1" {
	// 		withdrawDate = withdrawDates.FirstSemesterDate
	// 	} else {
	// 		withdrawDate = withdrawDates.SecondSemesterDate
	// 	}

	// 	localPools, err := on_chain.GetOnChainObjects[entities.LocalPool](on_chain.GetOnChainObjectsRequest{
	// 		Client:    client,
	// 		ObjectIds: pool.LocalPools,
	// 		ErrLogger: c.errLogger,
	// 	}, ctx)
	// 	if err != nil {
	// 		return response.UrlAPIResponse{}, err
	// 	}

	// 	var leaders []string
	// 	for _, localPool := range localPools {
	// 		if localPool.Region == child.Region {
	// 			leaders = localPool.Mods
	// 			break
	// 		}
	// 	}

	// 	if err := c.leaderNotiRepo.CreateNoti(entities.LeaderNoti{
	// 		ID:                      util.GenerateId(),
	// 		NeedID:                  id,
	// 		NeedType:                string(entities.BOOKS_NEED_PURPOSE),
	// 		ChildID:                 need.ChildID,
	// 		Region:                  child.Region,
	// 		AssignedLeaders:         leaders,
	// 		ExpectedWithdrawPeriods: []string{withdrawDate + "/" + need.Year},
	// 		Content:                 fmt.Sprintf("Withdraw books need semester %s for child %s", need.Semster, util.FormatAddress(child.ID.ID)),
	// 		CreatedAt: curTime,
	// 		UpdatedAt: curTime,
	// 	}, ctx); err != nil {
	// 		return response.UrlAPIResponse{}, err
	// 	}
	// }

	var donationId string = util.GenerateId()
	var curTime time.Time = time.Now()
	if err := c.donationRepo.CreateDonation(entities.OffChainDonation{
		ID:        donationId,
		Purpose:   string(entities.BOOKS_NEED_PURPOSE),
		Target:    id,
		CreatedAt: curTime,
	}, ctx); err != nil {
		return response.UrlAPIResponse{}, err
	}

	var description string = fmt.Sprintf("Support Books Need Semester %s - %s", need.Semster, need.Year)
	return response.UrlAPIResponse{
			Url: data.CheckoutUrl,
		}, c.paymentRepo.CreatePayment(entities.Payment{
			ID:            paymentId,
			Actor:         ctx.Value("address").(string),
			ProfileID:     profile.ID,
			DonationID:    &donationId,
			IsDonateTx:    true,
			TransactionId: fmt.Sprint(orderCode),
			Amount:        amount,
			Currency:      shared.VIETNAMDONG_CURRENCY,
			Status:        payment_pending_status,
			Method:        shared.PAYMENT_PAYOS_METHOD,
			Message:       description,
			ExpiredAt:     time.Unix(int64(*data.ExpiredAt), 0),
			CreatedAt:     curTime,
			UpdatedAt:     curTime,
		}, ctx)
}

// SupportMealNeed implements business.IChildService.
func (c *childService) SupportMealNeed(id string, req request.SupportMealNeadRequest, ctx context.Context) (response.UrlAPIResponse, error) {
	profile, err := c.profileRepo.GetProfile(ctx.Value("sub").(string), ctx)
	if err != nil {
		return response.UrlAPIResponse{}, err
	}

	if profile == nil || profile.IdentityCode == nil {
		return response.UrlAPIResponse{}, errors.New(noti.PROFILE_EMPTY_MESSAGE)
	}

	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if !util.IsValidSuiAddressStrict(id) {
		return response.UrlAPIResponse{}, genericErr
	}

	var client = c.clients[constant.SuiTestnet]
	need, err := on_chain.GetOnChainObject[entities.MealNeed](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  id,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		c.errLogger.Println("Fail at get object")
		return response.UrlAPIResponse{}, err
	}

	if need == nil {
		return response.UrlAPIResponse{}, genericErr
	}

	var curTime time.Time = time.Now()
	var lastDuration = need.Durations[len(need.Durations)-1]
	var endPeriod time.Time = util.RawDateToTime(lastDuration.Fields.EndPeriod)
	var rawExpectedStart, rawExpectedEnd string
	var nextStartPeriod time.Time
	if curTime.Before(endPeriod) { // Donate time: 1/1/2026 | Last supported: 15/7/2026
		nextStartPeriod = endPeriod.AddDate(0, 0, 2)
	} else {
		nextStartPeriod = curTime.AddDate(0, 0, 2)
	}

	var nextYear int = curTime.Year() + 1
	var rawMaxSupportedEndPeriod string = fmt.Sprintf("15/01/%d", nextYear)
	var nextEndPeriod time.Time = nextStartPeriod.AddDate(0, req.Months, 0)
	if nextEndPeriod.After(util.RawDateToTime(rawMaxSupportedEndPeriod)) {
		// Support 6 months -> 16/1/2027 -> Deny
		c.errLogger.Println("Raw last end period:", lastDuration.Fields.EndPeriod)
		c.errLogger.Println("Last end period:", endPeriod)
		c.errLogger.Println("Next year:", nextYear)
		c.errLogger.Println("Next start period:", nextStartPeriod)
		c.errLogger.Println("Next end period:", nextEndPeriod)
		return response.UrlAPIResponse{}, errors.New(noti.MEAL_NEED_SUPPORT_DURATION_OUT_RANGE_MESSAGE)
	}

	rawExpectedStart = util.TimeToRawDate(nextStartPeriod)
	rawExpectedEnd = util.TimeToRawDate(nextEndPeriod)

	var paymentId string = util.GenerateId()
	var orderCode int = util.GenerateNumber()
	var callbackUrl string = os.Getenv(payment.PAYMENT_CALLBACK_URL) + paymentId
	value, _ := strconv.ParseInt(need.Value, 10, 64)
	var amount int64 = value * int64(req.Months)
	var paymentDescription string = entities.MEAL_NEED_PAYMENT_DESCRIPTION.GenerateSupportPaymentDescription()
	data, err := payos.CreatePaymentLink(payos.CheckoutRequestType{
		OrderCode:   int64(orderCode),
		Amount:      int(amount),
		Description: paymentDescription,
		ReturnUrl:   callbackUrl,
		CancelUrl:   callbackUrl,
	})
	if err != nil {
		c.errLogger.Println("Err: ", err.Error())
		c.errLogger.Println("Fail at create payos")
		return response.UrlAPIResponse{}, errors.New(noti.INTERNALL_ERR_MSG)
	}

	// var expectedWithdrawDate time.Time = nextEndPeriod.AddDate(0, 0, -1)
	// leaderNoti, err := c.leaderNotiRepo.GetNotiByMealNeed(id, ctx)
	// if err != nil {
	// 	return response.UrlAPIResponse{}, err
	// }

	// child, err := on_chain.GetOnChainObject[entities.Child](on_chain.GetOnChainObjectRequest{
	// 	Client:    client,
	// 	ObjectId:  need.ChildID,
	// 	ErrLogger: c.errLogger,
	// }, ctx)
	// if err != nil {
	// 	return response.UrlAPIResponse{}, err
	// }

	// var expectedWithdrawDates []string
	// for i := 0; i < req.Months; i++ {
	// 	var withdrawDate time.Time = expectedWithdrawDate.AddDate(0, i, 0)
	// 	expectedWithdrawDates = append(expectedWithdrawDates, util.TimeToRawDate(withdrawDate))
	// }

	// if leaderNoti != nil {
	// 	leaderNoti.ExpectedWithdrawPeriods = append(leaderNoti.ExpectedWithdrawPeriods, expectedWithdrawDates...)
	// 	if err := c.leaderNotiRepo.UpdateNoti(*leaderNoti, ctx); err != nil {
	// 		return response.UrlAPIResponse{}, err
	// 	}
	// } else {
	// 	pool, err := on_chain.GetOnChainObject[entities.MainPool](on_chain.GetOnChainObjectRequest{
	// 		Client:    client,
	// 		ObjectId:  os.Getenv(env.POOL_ID),
	// 		ErrLogger: c.errLogger,
	// 	}, ctx)
	// 	if err != nil {
	// 		return response.UrlAPIResponse{}, err
	// 	}

	// 	localPools, err := on_chain.GetOnChainObjects[entities.LocalPool](on_chain.GetOnChainObjectsRequest{
	// 		Client:    client,
	// 		ObjectIds: pool.LocalPools,
	// 		ErrLogger: c.errLogger,
	// 	}, ctx)
	// 	if err != nil {
	// 		return response.UrlAPIResponse{}, err
	// 	}

	// 	var leaders []string
	// 	for _, localPool := range localPools {
	// 		if localPool.Region == child.Region {
	// 			leaders = localPool.Mods
	// 			break
	// 		}
	// 	}

	// 	if err := c.leaderNotiRepo.CreateNoti(entities.LeaderNoti{
	// 		ID:                      util.GenerateId(),
	// 		NeedID:                  id,
	// 		ChildID:                 need.ChildID,
	// 		Region:                  child.Region,
	// 		AssignedLeaders:         leaders,
	// 		ExpectedWithdrawPeriods: expectedWithdrawDates,
	// 		//Content:                 fmt.Sprintf("Withdraw meal need for child %s", util.FormatAddress(child.ID.ID)),
	// 		CreatedAt: curTime,
	// 		UpdatedAt: curTime,
	// 	}, ctx); err != nil {
	// 		return response.UrlAPIResponse{}, err
	// 	}
	// }

	// manageObj, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
	// 	Client:    client,
	// 	ObjectId:  os.Getenv(env.PACKAGE_ID),
	// 	ErrLogger: c.errLogger,
	// }, ctx)
	// if err != nil {
	// 	return response.UrlAPIResponse{}, err
	// }

	// volunteers, err := on_chain.GetOnChainObjects[entities.StaffNft](on_chain.GetOnChainObjectsRequest{
	// 	Client:    client,
	// 	ObjectIds: manageObj.VolunteerNfts,
	// 	ErrLogger: c.errLogger,
	// }, ctx)
	// if err != nil {
	// 	return response.UrlAPIResponse{}, err
	// }

	// var volunteerAddresses []string
	// for i, volunteer := range volunteers {
	// 	if volunteer.Region == child.Region {
	// 		volunteerAddresses = append(volunteerAddresses, manageObj.VolunteerIds[i])
	// 	}
	// }

	// if err := c.volunteerNotiRepo.CreateNoti(entities.VolunteerNoti{
	// 	ID:                 util.GenerateId(),
	// 	ChildID:            need.ChildID,
	// 	Region:             child.Region,
	// 	AssginedVolunteers: volunteerAddresses,
	// 	Content:            fmt.Sprintf("Provide meal for child %s from %s to %s", util.FormatAddress(child.ID.ID), rawExpectedStart, rawExpectedEnd),
	// 	StartPeriod:        nextStartPeriod,
	// 	EndPeriod:          nextEndPeriod,
	// }, ctx); err != nil {
	// 	return response.UrlAPIResponse{}, err
	// }

	var mealSupportDurationId string = util.GenerateId()
	if err := c.mealDurationRepo.CreateMealSupportDuration(entities.OffChainMealSupportDuration{
		ID:          mealSupportDurationId,
		StartPeriod: rawExpectedStart,
		EndPeriod:   rawExpectedEnd,
	}, ctx); err != nil {
		return response.UrlAPIResponse{}, err
	}

	var donationId string = util.GenerateId()
	if err := c.donationRepo.CreateDonation(entities.OffChainDonation{
		ID:             donationId,
		Purpose:        string(entities.MEAL_NEED_PURPOSE),
		Target:         id,
		MealDurationID: &mealSupportDurationId,
		CreatedAt:      curTime,
	}, ctx); err != nil {
		return response.UrlAPIResponse{}, err
	}

	var description string = fmt.Sprintf("Support Meal Need %s - %s for child", rawExpectedStart, rawExpectedEnd)
	return response.UrlAPIResponse{
			Url: data.CheckoutUrl,
		}, c.paymentRepo.CreatePayment(entities.Payment{
			ID:            paymentId,
			Actor:         ctx.Value("address").(string),
			ProfileID:     profile.ID,
			DonationID:    &donationId,
			IsDonateTx:    true,
			TransactionId: fmt.Sprint(orderCode),
			Amount:        amount,
			Currency:      shared.VIETNAMDONG_CURRENCY,
			Status:        payment_pending_status,
			Method:        shared.PAYMENT_PAYOS_METHOD,
			Message:       description,
			ExpiredAt:     time.Unix(int64(*data.ExpiredAt), 0),
			CreatedAt:     curTime,
			UpdatedAt:     curTime,
		}, ctx)
}

// SupportSpecialNeed implements business.IChildService.
func (c *childService) SupportSpecialNeed(id string, req request.SupportSpecialNeedRequest, ctx context.Context) (response.UrlAPIResponse, error) {
	profile, err := c.profileRepo.GetProfile(ctx.Value("sub").(string), ctx)
	if err != nil {
		return response.UrlAPIResponse{}, err
	}

	if profile == nil || profile.IdentityCode == nil {
		return response.UrlAPIResponse{}, errors.New(noti.PROFILE_EMPTY_MESSAGE)
	}

	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if !util.IsValidSuiAddressStrict(id) {
		return response.UrlAPIResponse{}, genericErr
	}

	var client = c.clients[constant.SuiTestnet]
	campaign, err := on_chain.GetOnChainObject[entities.SpecialNeedCampaign](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  id,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return response.UrlAPIResponse{}, err
	}

	if campaign == nil {
		return response.UrlAPIResponse{}, genericErr
	}

	target, _ := strconv.ParseInt(campaign.Target, 10, 64)
	totalDonations, _ := strconv.ParseInt(campaign.TotalDonated, 10, 64)
	if req.Amount > target-totalDonations {
		return response.UrlAPIResponse{}, errors.New(noti.SUPPORT_SURPASS_CAMPAIGN_TARGET_MESSAGE)
	}

	var paymentId string = util.GenerateId()
	var orderCode int = util.GenerateNumber()
	var callbackUrl string = os.Getenv(payment.PAYMENT_CALLBACK_URL) + paymentId
	var paymentDescription string = entities.SPECIAL_NEED_CAMPAIGN_PAYMENT_DESCRIPTION.GenerateSupportPaymentDescription()
	data, err := payos.CreatePaymentLink(payos.CheckoutRequestType{
		OrderCode:   int64(orderCode),
		Amount:      int(req.Amount),
		Description: paymentDescription,
		ReturnUrl:   callbackUrl,
		CancelUrl:   callbackUrl,
	})
	if err != nil {
		c.errLogger.Println("Err: ", err.Error())
		return response.UrlAPIResponse{}, errors.New(noti.INTERNALL_ERR_MSG)
	}

	var donationId string = util.GenerateId()
	var curTime time.Time = time.Now()
	if err := c.donationRepo.CreateDonation(entities.OffChainDonation{
		ID:        donationId,
		Purpose:   string(entities.SPECIAL_NEED_PURPOSE),
		Target:    id,
		CreatedAt: curTime,
	}, ctx); err != nil {
		return response.UrlAPIResponse{}, err
	}

	return response.UrlAPIResponse{
			Url: data.CheckoutUrl,
		}, c.paymentRepo.CreatePayment(entities.Payment{
			ID:            paymentId,
			Actor:         ctx.Value("address").(string),
			ProfileID:     profile.ID,
			DonationID:    &donationId,
			IsDonateTx:    true,
			TransactionId: fmt.Sprint(orderCode),
			Amount:        req.Amount,
			Currency:      shared.VIETNAMDONG_CURRENCY,
			Status:        payment_pending_status,
			Method:        shared.PAYMENT_PAYOS_METHOD,
			Message:       strings.TrimSpace(req.Description),
			ExpiredAt:     time.Unix(int64(*data.ExpiredAt), 0),
			CreatedAt:     curTime,
			UpdatedAt:     curTime,
		}, ctx)
}

// VoteSpecialNeedProposal implements business.IChildService.
func (c *childService) VoteSpecialNeedProposal(id string, req request.VoteRequest, ctx context.Context) (response.BuildTransactionResponse, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if !util.IsValidSuiAddressStrict(id) {
		return response.BuildTransactionResponse{}, genericErr
	}

	var client = c.clients[constant.SuiTestnet]
	proposal, err := on_chain.GetOnChainObject[entities.SpecialNeedProposal](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  id,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	if proposal == nil {
		return response.BuildTransactionResponse{}, genericErr
	}

	var sender string = ctx.Value("address").(string)
	if proposal.Creator == sender {
		return response.BuildTransactionResponse{}, errors.New(noti.OWNER_VOTE_WARN_MSG)
	}

	closedAt, _ := strconv.ParseInt(proposal.ClosedAt, 10, 64)
	if time.Now().After(util.MilliSecToTime(closedAt)) {
		return response.BuildTransactionResponse{}, errors.New(noti.REQUEST_CLOSED_MESSAGE)
	}

	if slices.Contains(proposal.Approvers, sender) || slices.Contains(proposal.Refusers, sender) {
		return response.BuildTransactionResponse{}, errors.New(noti.ALREADY_VOTE_MESSAGE)
	}

	var donorModule = on_chain.InitializeModuleDonor()
	nfts, _ := on_chain.GetOnChainOwnedObjects[entities.Donor](on_chain.GetOnChainOwnedObjectsRequest{
		Client:       client,
		OwnerAddress: sender,
		StructType:   fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), donorModule.GetModule(), donorModule.GetDonorNftStruct()),
		ErrLogger:    c.errLogger,
	}, ctx)
	if nfts == nil || len(nfts) == 0 {
		return response.BuildTransactionResponse{}, errors.New(noti.HAVE_TO_DONATE_TO_VOTE)
	}

	var refuseReason string = strings.TrimSpace(req.RefuseReason)
	if refuseReason == "" {
		refuseReason = "Refuse"
	}

	var needModule = on_chain.InitializeModuleNeed()
	txBytes, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
		Client:    client,
		Sender:    sender,
		Module:    needModule.GetModule(),
		Function:  needModule.GetFunctionVoteSpecialNeedProposal(),
		ErrLogger: c.errLogger,
		Arguments: needModule.ToVoteSpecialNeedProposalArguments(on_chain.VoteSpecialNeedProposalArguments{
			ProposalID:   id,
			DonorNft:     nfts[0].ID.ID,
			IsApprove:    req.IsVoteYes,
			RefuseReason: refuseReason,
		}),
	}, ctx)

	return response.BuildTransactionResponse{
		TxBytes: txBytes,
	}, err
}

// UpdateBooksNeed implements business.IChildService.
func (c *childService) UpdateBooksNeed(req request.UpdateChildNeedRequest, ctx context.Context) (response.BuildTransactionResponse, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if !util.IsValidSuiAddressStrict(req.ChildID) || !util.IsValidSuiAddressStrict(req.NeedID) {
		return response.BuildTransactionResponse{}, genericErr
	}

	var client = c.clients[constant.SuiTestnet]
	child, err := on_chain.GetOnChainObject[entities.Child](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  req.ChildID,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	if child == nil || !slices.Contains(child.BooksNeeds, req.NeedID) {
		return response.BuildTransactionResponse{}, genericErr
	}

	var sender string = ctx.Value("address").(string)
	var staffModule = on_chain.InitializeModuleStaff()
	staffNfts, err := on_chain.GetOnChainOwnedObjects[entities.StaffNft](on_chain.GetOnChainOwnedObjectsRequest{
		Client:       client,
		OwnerAddress: sender,
		StructType:   fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), staffModule.GetModule(), staffModule.GetStaffNftObjectStruct()),
		ErrLogger:    c.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	var genericRightErr error = errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	if staffNfts == nil || len(staffNfts) == 0 {
		return response.BuildTransactionResponse{}, genericRightErr
	}

	var leaderNftId string
	for _, nft := range staffNfts {
		if nft.Role == local_leader_role && nft.Region == child.Region {
			leaderNftId = nft.ID.ID
			break
		}
	}

	if leaderNftId == "" {
		return response.BuildTransactionResponse{}, genericRightErr
	}

	if req.Value == nil {
		return response.BuildTransactionResponse{}, nil
	}

	if *req.Value < 10_000 {
		return response.BuildTransactionResponse{}, errors.New(noti.NEED_VALUE_INVALID_WARN_MSG)
	}

	need, err := on_chain.GetOnChainObject[entities.BooksNeed](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  req.NeedID,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	var curTime time.Time = time.Now()
	if need.IsUpdated {
		if slices.Contains(need.YearChanges, fmt.Sprint(curTime.Year())) {
			return response.BuildTransactionResponse{}, errors.New(noti.CHILD_NEED_UPDATED_MESSAGE)
		}

		editDates, err := on_chain.GetOnChainObject[entities.EditNeedDates](on_chain.GetOnChainObjectRequest{
			Client:    client,
			ObjectId:  os.Getenv(env.EDIT_BOOKS_NEED_DATES_ID),
			ErrLogger: c.errLogger,
		}, ctx)
		if err != nil {
			return response.BuildTransactionResponse{}, err
		}

		var startDate time.Time = util.ToStartOfDate(util.RawDateToTime(fmt.Sprintf("%s/%d", editDates.StartDate, curTime.Year())))
		var endDate time.Time = util.ToEndOfDate(util.RawDateToTime(fmt.Sprintf("%s/%d", editDates.EndDate, curTime.Year())))
		if curTime.Before(startDate) || curTime.After(endDate) {
			return response.BuildTransactionResponse{}, errors.New(noti.NOTE_UPDATE_CHILD_NEED_DATE_MESSAGE)
		}
	}

	var childModule = on_chain.InitializeModuleChild()
	txBytes, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
		Client:    client,
		Sender:    sender,
		Module:    childModule.GetModule(),
		Function:  childModule.GetFunctionUpdateChildBooksNeed(),
		ErrLogger: c.errLogger,
		Arguments: childModule.ToUpdateChildNeedArguments(on_chain.UpdateChildNeedArguments{
			StaffNft: leaderNftId,
			ChildID:  req.ChildID,
			NeedID:   req.NeedID,
			Year:     curTime.Year(),
			Value:    *req.Value,
		}),
	}, ctx)

	return response.BuildTransactionResponse{
		TxBytes: txBytes,
	}, err
}

// UpdateHealthInsuranceNeed implements business.IChildService.
func (c *childService) UpdateHealthInsuranceNeed(req request.UpdateChildNeedRequest, ctx context.Context) (response.BuildTransactionResponse, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if !util.IsValidSuiAddressStrict(req.ChildID) || !util.IsValidSuiAddressStrict(req.NeedID) {
		return response.BuildTransactionResponse{}, genericErr
	}

	var client = c.clients[constant.SuiTestnet]
	child, err := on_chain.GetOnChainObject[entities.Child](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  req.ChildID,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	if child == nil || child.HealthInsuranceNeed != req.NeedID {
		return response.BuildTransactionResponse{}, genericErr
	}

	var sender string = ctx.Value("address").(string)
	var staffModule = on_chain.InitializeModuleStaff()
	staffNfts, err := on_chain.GetOnChainOwnedObjects[entities.StaffNft](on_chain.GetOnChainOwnedObjectsRequest{
		Client:       client,
		OwnerAddress: sender,
		StructType:   fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), staffModule.GetModule(), staffModule.GetStaffNftObjectStruct()),
		ErrLogger:    c.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	var genericRightErr error = errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	if staffNfts == nil || len(staffNfts) == 0 {
		return response.BuildTransactionResponse{}, genericRightErr
	}

	var leaderNftId string
	for _, nft := range staffNfts {
		if nft.Role == local_leader_role && nft.Region == child.Region {
			leaderNftId = nft.ID.ID
			break
		}
	}

	if leaderNftId == "" {
		return response.BuildTransactionResponse{}, genericRightErr
	}

	if req.Value == nil {
		return response.BuildTransactionResponse{}, nil
	}

	if *req.Value < 10_000 {
		return response.BuildTransactionResponse{}, errors.New(noti.NEED_VALUE_INVALID_WARN_MSG)
	}

	need, err := on_chain.GetOnChainObject[entities.HealthInsuranceNeed](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  req.NeedID,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	var curTime time.Time = time.Now()
	if need.IsUpdated {
		if slices.Contains(need.YearChanges, fmt.Sprint(curTime.Year())) {
			return response.BuildTransactionResponse{}, errors.New(noti.CHILD_NEED_UPDATED_MESSAGE)
		}

		editDates, err := on_chain.GetOnChainObject[entities.EditNeedDates](on_chain.GetOnChainObjectRequest{
			Client:    client,
			ObjectId:  os.Getenv(env.EDIT_HEALTH_INSURANCE_NEED_DATES_ID),
			ErrLogger: c.errLogger,
		}, ctx)
		if err != nil {
			return response.BuildTransactionResponse{}, err
		}

		var startDate time.Time = util.ToStartOfDate(util.RawDateToTime(fmt.Sprintf("%s/%d", editDates.StartDate, curTime.Year())))
		var endDate time.Time = util.ToEndOfDate(util.RawDateToTime(fmt.Sprintf("%s/%d", editDates.EndDate, curTime.Year())))
		if curTime.Before(startDate) || curTime.After(endDate) {
			return response.BuildTransactionResponse{}, errors.New(noti.NOTE_UPDATE_CHILD_NEED_DATE_MESSAGE)
		}
	}

	var childModule = on_chain.InitializeModuleChild()
	txBytes, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
		Client:    client,
		Sender:    sender,
		Module:    childModule.GetModule(),
		Function:  childModule.GetFunctionUpdateChildHealthInsuranceNeed(),
		ErrLogger: c.errLogger,
		Arguments: childModule.ToUpdateChildNeedArguments(on_chain.UpdateChildNeedArguments{
			StaffNft: leaderNftId,
			ChildID:  req.ChildID,
			NeedID:   req.NeedID,
			Year:     curTime.Year(),
			Value:    *req.Value,
		}),
	}, ctx)

	return response.BuildTransactionResponse{
		TxBytes: txBytes,
	}, err
}

// UpdateMealNeed implements business.IChildService.
func (c *childService) UpdateMealNeed(req request.UpdateChildNeedRequest, ctx context.Context) (response.BuildTransactionResponse, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if !util.IsValidSuiAddressStrict(req.ChildID) || !util.IsValidSuiAddressStrict(req.NeedID) {
		return response.BuildTransactionResponse{}, genericErr
	}

	var client = c.clients[constant.SuiTestnet]
	child, err := on_chain.GetOnChainObject[entities.Child](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  req.ChildID,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	if child == nil || child.MealNeed != req.NeedID {
		return response.BuildTransactionResponse{}, genericErr
	}

	var sender string = ctx.Value("address").(string)
	var staffModule = on_chain.InitializeModuleStaff()
	staffNfts, err := on_chain.GetOnChainOwnedObjects[entities.StaffNft](on_chain.GetOnChainOwnedObjectsRequest{
		Client:       client,
		OwnerAddress: sender,
		StructType:   fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), staffModule.GetModule(), staffModule.GetStaffNftObjectStruct()),
		ErrLogger:    c.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	var genericRightErr error = errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	if staffNfts == nil || len(staffNfts) == 0 {
		return response.BuildTransactionResponse{}, genericRightErr
	}

	var leaderNftId string
	for _, nft := range staffNfts {
		if nft.Role == local_leader_role && nft.Region == child.Region {
			leaderNftId = nft.ID.ID
			break
		}
	}

	if leaderNftId == "" {
		return response.BuildTransactionResponse{}, genericRightErr
	}

	if req.Value == nil {
		return response.BuildTransactionResponse{}, nil
	}

	if *req.Value < 10_000 {
		return response.BuildTransactionResponse{}, errors.New(noti.NEED_VALUE_INVALID_WARN_MSG)
	}

	need, err := on_chain.GetOnChainObject[entities.MealNeed](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  req.NeedID,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	var curTime time.Time = time.Now()
	if need.IsUpdated {
		var rawCurYear string = fmt.Sprintf("%d", curTime.Year())
		if need.Year == rawCurYear {
			return response.BuildTransactionResponse{}, errors.New(noti.CHILD_NEED_UPDATED_MESSAGE)
		}

		editDates, err := on_chain.GetOnChainObject[entities.EditNeedDates](on_chain.GetOnChainObjectRequest{
			Client:    client,
			ObjectId:  os.Getenv(env.EDIT_BOOKS_NEED_DATES_ID),
			ErrLogger: c.errLogger,
		}, ctx)
		if err != nil {
			return response.BuildTransactionResponse{}, err
		}

		var startDate time.Time = util.ToStartOfDate(util.RawDateToTime(fmt.Sprintf("%s/%s", editDates.StartDate, rawCurYear)))
		var endDate time.Time = util.ToEndOfDate(util.RawDateToTime(fmt.Sprintf("%s/%s", editDates.EndDate, rawCurYear)))
		if curTime.Before(startDate) || curTime.After(endDate) {
			return response.BuildTransactionResponse{}, errors.New(noti.NOTE_UPDATE_CHILD_NEED_DATE_MESSAGE)
		}
	}

	var childModule = on_chain.InitializeModuleChild()
	txBytes, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
		Client:    client,
		Sender:    sender,
		Module:    childModule.GetModule(),
		Function:  childModule.GetFunctionUpdateChildMealNeed(),
		ErrLogger: c.errLogger,
		Arguments: childModule.ToUpdateChildNeedArguments(on_chain.UpdateChildNeedArguments{
			StaffNft: leaderNftId,
			ChildID:  req.ChildID,
			NeedID:   req.NeedID,
			Year:     curTime.Year(),
			Value:    *req.Value,
		}),
	}, ctx)

	return response.BuildTransactionResponse{
		TxBytes: txBytes,
	}, err
}

func (c *childService) getGetChildrenRediskey(req request.GetChildrenRequest) string {
	var keyword string = "empty"
	if req.Keyword != "" {
		keyword = req.Keyword
	}

	var region string = "empty"
	if req.Region != "" {
		region = req.Region
	}

	var yob string = "empty"
	if req.YearOfBirth != nil {
		yob = fmt.Sprintf("%d", *req.YearOfBirth)
	}

	var gender string = "empty"
	if req.Gender != "" {
		gender = req.Gender
	}

	return fmt.Sprintf("child:kw:%s:r:%s:y:%s:o:%s:g:%s:s:%d:p:%d",
		keyword, region, yob, req.SortOrder, gender, req.PageSize, req.Page)
}

func (c *childService) getGetChildRediskey(id string) string {
	return fmt.Sprintf("child:%s", id)
}
