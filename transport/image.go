package transport

import (
	"raise-child/business"
	action_type "raise-child/constants/action_type"
	"raise-child/model/dtos/response"
	"raise-child/util"

	"github.com/gin-gonic/gin"
)

// GetPresignedUrl godoc
// @Summary      Presign before upload image to Cloudinary
// @Description  Presign before upload image to Cloudinary
// @Tags         image
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.PresignedUrlResponse
// @Failure      400  {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      403  {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500  {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /images/presigned-url [get]
func GetPresignedUrl(ctx *gin.Context) {
	var service = business.GenerateImageService()

	res, err := service.GetPresignedUrl(ctx)

	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}
