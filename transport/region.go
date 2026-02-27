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

// GetSupportedRegionProposals godoc
// @Summary      Get list of supported region proposals
// @Description  Retrieves a list of supported region proposals based on filter criteria
// @Tags         regions
// @Accept       json
// @Produce      json
// @Param        request  query     request.GetSupportedRegionProposalsRequest  true  "Filter Criteria"
// @Success      200  {object}  response.PaginationDataResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /regions/supported-proposal [get]
func GetSupportedRegionProposals(ctx *gin.Context) {
	var request request.GetSupportedRegionProposalsRequest
	if ctx.ShouldBindQuery(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateRegionService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetSupportedRegionProposals(request, ctx)

	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetUserSupportedRegionProposals godoc
// @Summary      Get list of supported region proposals from a user
// @Description  Retrieves a list of supported region proposals based on filter criteria from a user
// @Tags         regions
// @Accept       json
// @Produce      json
// @Param        id       path      string  true  "User Wallet Address"
// @Param        request  query     request.GetSupportedRegionProposalsRequest  true  "Filter Criteria"
// @Success      200  {object}  response.PaginationDataResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /regions/supported-proposal/user/{id} [get]
func GetUserSupportedRegionProposals(ctx *gin.Context) {
	var request request.GetSupportedRegionProposalsRequest
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

	res, err := service.GetSupportedRegionProposals(request, ctx)

	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetSupportedRegionProposal godoc
// @Summary      Get a supported region proposal detail
// @Description  Retrieves a supported region proposal detailed information
// @Tags         regions
// @Accept       json
// @Produce      json
// @Param        id       path      string  true  "Supported Region Proposal ID"
// @Success      200  {object}  entities.SupportedRegionProposal
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /regions/supported-proposal/{id} [get]
func GetSupportedRegionProposal(ctx *gin.Context) {
	service, err := business.GenerateRegionService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetSupportedRegionProposal(ctx.Param("id"), ctx)

	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// CreateSupportedRegionProposal godoc
// @Summary      Create a new supported region proposal
// @Description  Submit a new supported region proposal.
// @Tags         regions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      request.CreateSupportedRegionProposalsRequest  true  "Supported Region Request Body"
// @Success      201      {object}  entities.SupportedRegionProposal
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /regions/supported-proposal [post]
func CreateSupportedRegionProposal(ctx *gin.Context) {
	var request request.CreateSupportedRegionProposalsRequest
	if ctx.ShouldBindJSON(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateRegionService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.CreateSupportedRegionProposal(request, ctx)

	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.CREATE_ACTION,
	})
}
