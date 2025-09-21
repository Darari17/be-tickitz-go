package middlewares

import (
	"net/http"
	"slices"

	"github.com/Darari17/be-tickitz-full/internal/dtos"
	"github.com/Darari17/be-tickitz-full/pkg"
	"github.com/gin-gonic/gin"
)

func Access(roles ...string) func(*gin.Context) {
	return func(ctx *gin.Context) {
		claims, isExist := ctx.Get("claims")
		if !isExist {
			ctx.AbortWithStatusJSON(http.StatusForbidden, dtos.Response{
				Code:    http.StatusForbidden,
				Success: false,
				Message: "Please log in again",
			})
			return
		}

		user, ok := claims.(*pkg.Claims)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, dtos.Response{
				Code:    http.StatusInternalServerError,
				Success: false,
				Message: "Internal Server Error",
			})
			return
		}

		if !slices.Contains(roles, user.Role) {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, dtos.Response{
				Code:    http.StatusInternalServerError,
				Success: false,
				Message: "You are not authorized to access this resource",
			})
			return
		}

		ctx.Next()
	}
}
