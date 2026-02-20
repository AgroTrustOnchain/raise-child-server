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

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/block-vision/sui-go-sdk/utils"
)

type notiService struct {
	volunteerNotiRepo i_repository.IVolunteerNotiRepository
	leaderNotiRepo    i_repository.ILeaderNotiRepository
	clients           map[string]sui.ISuiAPI
	errLogger         *log.Logger
}

func InitializeNotiService(db *sql.DB, errLogger *log.Logger) business.INotiService {
	return &notiService{
		volunteerNotiRepo: repository.InitializeVolunteerNotiRepository(db, errLogger),
		leaderNotiRepo:    repository.InitializeLeaderNotiRepository(db, errLogger),
		clients:           _networkAliases,
		errLogger:         errLogger,
	}
}

func GenerateNotiService() (business.INotiService, error) {
	var errLogger = util.GetLogConfig(shared.ERROR_LEVEL)

	cnn, err := db.ConnectDB(errLogger, db.InitializePostgreSQL())
	if err != nil {
		return nil, err
	}

	return InitializeNotiService(cnn, errLogger), nil
}

// GetCurrentWalletNotis implements business.INotiService.
func (n *notiService) GetCurrentWalletNotis(wallet string, req request.GetNotisRequest, ctx context.Context) (response.PaginationDataResponse, error) {
	if !utils.IsValidSuiAddress(models.SuiAddress(wallet)) {
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

	if req.Page < 1 {
		req.Page = 1
	}

	if req.PageSize < 1 {
		req.PageSize = default_page_size
	}

	var role string = ctx.Value("role").(string)
	var data []response.NotiResponse
	var isFound bool = false
	for _, nft := range nfts {
		if nft.Role == role {
			if role == volunteer_role {
				notis, err := n.volunteerNotiRepo.GetCurrentVolunteerNotis(req, wallet, ctx)
				if err != nil {
					return response.PaginationDataResponse{}, err
				}

				if notis != nil && len(notis) > 0 {
					var notiType string = "Volunteer Notification"
					for _, noti := range notis {
						data = append(data, response.NotiResponse{
							ID:      noti.ID,
							Content: noti.Content,
							Type:    notiType,
						})
					}
				}

				isFound = true
				break
			} else if role == local_leader_role {
				notis, err := n.leaderNotiRepo.GetCurrentLeaderNotis(req, wallet, ctx)
				if err != nil {
					return response.PaginationDataResponse{}, err
				}

				if notis != nil && len(notis) > 0 {
					var notiType string = "Leader Notification"
					for _, noti := range notis {
						data = append(data, response.NotiResponse{
							ID:      noti.ID,
							Content: noti.Content,
							Type:    notiType,
						})
					}
				}

				isFound = true
				break
			}
		}
	}

	if !isFound {
		return response.PaginationDataResponse{}, errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	}

	return response.PaginationDataResponse{
		Data: data,
		Page: req.Page,
	}, nil
}
