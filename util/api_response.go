package util

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	action_type "raise-child/constants/action_type"
	"raise-child/constants/noti"
	"raise-child/model/dtos/response"

	"github.com/gin-gonic/gin"
)

func ProcessResponse(data response.APIResponse) {
	if data.ErrMsg != nil {
		processFailResponse(data.ErrMsg, data.Context)
		return
	}

	if data.PostType != action_type.NON_POST {
		processSuccessPostReponse(data.Data2, data.PostType, data.Context)
		return
	}

	processSuccessResponse(data.Data1, data.Context)
}

func GenerateInvalidRequestAndSystemProblemModel(ctx *gin.Context, err error) response.APIResponse {
	var errMsg error = err
	if errMsg == nil {
		errMsg = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	return response.APIResponse{
		ErrMsg:   errMsg,
		Context:  ctx,
		PostType: action_type.NON_POST,
	}
}

func GetUnAuthBodyResponse(ctx *gin.Context) response.APIResponse {
	return response.APIResponse{
		ErrMsg:  errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG),
		Context: ctx,
	}
}

func GetRateLimitRequestBodyResponse(ctx *gin.Context) response.APIResponse {
	return response.APIResponse{
		ErrMsg:  errors.New(noti.TOO_MANY_REQUESTS_WARN_MSG),
		Context: ctx,
	}
}

func ProcessLoginResponse(data response.APIResponse) {
	if data.ErrMsg != nil {
		processFailResponse(data.ErrMsg, data.Context)
		return
	}

	var stringRes1 string = fmt.Sprint(data.Data1)
	var stringRes2 string = fmt.Sprint(data.Data2)

	var singleResponse response.MessageAPIResponse
	singleResponse.Message = stringRes2

	switch stringRes1 {
	case action_type.ACTIVATE_TYPE:
		data.Context.IndentedJSON(http.StatusContinue, singleResponse)
	case action_type.VERIFY_TYPE:
		data.Context.IndentedJSON(http.StatusContinue, singleResponse)
	case action_type.REDIRECT:
		processRedirectResponse(stringRes2, data.Context)
	default:
		data.Context.IndentedJSON(http.StatusOK, response.LoginSuccessResponse{
			AccessToken:  stringRes1,
			RefreshToken: stringRes2,
		})
	}
}

func ProcessRefreshTokenResponse(data response.APIResponse) {
	if data.ErrMsg != nil {
		if errMsg := data.ErrMsg.Error(); errMsg != "" {
			data.Context.IndentedJSON(http.StatusUnauthorized, response.MessageAPIResponse{
				Message: errMsg,
			})

			return
		}

	}
	data.Context.IndentedJSON(http.StatusOK, response.RefreshTokenSuccessResponse{
		AccessToken: fmt.Sprint(data.Data1),
	})
}

func processFailResponse(err error, ctx *gin.Context) {
	var errCode int

	switch err.Error() {
	case noti.INTERNALL_ERR_MSG:
		errCode = http.StatusInternalServerError
	case noti.GENERIC_RIGHT_ACCESS_WARN_MSG:
		errCode = http.StatusForbidden
	case noti.TOO_MANY_REQUESTS_WARN_MSG:
		errCode = http.StatusTooManyRequests
	default:
		errCode = http.StatusBadRequest
	}

	if isErrorTypeOfUndefined(err) {
		errCode = http.StatusNotFound
	}

	ctx.IndentedJSON(errCode, response.MessageAPIResponse{
		Message: err.Error(),
	})
}

func processSuccessPostReponse(res interface{}, postType string, ctx *gin.Context) {
	switch postType {
	case action_type.REDIRECT:
		processRedirectResponse(fmt.Sprint(res), ctx)
	case action_type.INFORM:
		processInformResponse(res, ctx)
	case action_type.CREATE_ACTION:
		processCreateResponse(res, ctx)
	default:
		ctx.IndentedJSON(http.StatusOK, response.MessageAPIResponse{
			Message: "success",
		})
	}
}

func processRedirectResponse(redirectUrl string, ctx *gin.Context) {
	ctx.Redirect(http.StatusSeeOther, redirectUrl)
}

func processInformResponse(message interface{}, ctx *gin.Context) {
	ctx.IndentedJSON(http.StatusOK, response.MessageAPIResponse{
		Message: fmt.Sprint(message),
	})
}

func processCreateResponse(data interface{}, ctx *gin.Context) {
	ctx.IndentedJSON(http.StatusCreated, data)
}

func processSuccessResponse(data interface{}, ctx *gin.Context) {
	ctx.IndentedJSON(http.StatusOK, data)
}

func isErrorTypeOfUndefined(err error) bool {
	return strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "undefined")
}
