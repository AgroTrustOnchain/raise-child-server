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

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/sui"
)

type notiService struct {
	leaderNotiRepo i_repository.ILeaderNotiRepository
	redisCache     cache.IRedisCache
	clients        map[string]sui.ISuiAPI
	errLogger      *log.Logger
}

func initializeNotiService(
	leaderNotiRepo i_repository.ILeaderNotiRepository,
	clients map[string]sui.ISuiAPI,
	errLogger *log.Logger,
) business.INotiService {
	return &notiService{
		leaderNotiRepo: leaderNotiRepo,
		redisCache:     cache.InitializeRedisCache(),
		clients:        clients,
		errLogger:      errLogger,
	}
}

func GenerateNotiService() (business.INotiService, error) {
	var errLogger = util.GetLogConfig(shared.ERROR_LEVEL)

	cnn, err := db.ConnectDB(errLogger, db.InitializePostgreSQL())
	if err != nil {
		return nil, err
	}

	return initializeNotiService(
		repository.InitializeLeaderNotiRepository(cnn, errLogger),
		_networkAliases,
		errLogger,
	), nil
}

// GetCurrentWalletNotis implements business.INotiService.
func (n *notiService) GetCurrentWalletNotis(wallet string, req request.GetNotisRequest, ctx context.Context) (response.PaginationDataResponse, error) {
	if !util.IsValidSuiAddressStrict(wallet) {
		return response.PaginationDataResponse{}, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	var module = on_chain.InitializeModuleStaff()
	nfts, err := on_chain.GetOnChainOwnedObjects[entities.StaffNft](on_chain.GetOnChainOwnedObjectsRequest{
		Client:       n.clients[constant.SuiTestnet],
		OwnerAddress: wallet,
		StructType:   fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), module.GetModule(), module.GetStaffNftObjectStruct()),
		ErrLogger:    n.errLogger,
	}, ctx)
	if err != nil {
		return response.PaginationDataResponse{}, err
	}

	var genericRightErr error = errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	if nfts == nil || len(nfts) == 0 {
		return response.PaginationDataResponse{}, genericRightErr
	}

	var isLeader bool = false
	for _, nft := range nfts {
		if nft.Role == local_leader_role {
			isLeader = true
			break
		}
	}

	if !isLeader {
		return response.PaginationDataResponse{}, genericRightErr
	}

	if req.Page < 1 {
		req.Page = 1
	}

	if req.PageSize < 1 {
		req.PageSize = default_page_size
	}

	var res response.PaginationDataResponse
	var redisKey string = n.getGetCurrentWalletNotisRedisKey(wallet, req)
	if n.redisCache.Get(redisKey, &res, ctx) {
		return res, nil
	}

	data, err := n.leaderNotiRepo.GetCurrentLeaderNotis(req, wallet, ctx)
	var amount int
	if data == nil || len(data) == 0 {
		amount = 0
	} else {
		amount = len(data)
	}

	return response.PaginationDataResponse{
		Data:   data,
		Amount: amount,
		Page:   req.Page,
	}, nil
}

func (n *notiService) getGetCurrentWalletNotisRedisKey(wallet string, req request.GetNotisRequest) string {
	return fmt.Sprintf("noti:of:%s:s:%d:p:%d", wallet, req.PageSize, req.Page)
}
