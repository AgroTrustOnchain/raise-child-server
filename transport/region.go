package transport

import (
	"raise-child/business"
	action_type "raise-child/constants/action_type"
	"raise-child/model/dtos/response"
	"raise-child/util"

	"github.com/gin-gonic/gin"
)

// GetRegions godoc
// @Summary      Get list of regions
// @Description  Retrieves a list of all available regions.
// @Tags         regions
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.RegionsResponse
// @Router       /regions [get]
func GetRegions(ctx *gin.Context) {
	util.ProcessResponse(response.APIResponse{
		Data1:    business.GenerateRegionService().GetRegions(),
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}
