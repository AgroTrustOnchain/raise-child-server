package business

import (
	"context"
	"errors"
	"log"
	"os"
	"raise-child/interfaces/business"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/util"
	on_chain "raise-child/util/on_chain"

	"raise-child/constants/env"
	"raise-child/constants/noti"
	internal_sui "raise-child/constants/on-chain/sui"
	"raise-child/constants/shared"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/sui"
)

type onChainService struct {
	clients   map[string]sui.ISuiAPI
	errLogger *log.Logger
}

// Money actions
const (
	withdraw_action string = "withdraw"
	donate_action   string = "action"
)

func InitializeOnChainService(clients map[string]sui.ISuiAPI, errLogger *log.Logger) business.IOnChainService {
	return &onChainService{
		clients:   clients,
		errLogger: errLogger,
	}
}

func GenerateOnChainService() (business.IOnChainService, error) {
	return InitializeOnChainService(_networkAliases, util.GetLogConfig(shared.ERROR_LEVEL)), nil
}

// ExecuteTransaction implements business.IOnChainService.
func (o *onChainService) ExecuteTransaction(req request.ExecuteTransactionRequest, ctx context.Context) error {
	_, err := on_chain.ExecuteTransaction(on_chain.ExecuteTransactionRequest{
		Client:    o.clients[constant.SuiTestnet],
		TxBytes:   req.TxBytes,
		Signature: []string{req.Signature},
		ErrLogger: o.errLogger,
	}, ctx)

	// todo: save transction record to db

	return err
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

	// Not admin
	if req.Sender != os.Getenv(env.ADMIN) {
		return response.BuildTransactionResponse{}, errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	}

	var module = on_chain.InitializeModuleManage()
	txBytes, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
		Client:   o.clients[constant.SuiTestnet],
		Sender:   req.Sender,
		Module:   module.GetModule(),
		Function: module.GetFunctionWithdrawSuiPool(),
		Arguments: []interface{}{
			os.Getenv(env.SUI_POOL_ID),
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
