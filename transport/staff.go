package transport

import (
	"raise-child/business"
	action_type "raise-child/constants/action_type"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/util"

	"github.com/gin-gonic/gin"
)

// GetStaff godoc
// @Summary      Get staff details
// @Description  Retrieves staff information by its unique ID
// @Tags         staffs
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Staff ID"
// @Success      200      {object}  response.StaffNftResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /staffs/{id} [get]
func GetStaff(ctx *gin.Context) {
	service, err := business.GenerateStaffService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetStaff(ctx.Param("id"), ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetStaffByOwnerWallet godoc
// @Summary      Get staff details by owner
// @Description  Get staff details by owner
// @Tags         staffs
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Staff Owner Wallet Address"
// @Success      200      {object}  response.StaffNftResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /staffs/owner/{id} [get]
func GetStaffByOwnerWallet(ctx *gin.Context) {
	service, err := business.GenerateStaffService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetStaffByOwnerWallet(ctx.Param("id"), ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetStaffs godoc
// @Summary      List staffs
// @Description  Retrieves a list of staffs based on filter criteria
// @Tags         staffs
// @Accept       json
// @Produce      json
// @Param        request  query     request.GetStaffsRequest  true  "Filter Criteria"
// @Success      200      {object}  response.PaginationDataResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /staffs [get]
func GetStaffs(ctx *gin.Context) {
	var request request.GetStaffsRequest
	if ctx.ShouldBindQuery(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateStaffService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetStaffs(request, ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}
