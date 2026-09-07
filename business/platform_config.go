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
	"github.com/block-vision/sui-go-sdk/sui"
)

type platformConfigService struct {
	platformConfigRepo i_repository.IPlatformConfigRepository
	clients            map[string]sui.ISuiAPI
	errLogger          *log.Logger
}

func initializePlatformConfigService(
	platformConfigRepo i_repository.IPlatformConfigRepository,
	clients map[string]sui.ISuiAPI,
	errLogger *log.Logger,
) business.IPlatformConfigService {
	return &platformConfigService{
		platformConfigRepo: platformConfigRepo,
		clients:            clients,
		errLogger:          errLogger,
	}
}

func GeneratePlatformConfigService() (business.IPlatformConfigService, error) {
	var errLogger = util.GetLogConfig(shared.ERROR_LEVEL)

	cnn, err := db.ConnectDB(errLogger, db.InitializePostgreSQL())
	if err != nil {
		return nil, err
	}

	return initializePlatformConfigService(
		repository.InitializePlatformConfigRepository(cnn, errLogger),
		_networkAliases,
		errLogger,
	), nil
}

// GetConfig implements business.IPlatformConfigService.
func (p *platformConfigService) GetConfig(id string, isNumericConfig bool, ctx context.Context) (*entities.PlatformConfig, error) {
	return p.platformConfigRepo.GetConfig(id, p.getPlatformConfigTable(isNumericConfig), ctx)
}

// GetConfigs implements business.IPlatformConfigService.
func (p *platformConfigService) GetConfigs(req request.GetPlatformConfigsRequest, isNumericConfig bool, ctx context.Context) (response.PaginationDataResponse, error) {
	if req.ActorAddress != "" {
		if !util.IsValidSuiAddressStrict(req.ActorAddress) {
			return response.PaginationDataResponse{}, errors.New(noti.GENERIC_ERROR_WARN_MSG)
		}
	}

	req.Keyword = strings.TrimSpace(req.Keyword)
	if req.PageSize < 1 {
		req.PageSize = default_page_size
	}

	if req.Page < 1 {
		req.Page = 1
	}

	data, pages, err := p.platformConfigRepo.GetConfigs(req, p.getPlatformConfigTable(isNumericConfig), ctx)
	var amount int
	if len(data) == 0 {
		amount = 0
	} else {
		amount = len(data)
	}

	return response.PaginationDataResponse{
		Data:       data,
		Amount:     amount,
		Page:       req.Page,
		TotalPages: pages,
	}, err
}

// UpdateConfig implements business.IPlatformConfigService.
func (p *platformConfigService) UpdateConfig(id string, req request.UpdatePlatformConfigRequest, isNumericConfig bool, ctx context.Context) error {
	if req.Value == nil && req.Description == "" {
		return nil
	}

	var table string = p.getPlatformConfigTable(isNumericConfig)
	res, err := p.platformConfigRepo.GetConfig(id, table, ctx)
	if err != nil {
		return err
	}

	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if res == nil {
		return genericErr
	}

	manage, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    p.clients[constant.SuiTestnet],
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: p.errLogger,
	}, ctx)
	if err != nil {
		return err
	}

	if manage == nil {
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	var sender string = ctx.Value("address").(string)
	if !slices.Contains(manage.AdminIds, sender) {
		return errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	}

	if req.Value != nil {
		switch res.Value.(type) {
		case int, int64, float64:
			switch req.Value.(type) {
			case int, int64, float64:
				var ageMaxLimitConfig string = res.GetChildAgeMaxLimitConfigKey()
				var ageMinLimitConfig string = res.GetChildAgeMinLimitConfigKey()
				var minRegionStaffsConfig string = res.GetMinRegionStaffsConfigKey()
				if res.Key == ageMaxLimitConfig || res.Key == ageMinLimitConfig {
					if req.Value.(int) < 0 {
						return errors.New(noti.NEGATIVE_AGE_WARN_MSG)
					}

					switch res.Key {
					case ageMaxLimitConfig:
						child_max_age_accepted = req.Value.(int64)
					case ageMinLimitConfig:
						child_min_age_accepted = req.Value.(int64)
					case minRegionStaffsConfig:
						min_region_staffs = req.Value.(int64)
					}
				}
			default:
				return genericErr
			}
		case string:
			if _, ok := req.Value.(string); !ok {
				return genericErr
			}
		}

		res.Value = req.Value
	}

	if req.Description != "" {
		res.Description = strings.TrimSpace(req.Description)
	}

	var profileId string = ctx.Value("sub").(string)
	res.ActorProfileID = &profileId
	res.ActorAddress = &sender
	res.UpdatedAt = time.Now()

	return p.platformConfigRepo.UpdateConfig(*res, table, ctx)
}

func (p *platformConfigService) getPlatformConfigTable(isNumericConfig bool) string {
	var table string
	if isNumericConfig {
		var x *entities.NumericConfig
		table = x.GetTable()
	} else {
		var x *entities.StringConfig
		table = x.GetTable()
	}

	return table
}
