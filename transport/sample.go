package transport

import (
	"log"
	actiontype "raise-child/constants/action_type"
	"raise-child/model/dtos/response"
	"raise-child/util"
	on_chain "raise-child/util/on_chain"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/gin-gonic/gin"
)

var client = sui.NewSuiClient(constant.BvTestnetEndpoint)
var address string = ""

func stringToByteArray(s string) []interface{} {
	bytes := []byte(s)
	result := make([]interface{}, len(bytes))
	for i, b := range bytes {
		result[i] = int(b) // Convert to int, not uint8
	}
	return result
}

func BuildCampaignTransaction(ctx *gin.Context) {
	//tx := transaction.NewTransaction()
	var text = "Test Build Campaign on Go"
	// nameEncoded := base64.StdEncoding.EncodeToString([]byte(text))
	// descEncoded := base64.StdEncoding.EncodeToString([]byte(text))
	txtBytes, _ := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
		Client:   client,
		Sender:   address,
		Module:   "fundraising",
		Function: "create_campaign",
		Arguments: []interface{}{
			text,
			text,
		},
		ErrLogger: log.Default(),
	}, ctx)

	type BuildTxResponse struct {
		TxBytes string `json:"tx_bytes"`
	}

	util.ProcessResponse(response.APIResponse{
		Data1:    BuildTxResponse{TxBytes: txtBytes},
		Data2:    BuildTxResponse{TxBytes: txtBytes},
		Context:  ctx,
		PostType: actiontype.NON_POST,
	})
}

// func ExecuteCampaignTransaction(ctx *gin.Context) {
// 	var request request.ExecuteTransactionRequest
// 	ctx.ShouldBindJSON(&request)

// 	log.Println("TxBytes: ", request.TxBytes)
// 	log.Println("Signature:", request.Signature)

// 	err := on_chain.ExecuteTransaction(on_chain.ExecuteTransactionRequest{
// 		Client:    client,
// 		TxBytes:   request.TxBytes,
// 		Signature: []string{request.Signature},
// 		ErrLogger: log.Default(),
// 	}, ctx)

// 	if err != nil {
// 		log.Println("Transaction result:" + err.Error())
// 	} else {
// 		log.Println("Transaction executed successfully")
// 	}

// }
