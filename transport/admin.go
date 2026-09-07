package transport

import (
	"raise-child/business"
	action_type "raise-child/constants/action_type"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/util"

	"github.com/gin-gonic/gin"
)

// GetAdmins godoc
// @Summary      List admins
// @Description  Retrieves a list of admins based on filter criteria
// @Tags         admins
// @Accept       json
// @Produce      json
// @Param        request  query     request.GetAdminsRequest  true  "Filter Criteria"
// @Success      200      {object}  response.PaginationDataResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /admins [get]
func GetAdmins(ctx *gin.Context) {
	var request request.GetAdminsRequest
	if ctx.ShouldBindQuery(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateAdminService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetAdmins(request, ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetAdmin godoc
// @Summary      Retrieve an admin
// @Description  Retrieve an admin
// @Tags         admins
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Admin NFT ID"
// @Success      200      {object}  response.AdminNftResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /admins/{id} [get]
func GetAdmin(ctx *gin.Context) {
	service, err := business.GenerateAdminService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetAdmin(ctx.Param("id"), ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetAdminByOwner godoc
// @Summary      Retrieve an admin by owner address
// @Description  Retrieve an admin by owner address
// @Tags         admins
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Admin Address"
// @Success      200      {object}  response.AdminNftResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /admins/user/{id} [get]
func GetAdminByOwner(ctx *gin.Context) {
	service, err := business.GenerateAdminService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetAdminByOwner(ctx.Param("id"), ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// UpdatePublisherInfo godoc
// @Summary      Update publisher information
// @Description  Updates the details of publisher (admin) based on the provided request body with just only 1 call permitted.
// @Tags         admins
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      request.UpdatePublisherInfoRequest  true  "Publisher Update Information"
// @Success      200      {object}  response.BuildTransactionResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /admins [post]
func UpdatePublisherInfo(ctx *gin.Context) {
	var request request.UpdatePublisherInfoRequest
	if ctx.ShouldBindJSON(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateAdminService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.UpdatePublisherInfo(request, ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}
