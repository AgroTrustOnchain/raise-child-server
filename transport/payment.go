package transport

import (
	"raise-child/business"
	action_type "raise-child/constants/action_type"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/util"

	"github.com/gin-gonic/gin"
)

// GetPayments godoc
// @Summary      List Payments
// @Description  Retrieves a list of Payments based on filter criteria
// @Tags         payment
// @Accept       json
// @Produce      json
// @Param        request  query     request.GetPaymentsRequest  true  "Filter Criteria"
// @Success      200      {object}  response.PaginationDataResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /payments [get]
func GetPayments(ctx *gin.Context) {
	var request request.GetPaymentsRequest
	if ctx.ShouldBindQuery(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GeneratePaymentService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetPayments(request, ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetPayment godoc
// @Summary      Retrieve a payment
// @Description  Retrieve a payment
// @Tags         payment
// @Accept       json
// @Produce      json
// @Param        id  path      string  true  "Payment ID"
// @Success      200      {object}  entities.Payment
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /payments/{id} [get]
func GetPayment(ctx *gin.Context) {
	service, err := business.GeneratePaymentService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetPayment(ctx.Param("id"), ctx)

	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// ApprovePayment godoc
// @Summary      Approve a Payment with method as Manual Bank and upload to Sui Blockchain
// @Description  Prepares and builds a transaction for approve and upload Payment with method as Manual Bank on-chain
// @Tags         payment
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Payment ID"
// @Success      200      {object}  response.MessageAPIResponse "Success"
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /payments/{id}/approve [post]
func ApprovePayment(ctx *gin.Context) {
	service, err := business.GeneratePaymentService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	util.ProcessResponse(response.APIResponse{
		ErrMsg:  service.ApprovePayment(ctx.Param("id"), ctx),
		Context: ctx,
	})
}

// RefusePayment godoc
// @Summary      Refuse Payment
// @Description  Refuse Payment
// @Tags         payment
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Payment ID"
// @Success      200      {object}  response.MessageAPIResponse "Success"
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /payments/{id}/refuse [post]
func RefusePayment(ctx *gin.Context) {
	service, err := business.GeneratePaymentService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	util.ProcessResponse(response.APIResponse{
		ErrMsg:  service.RefusePayment(ctx.Param("id"), ctx),
		Context: ctx,
	})
}

// Donate godoc
// @Summary      Process a donation
// @Description  Initiates a donation process, binds the request JSON, and returns payment or redirect information.
// @Tags         payment
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      request.DonateRequest  true  "Donation Details"
// @Success      200      {object}  response.PaymentUrlResponse
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

	res, err := service.DonateV2(request, ctx)

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
// @Description  Process payment callback
// @Tags         payment
// @Accept       json
// @Produce      json
// @Param        id   path      string                true  "Payment ID"
// @Success      200      {object}  response.MessageAPIResponse "Success"
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /payments/callback/{id} [get]
func CallbackTransaction(ctx *gin.Context) {
	service, err := business.GeneratePaymentService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	util.ProcessResponse(response.APIResponse{
		ErrMsg:  service.CallbackV2(ctx.Param("id"), ctx),
		Context: ctx,
	})
}

// CallbackWithAuthTransaction godoc
// @Summary      Process payment callback with authorization
// @Description  Process payment callback with authorization
// @Tags         payment
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string                true  "Payment ID"
// @Param        request  body      request.PaymentAuthCallbackRequest   true  "Auth Payment Callback Detail"
// @Success      200      {object}  response.MessageAPIResponse "Success"
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /payments/auth-callback/{id} [post]
func CallbackWithAuthTransaction(ctx *gin.Context) {
	var request request.PaymentAuthCallbackRequest
	if ctx.ShouldBindJSON(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GeneratePaymentService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	util.ProcessResponse(response.APIResponse{
		ErrMsg:  service.CallbackWithAuthV2(ctx.Param("id"), request, ctx),
		Context: ctx,
	})
}
