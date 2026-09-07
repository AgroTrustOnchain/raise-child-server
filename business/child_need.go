package business

import (
	"context"
	"errors"
	"log"
	"raise-child/constants/noti"
	"raise-child/constants/shared"
	"raise-child/interfaces/business"
	"raise-child/model/dtos/response"
	"raise-child/model/entities"
	"raise-child/util"
	on_chain "raise-child/util/on_chain"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/sui"
)

type childNeedService struct {
	clients   map[string]sui.ISuiAPI
	errLogger *log.Logger
}

func initializeChildNeedService(
	clients map[string]sui.ISuiAPI,
	errLogger *log.Logger,
) business.IChildNeedService {
	return &childNeedService{
		clients:   clients,
		errLogger: errLogger,
	}
}

func GenerateChildNeedService() (business.IChildNeedService, error) {
	return initializeChildNeedService(
		_networkAliases,
		util.GetLogConfig(shared.ERROR_LEVEL),
	), nil
}

// GetChildBooksNeed implements business.IChildNeedService.
func (c *childNeedService) GetChildBooksNeed(id string, ctx context.Context) (response.BooksNeedResponse, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if !util.IsValidSuiAddressStrict(id) {
		return response.BooksNeedResponse{}, genericErr
	}

	res, err := on_chain.GetOnChainObject[entities.BooksNeed](on_chain.GetOnChainObjectRequest{
		Client:    c.clients[constant.SuiTestnet],
		ObjectId:  id,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return response.BooksNeedResponse{}, err
	}

	if res == nil {
		return response.BooksNeedResponse{}, genericErr
	}

	return res.ToBooksNeedReponse(), nil
}

// GetChildHealthInsuranceNeed implements business.IChildNeedService.
func (c *childNeedService) GetChildHealthInsuranceNeed(id string, ctx context.Context) (response.HealthInsuranceNeedResponse, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if !util.IsValidSuiAddressStrict(id) {
		return response.HealthInsuranceNeedResponse{}, genericErr
	}

	res, err := on_chain.GetOnChainObject[entities.HealthInsuranceNeed](on_chain.GetOnChainObjectRequest{
		Client:    c.clients[constant.SuiTestnet],
		ObjectId:  id,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return response.HealthInsuranceNeedResponse{}, err
	}

	if res == nil {
		return response.HealthInsuranceNeedResponse{}, genericErr
	}

	return res.ToHealthInsuranceNeedReponse(), nil
}

// GetChildMealNeed implements business.IChildNeedService.
func (c *childNeedService) GetChildMealNeed(id string, ctx context.Context) (response.MealNeedResponse, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if !util.IsValidSuiAddressStrict(id) {
		return response.MealNeedResponse{}, genericErr
	}

	res, err := on_chain.GetOnChainObject[entities.MealNeed](on_chain.GetOnChainObjectRequest{
		Client:    c.clients[constant.SuiTestnet],
		ObjectId:  id,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return response.MealNeedResponse{}, err
	}

	if res == nil {
		return response.MealNeedResponse{}, genericErr
	}

	return res.ToMealNeedResponse(), nil
}
