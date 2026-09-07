package transport

import (
	"raise-child/business"
	action_type "raise-child/constants/action_type"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/util"

	"github.com/gin-gonic/gin"
)

// GetNumericPlatformConfigs godoc
// @Summary      Get numeric platform configs
// @Description  Get numeric platform configs
// @Tags         Platform Config
// @Accept       json
// @Produce      json
// @Param        request  query     request.GetPlatformConfigsRequest  true  "Filter Criteria"
// @Success      200      {object}  response.PaginationDataResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /platform-configs/numeric [get]
func GetNumericPlatformConfigs(ctx *gin.Context) {
	var request request.GetPlatformConfigsRequest
	if ctx.ShouldBindQuery(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GeneratePlatformConfigService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetConfigs(request, true, ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetStringPlatformConfigs godoc
// @Summary      Get string platform configs
// @Description  Get string platform configs
// @Tags         Platform Config
// @Accept       json
// @Produce      json
// @Param        request  query     request.GetPlatformConfigsRequest  true  "Filter Criteria"
// @Success      200      {object}  response.PaginationDataResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /platform-configs/string [get]
func GetStringPlatformConfigs(ctx *gin.Context) {
	var request request.GetPlatformConfigsRequest
	if ctx.ShouldBindQuery(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GeneratePlatformConfigService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetConfigs(request, false, ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetNumericPlatformConfig godoc
// @Summary      Get numeric platform config
// @Description  Get numeric platform config
// @Tags         Platform Config
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Platform Config ID"
// @Success      200      {object}  entities.PlatformConfig
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /platform-configs/numeric/{id} [get]
func GetNumericPlatformConfig(ctx *gin.Context) {
	service, err := business.GeneratePlatformConfigService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetConfig(ctx.Param("id"), true, ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetStringPlatformConfig godoc
// @Summary      Get string platform config
// @Description  Get string platform config
// @Tags         Platform Config
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Platform Config ID"
// @Success      200      {object}  entities.PlatformConfig
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /platform-configs/string/{id} [get]
func GetStringPlatformConfig(ctx *gin.Context) {
	service, err := business.GeneratePlatformConfigService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetConfig(ctx.Param("id"), false, ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// UpdateNumericPlatformConfig godoc
// @Summary      Update numeric platform config
// @Description  Update numeric platform config
// @Tags         Platform Config
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Platform Config ID"
// @Param        request  body     request.UpdatePlatformConfigRequest  true  "Update Platform Config Detail"
// @Success      200      {object}  response.MessageAPIResponse "success"
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      403  {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /platform-configs/numeric/{id} [put]
func UpdateNumericPlatformConfig(ctx *gin.Context) {
	var request request.UpdatePlatformConfigRequest
	if ctx.ShouldBindJSON(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GeneratePlatformConfigService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	util.ProcessResponse(response.APIResponse{
		ErrMsg:  service.UpdateConfig(ctx.Param("id"), request, true, ctx),
		Context: ctx,
	})
}

// UpdateStringPlatformConfig godoc
// @Summary      Update string platform config
// @Description  Update string platform config
// @Tags         Platform Config
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Platform Config ID"
// @Param        request  body     request.UpdatePlatformConfigRequest  true  "Update Platform Config Detail"
// @Success      200      {object}  response.MessageAPIResponse "success"
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      403  {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /platform-configs/string/{id} [put]
func UpdateStringPlatformConfig(ctx *gin.Context) {
	var request request.UpdatePlatformConfigRequest
	if ctx.ShouldBindJSON(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GeneratePlatformConfigService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	util.ProcessResponse(response.APIResponse{
		ErrMsg:   service.UpdateConfig(ctx.Param("id"), request, false, ctx),
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}
