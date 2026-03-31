package transport

import (
	"raise-child/business"
	action_type "raise-child/constants/action_type"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/util"

	"github.com/gin-gonic/gin"
)

// GetPendingWithdrawProposal godoc
// @Summary      Get a pending withdraw proposal
// @Description  Retrieve details of a specific pending withdraw proposal by its ID
// @Tags         Pending Withdrawal
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Pending Withdraw Proposal ID"
// @Success      200      {object}  entities.PendingWithdrawProposal
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /pending-withdraw-proposals/{id} [get]
func GetPendingWithdrawProposal(ctx *gin.Context) {
	service, err := business.GeneratePendingWithdrawProposalService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetPendingWithdrawProposal(ctx.Param("id"), ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetPendingWithdrawProposals godoc
// @Summary      List pending withdraw proposals
// @Description  Retrieve a list of pending withdraw proposals with optional query filters
// @Tags         Pending Withdrawal
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  query     request.GetPendingWithdrawProposalsRequest  true  "Filter Criteria"
// @Success      200      {object}  response.PaginationDataResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /pending-withdraw-proposals [get]
func GetPendingWithdrawProposals(ctx *gin.Context) {
	var request request.GetPendingWithdrawProposalsRequest
	if ctx.ShouldBindQuery(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GeneratePendingWithdrawProposalService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetPendingWithdrawProposals(request, ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// CreatePendingWithdrawProposal godoc
// @Summary      Create a new pending withdraw proposal
// @Description  Create a new pending withdraw proposal
// @Tags         Pending Withdrawal
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      request.CreatePendingWithdrawProposalRequest   true  "Pending Withdraw Proposal details (e.g., "withdraw amount", "description")"
// @Success      201  {object}  response.BuildTransactionResponse
// @Failure      400  {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401  {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500  {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /pending-withdraw-proposals [post]
func CreatePendingWithdrawProposal(ctx *gin.Context) {
	var request request.CreatePendingWithdrawProposalRequest
	if ctx.ShouldBindJSON(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GeneratePendingWithdrawProposalService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.CreatePendingWithdrawProposal(request, ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.CREATE_ACTION,
	})
}

// ApprovePendingWithdrawProposal godoc
// @Summary      Approve a pending withdraw proposal
// @Description  Approve a pending withdraw proposal to prepare and build a transaction to publish this proposal on chain
// @Tags         Pending Withdrawal
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id     path      string               true  "Oending Withdraw Proposal ID"
// @Success      200    {object}  response.BuildTransactionResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /pending-withdraw-proposals/{id}/approve [post]
func ApprovePendingWithdrawProposal(ctx *gin.Context) {
	service, err := business.GeneratePendingWithdrawProposalService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.ApprovePendingWithdrawProposal(ctx.Param("id"), ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// RefusePendingWithdrawProposal godoc
// @Summary      Refuse a pending withdraw proposal
// @Description  Refuse a pending withdraw proposal
// @Tags         Pending Withdrawal
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id     path      string               true  "Oending Withdraw Proposal ID"
// @Success      200    {object}  response.BuildTransactionResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /pending-withdraw-proposals/{id}/refuse [post]
func RefusePendingWithdrawProposal(ctx *gin.Context) {
	service, err := business.GeneratePendingWithdrawProposalService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	util.ProcessResponse(response.APIResponse{
		ErrMsg:   service.RefusePendingWithdrawProposal(ctx.Param("id"), ctx),
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}
