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
	if !util.IsValidSuiAddressStrict(id) {
		return response.TransactionResponse{}, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	res, err := on_chain.GetOnChainObject[entities.Transaction](on_chain.GetOnChainObjectRequest{
		Client:    t.clients[constant.SuiTestnet],
		ObjectId:  id,
		ErrLogger: t.errLogger,
	}, ctx)

	return res.ToTransactionResponse(), err
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

	req.SortOrder = util.StanderizeSortOrder(req.SortOrder)
	req.SortCriteria = util.StanderizeSortCriteria(req.SortCriteria)
	req.Keyword = util.StanderizeString(req.Keyword)
	if req.Page < 1 {
		req.Page = 1
	}

	if req.PageSize < 1 {
		req.PageSize = default_page_size
	}

	var res response.PaginationDataResponse
	var redisKey string = t.getGetTransactionRecordsRedisKey(req)
	if t.redisCache.Get(redisKey, &res, ctx) {
		return res, nil
	}

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
			return res, err
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

	if txs == nil || len(txs) == 0 {
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

	t.redisCache.Set(redisKey, res, time.Minute*5, ctx)

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

	return fmt.Sprintf("tx_record:kw:%s:pool:%s:of:%s:type:%s:min:%s:max:%s:%s:sc:%s:o:%s:s:%d:p:%d",
		keyword, poolId, actor, actionType, minAmount, maxAmount, sortCriteria, req.SortOrder, req.PageSize, req.Page)
}
