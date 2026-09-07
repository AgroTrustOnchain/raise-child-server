package transport

import (
	"raise-child/business"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/util"

	"github.com/gin-gonic/gin"
)

// UpdateChildEditMealNeedDates godoc
// @Summary      Update child meal need edit dates
// @Description  Prepares and executes a transaction for updating child meal need edit dates on-chain
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

	util.ProcessResponse(response.APIResponse{
		ErrMsg:  service.UpdateChildEditMealNeedDates(request, ctx),
		Context: ctx,
	})
}

// UpdateChildEditBooksNeedDates godoc
// @Summary      Update child books need edit dates
// @Description  Prepares and executes a transaction for updating child books need edit dates on-chain
// @Tags         configs
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      request.UpdateChildEditNeedDatesRequest   true  "Update Child Books Need Edit Dates Detail"
// @Success      200      {object}  response.MessageAPIResponse "Success"
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

	util.ProcessResponse(response.APIResponse{
		ErrMsg:  service.UpdateChildEditBooksNeedDates(request, ctx),
		Context: ctx,
	})
}

// UpdateChildEditHealthInsuranceNeedDates godoc
// @Summary      Update child health insurance need edit dates
// @Description  Prepares and executes a transaction for updating child health insurance need edit dates on-chain
// @Tags         configs
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      request.UpdateChildEditNeedDatesRequest   true  "Update Child Health Insurance Need Edit Dates Detail"
// @Success      200      {object}  response.MessageAPIResponse "Success"
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

	util.ProcessResponse(response.APIResponse{
		ErrMsg:  service.UpdateChildEditHealthInsuranceNeedDates(request, ctx),
		Context: ctx,
	})
}

// EditSpecialNeedDao godoc
// @Summary      Update child special need dao
// @Description  Prepares and executes a transaction for updating child special need dao on-chain
// @Tags         configs
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      request.EditDaoRequest   true  "Edit Child Special Need DAO Request Detail"
// @Success      200      {object}  response.MessageAPIResponse "Success"
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /configs/child-special-need-dao [put]
func EditSpecialNeedDao(ctx *gin.Context) {
	var request request.EditDaoRequest
	if ctx.ShouldBindJSON(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateConfigService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	util.ProcessResponse(response.APIResponse{
		ErrMsg:  service.EditSpecialNeedDao(request, ctx),
		Context: ctx,
	})
}
