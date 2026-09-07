package transport

import (
	"log"
	actiontype "raise-child/constants/action_type"
	"raise-child/model/dtos/response"
	"raise-child/util"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/signer"
	"github.com/gin-gonic/gin"
)

var packageMintID string = ""
var capID string = ""
var privateKey string = ""
var addresses []string = []string{
	address, // Original
	"",
	"",
}

func MintCap(ctx *gin.Context) {
	signer, err := signer.NewSignerWithSecretKey(privateKey)
	if err != nil {
		log.Println("Error create signer: " + err.Error())
		util.ProcessResponse(response.APIResponse{
			ErrMsg:   err,
			Context:  ctx,
			PostType: actiontype.NON_POST,
		})
		return
	}

	log.Println("Signer address:" + string(signer.Address))

	txnData, err := client.MoveCall(ctx, models.MoveCallRequest{
		Signer:          addresses[1],
		PackageObjectId: packageMintID,
		Module:          "mint",
		Function:        "mint_cap",
		TypeArguments:   []interface{}{},
		Arguments: []interface{}{
			capID,
			"0x917c2fdd5d5703f98ee74a4dc3a8a1c974254483fa7ca7309e1a340eb06b09e6",
		},
		GasBudget: "20000000",
	})
	if err != nil {
		log.Println("Error create txn meta data: " + err.Error())
		util.ProcessResponse(response.APIResponse{
			ErrMsg:   err,
			Context:  ctx,
			PostType: actiontype.NON_POST,
		})
		return
	}

	_, errRess := client.SignAndExecuteTransactionBlock(ctx, models.SignAndExecuteTransactionBlockRequest{
		TxnMetaData: txnData,
		PriKey:      signer.PriKey,
		RequestType: "WaitForLocalExecution",
		Options: models.SuiTransactionBlockOptions{
			ShowEffects: true,
		},
	})

	if errRess != nil {
		log.Println("Error sign and execute tx: " + errRess.Error())
		util.ProcessResponse(response.APIResponse{
			ErrMsg:   errRess,
			Context:  ctx,
			PostType: actiontype.NON_POST,
		})
		return
	}

	util.ProcessResponse(response.APIResponse{
		Data1:    response.MessageAPIResponse{Message: "Mint Cap Success"},
		Context:  ctx,
		PostType: actiontype.NON_POST,
	})
}

func MintCaps(ctx *gin.Context) {
	signer, err := signer.NewSignerWithSecretKey(privateKey)
	if err != nil {
		log.Println("Error create signer: " + err.Error())
		util.ProcessResponse(response.APIResponse{
			ErrMsg:   err,
			Context:  ctx,
			PostType: actiontype.NON_POST,
		})
		return
	}

	var txRequestParams []models.RPCTransactionRequestParams
	for i := 1; i < 3; i++ {
		txRequestParams = append(txRequestParams, models.RPCTransactionRequestParams{
			MoveCallRequestParams: &models.MoveCallRequest{
				Signer:          address,
				PackageObjectId: packageMintID,
				Module:          "mint",
				Function:        "mint_cap",
				TypeArguments:   []interface{}{},
				Arguments: []interface{}{
					capID,
					addresses[i],
				},
				GasBudget: "20000000",
			},
		})
	}

	res, err := client.BatchTransaction(ctx, models.BatchTransactionRequest{
		Signer:                         address,
		RPCTransactionRequestParams:    txRequestParams,
		GasBudget:                      "20000000",
		SuiTransactionBlockBuilderMode: "Commit",
	})
	if err != nil {
		log.Println("Error batch tx: " + err.Error())
		util.ProcessResponse(response.APIResponse{
			ErrMsg:   err,
			Context:  ctx,
			PostType: actiontype.NON_POST,
		})
		return
	}

	_, errRess := client.SignAndExecuteTransactionBlock(ctx, models.SignAndExecuteTransactionBlockRequest{
		TxnMetaData: models.TxnMetaData{
			Gas:          res.Gas,
			InputObjects: res.InputObjects,
			TxBytes:      res.TxBytes,
		},
		PriKey:      signer.PriKey,
		RequestType: "WaitForLocalExecution",
		Options: models.SuiTransactionBlockOptions{
			ShowEffects: true,
		},
	})

	if errRess != nil {
		log.Println("Error sign and execute tx: " + errRess.Error())
		util.ProcessResponse(response.APIResponse{
			ErrMsg:   errRess,
			Context:  ctx,
			PostType: actiontype.NON_POST,
		})
		return
	}

	util.ProcessResponse(response.APIResponse{
		Data1:    response.MessageAPIResponse{Message: "Mint Cap Success"},
		Context:  ctx,
		PostType: actiontype.NON_POST,
	})
}
