package transport

import (
	"raise-child/business"
	action_type "raise-child/constants/action_type"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/util"

	"github.com/gin-gonic/gin"
)

// GetGift godoc
// @Summary      Get a gift information
// @Description  Retrieve details of a specific gift by its ID
// @Tags         gift
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Gift ID"
// @Success      200      {object}  response.GiftResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /gifts/{id} [get]
func GetGift(ctx *gin.Context) {
	service, err := business.GenerateGiftService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetGift(ctx.Param("id"), ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetGiftsOfChild godoc
// @Summary      List gifts of a child
// @Description  Retrieve a list of gifts sent to a child with optional query filters
// @Tags         gift
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Child ID"
// @Param        request  query     request.GetGiftsRequest  true  "Filter Criteria"
// @Success      200      {object}  response.PaginationDataResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /gifts/child/{id} [get]
func GetGiftsOfChild(ctx *gin.Context) {
	var request request.GetGiftsRequest
	if ctx.ShouldBindQuery(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateGiftService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetGiftsOfChild(ctx.Param("id"), request, ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// CreateGift godoc
// @Summary      Create a new gift for a child
// @Description  Prepare and create a new gift for a child on-chain
// @Tags         gift
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      request.CreateGiftRequest   true  "Create gift details"
// @Success      201  {object}  response.MessageAPIResponse "Success"
// @Failure      400  {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      403  {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500  {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /gifts [post]
func CreateGift(ctx *gin.Context) {
	var request request.CreateGiftRequest
	if ctx.ShouldBindQuery(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateGiftService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	util.ProcessResponse(response.APIResponse{
		ErrMsg:  service.CreateGift(request, ctx),
		Context: ctx,
	})
}

// ConfirmReceiveGift godoc
// @Summary      Confirm recieve a gift
// @Description  Confirm recieve a gift
// @Tags         gift
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Gift ID"
// @Param        request  body      request.ConfirmReceiveGiftRequest   true  "Confirm receive gift details"
// @Success      200  {object}  response.BuildTransactionResponse
// @Failure      400  {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401  {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500  {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /gifts/{id}/confirm [post]
func ConfirmReceiveGift(ctx *gin.Context) {
	var request request.ConfirmReceiveGiftRequest
	if ctx.ShouldBindQuery(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateGiftService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	util.ProcessResponse(response.APIResponse{
		ErrMsg:  service.ConfirmReceiveGift(ctx.Param("id"), request, ctx),
		Context: ctx,
	})
}
