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
// @Success      200      {object}  response.BuildTransactionResponse
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
