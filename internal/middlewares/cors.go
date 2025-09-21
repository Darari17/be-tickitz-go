package middlewares

import (
	"log"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware(ctx *gin.Context) {
	whitelist := []string{
		"http://127.0.0.1:5500",
		"http://127.0.0.1:3001",
		"http://localhost:5173",
	}
	origin := ctx.GetHeader("Origin")

	if slices.Contains(whitelist, origin) {
		ctx.Header("Access-Control-Allow-Origin", origin)
	} else {
		log.Printf("Origin is not in the whitelist: %s", origin)
	}

	ctx.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
	ctx.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept, Origin, X-Requested-With")
	ctx.Header("Access-Control-Allow-Credentials", "true")

	if ctx.Request.Method == http.MethodOptions {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}

	ctx.Next()
}
