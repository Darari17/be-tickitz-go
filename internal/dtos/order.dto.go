package dtos

type CreateOrderRequest struct {
	ScheduleID int      `json:"schedule_id" binding:"required" example:"8"`
	PaymentID  int      `json:"payment_id" binding:"required" example:"2"`
	FullName   string   `json:"fullname" binding:"required" example:"Farid Darari"`
	Email      string   `json:"email" binding:"required,email" example:"farid@example.com"`
	Phone      string   `json:"phone" binding:"required" example:"+628123456789"`
	SeatCodes  []string `json:"seat_codes" binding:"required,min=1" example:"[\"A1\",\"A2\"]"`
}
