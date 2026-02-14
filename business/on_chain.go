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

	"raise-child/constants/env"
	"raise-child/constants/noti"
	internal_sui "raise-child/constants/on-chain/sui"
	"raise-child/constants/shared"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/sui"
)

type onChainService struct {
	withdrawRepo i_repository.IOffChainWithdrawProposalRepository
	clients      map[string]sui.ISuiAPI
	errLogger    *log.Logger
}

// Money actions
const (
	withdraw_action string = "withdraw"
	donate_action   string = "action"
)

func InitializeOnChainService(db *sql.DB, errLogger *log.Logger) business.IOnChainService {
	return &onChainService{
		withdrawRepo: repository.InitializeOffChainWithdrawProposalRepository(db, errLogger),
		clients:      _networkAliases,
		errLogger:    errLogger,
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
	if req.Proposal != "" {
		proposal, err := o.withdrawRepo.GetOffChainWithdrawProposal(req.Proposal, ctx)
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

	if req.Proposal != "" {
		var events = res.Events
		if events == nil || len(events) == 0 {
			return genericErr
		}

		var module = on_chain.InitializeModulePool()
		var eventType string = fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), module.GetModule(), module.GetWithdrawProposalEventEmittedStruct())
		for _, event := range events {
			if event.Type == eventType {
				if onChainProposal, ok := event.ParsedJson["id"].(string); ok {
					return o.withdrawRepo.SetOnChainProposalIdAfterExecuteTx(req.Proposal, onChainProposal, ctx)
				}
			}
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
