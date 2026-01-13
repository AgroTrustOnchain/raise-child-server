package transport

import (
	"raise-child/business"
	action_type "raise-child/constants/action_type"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/util"

	"github.com/gin-gonic/gin"
)

// GetSponsor godoc
// @Summary      Get sponsor details
// @Description  Retrieves sponsor information by unique ID
// @Tags         sponsor
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Sponsor ID"
// @Success      200      {object}  response.SponsorResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /sponsors/{id} [get]
func GetSponsor(ctx *gin.Context) {
	service, err := business.GenerateSponsorService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetSponsor(ctx.Param("id"), ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetSponsors godoc
// @Summary      List sponsors
// @Description  Retrieves a list of sponsors based on filter criteria
// @Tags         sponsor
// @Accept       json
// @Produce      json
// @Param        request  query     request.GetSponsorsRequest  true  "Filter Criteria"
// @Success      200      {object}  response.PaginationDataResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /sponsors [get]
func GetSponsors(ctx *gin.Context) {
	var request request.GetSponsorsRequest
	if ctx.ShouldBindQuery(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateSponsorService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetSponsors(request, ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}
