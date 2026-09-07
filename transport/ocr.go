package transport

import (
	"raise-child/business"
	action_type "raise-child/constants/action_type"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/util"

	"github.com/gin-gonic/gin"
)

// GetExtractChildInfo godoc
// @Summary Extract child information from OCR input
// @Description Extract child information from an uploaded OCR payload
// @Tags OCR
// @Accept json
// @Produce json
// @Param request body request.ExtractChildUploadInfoRequest true "OCR request payload - second_guardian_id_card_url is optional"
// @Success 200 {object} response.ExtractChildUploadInfoResponse
// @Failure 400 {object} response.MessageAPIResponse
// @Failure 500 {object} response.MessageAPIResponse
// @Router /ocr/extract-child-info [post]
func GetExtractChildInfo(ctx *gin.Context) {
	var req request.ExtractChildUploadInfoRequest
	
	if err := ctx.ShouldBindJSON(&req); err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateOCRService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetExtractChildInfo(req, ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})

}
