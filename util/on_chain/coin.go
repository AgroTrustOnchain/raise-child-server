package onchain

import (
	"context"
	"errors"
	"fmt"
	"log"
	"raise-child/constants/noti"

	"strconv"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
)

const default_mist int64 = 1_000_000_000

func StandarizeToSuiMist(amount int) int64 {
	return int64(amount) * default_mist
}

func getOwnedCoinsMatchedAmount(client sui.ISuiAPI, owner, coinType string, amount int64, errLogger *log.Logger, ctx context.Context) ([]models.CoinData, error) {
	coins, err := client.SuiXGetCoins(ctx, models.SuiXGetCoinsRequest{
		Owner:    owner,
		CoinType: coinType,
	})

	if err != nil {
		errLogger.Println(noti.RETRIEVE_OWNED_COINS_ERR_MSG + err.Error())
		return nil, errors.New(noti.INTERNALL_ERR_MSG)
	}

	var res []models.CoinData
	for _, coin := range coins.Data {
		coinBalance, _ := strconv.ParseInt(coin.Balance, 10, 64)
		if coinBalance >= amount {
			res = append(res, coin)
		}
	}

	return res, nil
}

func splitCoin(client sui.ISuiAPI, sender, originalCoinId string, amount int64, errLogger *log.Logger, ctx context.Context) (string, error) {
	res, err := client.SplitCoin(ctx, models.SplitCoinRequest{
		Signer:       sender,
		CoinObjectId: originalCoinId,
		SplitAmounts: []string{fmt.Sprint(amount)},
		GasBudget:    fmt.Sprint(defaultGasBudget),
	})

	if err != nil {
		errLogger.Println(noti.SPLIT_COIN_ERR_MSG + err.Error())
		return "", errors.New(noti.INTERNALL_ERR_MSG)
	}

	return res.TxBytes, err
}
