package middleware

import (
	"net/http"
	"raise-child/business"
	"raise-child/constants/shared"
	"raise-child/util"
	"raise-child/util/security"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	Authorize                  = skipOnPreflight(authorize)
	AdminAuthorize             = skipOnPreflight(adminAuthorize)
	LeaderAuthorize            = skipOnPreflight(leaderAuthorize)
	ManagerRoleAuthorize       = skipOnPreflight(managerRoleAuthorize)
	StaffRoleAuthorize         = skipOnPreflight(staffRoleAuthorize)
	VolunteerRoleAuthorize     = skipOnPreflight(volunteerRoleAuthorize)
	GotRolesAuthorizeAuthorize = skipOnPreflight(gotRolesAuthorize)
	OptinalAuthorize           = skipOnPreflight(optinalAuthorize)
)

func skipOnPreflight(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Method == http.MethodOptions {
			ctx.Next()
			return
		}
		handler(ctx)
	}
}

func authorize(ctx *gin.Context) {
	var unAuthBodyResponse = util.GetUnAuthBodyResponse(ctx)
	var authHeader string = ctx.GetHeader("Authorization")
	if authHeader == "" {
		util.ProcessResponse(unAuthBodyResponse)
		ctx.Abort()
		return
	}

	//var token string = strings.TrimPrefix(authHeader, "Bearer ")

	address, sub, roles, exp, err := security.ExtractDataFromTokenV2(authHeader, util.GetLogConfig(shared.ERROR_LEVEL))
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

	if !business.IsWalletRegistered(sub) {
		util.ProcessResponse(unAuthBodyResponse)
		ctx.Abort()
		return
	}

	if !util.IsValidSuiAddressStrict(address) {
		util.ProcessResponse(unAuthBodyResponse)
		ctx.Abort()
		return
	}

	if len(roles) == 0 {
		util.ProcessResponse(unAuthBodyResponse)
		ctx.Abort()
		return
	}

	ctx.Set("address", address)
	ctx.Set("sub", sub)
	ctx.Set("roles", roles)
	ctx.Next()
}

func optinalAuthorize(ctx *gin.Context) {
	var authHeader string = ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.Next()
		return
	}

	//var token string = strings.TrimPrefix(authHeader, "Bearer ")

	address, sub, roles, exp, err := security.ExtractDataFromTokenV2(authHeader, util.GetLogConfig(shared.ERROR_LEVEL))
	if err != nil {
		ctx.Next()
		return
	}

	// Token expired
	if time.Now().After(exp) {
		ctx.Next()
		return
	}

	if !business.IsWalletRegistered(sub) {
		ctx.Next()
		return
	}

	if !util.IsValidSuiAddressStrict(address) {
		ctx.Next()
		return
	}

	if len(roles) == 0 {
		ctx.Next()
		return
	}

	ctx.Set("address", address)
	ctx.Set("sub", sub)
	ctx.Set("roles", roles)
	ctx.Next()
}

func adminAuthorize(ctx *gin.Context) {
	val, exists := ctx.Get("roles")
	if !exists {
		util.ProcessResponse(util.GetUnAuthBodyResponse(ctx))
		ctx.Abort()
		return
	}

	roles, ok := val.([]string)
	if !ok {
		util.ProcessResponse(util.GetUnAuthBodyResponse(ctx))
		ctx.Abort()
		return
	}

	if !slices.Contains(roles, "Admin") {
		util.ProcessResponse(util.GetUnAuthBodyResponse(ctx))
		ctx.Abort()
		return
	}

	ctx.Next()
}

func leaderAuthorize(ctx *gin.Context) {
	val, exists := ctx.Get("roles")
	if !exists {
		util.ProcessResponse(util.GetUnAuthBodyResponse(ctx))
		ctx.Abort()
		return
	}

	roles, ok := val.([]string)
	if !ok {
		util.ProcessResponse(util.GetUnAuthBodyResponse(ctx))
		ctx.Abort()
		return
	}

	if !slices.Contains(roles, "Local Leader") {
		util.ProcessResponse(util.GetUnAuthBodyResponse(ctx))
		ctx.Abort()
		return
	}

	ctx.Next()
}

func managerRoleAuthorize(ctx *gin.Context) {
	val, exists := ctx.Get("roles")
	if !exists {
		util.ProcessResponse(util.GetUnAuthBodyResponse(ctx))
		ctx.Abort()
		return
	}

	roles, ok := val.([]string)
	if !ok {
		util.ProcessResponse(util.GetUnAuthBodyResponse(ctx))
		ctx.Abort()
		return
	}

	if !slices.Contains(roles, "Admin") && !slices.Contains(roles, "Local Leader") {
		util.ProcessResponse(util.GetUnAuthBodyResponse(ctx))
		ctx.Abort()
		return
	}

	ctx.Next()
}

func staffRoleAuthorize(ctx *gin.Context) {
	val, exists := ctx.Get("roles")
	if !exists {
		util.ProcessResponse(util.GetUnAuthBodyResponse(ctx))
		ctx.Abort()
		return
	}

	roles, ok := val.([]string)
	if !ok {
		util.ProcessResponse(util.GetUnAuthBodyResponse(ctx))
		ctx.Abort()
		return
	}

	if !slices.Contains(roles, "Volunteer") && !slices.Contains(roles, "Local Leader") {
		util.ProcessResponse(util.GetUnAuthBodyResponse(ctx))
		ctx.Abort()
		return
	}

	ctx.Next()
}

func volunteerRoleAuthorize(ctx *gin.Context) {
	val, exists := ctx.Get("roles")
	if !exists {
		util.ProcessResponse(util.GetUnAuthBodyResponse(ctx))
		ctx.Abort()
		return
	}

	roles, ok := val.([]string)
	if !ok {
		util.ProcessResponse(util.GetUnAuthBodyResponse(ctx))
		ctx.Abort()
		return
	}

	if !slices.Contains(roles, "Volunteer") {
		util.ProcessResponse(util.GetUnAuthBodyResponse(ctx))
		ctx.Abort()
		return
	}

	ctx.Next()
}

func gotRolesAuthorize(ctx *gin.Context) {
	val, exists := ctx.Get("roles")
	if !exists {
		util.ProcessResponse(util.GetUnAuthBodyResponse(ctx))
		ctx.Abort()
		return
	}

	roles, ok := val.([]string)
	if !ok {
		util.ProcessResponse(util.GetUnAuthBodyResponse(ctx))
		ctx.Abort()
		return
	}

	if !slices.Contains(roles, "Admin") && !slices.Contains(roles, "Local Leader") && !slices.Contains(roles, "Volunteer") {
		util.ProcessResponse(util.GetUnAuthBodyResponse(ctx))
		ctx.Abort()
		return
	}

	ctx.Next()
}
