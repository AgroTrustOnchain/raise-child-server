package transport

import (
	"raise-child/business"
	action_type "raise-child/constants/action_type"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/util"

	"github.com/gin-gonic/gin"
)

// GetWithdrawProposal godoc
// @Summary      Get a withdraw proposal
// @Description  Retrieve details of a specific withdraw proposal by its ID
// @Tags         Withdrawal
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Withdraw Proposal ID"
// @Success      200      {object}  response.WithdrawProposalResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /withdraw-proposals/{id} [get]
func GetWithdrawProposal(ctx *gin.Context) {
	service, err := business.GenerateWithdrawProposalService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetWithdrawProposal(ctx.Param("id"), ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetWithdrawProposals godoc
// @Summary      List withdraw proposals
// @Description  Retrieve a list of withdraw proposals with optional query filters
// @Tags         Withdrawal
// @Accept       json
// @Produce      json
// @Param        request  query     request.GetWithdrawProposalsRequest  true  "Filter Criteria"
// @Success      200      {object}  response.PaginationDataResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /withdraw-proposals [get]
func GetWithdrawProposals(ctx *gin.Context) {
	var request request.GetWithdrawProposalsRequest
	if ctx.ShouldBindQuery(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateWithdrawProposalService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetWithdrawProposals(request, ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// VoteWithdrawProposal godoc
// @Summary      Vote on a withdraw proposal
// @Description  Submit a vote for a specific withdraw proposal. Note: This uses query parameters for the vote data.
// @Tags         Withdrawal
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id     path      string               true  "Withdraw Proposal ID"
// @Param        query  query     request.VoteRequest  true  "Vote Action Details"
// @Success      200    {object}  response.BuildTransactionResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /withdraw-proposals/{id}/vote [post]
func VoteWithdrawProposal(ctx *gin.Context) {
	var request request.VoteRequest
	if ctx.ShouldBindQuery(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateWithdrawProposalService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.VoteWithdrawProposal(ctx.Param("id"), request, ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// ConfirmWithdrawProposal godoc
// @Summary      Confirm a withdraw proposal
// @Description  Finalize and confirm a specific withdraw proposal by ID
// @Tags         Withdrawal
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Withdraw Proposal ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401  {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500  {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /withdraw-proposals/{id}/confirm [post]
func ConfirmWithdrawProposal(ctx *gin.Context) {
	service, err := business.GenerateWithdrawProposalService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.ConfirmWithdrawProposal(ctx.Param("id"), ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// ConfirmMainPoolWithdrawProposal godoc
// @Summary      Confirm a withdraw proposal from main pool
// @Description  Finalize and confirm a specific withdraw proposal by ID from main pool
// @Tags         Withdrawal
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Withdraw Proposal ID"
// @Param        id   query      string  true  "Bank Transaction Capture Image Blob ID"
// @Success      200  {object}  response.BuildTransactionResponse
// @Failure      400  {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401  {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500  {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /withdraw-proposals/{id}/main-pool-confirm [post]
func ConfirmMainPoolWithdrawProposal(ctx *gin.Context) {
	service, err := business.GenerateWithdrawProposalService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.ConfirmMainPoolWithdrawProposal(ctx.Param("id"), ctx.Query("blobId"), ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// CreateWithdrawProposal godoc
// @Summary      Create a new withdraw proposal
// @Description  Prepare and create a new withdraw proposal on-chain
// @Tags         Withdrawal
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      request.CreateWithdrawProposalRequest   true  "Withdraw Proposal details (e.g., "withdraw amount", "description")"
// @Success      201  {object}  response.BuildTransactionResponse
// @Failure      400  {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401  {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500  {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /withdraw-proposals [post]
func CreateWithdrawProposal(ctx *gin.Context) {
	var request request.CreateWithdrawProposalRequest
	if ctx.ShouldBindQuery(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateWithdrawProposalService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.CreateWithdrawProposal(request, ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.CREATE_ACTION,
	})
}
