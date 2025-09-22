package routers

import (
	"net/http"

	docs "github.com/Darari17/be-tickitz-full/docs"
	"github.com/Darari17/be-tickitz-full/internal/dtos"
	"github.com/Darari17/be-tickitz-full/internal/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(db *pgxpool.Pool, rdb *redis.Client) *gin.Engine {
	router := gin.Default()
	router.Use(middlewares.CORSMiddleware)

	initAuthRouter(router, db)
	initMovieRouter(router, db, rdb)
	initOrderRouter(router, db)
	initAdminRoutes(router, db)

	router.Static("/img", "public")

	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, dtos.Response{
			Code:    http.StatusNotFound,
			Success: false,
			Message: "You wrong route",
			Data:    nil,
		})
	})

	return router
}
