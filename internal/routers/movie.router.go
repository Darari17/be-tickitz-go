package routers

import (
	"github.com/Darari17/be-tickitz-full/internal/controllers"
	"github.com/Darari17/be-tickitz-full/internal/repositories"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func initMovieRouter(router *gin.Engine, db *pgxpool.Pool, redis *redis.Client) {
	movieRepo := repositories.NewMovieRepository(db, redis)
	movieHandler := controllers.NewMovieController(movieRepo)

	movies := router.Group("/movies")
	movies.GET("/upcoming", movieHandler.GetUpcomingMovies)
	movies.GET("/popular", movieHandler.GetPopularMovies)
	movies.GET("", movieHandler.GetAllMovies)
	movies.GET("/:id", movieHandler.GetMovieDetail)
	movies.GET("/genres", movieHandler.GetGenres)
}
