package business

import (
	"context"
	"log"
	"raise-child/constants/shared"
	"raise-child/interfaces/business"
	i_repository "raise-child/interfaces/repository"
	"raise-child/model/entities"
	"raise-child/repository"
	"raise-child/util"
	"raise-child/util/db"
	on_chain "raise-child/util/on_chain"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/sui"
)

type backgroundService struct {
	registrationRequestRepo i_repository.IRegistrationRequestRepository
	centerRequestRepo       i_repository.ICenterRequestRepository
	uploadChildRequestRepo  i_repository.IUploadChildRequestRepository
	clients                 map[string]sui.ISuiAPI
	errLogger               *log.Logger
}

func initializeBackgroundService(
	registrationRequestRepo i_repository.IRegistrationRequestRepository,
	centerRequestRepo i_repository.ICenterRequestRepository,
	uploadChildRequestRepo i_repository.IUploadChildRequestRepository,
	clients map[string]sui.ISuiAPI,
	errLogger *log.Logger,
) business.IBackgroundService {
	return &backgroundService{
		registrationRequestRepo: registrationRequestRepo,
		centerRequestRepo:       centerRequestRepo,
		uploadChildRequestRepo:  uploadChildRequestRepo,
		clients:                 clients,
		errLogger:               errLogger,
	}
}

func GenerateBackgroundService() (business.IBackgroundService, error) {
	var errLogger = util.GetLogConfig(shared.ERROR_LEVEL)

	cnn, err := db.ConnectDB(errLogger, db.InitializePostgreSQL())
	if err != nil {
		return nil, err
	}

	return initializeBackgroundService(
		repository.InitializeRegistrationRequestRepo(cnn, errLogger),
		repository.InitializeCenterRequestRepository(cnn, errLogger),
		repository.InitializeUploadChildRequestRepo(cnn, errLogger),
		_networkAliases,
		errLogger,
	), nil
}

// ProcessBackgroundCenterRequests implements business.IBackgroundService.
func (b *backgroundService) ProcessBackgroundCenterRequests(ctx context.Context) {
	pendingRes, approvedRes, err := b.centerRequestRepo.GetPendingRequests(ctx)
	if err != nil {
		return
	}

	var refusedReqs []entities.BackgroundRecord
	for _, req := range pendingRes {
		var rate float32 = float32(len(req.Approvers)) / float32(len(req.Approvers)+len(req.Refusers))
		if rate >= approve_rate_limit {
			approvedRes = append(approvedRes, req)
		} else {
			refusedReqs = append(refusedReqs, req)
		}
	}

	b.centerRequestRepo.SetRefusedStatuses(refusedReqs, ctx)

	var modules []string
	var functions []string
	var args [][]interface{}
	var module = on_chain.InitializeModuleManage()
	for i, req := range approvedRes {
		modules[i] = module.GetModule()
		functions[i] = module.GetFunctionMintUploadCenterCap()
		args = append(args, module.ToMintCapArguments(on_chain.MintCapArguments{
			Recipient: req.Sender,
		}))
	}

	if err := on_chain.BuildMultiBackgroundTransactions(on_chain.BuildMultiBackgroundTransactionsRequest{
		Client:    b.clients[constant.SuiTestnet],
		Modules:   modules,
		Functions: functions,
		Arguments: args,
		ErrLogger: b.errLogger,
	}, ctx); err != nil {
		return
	}

	b.centerRequestRepo.SetApprovedStatuses(approvedRes, ctx)
}

// ProcessBackgroundRegistrationRequests implements business.IBackgroundService.
func (b *backgroundService) ProcessBackgroundRegistrationRequests(ctx context.Context) {
	pendingRes, approvedRes, err := b.registrationRequestRepo.GetPendingRequests(ctx)
	if err != nil {
		return
	}

	var refusedReqs []entities.BackgroundRecord
	for _, req := range pendingRes {
		var rate float32 = float32(len(req.Approvers)) / float32(len(req.Approvers)+len(req.Refusers))
		if rate >= approve_rate_limit {
			approvedRes = append(approvedRes, req)
		} else {
			refusedReqs = append(refusedReqs, req)
		}
	}

	b.registrationRequestRepo.SetRefusedStatuses(refusedReqs, ctx)

	var modules []string
	var functions []string
	var args [][]interface{}
	var module = on_chain.InitializeModuleManage()
	for i, req := range approvedRes {
		modules[i] = module.GetModule()
		switch req.Role {
		case admin_role:
			functions[i] = module.GetFunctionMintRegisterAdminCap()
		case local_leader_role:
			functions[i] = module.GetFunctionMintRegisterLeaderCap()
		case volunteer_role:
			functions[i] = module.GetFunctionMintRegisterVolunteerCap()
		}
		args = append(args, module.ToMintCapArguments(on_chain.MintCapArguments{
			Recipient: req.Sender,
		}))
	}

	if err := on_chain.BuildMultiBackgroundTransactions(on_chain.BuildMultiBackgroundTransactionsRequest{
		Client:    b.clients[constant.SuiTestnet],
		Modules:   modules,
		Functions: functions,
		Arguments: args,
		ErrLogger: b.errLogger,
	}, ctx); err != nil {
		return
	}

	b.registrationRequestRepo.SetApprovedStatuses(approvedRes, ctx)
}

// ProcessBackgroundUploadChildRequests implements business.IBackgroundService.
func (b *backgroundService) ProcessBackgroundUploadChildRequests(ctx context.Context) {
	// pendingRes, approvedRes, err := b.uploadChildRequestRepo.GetPendingRequests(ctx)
	// if err != nil {
	// 	return
	// }

	// var refusedReqs []entities.BackgroundRecord
	// for _, req := range pendingRes {
	// 	var rate float32 = float32(len(req.Approvers)) / float32(len(req.Approvers)+len(req.Refusers))
	// 	if rate >= approve_rate_limit {
	// 		approvedRes = append(approvedRes, req)
	// 	} else {
	// 		refusedReqs = append(refusedReqs, req)
	// 	}
	// }

	// b.uploadChildRequestRepo.SetRefusedStatuses(refusedReqs, ctx)

	// var modules []string
	// var functions []string
	// var args [][]interface{}
	// var module = on_chain.InitializeModuleManage()
	// for i, req := range approvedRes {
	// 	modules[i] = module.GetModule()
	// 	functions[i] = module.GetFunctionMintUploadCenterCap()
	// 	args = append(args, module.ToMintCapArguments(on_chain.MintCapArguments{
	// 		Recipient: req.Sender,
	// 	}))
	// }

	// if err := on_chain.BuildMultiBackgroundTransactions(on_chain.BuildMultiBackgroundTransactionsRequest{
	// 	Client:    b.clients[constant.SuiTestnet],
	// 	Modules:   modules,
	// 	Functions: functions,
	// 	Arguments: args,
	// 	ErrLogger: b.errLogger,
	// }, ctx); err != nil {
	// 	return
	// }

	// b.uploadChildRequestRepo.SetApprovedStatuses(approvedRes, ctx)
}
