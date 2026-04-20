package business

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"raise-child/constants/env"
	"raise-child/constants/shared"
	"raise-child/interfaces/business"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/model/entities"
	"raise-child/util"
	"raise-child/util/cache"
	on_chain "raise-child/util/on_chain"
	"sort"
	"strings"
	"time"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/sui"
)

type centerService struct {
	redisCache cache.IRedisCache
	clients    map[string]sui.ISuiAPI
	errLogger  *log.Logger
}

func initializeCenterService(clients map[string]sui.ISuiAPI, errLogger *log.Logger) business.ICenterService {
	return &centerService{
		redisCache: cache.InitializeRedisCache(),
		clients:    clients,
		errLogger:  errLogger,
	}
}

func GenerateCenterService() (business.ICenterService, error) {
	var errLogger = util.GetLogConfig(shared.ERROR_LEVEL)
	return initializeCenterService(
		_networkAliases,
		errLogger,
	), nil
}

// GetCenters implements business.ICenterService.
func (c *centerService) GetCenters(req request.GetCentersRequest, ctx context.Context) (response.PaginationDataResponse, error) {
	req.Keyword = util.StandardizeString(req.Keyword)
	if req.Page < 1 {
		req.Page = 1
	}

	if req.PageSize < 1 {
		req.PageSize = default_page_size
	}

	var redisKey string = c.getGetCentersRedisKey(req)
	var res response.PaginationDataResponse
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

	centers, err := on_chain.GetOnChainObjects[entities.Center](on_chain.GetOnChainObjectsRequest{
		Client:    client,
		ObjectIds: manageObj.CreatedCenters,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return response.PaginationDataResponse{}, err
	}

	if centers == nil || len(centers) == 0 {
		return response.PaginationDataResponse{
			Page:   req.Page,
			Amount: 0,
		}, nil
	}

	var filteredCenters []entities.Center
	for _, center := range centers {
		if req.Keyword != "" {
			if !strings.Contains(strings.ToLower(center.Region), req.Keyword) && !strings.Contains(strings.ToLower(center.CenterAddress), req.Keyword) && !strings.Contains(strings.ToLower(center.CenterPhoneNumber), req.Keyword) { // Not matched
				continue
			}
		}

		filteredCenters = append(filteredCenters, center)
	}

	var skippedRecords int = (req.Page - 1) * req.PageSize
	if len(filteredCenters) <= skippedRecords {
		return response.PaginationDataResponse{}, nil
	}

	sort.Slice(filteredCenters, func(i, j int) bool {
		var updatedAt1 time.Time = util.RawDateToTime(filteredCenters[i].UploadedAt)
		var updatedAt2 time.Time = util.RawDateToTime(filteredCenters[j].UpdatedAt)
		if req.SortOrder == "DESC" {
			return updatedAt2.After(updatedAt1)
		}

		return updatedAt1.Before(updatedAt2)
	})

	var data []response.CenterCardMinimumResponse
	for i := skippedRecords; i < len(filteredCenters); i++ {
		data = append(data, filteredCenters[i].ToCenterCardMinimumResponse())
		if len(data) == req.PageSize {
			break
		}
	}

	res = response.PaginationDataResponse{
		Data:       data,
		Amount:     len(data),
		Page:       req.Page,
		TotalPages: int(math.Ceil(float64(len(filteredCenters)) / float64(req.PageSize))),
	}

	c.redisCache.Set(redisKey, res, time.Minute*5, ctx)

	return res, nil
}

func (c *centerService) getGetCentersRedisKey(req request.GetCentersRequest) string {
	var keyword string = "empty"
	if req.Keyword != "" {
		keyword = req.Keyword
	}

	return fmt.Sprintf("center:kw:%s:o:%s:s:%d:p:%d",
		keyword, req.SortOrder, req.PageSize, req.Page,
	)
}
