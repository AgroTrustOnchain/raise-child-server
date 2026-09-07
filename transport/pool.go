package transport

import (
	"raise-child/business"
	action_type "raise-child/constants/action_type"
	"raise-child/model/dtos/response"
	"raise-child/util"

	"github.com/gin-gonic/gin"
)

// GetLeaderPool godoc
// @Summary      Get pool detail of leader's region
// @Description  Get pool detail of leader's region
// @Tags         pools
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Leader Wallet Address"
// @Success      200      {object}  response.PoolResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /pools/leader/{id} [get]
func GetLeaderPool(ctx *gin.Context) {
	service, err := business.GeneratePoolService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetLeaderPool(ctx.Param("id"), ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}
