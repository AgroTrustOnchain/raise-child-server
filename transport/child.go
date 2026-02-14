package transport

import (
	"raise-child/business"
	action_type "raise-child/constants/action_type"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/util"

	"github.com/gin-gonic/gin"
)

// GetChildren godoc
// @Summary      List children
// @Description  Retrieves a list of children based on filter criteria
// @Tags         children
// @Accept       json
// @Produce      json
// @Param        request  query     request.GetChildrenRequest  true  "Filter Criteria"
// @Success      200      {object}  response.PaginationDataResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /children [get]
func GetChildren(ctx *gin.Context) {
	var request request.GetChildrenRequest
	if ctx.ShouldBindQuery(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateChildService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetChildren(request, ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetChild godoc
// @Summary      Get child details
// @Description  Retrieves child information by its unique ID
// @Tags         children
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Child ID"
// @Success      200      {object}  response.ChildResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /children/{id} [get]
func GetChild(ctx *gin.Context) {
	service, err := business.GenerateChildService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetChild(ctx.Param("id"), ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// UploadChild godoc
// @Summary      Upload a new child to Sui Blockchain
// @Description  Prepares and builds a transaction for uploading new child information (e.g., "height", "weight") on-chain
// @Tags         children
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      request.UploadChildRequest   true  "Child Information (e.g. "first name", "last name", "gender")"
// @Success      201      {object}  response.BuildTransactionResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /children [post]
func UploadChild(ctx *gin.Context) {
	var request request.UploadChildRequest
	if ctx.ShouldBindJSON(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateChildService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.UploadChild(request, ctx)

	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// AddChildStringMetadata godoc
// @Summary      Add a new string metadata to a child object
// @Description  Prepares and builds a transaction for adding new field information to a child information (e.g., height, weight) on-chain
// @Tags         children
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      request.AddChildStringMetadaRequest   true  "Add child string metadata details (e.g., key, value)"
// @Success      200      {object}  response.BuildTransactionResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /children/metadata/string/{id} [put]
func AddChildStringMetadata(ctx *gin.Context) {
	var request request.AddChildStringMetadaRequest
	if ctx.ShouldBindJSON(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateChildService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.AddStringMetada(ctx.Param("id"), request, ctx)

	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// AddChildNumberMetadata godoc
// @Summary      Add a new number metadata to a child object
// @Description  Prepares and builds a transaction for adding new field information to a child information (e.g., height, weight) on-chain
// @Tags         children
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      request.AddChildStringMetadaRequest   true  "Add child number metadata details (e.g., key, value)"
// @Success      200      {object}  response.BuildTransactionResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /children/metadata/number/{id} [put]
func AddChildNumberMetadata(ctx *gin.Context) {
	var request request.AddChildNumberMetadaRequest
	if ctx.ShouldBindJSON(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateChildService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.AddNumberMetada(ctx.Param("id"), request, ctx)

	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// SupportBooksNeed godoc
// @Summary      Support books need for a child
// @Description  Prepares and builds a transaction for supporting books need for a child on-chain
// @Tags         children
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Books Need ID"
// @Success      200      {object}  response.UrlAPIResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /children/books-need/{id}/support [post]
func SupportBooksNeed(ctx *gin.Context) {
	service, err := business.GenerateChildService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.SupportBooksNeed(ctx.Param("id"), ctx)

	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// SupportMealNeed godoc
// @Summary      Support meal need for a child
// @Description  Prepares and builds a transaction for supporting meal need for a child on-chain
// @Tags         children
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Meal Need ID"
// @Param        request  body      request.SupportMealNeadRequest   true  "Support Child Meal Need Detail"
// @Success      200      {object}  response.UrlAPIResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /children/meal-need/{id}/support [post]
func SupportMealNeed(ctx *gin.Context) {
	var request request.SupportMealNeadRequest
	if ctx.ShouldBindJSON(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateChildService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.SupportMealNeed(ctx.Param("id"), request, ctx)

	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// SupportSpecialNeed godoc
// @Summary      Support special need campaign of a child
// @Description  Prepares and builds a transaction for supporting special need campaign of a child on-chain
// @Tags         children
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Campaign ID"
// @Param        request  body      request.SupportSpecialNeedRequest   true  "Support Child Special Need Campaign Detail"
// @Success      200      {object}  response.UrlAPIResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /children/special-need/{id}/support [post]
func SupportSpecialNeed(ctx *gin.Context) {
	var request request.SupportSpecialNeedRequest
	if ctx.ShouldBindJSON(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateChildService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.SupportSpecialNeed(ctx.Param("id"), request, ctx)

	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// CreateBooksNeedWithdrawProposal godoc
// @Summary      Create books need withdraw proposal
// @Description  Prepares and builds a transaction for creating a new withdraw proposal from child's books need on-chain
// @Tags         children
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      request.CreateNormalNeedWithdrawProposalRequest   true  "Create Books Need Withdraw Proposal Detail"
// @Success      200      {object}  response.BuildTransactionResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /children/books-need/withdraw-proposal [post]
func CreateBooksNeedWithdrawProposal(ctx *gin.Context) {
	var request request.CreateNormalNeedWithdrawProposalRequest
	if ctx.ShouldBindJSON(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateChildService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.CreateBooksNeedWithdrawProposal(request, ctx)

	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// CreateMealNeedWithdrawProposal godoc
// @Summary      Create meal need withdraw proposal
// @Description  Prepares and builds a transaction for creating a new withdraw proposal from child's meal need on-chain
// @Tags         children
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      request.CreateNormalNeedWithdrawProposalRequest   true  "Create Meal Need Withdraw Proposal Detail"
// @Success      200      {object}  response.BuildTransactionResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /children/meal-need/withdraw-proposal [post]
func CreateMealNeedWithdrawProposal(ctx *gin.Context) {
	var request request.CreateNormalNeedWithdrawProposalRequest
	if ctx.ShouldBindJSON(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateChildService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.CreateMealNeedWithdrawProposal(request, ctx)

	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// CreateSpecialNeedWithdrawProposal godoc
// @Summary      Create special need campaign withdraw proposal
// @Description  Prepares and builds a transaction for creating a new withdraw proposal from child's special need campaign on-chain
// @Tags         children
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      request.CreateSpecialNeedWithdrawProposalRequest   true  "Create Special Need Campaign Withdraw Proposal Detail"
// @Success      200      {object}  response.BuildTransactionResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /children/special-need/withdraw-proposal [post]
func CreateSpecialNeedWithdrawProposal(ctx *gin.Context) {
	var request request.CreateSpecialNeedWithdrawProposalRequest
	if ctx.ShouldBindJSON(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateChildService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.CreateSpecialNeedWithdrawProposal(request, ctx)

	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// CreateSpecialNeedProposal godoc
// @Summary      Create special need proposal
// @Description  Prepares and builds a transaction for creating a new special need proposal for a child on-chain
// @Tags         children
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      request.CreateSpecialNeedProposalRequest   true  "Create Special Need Campaign Withdraw Proposal Detail"
// @Success      200      {object}  response.BuildTransactionResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /children/special-need/proposal [post]
func CreateSpecialNeedProposal(ctx *gin.Context) {
	var request request.CreateSpecialNeedProposalRequest
	if ctx.ShouldBindJSON(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateChildService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.CreateSpecialNeedProposal(request, ctx)

	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// ConfirmSpecialNeedProposal godoc
// @Summary      Confirm a special need proposal
// @Description  Prepares and builds a transaction for confirming a special need proposal for a child if accepted to create a new campaign on-chain
// @Tags         children
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Special Need Proposal ID"
// @Success      200      {object}  response.BuildTransactionResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /children/special-need/proposal/{id}/confirm [post]
func ConfirmSpecialNeedProposal(ctx *gin.Context) {
	service, err := business.GenerateChildService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.ConfirmSpecialNeedProposal(ctx.Param("id"), ctx)

	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// VoteSpecialNeedProposal godoc
// @Summary      Vote a special need proposal
// @Description  Prepares and builds a transaction for voting a special need proposal of a child on-chain
// @Tags         children
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Special Need Proposal ID"
// @Param        request  body      request.VoteRequest   true  "Voting Detail"
// @Success      200      {object}  response.BuildTransactionResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /children/special-need/proposal/{id}/vote [post]
func VoteSpecialNeedProposal(ctx *gin.Context) {
	var request request.VoteRequest
	if ctx.ShouldBindJSON(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateChildService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.VoteSpecialNeedProposal(ctx.Param("id"), request, ctx)

	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}
