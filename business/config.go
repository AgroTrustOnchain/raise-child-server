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
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/model/entities"
	"raise-child/util"
	on_chain "raise-child/util/on_chain"
	"slices"
	"strings"
	"time"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/sui"
)

type configService struct {
	clients   map[string]sui.ISuiAPI
	errLogger *log.Logger
}

func initializeConfigService(clients map[string]sui.ISuiAPI, errLogger *log.Logger) business.IConfigService {
	return &configService{
		clients:   clients,
		errLogger: errLogger,
	}
}

func GenerateConfigService() (business.IConfigService, error) {
	var errLogger = util.GetLogConfig(shared.ERROR_LEVEL)
	return initializeConfigService(_networkAliases, errLogger), nil
}

// UpdateChildEditBooksNeedDates implements business.IConfigService.
func (c *configService) UpdateChildEditBooksNeedDates(req request.UpdateChildEditNeedDatesRequest, ctx context.Context) (response.BuildTransactionResponse, error) {
	var client = c.clients[constant.SuiTestnet]
	manageObj, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	var sender string = ctx.Value("address").(string)
	if !slices.Contains(manageObj.AdminIds, sender) {
		return response.BuildTransactionResponse{}, errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	}

	editNeedDates, err := on_chain.GetOnChainObject[entities.EditNeedDates](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.EDIT_BOOKS_NEED_DATES_ID),
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	var rawStartDate string = strings.TrimSpace(req.StartDate)
	if rawStartDate == "" {
		rawStartDate = editNeedDates.StartDate
	}

	var curTime time.Time = time.Now()
	var startDate time.Time = util.ToStartOfDate(util.RawDateToTime(fmt.Sprintf("%s/%d", rawStartDate, curTime.Year())))
	var originalStartDate time.Time = util.ToStartOfDate(util.RawDateToTime(fmt.Sprintf("%s/%d", editNeedDates.StartDate, curTime.Year())))
	var invalidDateErr error = errors.New(noti.INVALID_DATE_MESSAGE)
	if startDate.IsZero() || startDate.Before(originalStartDate.AddDate(0, 0, -15)) || startDate.After(originalStartDate.AddDate(0, 0, 15)) { // Changes must be around 15 days from original date
		return response.BuildTransactionResponse{}, invalidDateErr
	}

	var rawEndDate string = strings.TrimSpace(req.StartDate)
	if rawEndDate == "" {
		rawEndDate = editNeedDates.EndDate
	}

	if rawStartDate == editNeedDates.StartDate && rawEndDate == editNeedDates.EndDate {
		return response.BuildTransactionResponse{}, nil
	}

	var endDate time.Time = util.ToEndOfDate(util.RawDateToTime(fmt.Sprintf("%s/%d", rawEndDate, curTime.Year())))
	if endDate.IsZero() || endDate.Before(startDate) {
		return response.BuildTransactionResponse{}, invalidDateErr
	}

	var needModule = on_chain.InitializeModuleNeed()
	txBytes, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
		Client:    client,
		Sender:    sender,
		Module:    needModule.GetModule(),
		Function:  needModule.GetFunctionEditUpdateBooksNeedDates(),
		ErrLogger: c.errLogger,
		Arguments: needModule.ToEditUpdateNeedDatesArguments(on_chain.EditUpdateNeedDatesArguments{
			EditDatesID: editNeedDates.ID.ID,
			StartDate:   rawStartDate,
			EndDate:     rawEndDate,
		}),
	}, ctx)

	return response.BuildTransactionResponse{
		TxBytes: txBytes,
	}, err
}

// UpdateChildEditHealthInsuranceNeedDates implements business.IConfigService.
func (c *configService) UpdateChildEditHealthInsuranceNeedDates(req request.UpdateChildEditNeedDatesRequest, ctx context.Context) (response.BuildTransactionResponse, error) {
	var client = c.clients[constant.SuiTestnet]
	manageObj, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	var sender string = ctx.Value("address").(string)
	if !slices.Contains(manageObj.AdminIds, sender) {
		return response.BuildTransactionResponse{}, errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	}

	editNeedDates, err := on_chain.GetOnChainObject[entities.EditNeedDates](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.EDIT_HEALTH_INSURANCE_NEED_DATES_ID),
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	var rawStartDate string = strings.TrimSpace(req.StartDate)
	if rawStartDate == "" {
		rawStartDate = editNeedDates.StartDate
	}

	var curTime time.Time = time.Now()
	var startDate time.Time = util.ToStartOfDate(util.RawDateToTime(fmt.Sprintf("%s/%d", rawStartDate, curTime.Year())))
	var originalStartDate time.Time = util.ToStartOfDate(util.RawDateToTime(fmt.Sprintf("%s/%d", editNeedDates.StartDate, curTime.Year())))
	var invalidDateErr error = errors.New(noti.INVALID_DATE_MESSAGE)
	if startDate.IsZero() || startDate.Before(originalStartDate.AddDate(0, 0, -15)) || startDate.After(originalStartDate.AddDate(0, 0, 15)) { // Changes must be around 15 days from original date
		return response.BuildTransactionResponse{}, invalidDateErr
	}

	var rawEndDate string = strings.TrimSpace(req.StartDate)
	if rawEndDate == "" {
		rawEndDate = editNeedDates.EndDate
	}

	if rawStartDate == editNeedDates.StartDate && rawEndDate == editNeedDates.EndDate {
		return response.BuildTransactionResponse{}, nil
	}

	var endDate time.Time = util.ToEndOfDate(util.RawDateToTime(fmt.Sprintf("%s/%d", rawEndDate, curTime.Year())))
	if endDate.IsZero() || endDate.Before(startDate) {
		return response.BuildTransactionResponse{}, invalidDateErr
	}

	var needModule = on_chain.InitializeModuleNeed()
	txBytes, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
		Client:    client,
		Sender:    sender,
		Module:    needModule.GetModule(),
		Function:  needModule.GetFunctionEditUpdateHealthInsuranceNeedDates(),
		ErrLogger: c.errLogger,
		Arguments: needModule.ToEditUpdateNeedDatesArguments(on_chain.EditUpdateNeedDatesArguments{
			EditDatesID: editNeedDates.ID.ID,
			StartDate:   rawStartDate,
			EndDate:     rawEndDate,
		}),
	}, ctx)

	return response.BuildTransactionResponse{
		TxBytes: txBytes,
	}, err
}

// UpdateChildEditMealNeedDates implements business.IConfigService.
func (c *configService) UpdateChildEditMealNeedDates(req request.UpdateChildEditNeedDatesRequest, ctx context.Context) (response.BuildTransactionResponse, error) {
	var client = c.clients[constant.SuiTestnet]
	manageObj, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	var sender string = ctx.Value("address").(string)
	if !slices.Contains(manageObj.AdminIds, sender) {
		return response.BuildTransactionResponse{}, errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	}

	editNeedDates, err := on_chain.GetOnChainObject[entities.EditNeedDates](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.EDIT_MEAL_NEED_DATES_ID),
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	var rawStartDate string = strings.TrimSpace(req.StartDate)
	if rawStartDate == "" {
		rawStartDate = editNeedDates.StartDate
	}

	var curTime time.Time = time.Now()
	var startDate time.Time = util.ToStartOfDate(util.RawDateToTime(fmt.Sprintf("%s/%d", rawStartDate, curTime.Year())))
	var originalStartDate time.Time = util.ToStartOfDate(util.RawDateToTime(fmt.Sprintf("%s/%d", editNeedDates.StartDate, curTime.Year())))
	var invalidDateErr error = errors.New(noti.INVALID_DATE_MESSAGE)
	if startDate.IsZero() || startDate.Before(originalStartDate.AddDate(0, 0, -15)) || startDate.After(originalStartDate.AddDate(0, 0, 15)) { // Changes must be around 15 days from original date
		return response.BuildTransactionResponse{}, invalidDateErr
	}

	var rawEndDate string = strings.TrimSpace(req.StartDate)
	if rawEndDate == "" {
		rawEndDate = editNeedDates.EndDate
	}

	if rawStartDate == editNeedDates.StartDate && rawEndDate == editNeedDates.EndDate {
		return response.BuildTransactionResponse{}, nil
	}

	var endDate time.Time = util.ToEndOfDate(util.RawDateToTime(fmt.Sprintf("%s/%d", rawEndDate, curTime.Year())))
	if endDate.IsZero() || endDate.Before(startDate) {
		return response.BuildTransactionResponse{}, invalidDateErr
	}

	var needModule = on_chain.InitializeModuleNeed()
	txBytes, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
		Client:    client,
		Sender:    sender,
		Module:    needModule.GetModule(),
		Function:  needModule.GetFunctionEditUpdateMealNeedDates(),
		ErrLogger: c.errLogger,
		Arguments: needModule.ToEditUpdateNeedDatesArguments(on_chain.EditUpdateNeedDatesArguments{
			EditDatesID: editNeedDates.ID.ID,
			StartDate:   rawStartDate,
			EndDate:     rawEndDate,
		}),
	}, ctx)

	return response.BuildTransactionResponse{
		TxBytes: txBytes,
	}, err
}
