package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Darari17/be-tickitz-full/internal/dtos"
	"github.com/Darari17/be-tickitz-full/internal/models"
	"github.com/Darari17/be-tickitz-full/internal/repositories"
	"github.com/Darari17/be-tickitz-full/internal/utils"
	"github.com/gin-gonic/gin"
)

type OrderController struct {
	orderRepo *repositories.OrderRepo
}

func NewOrderController(or *repositories.OrderRepo) *OrderController {
	return &OrderController{orderRepo: or}
}

//
// ========================
// Orders
// ========================
//

// CreateOrder godoc
// @Summary Create a new order
// @Description Create a new order and save to database
// @Tags Orders
// @Accept json
// @Produce json
// @Param order body dtos.CreateOrderRequest true "Order Data"
// @Success 201 {object} dtos.Response
// @Failure 400 {object} dtos.Response
// @Failure 401 {object} dtos.Response
// @Failure 500 {object} dtos.Response
// @Router /orders [post]
// CreateOrder godoc
func (oc *OrderController) CreateOrder(ctx *gin.Context) {
	user, err := utils.GetUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dtos.Response{
			Code:    http.StatusUnauthorized,
			Success: false,
			Message: "Unauthorized: " + err.Error(),
		})
		return
	}

	var req dtos.CreateOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	// 🔑 Konversi seat codes -> seat IDs
	seatIDs, err := oc.orderRepo.GetSeatIDsByCodes(ctx.Request.Context(), req.SeatCodes)
	if err != nil {
		log.Println("GetSeatIDsByCodes error:", err)
		ctx.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to map seat codes",
		})
		return
	}

	order := &models.Order{
		UserID:     user.ID,
		ScheduleID: req.ScheduleID,
		PaymentID:  req.PaymentID,
		FullName:   req.FullName,
		Email:      req.Email,
		Phone:      req.Phone,
	}

	createdOrder, err := oc.orderRepo.CreateOrder(ctx.Request.Context(), order, seatIDs)
	if err != nil {
		log.Println("CreateOrder error:", err)
		ctx.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to create order",
		})
		return
	}

	ctx.JSON(http.StatusCreated, dtos.Response{
		Code:    http.StatusCreated,
		Success: true,
		Data: map[string]interface{}{
			"order_id": createdOrder.ID,
			"qr_code":  createdOrder.QRCode,
			"seats":    createdOrder.Seats,
		},
	})
}

// GetSchedules godoc
// @Summary Get schedules by movie ID
// @Description Retrieve all schedules for a movie
// @Tags Orders
// @Produce json
// @Param movie_id query int true "Movie ID"
// @Success 200 {object} dtos.Response
// @Failure 400 {object} dtos.Response
// @Failure 500 {object} dtos.Response
// @Router /orders/schedules [get]
func (oc *OrderController) GetSchedules(ctx *gin.Context) {
	movieIDStr := ctx.Query("movie_id")
	if movieIDStr == "" {
		ctx.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "movie_id is required",
		})
		return
	}

	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "invalid movie_id",
		})
		return
	}

	schedules, err := oc.orderRepo.GetSchedules(ctx.Request.Context(), movieID)
	if err != nil {
		log.Println("GetSchedules error:", err)
		ctx.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to fetch schedules",
		})
		return
	}

	if schedules == nil {
		schedules = []models.Schedule{}
	}

	ctx.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    schedules,
	})
}

// GetAvailableSeats godoc
// @Summary Get available seats
// @Description Retrieve available seats for a schedule
// @Tags Orders
// @Produce json
// @Param schedule_id query int true "Schedule ID"
// @Success 200 {object} dtos.Response
// @Failure 400 {object} dtos.Response
// @Failure 500 {object} dtos.Response
// @Router /orders/seats [get]
func (oc *OrderController) GetAvailableSeats(ctx *gin.Context) {
	scheduleIDStr := ctx.Query("schedule_id")
	if scheduleIDStr == "" {
		ctx.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "schedule_id is required",
		})
		return
	}

	scheduleID, err := strconv.Atoi(scheduleIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "invalid schedule_id",
		})
		return
	}

	seats, err := oc.orderRepo.GetAvailableSeats(ctx.Request.Context(), scheduleID)
	if err != nil {
		log.Println("GetAvailableSeats error:", err)
		ctx.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to fetch seats",
		})
		return
	}

	ctx.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    seats,
	})
}

// GetTransactionDetail godoc
// @Summary Get transaction detail
// @Description Retrieve detail of a transaction
// @Tags Orders
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} dtos.Response
// @Failure 400 {object} dtos.Response
// @Failure 500 {object} dtos.Response
// @Router /orders/{id} [get]
func (oc *OrderController) GetTransactionDetail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	orderID, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "invalid order id",
		})
		return
	}

	detail, err := oc.orderRepo.GetTransactionDetail(ctx.Request.Context(), orderID)
	if err != nil {
		log.Println("GetTransactionDetail error:", err)
		ctx.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to fetch order detail",
		})
		return
	}

	ctx.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    detail,
	})
}

// GetOrderHistory godoc
// @Summary Get order history
// @Description Retrieve order history for current user
// @Tags Orders
// @Produce json
// @Success 200 {object} dtos.Response
// @Failure 401 {object} dtos.Response
// @Failure 500 {object} dtos.Response
// @Router /orders/history [get]
func (oc *OrderController) GetOrderHistory(ctx *gin.Context) {
	user, err := utils.GetUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dtos.Response{
			Code:    http.StatusUnauthorized,
			Success: false,
			Message: "Unauthorized: " + err.Error(),
		})
		return
	}

	history, err := oc.orderRepo.GetOrderHistory(ctx.Request.Context(), user.ID)
	if err != nil {
		log.Println("GetOrderHistory error:", err)
		ctx.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to fetch history",
		})
		return
	}

	ctx.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    history,
	})
}

// GetPayments godoc
// @Summary Get all payment methods
// @Description Retrieve list of available payment methods
// @Tags Orders
// @Produce json
// @Success 200 {object} dtos.Response
// @Failure 500 {object} dtos.Response
// @Router /orders/payments [get]
func (oc *OrderController) GetPayments(ctx *gin.Context) {
	payments, err := oc.orderRepo.GetPayments(ctx.Request.Context())
	if err != nil {
		log.Println("GetPayments error:", err)
		ctx.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to fetch payments",
		})
		return
	}

	ctx.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    payments,
	})
}

// GetCinemas godoc
// @Summary Get all cinemas
// @Description Retrieve list of available cinemas
// @Tags Orders
// @Produce json
// @Success 200 {object} dtos.Response
// @Failure 500 {object} dtos.Response
// @Router /orders/cinemas [get]
func (oc *OrderController) GetCinemas(ctx *gin.Context) {
	cinemas, err := oc.orderRepo.GetCinemas(ctx.Request.Context())
	if err != nil {
		log.Println("GetCinemas error:", err)
		ctx.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to fetch cinemas",
		})
		return
	}

	ctx.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    cinemas,
	})
}

// GetLocations godoc
// @Summary Get all locations
// @Description Retrieve list of available locations
// @Tags Orders
// @Produce json
// @Success 200 {object} dtos.Response
// @Failure 500 {object} dtos.Response
// @Router /orders/locations [get]
func (oc *OrderController) GetLocations(ctx *gin.Context) {
	locations, err := oc.orderRepo.GetLocations(ctx.Request.Context())
	if err != nil {
		log.Println("GetLocations error:", err)
		ctx.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to fetch locations",
		})
		return
	}

	ctx.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    locations,
	})
}

// GetTimes godoc
// @Summary Get all times
// @Description Retrieve list of available movie times
// @Tags Orders
// @Produce json
// @Success 200 {object} dtos.Response
// @Failure 500 {object} dtos.Response
// @Router /orders/times [get]
func (oc *OrderController) GetTimes(ctx *gin.Context) {
	times, err := oc.orderRepo.GetTimes(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to fetch times",
		})
		return
	}

	ctx.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    times,
	})
}
