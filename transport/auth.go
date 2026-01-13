package transport

import (
	"raise-child/business"
	action_type "raise-child/constants/action_type"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/util"

	"github.com/gin-gonic/gin"
)

// Login godoc
// @Summary      User login
// @Description  Authenticates a user and returns access tokens/session data
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      request.LoginRequestV2  true  "Login credentials (e.g., address and sub)"
// @Success      200      {object}  response.MessageAPIResponse "Success"
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /auth/login [post]
func Login(ctx *gin.Context) {
	var request request.LoginRequestV2
	if ctx.ShouldBindJSON(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateAuthService()
	if ctx.ShouldBindJSON(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.LoginV2(request, ctx)

	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// Logout godoc
// @Summary      Logout user
// @Description  Invalidates the session for the specified user address
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200      {object}  response.MessageAPIResponse "Success"
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /auth/logout [post]
func Logout(ctx *gin.Context) {
	service, err := business.GenerateAuthService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	util.ProcessResponse(response.APIResponse{
		ErrMsg:   service.LogoutV2(ctx),
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetNonce godoc
// @Summary      Get login nonce
// @Description  Retrieves a unique nonce for the user to sign
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        address  path      string  true  "User wallet address"
// @Success      200      {object}  response.MessageAPIResponse "Success"
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /auth/nonce/{address} [get]
func GetNonce(ctx *gin.Context) {
	service, err := business.GenerateAuthService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetNonce(ctx.Param("address"), ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetSalt godoc
// @Summary      Get wallet salt
// @Description  Retrieves a unique salt for the user to login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        id  path      string  true  "User sub"
// @Success      200      {object}  response.MessageAPIResponse "Success"
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /auth/salt/{id} [get]
func GetSalt(ctx *gin.Context) {
	service, err := business.GenerateAuthService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetSalt(ctx.Param("id"), ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}
