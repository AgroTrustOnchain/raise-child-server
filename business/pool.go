package business

import (
	"context"
	"errors"
	"log"
	"os"
	"raise-child/constants/env"
	"raise-child/constants/noti"
	"raise-child/constants/shared"
	"raise-child/interfaces/business"
	"raise-child/model/dtos/response"
	"raise-child/model/entities"
	"raise-child/util"
	on_chain "raise-child/util/on_chain"
	"strconv"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/sui"
)

type poolService struct {
	clients   map[string]sui.ISuiAPI
	errLogger *log.Logger
}

func initializePoolService(
	clients map[string]sui.ISuiAPI,
	errLogger *log.Logger,
) business.IPoolService {
	return &poolService{
		clients:   clients,
		errLogger: errLogger,
	}
}

func GeneratePoolService() (business.IPoolService, error) {
	var errLogger = util.GetLogConfig(shared.ERROR_LEVEL)

	return initializePoolService(
		_networkAliases,
		errLogger,
	), nil
}

// GetLeaderPool implements business.IPoolService.
func (p *poolService) GetLeaderPool(id string, ctx context.Context) (response.PoolResponse, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if !util.IsValidSuiAddressStrict(id) {
		return response.PoolResponse{}, genericErr
	}

	var client = p.clients[constant.SuiTestnet]
	manage, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: p.errLogger,
	}, ctx)
	if err != nil {
		return response.PoolResponse{}, err
	}

	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)
	if manage == nil {
		return response.PoolResponse{}, internalErr
	}

	var leaderNftId string
	for i, leader := range manage.LocalLeaderIds {
		if leader == id {
			leaderNftId = manage.LocalLeaderNfts[i]
			break
		}
	}

	if leaderNftId == "" {
		return response.PoolResponse{}, errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	}

	leaderNft, err := on_chain.GetOnChainObject[entities.StaffNft](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  leaderNftId,
		ErrLogger: p.errLogger,
	}, ctx)
	if err != nil {
		return response.PoolResponse{}, err
	}

	if leaderNft == nil {
		return response.PoolResponse{}, internalErr
	}

	pool, err := on_chain.GetOnChainObject[entities.MainPool](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.POOL_ID),
		ErrLogger: p.errLogger,
	}, ctx)
	if err != nil {
		return response.PoolResponse{}, err
	}

	if pool == nil {
		return response.PoolResponse{}, internalErr
	}

	localPools, err := on_chain.GetOnChainObjects[entities.LocalPool](on_chain.GetOnChainObjectsRequest{
		Client:    client,
		ObjectIds: pool.LocalPools,
		ErrLogger: p.errLogger,
	}, ctx)
	if err != nil {
		return response.PoolResponse{}, err
	}

	if localPools == nil {
		return response.PoolResponse{}, internalErr
	}

	for _, localPool := range localPools {
		if localPool.Region == leaderNft.Region {
			totalDonation, _ := strconv.ParseInt(localPool.TotalAmount, 10, 64)
			return response.PoolResponse{
				ID:            localPool.ID.ID,
				PoolName:      localPool.Region,
				TotalDonation: totalDonation,
			}, nil
		}
	}

	return response.PoolResponse{}, genericErr
}
