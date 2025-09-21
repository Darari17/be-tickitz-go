package routers

import (
	"github.com/Darari17/be-tickitz-full/internal/controllers"
	"github.com/Darari17/be-tickitz-full/internal/middlewares"
	"github.com/Darari17/be-tickitz-full/internal/repositories"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func initAuthRouter(router *gin.Engine, db *pgxpool.Pool) {
	authRepo := repositories.NewUserRepository(db)
	authHandler := controllers.NewUserController(authRepo)

	auth := router.Group("/auth")
	auth.POST("/login", authHandler.Login)
	auth.POST("/register", authHandler.Register)

	profile := router.Group("/profile", middlewares.RequiredToken, middlewares.Access("user"))
	profile.GET("", authHandler.GetProfile)
	profile.PATCH("", authHandler.UpdateProfile)
	profile.PATCH("/change-password", authHandler.ChangePassword)
	profile.PATCH("/change-avatar", authHandler.ChangeAvatar)
}
