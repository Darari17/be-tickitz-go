package routers

import (
	"github.com/Darari17/be-tickitz-full/internal/controllers"
	"github.com/Darari17/be-tickitz-full/internal/middlewares"
	"github.com/Darari17/be-tickitz-full/internal/repositories"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func AdminRoutes(r *gin.Engine, db *pgxpool.Pool) {
	adminRepo := repositories.NewAdminRepository(db)
	adminCtrl := controllers.NewAdminController(adminRepo)

	admin := r.Group("/admin", middlewares.RequiredToken, middlewares.Access("admin"))
	{
		admin.POST("/movies", adminCtrl.CreateMovie)
		admin.GET("/movies", adminCtrl.GetMovies)
		admin.GET("/movies/:id", adminCtrl.GetMovieByID)
		admin.PATCH("/movies/:id", adminCtrl.UpdateMovie)
		admin.DELETE("/movies/:id", adminCtrl.DeleteMovie)
	}
}
