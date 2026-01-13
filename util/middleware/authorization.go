package middleware

import (
	"raise-child/business"
	"raise-child/constants/shared"
	"raise-child/util"
	"raise-child/util/security"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Authorize(ctx *gin.Context) {
	var unAuthBodyResponse = util.GetUnAuthBodyResponse(ctx)
	var authHeader string = ctx.GetHeader("Authorization")
	if authHeader == "" {
		util.ProcessResponse(unAuthBodyResponse)
		ctx.Abort()
		return
	}

	var token string = strings.TrimPrefix(authHeader, "Bearer ")

	address, sub, exp, err := security.ExtractDataFromToken(token, util.GetLogConfig(shared.ERROR_LEVEL))
	if err != nil {
		util.ProcessResponse(unAuthBodyResponse)
		ctx.Abort()
		return
	}

	// Token expired
	if time.Now().After(exp) {
		util.ProcessResponse(unAuthBodyResponse)
		ctx.Abort()
		return
	}

	var wallets = business.GetWallets()
	if _, isExist := wallets[sub]; !isExist {
		util.ProcessResponse(unAuthBodyResponse)
		ctx.Abort()
		return
	}

	ctx.Set("address", address)
	ctx.Set("sub", sub)
	ctx.Next()
}
