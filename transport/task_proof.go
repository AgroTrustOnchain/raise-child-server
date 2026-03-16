package transport

import (
	"raise-child/business"
	action_type "raise-child/constants/action_type"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/util"

	"github.com/gin-gonic/gin"
)

// GetTaskProofs godoc
// @Summary      List task proofs
// @Description  Retrieve a list of task proofs with optional query filters
// @Tags         task-proofs
// @Accept       json
// @Produce      json
// @Param        request  query     request.GetTaskProofsRequest  true  "Filter Criteria"
// @Success      200      {object}  response.PaginationDataResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /task-proofs [get]
func GetTaskProofs(ctx *gin.Context) {
	var request request.GetTaskProofsRequest
	if ctx.ShouldBindQuery(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateTaskProofService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetTaskProofs(request, ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetTaskProof godoc
// @Summary      Get a task proof
// @Description  Retrieve details of a task proof by its ID
// @Tags         task-proofs
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Task Proof ID"
// @Success      200      {object}  response.PaginationDataResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /task-proofs/{id} [get]
func GetTaskProof(ctx *gin.Context) {
	service, err := business.GenerateTaskProofService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetTaskProof(ctx.Param("id"), ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// SubmitTaskProof godoc
// @Summary      Submit task proof for a task
// @Description  Submit task proof for a task
// @Tags         task-proofs
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Task ID"
// @Param        request  query     request.SubmitTaskProofRequest  true  "Submit Task Proof Detail"
// @Success      200      {object}  response.MessageAPIResponse "Sucsess"
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /task-proofs/task/{id}/submit [post]
func SubmitTaskProof(ctx *gin.Context) {
	var request request.SubmitTaskProofRequest
	if ctx.ShouldBindQuery(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateTaskProofService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	util.ProcessResponse(response.APIResponse{
		ErrMsg:   service.SubmitTaskProof(ctx.Param("id"), request, ctx),
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// ApproveTaskProof godoc
// @Summary      Approve a task proof
// @Description  Prepare and build a transaction for approve and upload a task proof on chain
// @Tags         task-proofs
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Task Proof ID"
// @Success      200      {object}  response.BuildTransactionResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /task-proofs/{id}/approve [post]
func ApproveTaskProof(ctx *gin.Context) {
	service, err := business.GenerateTaskProofService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.ApproveTaskProof(ctx.Param("id"), ctx)

	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// RefuseTaskProof godoc
// @Summary      Refuse a task proof
// @Description  Refuse a task proof
// @Tags         task-proofs
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Task Proof ID"
// @Success      200      {object}  response.MessageAPIResponse "Sucsess"
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /task-proofs/{id}/refuse [post]
func RefuseTaskProof(ctx *gin.Context) {
	service, err := business.GenerateTaskProofService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	util.ProcessResponse(response.APIResponse{
		ErrMsg:   service.RefuseTaskProof(ctx.Param("id"), ctx),
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}
