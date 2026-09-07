package onchain

import (
	"context"
	"errors"
	"log"
	"raise-child/constants/noti"
	internal_sui "raise-child/constants/on-chain/sui"
	"strconv"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/block-vision/sui-go-sdk/utils"
)

func FaucetTestnetBalance(client sui.ISuiAPI, address string, errLogger *log.Logger, ctx context.Context) error {
	// Invalid sui wallet
	if !utils.IsValidSuiAddress(models.SuiAddress(address)) {
		return errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)
	res, err := client.SuiXGetBalance(ctx, models.SuiXGetBalanceRequest{
		Owner:    address,
		CoinType: internal_sui.SUI_COIN_TYPE,
	})

	if err != nil {
		errLogger.Println(noti.GET_BALANCE_ERR_MSG + err.Error())
		return internalErr
	}

	mistBalance, _ := strconv.ParseFloat(res.TotalBalance, 64)
	if mistBalance > 0 { // Enough balance
		return nil
	}

	faucetHost, err := sui.GetFaucetHost(constant.SuiTestnet)
	if err != nil {
		errLogger.Println(noti.GET_FAUCET_HOST_ERR_MSG + err.Error())
		return internalErr
	}

	if err := sui.RequestSuiFromFaucet(faucetHost, address, map[string]string{}); err != nil {
		errLogger.Println(noti.FAUCET_ERR_MSG + err.Error())
		return internalErr
	}

	return nil
}
