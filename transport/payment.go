package transport

import (
	"raise-child/business"
	action_type "raise-child/constants/action_type"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/util"

	"github.com/gin-gonic/gin"
)

// Donate godoc
// @Summary      Process a donation
// @Description  Initiates a donation process, binds the request JSON, and returns payment or redirect information.
// @Tags         payment
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      request.DonateRequest  true  "Donation Details"
// @Success      200      {object}  response.UrlAPIResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /payments/donate [post]
func Donate(ctx *gin.Context) {
	var request request.DonateRequest
	if ctx.ShouldBindJSON(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GeneratePaymentService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.Donate(request, ctx)

	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// CallbackTransaction godoc
// @Summary      Process payment callback
// @Description  Handles the transaction callback/webhook from the payment provider to update transaction status by ID and build on-chain transaction.
// @Tags         payment
// @Accept       json
// @Produce      json
// @Param        id   path      string                true  "Payment ID"
// @Success      200  {object}  response.BuildTransactionResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /payments/callback/{id} [get]
func CallbackTransaction(ctx *gin.Context) {
	service, err := business.GeneratePaymentService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.Callback(ctx.Param("id"), ctx)

	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.REDIRECT,
	})
}

// CallbackWithAuthTransaction godoc
// @Summary      Process payment callback with authorization
// @Description  Handles the transaction callback/webhook from the payment provider to update transaction status by ID and build on-chain transaction.
// @Tags         payment
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string                true  "Payment ID"
// @Param        imageBlobId   query      string                true  "Captured Image Blob ID"
// @Success      200  {object}  response.BuildTransactionResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /payments/auth-callback/{id} [get]
func CallbackWithAuthTransaction(ctx *gin.Context) {
	service, err := business.GeneratePaymentService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.CallbackWithAuth(ctx.Param("id"), ctx.Query("imageBlobId"), ctx)

	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.REDIRECT,
	})
}
