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
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/model/entities"
	"raise-child/util"
	"raise-child/util/cache"
	on_chain "raise-child/util/on_chain"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/sui"
)

type transactionRecordService struct {
	redisCache cache.IRedisCache
	clients    map[string]sui.ISuiAPI
	errLogger  *log.Logger
}

func initializeTransactionRecordService(clients map[string]sui.ISuiAPI, errLogger *log.Logger) business.ITransactionRecordService {
	return &transactionRecordService{
		redisCache: cache.InitializeRedisCache(),
		clients:    clients,
		errLogger:  errLogger,
	}
}

func GenerateTransactionRecordService() (business.ITransactionRecordService, error) {
	var errLogger = util.GetLogConfig(shared.ERROR_LEVEL)
	return initializeTransactionRecordService(_networkAliases, errLogger), nil
}

// GetTransactionRecord implements business.ITransactionRecordService.
func (t *transactionRecordService) GetTransactionRecord(id string, ctx context.Context) (response.TransactionResponse, error) {
	// if !util.IsValidSuiAddressStrict(id) {
	// 	return response.TransactionResponse{}, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	// }

	// res, err := on_chain.GetOnChainObject[entities.Transaction](on_chain.GetOnChainObjectRequest{
	// 	Client:    t.clients[constant.SuiTestnet],
	// 	ObjectId:  id,
	// 	ErrLogger: t.errLogger,
	// }, ctx)

	// return res.ToTransactionResponse(), err

	// MOCK DATA
	for _, tx := range mockTxRecords {
		if tx.ID == id {
			return tx, nil
		}
	}

	return response.TransactionResponse{}, nil
}

// GetTransactionRecordsV2 implements business.ITransactionRecordService.
func (t *transactionRecordService) GetTransactionRecordsV2(req request.GetTransactionRecordsRequest, ctx context.Context) (response.GetTransactionRecordsResponse, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if req.Actor != "" {
		if !util.IsValidSuiAddressStrict(req.Actor) {
			return response.GetTransactionRecordsResponse{}, genericErr
		}
	}

	if req.PoolID != "" {
		if !util.IsValidSuiAddressStrict(req.PoolID) {
			return response.GetTransactionRecordsResponse{}, genericErr
		}
	}

	if req.MaxAmount != nil {
		if *req.MaxAmount < min_withdraw_proposal_amount_value {
			return response.GetTransactionRecordsResponse{}, nil
		}

		if req.MinAmount != nil {
			if *req.MaxAmount <= *req.MinAmount {
				return response.GetTransactionRecordsResponse{}, nil
			}
		}
	}

	req.SortOrder = util.StandardizeSortOrder(req.SortOrder)
	req.SortCriteria = util.StandardizeSortCriteria(req.SortCriteria)
	req.Keyword = util.StandardizeString(req.Keyword)
	if req.Page < 1 {
		req.Page = 1
	}

	if req.PageSize < 1 {
		req.PageSize = default_page_size
	}

	var res response.GetTransactionRecordsResponse

	var client = t.clients[constant.SuiTestnet]
	var txs []entities.Transaction
	var errRes error
	if req.Actor != "" {
		var recordModule = on_chain.InitializeModuleRecord()
		txs, errRes = on_chain.GetOnChainOwnedObjects[entities.Transaction](on_chain.GetOnChainOwnedObjectsRequest{
			Client:       client,
			OwnerAddress: req.Actor,
			StructType:   fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), recordModule.GetModule(), recordModule.GetTransactionRecordStruct()),
			ErrLogger:    t.errLogger,
		}, ctx)
	} else {
		manageObj, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
			Client:    client,
			ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
			ErrLogger: t.errLogger,
		}, ctx)
		if err != nil {
			return response.GetTransactionRecordsResponse{}, err
		}

		txs, errRes = on_chain.GetOnChainObjects[entities.Transaction](on_chain.GetOnChainObjectsRequest{
			Client:    client,
			ObjectIds: manageObj.TransactionRecords,
			ErrLogger: t.errLogger,
		}, ctx)
	}

	if errRes != nil {
		return res, errRes
	}

	if len(txs) == 0 {
		return res, nil
	}

	var poolName string
	if req.PoolID != "" && req.PoolID != os.Getenv(env.POOL_ID) {
		pool, _ := on_chain.GetOnChainObject[entities.LocalPool](on_chain.GetOnChainObjectRequest{
			Client:    client,
			ObjectId:  req.PoolID,
			ErrLogger: t.errLogger,
		}, ctx)

		if pool != nil {
			poolName = pool.Region

			totalDonation, _ := strconv.ParseInt(pool.TotalAmount, 10, 64)
			res.PoolID = req.PoolID
			res.PoolName = poolName
			res.PoolTotalDonation = totalDonation
		}
	} else {
		pool, _ := on_chain.GetOnChainObject[entities.MainPool](on_chain.GetOnChainObjectRequest{
			Client:    client,
			ObjectId:  os.Getenv(env.POOL_ID),
			ErrLogger: t.errLogger,
		}, ctx)

		if pool != nil {
			totalDonation, _ := strconv.ParseInt(pool.TotalAmount, 10, 64)
			res.PoolID = req.PoolID
			res.PoolName = "Main Pool"
			res.PoolTotalDonation = totalDonation
		}

		if req.PoolID != "" {
			poolName = "Main Pool"
		}
	}

	var filteredTxs []entities.Transaction
	for _, tx := range txs {
		if req.Keyword != "" {
			if !strings.Contains(tx.Message, req.Keyword) {
				continue
			}
		}

		if poolName != "" {
			if tx.PoolName != poolName {
				continue
			}
		}

		if req.ActionType != "" {
			if tx.ActionType != req.ActionType {
				continue
			}
		}

		amount, _ := strconv.ParseInt(tx.Amount, 10, 64)
		if req.MinAmount != nil {
			if amount < *req.MinAmount {
				continue
			}
		}

		if req.MaxAmount != nil {
			if amount > *req.MaxAmount {
				continue
			}
		}

		filteredTxs = append(filteredTxs, tx)
	}

	sort.Slice(filteredTxs, func(i, j int) bool {
		if req.SortCriteria == "amount" {
			amount1, _ := strconv.ParseInt(filteredTxs[i].Amount, 10, 64)
			amount2, _ := strconv.ParseInt(filteredTxs[j].Amount, 10, 64)
			if req.SortOrder == "DESC" {
				return amount2 > amount1
			}

			return amount2 < amount1
		}

		if req.SortOrder == "ASC" {
			return true
		}

		return false
	})

	var skippedRecords int = (req.Page - 1) * req.PageSize
	if len(filteredTxs) <= skippedRecords {
		return response.GetTransactionRecordsResponse{}, nil
	}

	var data []response.TransactionResponse
	for i := skippedRecords; i < len(filteredTxs); i++ {
		data = append(data, filteredTxs[i].ToTransactionResponse())
		if len(data) == req.PageSize {
			break
		}
	}

	res.Data = data
	res.Amount = len(data)
	res.Page = req.Page
	res.TotalPages = int(math.Ceil(float64(len(filteredTxs)) / float64(req.PageSize)))

	return res, nil
}

// GetTransactionRecords implements business.ITransactionRecordService.
func (t *transactionRecordService) GetTransactionRecords(req request.GetTransactionRecordsRequest, ctx context.Context) (response.PaginationDataResponse, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if req.Actor != "" {
		if !util.IsValidSuiAddressStrict(req.Actor) {
			return response.PaginationDataResponse{}, genericErr
		}
	}

	if req.PoolID != "" {
		if !util.IsValidSuiAddressStrict(req.PoolID) {
			return response.PaginationDataResponse{}, genericErr
		}
	}

	if req.MaxAmount != nil {
		if *req.MaxAmount < min_withdraw_proposal_amount_value {
			return response.PaginationDataResponse{}, nil
		}

		if req.MinAmount != nil {
			if *req.MaxAmount <= *req.MinAmount {
				return response.PaginationDataResponse{}, nil
			}
		}
	}

	req.SortOrder = util.StandardizeSortOrder(req.SortOrder)
	req.SortCriteria = util.StandardizeSortCriteria(req.SortCriteria)
	req.Keyword = util.StandardizeString(req.Keyword)
	if req.Page < 1 {
		req.Page = 1
	}

	if req.PageSize < 1 {
		req.PageSize = default_page_size
	}

	var res response.PaginationDataResponse

	var client = t.clients[constant.SuiTestnet]
	var txs []entities.Transaction
	var errRes error
	if req.Actor != "" {
		var recordModule = on_chain.InitializeModuleRecord()
		txs, errRes = on_chain.GetOnChainOwnedObjects[entities.Transaction](on_chain.GetOnChainOwnedObjectsRequest{
			Client:       client,
			OwnerAddress: req.Actor,
			StructType:   fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), recordModule.GetModule(), recordModule.GetTransactionRecordStruct()),
			ErrLogger:    t.errLogger,
		}, ctx)
	} else {
		manageObj, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
			Client:    client,
			ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
			ErrLogger: t.errLogger,
		}, ctx)
		if err != nil {
			return response.PaginationDataResponse{}, err
		}

		txs, errRes = on_chain.GetOnChainObjects[entities.Transaction](on_chain.GetOnChainObjectsRequest{
			Client:    client,
			ObjectIds: manageObj.TransactionRecords,
			ErrLogger: t.errLogger,
		}, ctx)
	}

	if errRes != nil {
		return res, errRes
	}

	if len(txs) == 0 {
		return res, nil
	}

	var poolName string
	if req.PoolID != "" {
		if req.PoolID == os.Getenv(env.POOL_ID) {
			poolName = "Main Pool"
		} else {
			pool, _ := on_chain.GetOnChainObject[entities.LocalPool](on_chain.GetOnChainObjectRequest{
				Client:    client,
				ObjectId:  req.PoolID,
				ErrLogger: t.errLogger,
			}, ctx)

			if pool != nil {
				poolName = pool.Region
			}
		}
	}

	var filteredTxs []entities.Transaction
	for _, tx := range txs {
		if req.Keyword != "" {
			if !strings.Contains(tx.Message, req.Keyword) {
				continue
			}
		}

		if poolName != "" {
			if tx.PoolName != poolName {
				continue
			}
		}

		if req.ActionType != "" {
			if tx.ActionType != req.ActionType {
				continue
			}
		}

		amount, _ := strconv.ParseInt(tx.Amount, 10, 64)
		if req.MinAmount != nil {
			if amount < *req.MinAmount {
				continue
			}
		}

		if req.MaxAmount != nil {
			if amount > *req.MaxAmount {
				continue
			}
		}

		filteredTxs = append(filteredTxs, tx)
	}

	sort.Slice(filteredTxs, func(i, j int) bool {
		if req.SortCriteria == "amount" {
			amount1, _ := strconv.ParseInt(filteredTxs[i].Amount, 10, 64)
			amount2, _ := strconv.ParseInt(filteredTxs[j].Amount, 10, 64)
			if req.SortOrder == "DESC" {
				return amount2 > amount1
			}

			return amount2 < amount1
		}

		if req.SortOrder == "ASC" {
			return true
		}

		return false
	})

	var skippedRecords int = (req.Page - 1) * req.PageSize
	if len(filteredTxs) <= skippedRecords {
		return response.PaginationDataResponse{}, nil
	}

	var data []response.TransactionResponse
	for i := skippedRecords; i < len(filteredTxs); i++ {
		data = append(data, filteredTxs[i].ToTransactionResponse())
		if len(data) == req.PageSize {
			break
		}
	}

	res = response.PaginationDataResponse{
		Data:       data,
		Amount:     len(data),
		Page:       req.Page,
		TotalPages: int(math.Ceil(float64(len(filteredTxs)) / float64(req.PageSize))),
	}

	return res, nil
}

func (t *transactionRecordService) getGetTransactionRecordsRedisKey(req request.GetTransactionRecordsRequest) string {
	var keyword string = "empty"
	if req.Keyword != "" {
		keyword = req.Keyword
	}

	var poolId string = "empty"
	if req.PoolID != "" {
		poolId = req.PoolID
	}

	var actor string = "empty"
	if req.Actor != "" {
		actor = req.Actor
	}

	var actionType string = "empty"
	if req.ActionType != "" {
		actionType = req.ActionType
	}

	var minAmount string = "empty"
	if req.MinAmount != nil {
		minAmount = fmt.Sprintf("%d", *req.MinAmount)
	}

	var maxAmount string = "empty"
	if req.MaxAmount != nil {
		maxAmount = fmt.Sprintf("%d", *req.MaxAmount)
	}

	var sortCriteria string = "empty"
	if req.SortCriteria != "" {
		sortCriteria = req.SortCriteria
	}

	return fmt.Sprintf("tx_record:kw:%s:pool:%s:of:%s:type:%s:min:%s:max:%s:sc:%s:o:%s:s:%d:p:%d",
		keyword, poolId, actor, actionType, minAmount, maxAmount, sortCriteria, req.SortOrder, req.PageSize, req.Page)
}

var mockTxRecords = getMockTxRecords()

func getMockTxRecords() []response.TransactionResponse {
	var res []response.TransactionResponse
	for i := 1; i <= 20; i++ {
		action := "Deposit"
		if i%2 == 0 {
			action = "Withdraw"
		}

		var tx = response.TransactionResponse{
			ID:           fmt.Sprintf("tx-uuid-%03d", i),
			ActorAddress: fmt.Sprintf("0xAccount%d", i*123),
			ActionType:   action,
			PoolName:     "Main Liquidity Pool",
			Amount:       int64(i * 100000),
			Message:      fmt.Sprintf("Transaction number %d", i),
			CoinType:     "VND",
			CreatedAt:    time.Now().Add(time.Duration(-i) * time.Hour), // Lùi lại i giờ
		}

		res = append(res, tx)
	}

	return res
}
