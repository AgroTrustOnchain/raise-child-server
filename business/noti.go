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

	manage, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    n.clients[constant.SuiTestnet],
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: n.errLogger,
	}, ctx)
	if err != nil {
		return response.PaginationDataResponse{}, err
	}

	if !slices.Contains(manage.LocalLeaderIds, wallet) {
		return response.PaginationDataResponse{}, errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	}

	if req.Page < 1 {
		req.Page = 1
	}

	if req.PageSize < 1 {
		req.PageSize = default_page_size
	}

	data, err := n.leaderNotiRepo.GetCurrentLeaderNotis(req, wallet, ctx)
	var amount int
	if len(data) == 0 {
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
