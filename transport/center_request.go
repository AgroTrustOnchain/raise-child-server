package transport

import (
	"raise-child/business"
	action_type "raise-child/constants/action_type"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/util"

	"github.com/gin-gonic/gin"
)

// GetCenterRequests godoc
// @Summary      Get a list of center requests
// @Description  Retrieve center requests with filtering based on query parameters
// @Tags         center
// @Accept       json
// @Produce      json
// @Param        request  query     request.GetCenterRequests  true  "Filter Criteria"
// @Success      200      {object}  response.PaginationDataResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /centers [get]
func GetCenterRequests(ctx *gin.Context) {
	var request request.GetCenterRequests
	if ctx.ShouldBindQuery(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateCenterRequestService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetRequests(request, ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetWalletCenterRequests godoc
// @Summary      Get Center requests by a wallet
// @Description  Retrieve specific center requests associated with a unique wallet address
// @Tags         center
// @Accept       json
// @Produce      json
// @Param        id       path      string  true  "User Wallet Address"
// @Success      200      {object}  response.PaginationDataResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /centers/user/{id} [get]
func GetWalletCenterRequests(ctx *gin.Context) {
	service, err := business.GenerateCenterRequestService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetWalletRequests(ctx.Param("id"), ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetCenterRequest godoc
// @Summary      Get center request details
// @Description  Retrieve the full details of a single Center request by its unique ID
// @Tags         center
// @Accept       json
// @Produce      json
// @Param        id       path      string  true  "Center Request ID (UUID)"
// @Success      200      {object}  response.PaginationDataResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /centers/{id} [get]
func GetCenterRequest(ctx *gin.Context) {
	service, err := business.GenerateCenterRequestService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetRequest(ctx.Param("id"), ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// CreateCenterRequest godoc
// @Summary      Create a new center request
// @Description  Submit a new center request
// @Tags         center
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      request.CreateCenterRequest  true  "Center Request Body"
// @Success      201      {object}  entities.CenterRequest
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /centers [post]
func CreateCenterRequest(ctx *gin.Context) {
	var request request.CreateCenterRequest
	if ctx.ShouldBindJSON(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateCenterRequestService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.CreateRequest(request, ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.CREATE_ACTION,
	})
}

// VoteCenterRequest godoc
// @Summary      Vote on a center request
// @Description  Submit an approval or refusal vote for a specific center request.
// @Tags         center
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                true  "Center Request ID"
// @Param        request  body      request.VoteRequest   true  "Vote Details"
// @Success      200      {object}  response.MessageAPIResponse "Success"
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /centers/{id}/vote [post]
func VoteCenterRequest(ctx *gin.Context) {
	var request request.VoteRequest
	if ctx.ShouldBindJSON(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateCenterRequestService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	util.ProcessResponse(response.APIResponse{
		ErrMsg:   service.VoteRequest(ctx.Param("id"), request, ctx),
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// ConfirmCenterRequest godoc
// @Summary      Confirm and upload center information to Sui Blockhain
// @Description  Prepares and builds a transaction for uploading center information on-chain
// @Tags         center
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string  true  "Center Request ID (UUID)"
// @Success      200      {object}  response.BuildTransactionResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /centers/{id}/confirm [post]
func ConfirmCenterRequest(ctx *gin.Context) {
	service, err := business.GenerateCenterRequestService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.ConfirmRequest(ctx.Param("id"), ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}
