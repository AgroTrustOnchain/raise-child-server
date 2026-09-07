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
	"raise-child/model/entities"
	"raise-child/util"
	"raise-child/util/cache"
	on_chain "raise-child/util/on_chain"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/sui"
)

type configService struct {
	redisCache cache.IRedisCache
	clients    map[string]sui.ISuiAPI
	errLogger  *log.Logger
}

func initializeConfigService(clients map[string]sui.ISuiAPI, errLogger *log.Logger) business.IConfigService {
	return &configService{
		redisCache: cache.InitializeRedisCache(),
		clients:    clients,
		errLogger:  errLogger,
	}
}

func GenerateConfigService() (business.IConfigService, error) {
	var errLogger = util.GetLogConfig(shared.ERROR_LEVEL)
	return initializeConfigService(_networkAliases, errLogger), nil
}

// // UpdateChildEditBooksNeedDates implements business.IConfigService.
// func (c *configService) UpdateChildEditBooksNeedDates(req request.UpdateChildEditNeedDatesRequest, ctx context.Context) (response.BuildTransactionResponse, error) {
// 	var client = c.clients[constant.SuiTestnet]
// 	var manageObj entities.Manage
// 	if !c.redisCache.Get(manageObj.GetRedisKey(), &manageObj, ctx) {
// 		res, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
// 			Client:    client,
// 			ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
// 			ErrLogger: c.errLogger,
// 		}, ctx)
// 		if err != nil {
// 			return response.BuildTransactionResponse{}, err
// 		}

// 		if res != nil {
// 			c.redisCache.Set(manageObj.GetRedisKey(), res, time.Minute, ctx)
// 			manageObj = *res
// 		}
// 	}

// 	var sender string = ctx.Value("address").(string)
// 	if !slices.Contains(manageObj.AdminIds, sender) {
// 		return response.BuildTransactionResponse{}, errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
// 	}

// 	editNeedDates, err := on_chain.GetOnChainObject[entities.EditNeedDates](on_chain.GetOnChainObjectRequest{
// 		Client:    client,
// 		ObjectId:  os.Getenv(env.EDIT_BOOKS_NEED_DATES_ID),
// 		ErrLogger: c.errLogger,
// 	}, ctx)
// 	if err != nil {
// 		return response.BuildTransactionResponse{}, err
// 	}

// 	var rawStartDate string = strings.TrimSpace(req.StartDate)
// 	if rawStartDate == "" {
// 		rawStartDate = editNeedDates.StartDate
// 	}

// 	var curTime time.Time = time.Now()
// 	var startDate time.Time = util.ToStartOfDate(util.RawDateToTime(fmt.Sprintf("%s/%d", rawStartDate, curTime.Year())))
// 	var originalStartDate time.Time = util.ToStartOfDate(util.RawDateToTime(fmt.Sprintf("%s/%d", editNeedDates.StartDate, curTime.Year())))
// 	var invalidDateErr error = errors.New(noti.INVALID_DATE_MESSAGE)
// 	if startDate.IsZero() || startDate.Before(originalStartDate.AddDate(0, 0, -15)) || startDate.After(originalStartDate.AddDate(0, 0, 15)) { // Changes must be around 15 days from original date
// 		return response.BuildTransactionResponse{}, invalidDateErr
// 	}

// 	var rawEndDate string = strings.TrimSpace(req.StartDate)
// 	if rawEndDate == "" {
// 		rawEndDate = editNeedDates.EndDate
// 	}

// 	if rawStartDate == editNeedDates.StartDate && rawEndDate == editNeedDates.EndDate {
// 		return response.BuildTransactionResponse{}, nil
// 	}

// 	var endDate time.Time = util.ToEndOfDate(util.RawDateToTime(fmt.Sprintf("%s/%d", rawEndDate, curTime.Year())))
// 	if endDate.IsZero() || endDate.Before(startDate) {
// 		return response.BuildTransactionResponse{}, invalidDateErr
// 	}

// 	var needModule = on_chain.InitializeModuleNeed()
// 	txBytes, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
// 		Client:    client,
// 		Sender:    sender,
// 		Module:    needModule.GetModule(),
// 		Function:  needModule.GetFunctionEditUpdateBooksNeedDates(),
// 		ErrLogger: c.errLogger,
// 		Arguments: needModule.ToEditUpdateNeedDatesArguments(on_chain.EditUpdateNeedDatesArguments{
// 			EditDatesID: editNeedDates.ID.ID,
// 			StartDate:   rawStartDate,
// 			EndDate:     rawEndDate,
// 		}),
// 	}, ctx)

// 	return response.BuildTransactionResponse{
// 		TxBytes: txBytes,
// 	}, err
// }

// // UpdateChildEditHealthInsuranceNeedDates implements business.IConfigService.
// func (c *configService) UpdateChildEditHealthInsuranceNeedDates(req request.UpdateChildEditNeedDatesRequest, ctx context.Context) (response.BuildTransactionResponse, error) {
// 	var client = c.clients[constant.SuiTestnet]
// 	var manageObj entities.Manage
// 	if !c.redisCache.Get(manageObj.GetRedisKey(), &manageObj, ctx) {
// 		res, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
// 			Client:    client,
// 			ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
// 			ErrLogger: c.errLogger,
// 		}, ctx)
// 		if err != nil {
// 			return response.BuildTransactionResponse{}, err
// 		}

// 		if res != nil {
// 			c.redisCache.Set(manageObj.GetRedisKey(), res, time.Minute, ctx)
// 			manageObj = *res
// 		}
// 	}

// 	var sender string = ctx.Value("address").(string)
// 	if !slices.Contains(manageObj.AdminIds, sender) {
// 		return response.BuildTransactionResponse{}, errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
// 	}

// 	editNeedDates, err := on_chain.GetOnChainObject[entities.EditNeedDates](on_chain.GetOnChainObjectRequest{
// 		Client:    client,
// 		ObjectId:  os.Getenv(env.EDIT_HEALTH_INSURANCE_NEED_DATES_ID),
// 		ErrLogger: c.errLogger,
// 	}, ctx)
// 	if err != nil {
// 		return response.BuildTransactionResponse{}, err
// 	}

// 	var rawStartDate string = strings.TrimSpace(req.StartDate)
// 	if rawStartDate == "" {
// 		rawStartDate = editNeedDates.StartDate
// 	}

// 	var curTime time.Time = time.Now()
// 	var startDate time.Time = util.ToStartOfDate(util.RawDateToTime(fmt.Sprintf("%s/%d", rawStartDate, curTime.Year())))
// 	var originalStartDate time.Time = util.ToStartOfDate(util.RawDateToTime(fmt.Sprintf("%s/%d", editNeedDates.StartDate, curTime.Year())))
// 	var invalidDateErr error = errors.New(noti.INVALID_DATE_MESSAGE)
// 	if startDate.IsZero() || startDate.Before(originalStartDate.AddDate(0, 0, -15)) || startDate.After(originalStartDate.AddDate(0, 0, 15)) { // Changes must be around 15 days from original date
// 		return response.BuildTransactionResponse{}, invalidDateErr
// 	}

// 	var rawEndDate string = strings.TrimSpace(req.StartDate)
// 	if rawEndDate == "" {
// 		rawEndDate = editNeedDates.EndDate
// 	}

// 	if rawStartDate == editNeedDates.StartDate && rawEndDate == editNeedDates.EndDate {
// 		return response.BuildTransactionResponse{}, nil
// 	}

// 	var endDate time.Time = util.ToEndOfDate(util.RawDateToTime(fmt.Sprintf("%s/%d", rawEndDate, curTime.Year())))
// 	if endDate.IsZero() || endDate.Before(startDate) {
// 		return response.BuildTransactionResponse{}, invalidDateErr
// 	}

// 	var needModule = on_chain.InitializeModuleNeed()
// 	txBytes, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
// 		Client:    client,
// 		Sender:    sender,
// 		Module:    needModule.GetModule(),
// 		Function:  needModule.GetFunctionEditUpdateHealthInsuranceNeedDates(),
// 		ErrLogger: c.errLogger,
// 		Arguments: needModule.ToEditUpdateNeedDatesArguments(on_chain.EditUpdateNeedDatesArguments{
// 			EditDatesID: editNeedDates.ID.ID,
// 			StartDate:   rawStartDate,
// 			EndDate:     rawEndDate,
// 		}),
// 	}, ctx)

// 	return response.BuildTransactionResponse{
// 		TxBytes: txBytes,
// 	}, err
// }

// // UpdateChildEditMealNeedDates implements business.IConfigService.
// func (c *configService) UpdateChildEditMealNeedDates(req request.UpdateChildEditNeedDatesRequest, ctx context.Context) (response.BuildTransactionResponse, error) {
// 	var client = c.clients[constant.SuiTestnet]
// 	var manageObj entities.Manage
// 	if !c.redisCache.Get(manageObj.GetRedisKey(), &manageObj, ctx) {
// 		res, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
// 			Client:    client,
// 			ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
// 			ErrLogger: c.errLogger,
// 		}, ctx)
// 		if err != nil {
// 			return response.BuildTransactionResponse{}, err
// 		}

// 		if res != nil {
// 			c.redisCache.Set(manageObj.GetRedisKey(), res, time.Minute, ctx)
// 			manageObj = *res
// 		}
// 	}
// 	var sender string = ctx.Value("address").(string)
// 	if !slices.Contains(manageObj.AdminIds, sender) {
// 		return response.BuildTransactionResponse{}, errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
// 	}

// 	editNeedDates, err := on_chain.GetOnChainObject[entities.EditNeedDates](on_chain.GetOnChainObjectRequest{
// 		Client:    client,
// 		ObjectId:  os.Getenv(env.EDIT_MEAL_NEED_DATES_ID),
// 		ErrLogger: c.errLogger,
// 	}, ctx)
// 	if err != nil {
// 		return response.BuildTransactionResponse{}, err
// 	}

// 	var rawStartDate string = strings.TrimSpace(req.StartDate)
// 	if rawStartDate == "" {
// 		rawStartDate = editNeedDates.StartDate
// 	}

// 	var curTime time.Time = time.Now()
// 	var startDate time.Time = util.ToStartOfDate(util.RawDateToTime(fmt.Sprintf("%s/%d", rawStartDate, curTime.Year())))
// 	var originalStartDate time.Time = util.ToStartOfDate(util.RawDateToTime(fmt.Sprintf("%s/%d", editNeedDates.StartDate, curTime.Year())))
// 	var invalidDateErr error = errors.New(noti.INVALID_DATE_MESSAGE)
// 	if startDate.IsZero() || startDate.Before(originalStartDate.AddDate(0, 0, -15)) || startDate.After(originalStartDate.AddDate(0, 0, 15)) { // Changes must be around 15 days from original date
// 		return response.BuildTransactionResponse{}, invalidDateErr
// 	}

// 	var rawEndDate string = strings.TrimSpace(req.StartDate)
// 	if rawEndDate == "" {
// 		rawEndDate = editNeedDates.EndDate
// 	}

// 	if rawStartDate == editNeedDates.StartDate && rawEndDate == editNeedDates.EndDate {
// 		return response.BuildTransactionResponse{}, nil
// 	}

// 	var endDate time.Time = util.ToEndOfDate(util.RawDateToTime(fmt.Sprintf("%s/%d", rawEndDate, curTime.Year())))
// 	if endDate.IsZero() || endDate.Before(startDate) {
// 		return response.BuildTransactionResponse{}, invalidDateErr
// 	}

// 	var needModule = on_chain.InitializeModuleNeed()
// 	txBytes, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
// 		Client:    client,
// 		Sender:    sender,
// 		Module:    needModule.GetModule(),
// 		Function:  needModule.GetFunctionEditUpdateMealNeedDates(),
// 		ErrLogger: c.errLogger,
// 		Arguments: needModule.ToEditUpdateNeedDatesArguments(on_chain.EditUpdateNeedDatesArguments{
// 			EditDatesID: editNeedDates.ID.ID,
// 			StartDate:   rawStartDate,
// 			EndDate:     rawEndDate,
// 		}),
// 	}, ctx)

// 	return response.BuildTransactionResponse{
// 		TxBytes: txBytes,
// 	}, err
// }

// EditChildMaxAgeInSupport implements business.IConfigService.
func (c *configService) EditChildMaxAgeInSupport(ctx context.Context) error {
	panic("unimplemented")
}

// EditSpecialNeedDao implements business.IConfigService.
func (c *configService) EditSpecialNeedDao(req request.EditDaoRequest, ctx context.Context) error {
	var client = c.clients[constant.SuiTestnet]
	manage, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return err
	}

	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)
	if manage == nil {
		return internalErr
	}

	var sender string = ctx.Value("address").(string)
	if !slices.Contains(manage.AdminIds, sender) {
		return errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	}

	dao, err := on_chain.GetOnChainObject[entities.DaoStruct](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.SPECIAL_NEED_DAO_ID),
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return err
	}

	if dao == nil {
		return internalErr
	}

	if req.MinRate == nil && req.MinVoters == nil {
		return nil
	}

	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if req.MinRate != nil {
		if *req.MinRate < 80 {
			return genericErr
		}
	} else {
		minRate, _ := strconv.Atoi(dao.MinApprovedRate)
		req.MinRate = &minRate
	}

	if req.MinVoters != nil {
		if *req.MinVoters < 2 {
			return genericErr
		}
	} else {
		minApprovers, _ := strconv.Atoi(dao.MinVoters)
		req.MinVoters = &minApprovers
	}

	var module = on_chain.InitializeModuleNeed()
	_, errRes := on_chain.ExecuteTransactionV2(on_chain.ExecuteTransactionRequestV2{
		Client:   client,
		Module:   module.GetModule(),
		Function: module.GetFunctionEditSpecialNeedProposalDao(),
		Arguments: module.ToEditSpecialNeedProposalDaoArguments(on_chain.EditSpecialNeedProposalDaoArguments{
			MinVoters: *req.MinRate,
			MinRate:   int64(*req.MinVoters) * 100,
			Sender:    sender,
		}),
		ErrLogger: c.errLogger,
	}, ctx)

	return errRes
}

// UpdateChildEditBooksNeedDates implements business.IConfigService.
func (c *configService) UpdateChildEditBooksNeedDates(req request.UpdateChildEditNeedDatesRequest, ctx context.Context) error {
	var client = c.clients[constant.SuiTestnet]
	manage, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return err
	}

	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)
	if manage == nil {
		return internalErr
	}

	var sender string = ctx.Value("address").(string)
	if !slices.Contains(manage.AdminIds, sender) {
		return errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	}

	editNeedDates, err := on_chain.GetOnChainObject[entities.EditNeedDates](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.EDIT_BOOKS_NEED_DATES_ID),
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return err
	}

	if editNeedDates == nil {
		return internalErr
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
		return invalidDateErr
	}

	var rawEndDate string = strings.TrimSpace(req.StartDate)
	if rawEndDate == "" {
		rawEndDate = editNeedDates.EndDate
	}

	if rawStartDate == editNeedDates.StartDate && rawEndDate == editNeedDates.EndDate {
		return nil
	}

	var endDate time.Time = util.ToEndOfDate(util.RawDateToTime(fmt.Sprintf("%s/%d", rawEndDate, curTime.Year())))
	if endDate.IsZero() || endDate.Before(startDate) {
		return invalidDateErr
	}

	var needModule = on_chain.InitializeModuleNeed()
	_, errRes := on_chain.ExecuteTransactionV2(on_chain.ExecuteTransactionRequestV2{
		Client:    client,
		Module:    needModule.GetModule(),
		Function:  needModule.GetFunctionEditUpdateBooksNeedDates(),
		ErrLogger: c.errLogger,
		Arguments: needModule.ToEditUpdateNeedDatesArguments(on_chain.EditUpdateNeedDatesArguments{
			EditDatesID: editNeedDates.ID.ID,
			StartDate:   rawStartDate,
			EndDate:     rawEndDate,
			Sender:      sender,
		}),
	}, ctx)

	return errRes
}

// UpdateChildEditHealthInsuranceNeedDates implements business.IConfigService.
func (c *configService) UpdateChildEditHealthInsuranceNeedDates(req request.UpdateChildEditNeedDatesRequest, ctx context.Context) error {
	var client = c.clients[constant.SuiTestnet]
	manage, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return err
	}

	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)
	if manage == nil {
		return internalErr
	}

	var sender string = ctx.Value("address").(string)
	if !slices.Contains(manage.AdminIds, sender) {
		return errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	}

	editNeedDates, err := on_chain.GetOnChainObject[entities.EditNeedDates](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.EDIT_HEALTH_INSURANCE_NEED_DATES_ID),
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return err
	}

	if editNeedDates == nil {
		return internalErr
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
		return invalidDateErr
	}

	var rawEndDate string = strings.TrimSpace(req.StartDate)
	if rawEndDate == "" {
		rawEndDate = editNeedDates.EndDate
	}

	if rawStartDate == editNeedDates.StartDate && rawEndDate == editNeedDates.EndDate {
		return nil
	}

	var endDate time.Time = util.ToEndOfDate(util.RawDateToTime(fmt.Sprintf("%s/%d", rawEndDate, curTime.Year())))
	if endDate.IsZero() || endDate.Before(startDate) {
		return invalidDateErr
	}

	var needModule = on_chain.InitializeModuleNeed()
	_, errRes := on_chain.ExecuteTransactionV2(on_chain.ExecuteTransactionRequestV2{
		Client:    client,
		Module:    needModule.GetModule(),
		Function:  needModule.GetFunctionEditUpdateHealthInsuranceNeedDates(),
		ErrLogger: c.errLogger,
		Arguments: needModule.ToEditUpdateNeedDatesArguments(on_chain.EditUpdateNeedDatesArguments{
			EditDatesID: editNeedDates.ID.ID,
			StartDate:   rawStartDate,
			EndDate:     rawEndDate,
			Sender:      sender,
		}),
	}, ctx)

	return errRes
}

// UpdateChildEditMealNeedDates implements business.IConfigService.
func (c *configService) UpdateChildEditMealNeedDates(req request.UpdateChildEditNeedDatesRequest, ctx context.Context) error {
	var client = c.clients[constant.SuiTestnet]
	manage, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return err
	}

	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)
	if manage == nil {
		return internalErr
	}

	var sender string = ctx.Value("address").(string)
	if !slices.Contains(manage.AdminIds, sender) {
		return errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	}

	editNeedDates, err := on_chain.GetOnChainObject[entities.EditNeedDates](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.EDIT_MEAL_NEED_DATES_ID),
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return err
	}

	if editNeedDates == nil {
		return internalErr
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
		return invalidDateErr
	}

	var rawEndDate string = strings.TrimSpace(req.StartDate)
	if rawEndDate == "" {
		rawEndDate = editNeedDates.EndDate
	}

	if rawStartDate == editNeedDates.StartDate && rawEndDate == editNeedDates.EndDate {
		return nil
	}

	var endDate time.Time = util.ToEndOfDate(util.RawDateToTime(fmt.Sprintf("%s/%d", rawEndDate, curTime.Year())))
	if endDate.IsZero() || endDate.Before(startDate) {
		return invalidDateErr
	}

	var needModule = on_chain.InitializeModuleNeed()
	_, errRes := on_chain.ExecuteTransactionV2(on_chain.ExecuteTransactionRequestV2{
		Client:    client,
		Module:    needModule.GetModule(),
		Function:  needModule.GetFunctionEditUpdateMealNeedDates(),
		ErrLogger: c.errLogger,
		Arguments: needModule.ToEditUpdateNeedDatesArguments(on_chain.EditUpdateNeedDatesArguments{
			EditDatesID: editNeedDates.ID.ID,
			StartDate:   rawStartDate,
			EndDate:     rawEndDate,
			Sender:      sender,
		}),
	}, ctx)

	return errRes
}
