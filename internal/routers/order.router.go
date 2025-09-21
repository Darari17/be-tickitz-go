package routers

import (
	"github.com/Darari17/be-tickitz-full/internal/controllers"
	"github.com/Darari17/be-tickitz-full/internal/middlewares"
	"github.com/Darari17/be-tickitz-full/internal/repositories"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func initOrderRouter(router *gin.Engine, db *pgxpool.Pool) {
	orderRepo := repositories.NewOrderRepo(db)
	orderController := controllers.NewOrderController(orderRepo)

	orderGroup := router.Group("/orders", middlewares.RequiredToken, middlewares.Access("user"))
	orderGroup.POST("", orderController.CreateOrder)
	orderGroup.GET("/history", orderController.GetOrderHistory)
	orderGroup.GET("/schedules", orderController.GetSchedules)
	orderGroup.GET("/seats", orderController.GetAvailableSeats)
	orderGroup.GET("/:id", orderController.GetTransactionDetail)

	orderGroup.GET("/payments", orderController.GetPayments)
	orderGroup.GET("/cinemas", orderController.GetCinemas)
	orderGroup.GET("/locations", orderController.GetLocations)
	orderGroup.GET("/times", orderController.GetTimes)

}
