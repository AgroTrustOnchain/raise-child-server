package transport

import (
	"raise-child/business"
	action_type "raise-child/constants/action_type"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/util"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetUploadChildRequests godoc
// @Summary      Get upload-child requests
// @Description  Retrieve a list of child upload requests based on query parameters.
// @Tags         Child Upload Request
// @Accept       json
// @Produce      json
// @Param        request  query     request.GetUploadChildRequests   true  "Filter Criteria"
// @Success      200      {object}  response.PaginationDataResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /child-upload-reqs [get]
func GetUploadChildRequests(ctx *gin.Context) {
	var request request.GetUploadChildRequests
	if ctx.ShouldBindQuery(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateUploadChildRequestService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetUploadChildRequests(request, ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetWalletUploadChildRequests godoc
// @Summary      Get upload-child requests by a wallet
// @Description  Retrieve child upload requests associated with a specific wallet and page number.
// @Tags         Child Upload Request
// @Accept       json
// @Produce      json
// @Param        id     path      string  true  "Wallet address"
// @Param        page   query     int     false "Page number for pagination" default(1)
// @Success      200      {object}  response.PaginationDataResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /child-upload-reqs/user/{id} [get]
func GetWalletUploadChildRequests(ctx *gin.Context) {
	service, err := business.GenerateUploadChildRequestService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	page, _ := strconv.Atoi(ctx.Query("page"))

	res, err := service.GetWalletUploadChildRequests(ctx.Param("id"), page, ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetUploadChildRequest godoc
// @Summary      Get upload-child request details
// @Description  Retrieve the full details of a single upload-child request by its unique ID
// @Tags         Child Upload Request
// @Accept       json
// @Produce      json
// @Param        id       path      string  true  "Upload-child Request ID (UUID)"
// @Success      200      {object}  response.UploadChildRequestResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /child-upload-reqs/{id} [get]
func GetUploadChildRequest(ctx *gin.Context) {
	service, err := business.GenerateUploadChildRequestService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetUploadChildRequest(ctx.Param("id"), ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// CreateUploadChildRequest godoc
// @Summary      Create a new upload-child request
// @Description  Submit a JSON payload to create a new child upload request in the system.
// @Tags         Child Upload Request
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      request.UploadChildRequest  true  "Upload-child request payload"
// @Success      201      {object}  entities.UploadChildRequest
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /child-upload-reqs [post]
func CreateUploadChildRequest(ctx *gin.Context) {
	var request request.UploadChildRequest
	if ctx.ShouldBindJSON(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateUploadChildRequestService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.CreateUploadChildRequest(request, ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.CREATE_ACTION,
	})
}

// ReviewUploadChildRequest godoc
// @Summary      Review an upload-child request
// @Description  Submit an approval or refusal vote for reviewing a specific child upload request using its ID.
// @Tags         Child Upload Request
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                true  "Upload Child Request ID"
// @Param        request  body      request.VoteRequest   true  "Voting details (e.g., vote type, comments)"
// @Success      200      {object}  response.MessageAPIResponse "Success"
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /child-upload-reqs/{id}/review [post]
func ReviewUploadChildRequest(ctx *gin.Context) {
	var request request.VoteRequest
	if ctx.ShouldBindJSON(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateUploadChildRequestService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	util.ProcessResponse(response.APIResponse{
		ErrMsg:  service.ReviewUploadChildRequest(ctx.Param("id"), request, ctx),
		Context: ctx,
	})
}
