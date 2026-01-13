package transport

import (
	"raise-child/business"
	action_type "raise-child/constants/action_type"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/util"

	"github.com/gin-gonic/gin"
)

// UploadProfile godoc
// @Summary      Upload an account personal information
// @Description  Upload an account personal information based on profile ID
// @Tags         profile
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                      true  "Profile ID"
// @Param        request  body      request.UploadProfileRequest true  "Profile Information"
// @Success      200      {object}  response.PersonalProfileResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /profiles/{id} [post]
func UploadProfile(ctx *gin.Context) {
	var request request.UploadProfileRequest
	if ctx.ShouldBindJSON(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateProfileService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.UploadProfile(ctx.Param("id"), request, ctx)

	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}
