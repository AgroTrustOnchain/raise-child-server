package transport

import (
	"raise-child/business"
	action_type "raise-child/constants/action_type"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/util"

	"github.com/gin-gonic/gin"
)

// GetCampaigns godoc
// @Summary      List Pool Campaigns
// @Description  Retrieves a list of Pool Campaigns based on filter criteria
// @Tags         pool-campaigns
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  query     request.GetCampaignsRequest  true  "Filter Criteria"
// @Success      200      {object}  response.PaginationDataResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /pool-campaigns [get]
func GetCampaigns(ctx *gin.Context) {
	var request request.GetCampaignsRequest
	if ctx.ShouldBindQuery(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateCampaignService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetCampaigns(request, ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetCampaign godoc
// @Summary      Get Pool Campaign Detail
// @Description  Retrieves Pool Campaign information by its unique ID
// @Tags         pool-campaigns
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Pool Campaign ID"
// @Success      200      {object}  response.OnChainCampaignResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /pool-campaigns/{id} [get]
func GetCampaign(ctx *gin.Context) {
	service, err := business.GenerateCampaignService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetCampaign(ctx.Param("id"), ctx)

	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// SupportCampaign godoc
// @Summary      Support a Pool Campaign
// @Description  Support a Pool Campaign
// @Tags         pool-campaigns
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Pool Campaign ID"
// @Param        request  body      request.SupportCampaignRequest   true  "Support Pool Campaign Detail"
// @Success      200      {object}  response.PaymentUrlResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /pool-campaigns/{id}/support [post]
func SupportCampaign(ctx *gin.Context) {
	var request request.SupportCampaignRequest
	if ctx.ShouldBindJSON(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateCampaignService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.SupportCampaign(ctx.Param("id"), request, ctx)

	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// CreateCampaignWithdrawProposal godoc
// @Summary      Create Pool Campaign Pending Withdraw Proposal
// @Description  Create Pool Campaign Pending Withdraw Proposal
// @Tags         pool-campaigns
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      request.CreateCampaignWithdrawProposalRequest   true  "Create Pool Campaign Pending Withdraw Proposal Detail"
// @Success      201      {object}  response.MessageAPIResponse "Success"
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /pool-campaigns/withdraw-proposal [post]
func CreateCampaignWithdrawProposal(ctx *gin.Context) {
	var request request.CreateCampaignWithdrawProposalRequest
	if ctx.ShouldBindJSON(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateCampaignService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	util.ProcessResponse(response.APIResponse{
		ErrMsg:  service.CreateCampaignWithdrawProposal(request, ctx),
		Context: ctx,
	})
}
