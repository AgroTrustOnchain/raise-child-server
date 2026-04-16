package transport

import (
	"raise-child/business"
	action_type "raise-child/constants/action_type"
	"raise-child/model/dtos/response"
	"raise-child/util"

	"github.com/gin-gonic/gin"
)

// GetChildBooksNeed godoc
// @Summary      Get books need of a child
// @Description  Retrieves books need of a child by its unique ID
// @Tags         child-needs
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Child Books Need ID"
// @Success      200      {object}  response.BooksNeedResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /child-needs/books-need/{id} [get]
func GetChildBooksNeed(ctx *gin.Context) {
	service, err := business.GenerateChildNeedService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetChildBooksNeed(ctx.Param("id"), ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetChildHealthInsuranceNeed godoc
// @Summary      Get health insurance need of a child
// @Description  Retrieves health insurance need of a child by its unique ID
// @Tags         child-needs
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Child Health Insurance Need ID"
// @Success      200      {object}  response.HealthInsuranceNeedResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /child-needs/health-insurance-need/{id} [get]
func GetChildHealthInsuranceNeed(ctx *gin.Context) {
	service, err := business.GenerateChildNeedService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetChildHealthInsuranceNeed(ctx.Param("id"), ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetChildMealNeed godoc
// @Summary      Get meal need of a child
// @Description  Retrieves meal need of a child by its unique ID
// @Tags         child-needs
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Child Meal Need ID"
// @Success      200      {object}  response.MealNeedResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /child-needs/meal-need/{id} [get]
func GetChildMealNeed(ctx *gin.Context) {
	service, err := business.GenerateChildNeedService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetChildMealNeed(ctx.Param("id"), ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}
