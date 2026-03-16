package transport

import (
	"raise-child/business"
	action_type "raise-child/constants/action_type"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/util"

	"github.com/gin-gonic/gin"
)

// GetRegions godoc
// @Summary      Get list of regions
// @Description  Retrieves a list of all available regions
// @Tags         regions
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.RegionsResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /regions [get]
func GetRegions(ctx *gin.Context) {
	service, err := business.GenerateRegionService()
	util.ProcessResponse(response.APIResponse{
		Data1:    service.GetRegions(),
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetSupportedRegionSuggestions godoc
// @Summary      Get list of supported region proposals
// @Description  Retrieves a list of supported region proposals based on filter criteria
// @Tags         regions
// @Accept       json
// @Produce      json
// @Param        request  query     request.GetSupportedRegionSuggestionsRequest  true  "Filter Criteria"
// @Success      200  {object}  response.PaginationDataResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /regions/supported-suggestions [get]
func GetSupportedRegionSuggestions(ctx *gin.Context) {
	var request request.GetSupportedRegionSuggestionsRequest
	if ctx.ShouldBindQuery(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateRegionService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetSupportedRegionSuggestions(request, ctx)

	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetWalletSupportedRegionSuggestions godoc
// @Summary      Get list of supported region proposals created by a user
// @Description  Retrieves a list of supported region proposals of a user based on filter criteria
// @Tags         regions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string  true  "User Wallet Address"
// @Param        request  query     request.GetSupportedRegionSuggestionsRequest  true  "Filter Criteria"
// @Success      200      {object}  response.PaginationDataResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /regions/user/{id}/supported-suggestions [get]
func GetWalletSupportedRegionSuggestions(ctx *gin.Context) {
	var request request.GetSupportedRegionSuggestionsRequest
	if ctx.ShouldBindQuery(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateRegionService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	request.CreatedBy = ctx.Param("id")

	res, err := service.GetWalletSupportedRegionSuggestions(request, ctx)

	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// AdminGetSupportedRegionSuggestions godoc
// @Summary      Admin gets a list of supported region proposals
// @Description  Retrieves a list of supported region proposals based on filter criteria as actor is admin
// @Tags         regions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200      {object}  response.PaginationDataResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /regions/admin/supported-suggestions [get]
func AdminGetSupportedRegionSuggestions(ctx *gin.Context) {
	var request request.GetSupportedRegionSuggestionsRequest
	if ctx.ShouldBindQuery(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	request.CreatedBy = ctx.Param("id")

	service, err := business.GenerateRegionService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.AdminGetSupportedRegionSuggestions(request, ctx)

	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetSupportedRegionSuggestion godoc
// @Summary      Get a supported region proposal detail
// @Description  Retrieves a supported region proposal detailed information
// @Tags         regions
// @Accept       json
// @Produce      json
// @Param        id       path      string  true  "Supported Region Proposal ID"
// @Success      200  {object}  entities.SupportedRegionSuggestion
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /regions/supported-suggestions/{id} [get]
func GetSupportedRegionSuggestion(ctx *gin.Context) {
	service, err := business.GenerateRegionService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetSupportedRegionSuggestion(ctx.Param("id"), ctx)

	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// CreateSupportedRegionSuggestion godoc
// @Summary      Create a new supported region proposal
// @Description  Submit a new supported region proposal.
// @Tags         regions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      request.CreateSupportedRegionSuggestionsRequest  true  "Supported Region Request Body"
// @Success      201      {object}  entities.SupportedRegionSuggestion
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /regions/supported-suggestions [post]
func CreateSupportedRegionSuggestion(ctx *gin.Context) {
	var request request.CreateSupportedRegionSuggestionsRequest
	if ctx.ShouldBindJSON(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateRegionService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.CreateSupportedRegionSuggestion(request, ctx)

	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.CREATE_ACTION,
	})
}

// ReviewRegionSuggestion godoc
// @Summary      Review supported region proposal
// @Description  Review a supported region proposal.
// @Tags         regions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      request.VoteRequest  true  "Review Region Suggestion Detail"
// @Success      200      {object}  response.MessageAPIResponse "Success"
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /regions/supported-suggestions/{id}/review [post]
func ReviewRegionSuggestion(ctx *gin.Context) {
	var request request.VoteRequest
	if ctx.ShouldBindJSON(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateRegionService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	util.ProcessResponse(response.APIResponse{
		ErrMsg:   service.ReviewRegionSuggestion(ctx.Param("id"), request, ctx),
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}
