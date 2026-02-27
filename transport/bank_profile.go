package transport

import (
	"raise-child/business"
	action_type "raise-child/constants/action_type"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/util"

	"github.com/gin-gonic/gin"
)

// GetBankProfileByOwner godoc
// @Summary      Get bank profile owned by a wallet
// @Description  Retrieve detailed bank profile information using the bank's owner wallet
// @Tags         banks
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Wallet address"
// @Success      200  {object}  response.BankProfileResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /banks/user/{id} [get]
func GetBankProfileByOwner(ctx *gin.Context) {
	service, err := business.GenerateBankProfileService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetBankProfileByOwner(ctx.Param("id"), ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetBankProfile godoc
// @Summary      Get bank profile by specific id
// @Description  Retrieve detailed bank profile information by its unique id
// @Tags         banks
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Bank Profile ID"
// @Success      200  {object}  response.BankProfileResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /banks/{id} [get]
func GetBankProfile(ctx *gin.Context) {
	service, err := business.GenerateBankProfileService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetBankProfile(ctx.Param("id"), ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// CreateBankProfile godoc
// @Summary      Create a new bank profile
// @Description  Submit bank details to create a new profile in the system
// @Tags         banks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      request.CreateBankProfileRequest  true  "Bank Profile Creation Data"
// @Success      201      {object}  entities.BankProfile
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /banks [post]
func CreateBankProfile(ctx *gin.Context) {
	var request request.CreateBankProfileRequest
	if ctx.ShouldBindJSON(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateBankProfileService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.CreateBankProfile(request, ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.CREATE_ACTION,
	})
}

// UpdateBankProfile godoc
// @Summary      Update an existing bank profile
// @Description  Modify the details of an existing bank profile by its ID
// @Tags         banks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                            true  "Bank Profile ID"
// @Param        request  body      request.UpdateBankProfileRequest  true  "Updated Bank Profile Data"
// @Success      200      {object}  entities.BankProfile
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /banks/{id} [put]
func UpdateBankProfile(ctx *gin.Context) {
	var request request.UpdateBankProfileRequest
	if ctx.ShouldBindJSON(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateBankProfileService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.UpdateBankProfile(ctx.Param("id"), request, ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}
