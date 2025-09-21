package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID         int        `db:"id" json:"id"`
	QRCode     string     `db:"qr_code" json:"qr_code"`
	UserID     uuid.UUID  `db:"users_id" json:"user_id"`
	ScheduleID int        `db:"schedules_id" json:"schedule_id"`
	PaymentID  int        `db:"payments_id" json:"payment_id"`
	FullName   string     `db:"fullname" json:"fullname"`
	Email      string     `db:"email" json:"email"`
	Phone      string     `db:"phone_number" json:"phone"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt  *time.Time `db:"updated_at" json:"updated_at"`
	Seats      []Seat     `db:"-" json:"seats"`
}

type Schedule struct {
	ID         int       `db:"id" json:"id"`
	MovieID    int       `db:"movies_id" json:"movie_id"`
	CinemaID   int       `db:"cinemas_id" json:"cinema_id"`
	TimeID     int       `db:"times_id" json:"time_id"`
	LocationID int       `db:"locations_id" json:"location_id"`
	Date       time.Time `db:"date" json:"date"`
}

// komposite pk
type OrderSeat struct {
	OrderID int `db:"orders_id" json:"order_id"`
	SeatID  int `db:"seats_id" json:"seat_id"`
}

type Seat struct {
	ID       int    `db:"id" json:"id"`
	SeatCode string `db:"seat_code" json:"seat_code"`
}

type PaymentMethod struct {
	ID   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type Cinema struct {
	ID   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type Location struct {
	ID   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type Time struct {
	ID   int    `db:"id" json:"id"`
	Time string `db:"time" json:"time"`
}

type OrderDetail struct {
	Order
	Movie       Movie     `json:"movie"`
	CinemaName  string    `json:"cinema_name"`
	Location    string    `json:"location"`
	TimeStr     string    `json:"time"`
	Date        time.Time `json:"date"`
	PaymentName string    `json:"payment"`
}
