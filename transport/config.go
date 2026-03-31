package transport

import (
	"raise-child/business"
	action_type "raise-child/constants/action_type"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/util"

	"github.com/gin-gonic/gin"
)

// UpdateChildEditMealNeedDates godoc
// @Summary      Update child meal need edit dates
// @Description  Prepares and builds a transaction for updating child meal need edit dates on-chain
// @Tags         configs
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      request.UpdateChildEditNeedDatesRequest   true  "Update Child Meal Need Edit Dates Detail"
// @Success      200      {object}  response.BuildTransactionResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /configs/meal-need-edit-dates [put]
func UpdateChildEditMealNeedDates(ctx *gin.Context) {
	var request request.UpdateChildEditNeedDatesRequest
	if ctx.ShouldBindJSON(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateConfigService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.UpdateChildEditMealNeedDates(request, ctx)

	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// UpdateChildEditBooksNeedDates godoc
// @Summary      Update child books need edit dates
// @Description  Prepares and builds a transaction for updating child books need edit dates on-chain
// @Tags         configs
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      request.UpdateChildEditNeedDatesRequest   true  "Update Child Books Need Edit Dates Detail"
// @Success      200      {object}  response.BuildTransactionResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /configs/books-need-edit-dates [put]
func UpdateChildEditBooksNeedDates(ctx *gin.Context) {
	var request request.UpdateChildEditNeedDatesRequest
	if ctx.ShouldBindJSON(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateConfigService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.UpdateChildEditBooksNeedDates(request, ctx)

	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// UpdateChildEditHealthInsuranceNeedDates godoc
// @Summary      Update child health insurance need edit dates
// @Description  Prepares and builds a transaction for updating child health insurance need edit dates on-chain
// @Tags         configs
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      request.UpdateChildEditNeedDatesRequest   true  "Update Child Health Insurance Need Edit Dates Detail"
// @Success      200      {object}  response.BuildTransactionResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /configs/health-insurance-need-edit-dates [put]
func UpdateChildEditHealthInsuranceNeedDates(ctx *gin.Context) {
	var request request.UpdateChildEditNeedDatesRequest
	if ctx.ShouldBindJSON(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateConfigService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.UpdateChildEditHealthInsuranceNeedDates(request, ctx)

	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}
