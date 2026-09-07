package transport

import (
	"raise-child/business"
	action_type "raise-child/constants/action_type"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/util"

	"github.com/gin-gonic/gin"
)

// GetPendingChildSpecialNeedProposals godoc
// @Summary      List Child Pending Special Need Proposals
// @Description  Retrieves a list of child pending special need proposals based on filter criteria
// @Tags         pending-special-needs
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  query     request.GetPendingChildSpecialNeedProposalsRequest  true  "Filter Criteria"
// @Success      200      {object}  response.PaginationDataResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /pending-special-needs [get]
func GetPendingChildSpecialNeedProposals(ctx *gin.Context) {
	var request request.GetPendingChildSpecialNeedProposalsRequest
	if ctx.ShouldBindQuery(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GeneratePendingChildSpecialNeedProposalService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetPendingChildSpecialNeedProposals(request, ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetPendingChildSpecialNeedProposal godoc
// @Summary      Get Child Pending Special Need Proposal Detail
// @Description  Retrieves child pending special need proposal information by its unique ID
// @Tags         pending-special-needs
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Pending Child Special Need Proposal ID"
// @Success      200      {object}  entities.PendingChildSpecialNeedProposal
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /pending-special-needs/{id} [get]
func GetPendingChildSpecialNeedProposal(ctx *gin.Context) {
	service, err := business.GeneratePendingChildSpecialNeedProposalService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetPendingChildSpecialNeedProposal(ctx.Param("id"), ctx)

	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// ApprovePendingChildSpecialNeedProposal godoc
// @Summary      Approve a Child Pending Special Need Proposal and upload to Sui Blockchain
// @Description  Prepares and executes a transaction for approve and upload child pending special need proposal on-chain
// @Tags         pending-special-needs
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Pending Child Special Need Proposal ID"
// @Success      200      {object}  response.MessageAPIResponse "Success"
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /pending-special-needs/{id}/approve [post]
func ApprovePendingChildSpecialNeedProposal(ctx *gin.Context) {
	service, err := business.GeneratePendingChildSpecialNeedProposalService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	util.ProcessResponse(response.APIResponse{
		ErrMsg:  service.ApprovePendingChildSpecialNeedProposal(ctx.Param("id"), ctx),
		Context: ctx,
	})
}

// RefusePendingChildSpecialNeedProposal godoc
// @Summary      Refuse a Child Pending Special Need Proposal
// @Description  Refuse a Child Pending Special Need Proposal
// @Tags         pending-special-needs
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Pending Child Special Need Proposal ID"
// @Success      200      {object}  response.MessageAPIResponse "Success"
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /pending-special-needs/{id}/refuse [post]
func RefusePendingChildSpecialNeedProposal(ctx *gin.Context) {
	service, err := business.GeneratePendingChildSpecialNeedProposalService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	util.ProcessResponse(response.APIResponse{
		ErrMsg:  service.RefusePendingChildSpecialNeedProposal(ctx.Param("id"), ctx),
		Context: ctx,
	})
}
