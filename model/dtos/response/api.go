package response

import "github.com/gin-gonic/gin"

type APIResponse struct {
	Data1    interface{}
	Data2    interface{}
	ErrMsg   error
	Context  *gin.Context
	PostType string
}

type LoginSuccessResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenSuccessResponse struct {
	AccessToken string `json:"access_token"`
}

type MessageAPIResponse struct {
	Message string `json:"message"`
}

type UrlAPIResponse struct {
	Url string `json:"url"`
}
