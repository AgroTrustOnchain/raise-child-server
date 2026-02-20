package transport

import (
	"raise-child/business"
	action_type "raise-child/constants/action_type"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/util"

	"github.com/gin-gonic/gin"
)

// GetWalletNotifications godoc
// @Summary      List notifications of a user
// @Description  Retrieves a list of notifications from a user based on filter criteria
// @Tags         notification
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  query     request.GetNotisRequest  true  "Filter Criteria"
// @Param        id   path      string  true  "Wallet Address"
// @Success      200      {object}  response.PaginationDataResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /notis/user/{id} [get]
func GetWalletNotifications(ctx *gin.Context) {
	var request request.GetNotisRequest
	if ctx.ShouldBindQuery(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateNotiService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetCurrentWalletNotis(ctx.Param("id"), request, ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}
