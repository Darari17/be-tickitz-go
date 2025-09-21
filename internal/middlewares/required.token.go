package middlewares

import (
	"log"
	"net/http"
	"strings"

	"github.com/Darari17/be-tickitz-full/internal/dtos"
	"github.com/Darari17/be-tickitz-full/pkg"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequiredToken(ctx *gin.Context) {
	bearerToken := ctx.GetHeader("Authorization")
	if bearerToken == "" || !strings.HasPrefix(bearerToken, "Bearer ") {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, dtos.Response{
			Code:    http.StatusUnauthorized,
			Success: false,
			Message: "Authentication required",
		})
		return
	}

	token := strings.TrimPrefix(bearerToken, "Bearer ")

	// cek token redis
	// ...

	claims := &pkg.Claims{}

	if err := claims.VerifyToken(token); err != nil {
		if strings.Contains(err.Error(), jwt.ErrTokenInvalidIssuer.Error()) {
			log.Println("JWT Error.\nCause: ", err.Error())
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dtos.Response{
				Code:    http.StatusUnauthorized,
				Success: false,
				Message: "Please log in again",
			})
			return
		}

		if strings.Contains(err.Error(), jwt.ErrTokenExpired.Error()) {
			log.Println("JWT Error.\nCause: ", err.Error())
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dtos.Response{
				Code:    http.StatusUnauthorized,
				Success: false,
				Message: "Please log in again",
			})
			return
		}

		log.Println("Internal Server Error.\nCause: ", err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Internal Server Error",
		})
		return
	}

	ctx.Set("claims", claims)
	ctx.Next()

}
