package transport

import (
	"raise-child/business"
	action_type "raise-child/constants/action_type"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/util"

	"github.com/gin-gonic/gin"
)

// GetPendingCampaigns godoc
// @Summary      List Pending Pool Campaigns
// @Description  Retrieves a list of Pending Pool Campaigns based on filter criteria
// @Tags         pending-pool-campaigns
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  query     request.GetPendingCampaignsRequest  true  "Filter Criteria"
// @Success      200      {object}  response.PaginationDataResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /pending-pool-campaigns [get]
func GetPendingCampaigns(ctx *gin.Context) {
	var request request.GetPendingCampaignsRequest
	if ctx.ShouldBindQuery(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GeneratePendingCampaignService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetPendingCampaigns(request, ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetPendingCampaign godoc
// @Summary      Get Pending Pool Campaign Detail
// @Description  Retrieves Pending Pool Campaign information by its unique ID
// @Tags         pending-pool-campaigns
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Pending Pool Campaign ID"
// @Success      200      {object}  entities.PendingCampaign
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /pending-pool-campaigns/{id} [get]
func GetPendingCampaign(ctx *gin.Context) {
	service, err := business.GeneratePendingCampaignService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetPendingCampaign(ctx.Param("id"), ctx)

	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// ApprovePendingCampaign godoc
// @Summary      Approve a Pending Pool Campaign and upload to Sui Blockchain
// @Description  Prepares and builds a transaction for approve and upload Pending Pool Campaign on-chain
// @Tags         pending-pool-campaigns
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Pending Pool Campaign ID"
// @Success      200      {object}  response.BuildTransactionResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /pending-pool-campaigns/{id}/approve [post]
func ApprovePendingCampaign(ctx *gin.Context) {
	service, err := business.GeneratePendingCampaignService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.ApprovePendingCampaign(ctx.Param("id"), ctx)

	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// RefusePendingCampaign godoc
// @Summary      Refuse Pending Pool Campaign
// @Description  Refuse Pending Pool Campaign
// @Tags         pending-pool-campaigns
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Pending Pool Campaign ID"
// @Success      200      {object}  response.MessageAPIResponse "Success"
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /pending-pool-campaigns/{id}/refuse [post]
func RefusePendingCampaign(ctx *gin.Context) {
	service, err := business.GeneratePendingCampaignService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	util.ProcessResponse(response.APIResponse{
		ErrMsg:   service.RefusePendingCampaign(ctx.Param("id"), ctx),
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// CreatePendingCampaign godoc
// @Summary      Create Pending Pool Campaign
// @Description  Create Pending Pool Campaign and wait for admin to review it
// @Tags         pending-pool-campaigns
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      request.CreatePendingCampaignRequest   true  "Create Pending Pool Campaign Detail"
// @Success      201      {object}  entities.PendingCampaign
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /pending-pool-campaigns [post]
func CreatePendingCampaign(ctx *gin.Context) {
	var request request.CreatePendingCampaignRequest
	if ctx.ShouldBindJSON(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GeneratePendingCampaignService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.CreatePendingCampaign(request, ctx)

	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.CREATE_ACTION,
	})
}
