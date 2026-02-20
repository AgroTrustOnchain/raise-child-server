package business

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"raise-child/interfaces/business"
	i_repository "raise-child/interfaces/repository"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/repository"
	"raise-child/util"
	"raise-child/util/db"
	on_chain "raise-child/util/on_chain"
	"time"

	"raise-child/constants/env"
	"raise-child/constants/noti"
	internal_sui "raise-child/constants/on-chain/sui"
	"raise-child/constants/shared"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/sui"
)

type onChainService struct {
	withdrawRepo      i_repository.IOffChainWithdrawProposalRepository
	centerRepo        i_repository.ICenterRequestRepository
	uploadChildRepo   i_repository.IUploadChildRequestRepository
	registrationRepo  i_repository.IRegistrationRequestRepository
	volunteerNotiRepo i_repository.IVolunteerNotiRepository
	leaderNotiRepo    i_repository.ILeaderNotiRepository
	clients           map[string]sui.ISuiAPI
	errLogger         *log.Logger
}

// GetCurrentWalletNotis implements business.INotiService.
func (o *onChainService) GetCurrentWalletNotis(wallet string, req request.GetNotisRequest, ctx context.Context) (response.PaginationDataResponse, error) {
	panic("unimplemented")
}

// Money actions
const (
	withdraw_action string = "withdraw"
	donate_action   string = "action"
)

func InitializeOnChainService(db *sql.DB, errLogger *log.Logger) business.IOnChainService {
	return &onChainService{
		withdrawRepo:      repository.InitializeOffChainWithdrawProposalRepository(db, errLogger),
		centerRepo:        repository.InitializeCenterRequestRepository(db, errLogger),
		uploadChildRepo:   repository.InitializeUploadChildRequestRepo(db, errLogger),
		registrationRepo:  repository.InitializeRegistrationRequestRepo(db, errLogger),
		volunteerNotiRepo: repository.InitializeVolunteerNotiRepository(db, errLogger),
		leaderNotiRepo:    repository.InitializeLeaderNotiRepository(db, errLogger),
		clients:           _networkAliases,
		errLogger:         errLogger,
	}
}

func GenerateOnChainService() (business.IOnChainService, error) {
	var errLogger = util.GetLogConfig(shared.ERROR_LEVEL)

	cnn, err := db.ConnectDB(errLogger, db.InitializePostgreSQL())
	if err != nil {
		return nil, err
	}

	return InitializeOnChainService(cnn, errLogger), nil
}

// ExecuteTransaction implements business.IOnChainService.
func (o *onChainService) ExecuteTransaction(req request.ExecuteTransactionRequest, ctx context.Context) error {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	var curTime time.Time = time.Now()
	if req.ProposalID != "" {
		proposal, err := o.withdrawRepo.GetOffChainWithdrawProposal(req.ProposalID, ctx)
		if err != nil {
			return err
		}

		if proposal.ProposalID != "" {
			return genericErr
		}
	}

	res, err := on_chain.ExecuteTransaction(on_chain.ExecuteTransactionRequest{
		Client:    o.clients[constant.SuiTestnet],
		TxBytes:   req.TxBytes,
		Signature: []string{req.Signature},
		ErrLogger: o.errLogger,
	}, ctx)
	if err != nil {
		return err
	}

	if req.ProposalID != "" {
		var events = res.Events
		if events == nil || len(events) == 0 {
			return genericErr
		}

		var module = on_chain.InitializeModulePool()
		var eventType string = fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), module.GetModule(), module.GetWithdrawProposalEventEmittedStruct())
		for _, event := range events {
			if event.Type == eventType {
				if onChainProposal, ok := event.ParsedJson["id"].(string); ok {
					o.withdrawRepo.SetOnChainProposalIdAfterExecuteTx(req.ProposalID, onChainProposal, ctx)
					break
				}
			}
		}
	} else if req.CenterReq != "" {
		req, err := o.centerRepo.GetRequest(req.CenterReq, ctx)
		if err != nil {
			return err
		}

		if req.IsConfirmRegister {
			return genericErr
		}

		req.IsConfirmRegister = true
		req.UpdatedAt = curTime

		o.centerRepo.UpdateRegistrationRequest(*req, ctx)
	} else if req.UploadChildReq != "" {
		req, err := o.uploadChildRepo.GetUploadChildRequest(req.UploadChildReq, ctx)
		if err != nil {
			return err
		}

		if req.IsConfirmUpload {
			return genericErr
		}

		req.IsConfirmUpload = true
		req.UpdatedAt = curTime

		o.uploadChildRepo.UpdateUploadChildRequest(*req, ctx)
	} else if req.RegistraionReq != "" {
		req, err := o.registrationRepo.GetRegistrationRequest(req.RegistraionReq, ctx)
		if err != nil {
			return err
		}

		if req.IsConfirmRegister {
			return genericErr
		}

		req.IsConfirmRegister = true
		req.UpdatedAt = curTime

		o.registrationRepo.UpdateRegistrationRequest(*req, ctx)

		if req.RegisterRole == volunteer_role {
			o.volunteerNotiRepo.AssignVolunteer(req.CreatedBy, req.Region, ctx)
		} else if req.RegisterRole == local_leader_role {
			o.leaderNotiRepo.AssignLeader(req.CreatedBy, req.Region, ctx)
		}
	}

	return nil
}

// BuildMoneyTransaction implements business.IOnChainService.
func (o *onChainService) BuildMoneyTransaction(req request.MoneyActionRequest, ctx context.Context) (response.BuildTransactionResponse, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)

	// Not logged in
	if info := getSecurityInfo(req.Sender); !info.isLoggedIn {
		return response.BuildTransactionResponse{}, genericErr
	}

	// Invalid action
	if req.ActionType != withdraw_action && req.ActionType != donate_action {
		return response.BuildTransactionResponse{}, genericErr
	}

	if req.CoinType == "" {
		req.CoinType = internal_sui.SUI_COIN_TYPE
	}

	if req.Message == "" {
		req.Message = req.ActionType
	}

	if req.ActionType == donate_action {
		txBytes, err := on_chain.BuildDonateTransaction(on_chain.DonateTransactionRequest{
			Client:    o.clients[constant.SuiTestnet],
			Sender:    req.Sender,
			CoinType:  req.CoinType,
			Amount:    on_chain.StandarizeToSuiMist(req.Amount),
			Message:   req.Message,
			ErrLogger: o.errLogger,
		}, ctx)

		return response.BuildTransactionResponse{
			TxBytes: txBytes,
		}, err
	}

	// // Not admin
	// if req.Sender != os.Getenv(env.ADMIN) {
	// 	return response.BuildTransactionResponse{}, errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	// }

	var module = on_chain.InitializeModuleManage()
	txBytes, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
		Client:   o.clients[constant.SuiTestnet],
		Sender:   req.Sender,
		Module:   module.GetModule(),
		Function: module.GetFunctionWithdrawSuiPool(),
		Arguments: []interface{}{
			os.Getenv(env.POOL_ID),
			on_chain.StandarizeToSuiMist(req.Amount),
			internal_sui.CLOCK_OBJECT_ID,
			req.Message,
		},
		ErrLogger: o.errLogger,
	}, ctx)

	return response.BuildTransactionResponse{
		TxBytes: txBytes,
	}, err
}
