package onchain

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"raise-child/constants/env"
	"raise-child/constants/noti"
	"strconv"

	internal_sui "raise-child/constants/on-chain/sui"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/signer"
	"github.com/block-vision/sui-go-sdk/sui"
)

type BuildTransactionRequest struct {
	Client    sui.ISuiAPI
	Sender    string
	Module    string
	Function  string
	ErrLogger *log.Logger
	Arguments []interface{}
}

type BuildMultiTransactionsRequest struct {
	Sender string
	BuildMultiBackgroundTransactionsRequest
}

type BuildMultiBackgroundTransactionsRequest struct {
	Client    sui.ISuiAPI
	Modules   []string
	Functions []string
	Arguments [][]interface{}
	ErrLogger *log.Logger
}

type SplitCoinRequest struct {
	Client    sui.ISuiAPI
	Sender    string
	CoinType  string
	Amount    int64
	ErrLogger *log.Logger
}

type ExecuteTransactionRequest struct {
	Client    sui.ISuiAPI
	TxBytes   string
	Signature []string
	ErrLogger *log.Logger
}

type DonateTransactionRequest struct {
	Client    sui.ISuiAPI
	Sender    string
	CoinType  string
	Amount    int64
	Message   string
	ErrLogger *log.Logger
}

type ExecuteTransactionRequestV2 struct {
	Client    sui.ISuiAPI
	Module    string
	Function  string
	Arguments []interface{}
	ErrLogger *log.Logger
}

const (
	defaultGasBudget   int    = 100_000_000
	defaultRequestType string = "WaitForLocalExecution"
)

func BuildTransaction(req BuildTransactionRequest, ctx context.Context) (string, error) {
	res, err := req.Client.MoveCall(ctx, models.MoveCallRequest{
		Signer:          req.Sender,
		PackageObjectId: os.Getenv(env.PACKAGE_ID),
		Module:          req.Module,
		Function:        req.Function,
		TypeArguments:   []interface{}{},
		Arguments:       req.Arguments,
		GasBudget:       fmt.Sprint(defaultGasBudget),
	})

	if err != nil {
		req.ErrLogger.Println(noti.BUILDING_TX_ERR_MSG + err.Error())
		return "", errors.New(noti.INTERNALL_ERR_MSG)
	}

	return res.TxBytes, nil
}

func BuildMultiBackgroundTransactions(req BuildMultiBackgroundTransactionsRequest, ctx context.Context) error {
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)
	signer, err := signer.NewSignerWithSecretKey(os.Getenv(env.PUBLISHER_PRIVATE_KEY))
	if err != nil {
		req.ErrLogger.Println(err.Error())
		return internalErr
	}

	var packageId string = os.Getenv(env.PACKAGE_ID)
	var rawGasBudget string = fmt.Sprint(defaultGasBudget)
	var sender string = os.Getenv(env.PUBLISHER_ADDRESS)
	var txRequestParams []models.RPCTransactionRequestParams

	for i := 0; i < len(req.Functions); i++ {
		txRequestParams = append(txRequestParams, models.RPCTransactionRequestParams{
			MoveCallRequestParams: &models.MoveCallRequest{
				Signer:          sender,
				PackageObjectId: packageId,
				Module:          req.Modules[i],
				Function:        req.Functions[i],
				TypeArguments:   []interface{}{},
				Arguments:       req.Arguments[i],
				GasBudget:       rawGasBudget,
			},
		})
	}

	res, err := req.Client.BatchTransaction(ctx, models.BatchTransactionRequest{
		Signer:                         sender,
		RPCTransactionRequestParams:    txRequestParams,
		GasBudget:                      "20000000",
		SuiTransactionBlockBuilderMode: "Commit",
	})
	if err != nil {
		req.ErrLogger.Println(err.Error())
		return internalErr
	}

	if _, err := req.Client.SignAndExecuteTransactionBlock(ctx, models.SignAndExecuteTransactionBlockRequest{
		TxnMetaData: models.TxnMetaData{
			Gas:          res.Gas,
			InputObjects: res.InputObjects,
			TxBytes:      res.TxBytes,
		},
		PriKey:      signer.PriKey,
		RequestType: defaultRequestType,
		Options: models.SuiTransactionBlockOptions{
			ShowEffects: true,
		},
	}); err != nil {
		req.ErrLogger.Println(err.Error())
		return internalErr
	}

	return nil
}

func BuildMultiTransactions(req BuildMultiTransactionsRequest, ctx context.Context) (string, error) {
	var packageId string = os.Getenv(env.PACKAGE_ID)
	var rawGasBudget string = fmt.Sprint(defaultGasBudget)
	var txRequestParams []models.RPCTransactionRequestParams

	for i := 0; i < len(req.Functions); i++ {
		txRequestParams = append(txRequestParams, models.RPCTransactionRequestParams{
			MoveCallRequestParams: &models.MoveCallRequest{
				Signer:          req.Sender,
				PackageObjectId: packageId,
				Module:          req.Modules[i],
				Function:        req.Functions[i],
				TypeArguments:   []interface{}{},
				Arguments:       req.Arguments[i],
				GasBudget:       rawGasBudget,
			},
		})
	}

	res, err := req.Client.BatchTransaction(ctx, models.BatchTransactionRequest{
		Signer:                         req.Sender,
		RPCTransactionRequestParams:    txRequestParams,
		GasBudget:                      rawGasBudget,
		SuiTransactionBlockBuilderMode: "Commit",
	})

	if err != nil {
		req.ErrLogger.Println(noti.BATCHING_TX_ERR_MSG + err.Error())
		return "", errors.New(noti.INTERNALL_ERR_MSG)
	}

	return res.TxBytes, nil
}

func BuildDonateTransaction(req DonateTransactionRequest, ctx context.Context) (string, error) {
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	coins, err := getOwnedCoinsMatchedAmount(req.Client, req.Sender, req.CoinType, req.Amount, req.ErrLogger, ctx)
	if err != nil {
		req.ErrLogger.Println(err.Error())
		return "", internalErr
	}

	// Not enough balance
	if len(coins) == 0 {
		return "", errors.New(noti.NOT_ENOUGH_BALANCE_MESSAGE)
	}

	var targetCoin models.CoinData
	for _, coin := range coins {
		coinBalance, _ := strconv.ParseInt(coin.Balance, 10, 64)
		if coinBalance == req.Amount {
			targetCoin = coin
			break
		}
	}

	// No coin matched exactly requested amount
	if targetCoin.CoinObjectId == "" {
		res, err := splitCoin(req.Client, req.Sender, coins[0].CoinObjectId, req.Amount, req.ErrLogger, ctx)
		return res, err
	}

	var module = InitializeModuleManage()
	return BuildTransaction(BuildTransactionRequest{
		Client:    req.Client,
		Sender:    req.Sender,
		Module:    module.GetModule(),
		Function:  module.GetFunctionDonateSuiPool(),
		ErrLogger: req.ErrLogger,
		Arguments: []interface{}{
			os.Getenv(env.POOL_ID),
			targetCoin.CoinObjectId,
			internal_sui.CLOCK_OBJECT_ID,
			req.Message,
		},
	}, ctx)
}

func ExecuteTransaction(req ExecuteTransactionRequest, ctx context.Context) (models.SuiTransactionBlockResponse, error) {
	//req.Client.SignAndExecuteTransactionBlock(ctx, models.SignAndExecuteTransactionBlockRequest{})
	res, err := req.Client.SuiExecuteTransactionBlock(ctx, models.SuiExecuteTransactionBlockRequest{
		TxBytes:     req.TxBytes,
		Signature:   req.Signature,
		RequestType: defaultRequestType,
		Options: models.SuiTransactionBlockOptions{
			ShowEffects: true,
			ShowEvents:  true,
		},
	})

	if err != nil {
		req.ErrLogger.Println(noti.EXECUTING_TX_ERR_MSG + err.Error())
		err = errors.New(noti.INTERNALL_ERR_MSG)
	}

	return res, err
}

func ExecuteTransactionV2(req ExecuteTransactionRequestV2, ctx context.Context) (models.SuiTransactionBlockResponse, error) {
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)
	signer, err := signer.NewSignerWithSecretKey(os.Getenv(env.PUBLISHER_PRIVATE_KEY))
	if err != nil {
		req.ErrLogger.Println(noti.EXECUTING_TX_ERR_MSG + err.Error())
		return models.SuiTransactionBlockResponse{}, internalErr
	}

	txnData, err := req.Client.MoveCall(ctx, models.MoveCallRequest{
		Signer:          os.Getenv(env.PUBLISHER_ADDRESS),
		PackageObjectId: os.Getenv(env.PACKAGE_ID),
		Module:          req.Module,
		Function:        req.Function,
		TypeArguments:   []interface{}{},
		Arguments:       req.Arguments,
		GasBudget:       fmt.Sprintf("%d", defaultGasBudget),
	})
	if err != nil {
		req.ErrLogger.Println(noti.EXECUTING_TX_ERR_MSG + err.Error())
		return models.SuiTransactionBlockResponse{}, internalErr
	}

	res, err := req.Client.SignAndExecuteTransactionBlock(ctx, models.SignAndExecuteTransactionBlockRequest{
		TxnMetaData: txnData,
		PriKey:      signer.PriKey,
		RequestType: defaultRequestType,
		Options: models.SuiTransactionBlockOptions{
			ShowEffects: true,
			ShowEvents:  true,
		},
	})
	if err != nil {
		req.ErrLogger.Println(noti.EXECUTING_TX_ERR_MSG + err.Error())
		return models.SuiTransactionBlockResponse{}, internalErr
	}

	return res, nil
}
