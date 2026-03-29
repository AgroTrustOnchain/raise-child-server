package transport

import (
	"raise-child/business"
	action_type "raise-child/constants/action_type"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/util"

	"github.com/gin-gonic/gin"
)

// GetTransactionRecords godoc
// @Summary      Get a list of transaction records
// @Description  Get a list of transaction records
// @Tags         tx-records
// @Accept       json
// @Produce      json
// @Param        request  body      request.GetTransactionRecordsRequest true  "Filter criteria"
// @Success      200      {object}  response.PaginationDataResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /tx-records [get]
func GetTransactionRecords(ctx *gin.Context) {
	var request request.GetTransactionRecordsRequest
	if ctx.ShouldBindQuery(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateTransactionRecordService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetTransactionRecords(request, ctx)

	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetTransactionRecord godoc
// @Summary      Get a transaction record specific information
// @Description  Get a transaction record specific information by its unique ID
// @Tags         tx-records
// @Accept       json
// @Produce      json
// @Param        id       path      string                      true  "Transaction Record ID"
// @Success      200      {object}  response.TransactionResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /tx-records/{id} [get]
func GetTransactionRecord(ctx *gin.Context) {
	service, err := business.GenerateTransactionRecordService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetTransactionRecord(ctx.Param("id"), ctx)

	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}
