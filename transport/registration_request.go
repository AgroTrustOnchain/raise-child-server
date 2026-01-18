package transport

import (
	"raise-child/business"
	action_type "raise-child/constants/action_type"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/util"

	"github.com/gin-gonic/gin"
)

// GetRegistrationRequests godoc
// @Summary      Get a list of registration requests
// @Description  Retrieve registration requests with filtering based on query parameters
// @Tags         registration
// @Accept       json
// @Produce      json
// @Param        request  query     request.GetRegistrationRequests  true  "Filter Criteria"
// @Success      200      {object}  response.PaginationDataResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /registrations [get]
func GetRegistrationRequests(ctx *gin.Context) {
	var request request.GetRegistrationRequests
	if ctx.ShouldBindQuery(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateRegistrationRequestService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetRegistrationRequests(request, ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetWalletRegistrationRequests godoc
// @Summary      Get registration requests by a wallet
// @Description  Retrieve specific registration requests associated with a unique wallet address
// @Tags         registration
// @Accept       json
// @Produce      json
// @Param        id       path      string  true  "User Wallet Address" default(0x...)
// @Success      200      {object}  response.PaginationDataResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /registrations/user/{id} [get]
func GetWalletRegistrationRequests(ctx *gin.Context) {
	service, err := business.GenerateRegistrationRequestService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetWalletRegistrationRequests(ctx.Param("id"), ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetRegistrationRequest godoc
// @Summary      Get registration request details
// @Description  Retrieve the full details of a single registration request by its unique ID
// @Tags         registration
// @Accept       json
// @Produce      json
// @Param        id       path      string  true  "Registration Request ID (UUID)"
// @Success      200      {object}  response.PaginationDataResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /registrations/{id} [get]
func GetRegistrationRequest(ctx *gin.Context) {
	service, err := business.GenerateRegistrationRequestService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetRegistrationRequest(ctx.Param("id"), ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// CreateRegistrationRequest godoc
// @Summary      Create a new registration request
// @Description  Submit a new registration request. The phone number should be in a valid format (e.g., 090...) and email must be valid.
// @Tags         registration
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      request.CreateRegistrationRequest  true  "Registration Request Body"
// @Success      201      {object}  entities.RegistrationRequest
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /registrations [post]
func CreateRegistrationRequest(ctx *gin.Context) {
	var request request.CreateRegistrationRequest
	if ctx.ShouldBindJSON(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateRegistrationRequestService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.CreateRegistrationRequest(request, ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.CREATE_ACTION,
	})
}

// VoteRegistrationRequest godoc
// @Summary      Vote on a registration request
// @Description  Submit an approval or refusal vote for a specific registration request.
// @Tags         registration
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                true  "Registration Request ID"
// @Param        request  body      request.VoteRequest   true  "Vote Details"
// @Success      200      {object}  response.MessageAPIResponse "Success"
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /registrations/{id}/vote [post]
func VoteRegistrationRequest(ctx *gin.Context) {
	var request request.VoteRequest
	if ctx.ShouldBindJSON(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateRegistrationRequestService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	util.ProcessResponse(response.APIResponse{
		ErrMsg:   service.VoteRegistrationRequest(ctx.Param("id"), request, ctx),
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// ConfirmRegistrationRequest godoc
// @Summary      Confirm and register a staff role uploaded to Sui Blockhain
// @Description  Prepares and builds a transaction for registering new staff information on-chain
// @Tags         registration
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string  true  "Registration Request ID (UUID)"
// @Success      200      {object}  response.BuildTransactionResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /registrations/{id}/confirm [post]
func ConfirmRegistrationRequest(ctx *gin.Context) {
	service, err := business.GenerateRegistrationRequestService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.ConfirmRegistrationRequest(ctx.Param("id"), ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}
