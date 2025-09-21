package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Darari17/be-tickitz-full/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepo struct {
	db *pgxpool.Pool
}

func NewOrderRepo(db *pgxpool.Pool) *OrderRepo {
	return &OrderRepo{db: db}
}

//
// ========================
// Orders
// ========================
//

func (or *OrderRepo) CreateOrder(ctx context.Context, order *models.Order, seatIDs []int) (*models.Order, error) {
	tx, err := or.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// hardcode QR Code jika kosong
	if order.QRCode == "" {
		order.QRCode = fmt.Sprintf("QR-%d", time.Now().Unix())
	}

	query := `
        INSERT INTO orders (qr_code, users_id, schedules_id, payments_id, fullname, email, phone_number, created_at)
        VALUES ($1,$2,$3,$4,$5,$6,$7,NOW())
        RETURNING id, created_at
    `
	err = tx.QueryRow(ctx, query,
		order.QRCode, order.UserID, order.ScheduleID, order.PaymentID,
		order.FullName, order.Email, order.Phone,
	).Scan(&order.ID, &order.CreatedAt)
	if err != nil {
		return nil, err
	}

	for _, seatID := range seatIDs {
		_, err := tx.Exec(ctx, `INSERT INTO order_seats (orders_id, seats_id) VALUES ($1,$2)`, order.ID, seatID)
		if err != nil {
			return nil, err
		}
	}

	rows, err := tx.Query(ctx, `SELECT id, seat_code FROM seats WHERE id = ANY($1)`, seatIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var seat models.Seat
		if err := rows.Scan(&seat.ID, &seat.SeatCode); err != nil {
			return nil, err
		}
		order.Seats = append(order.Seats, seat)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return order, nil
}

func (or *OrderRepo) GetSeatIDsByCodes(ctx context.Context, seatCodes []string) ([]int, error) {
	rows, err := or.db.Query(ctx, `SELECT id FROM seats WHERE seat_code = ANY($1)`, seatCodes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (or *OrderRepo) GetSchedules(ctx context.Context, movieID int) ([]models.Schedule, error) {
	rows, err := or.db.Query(ctx, `
		SELECT id, movies_id, cinemas_id, times_id, locations_id, date
FROM schedules 
WHERE movies_id=$1


	`, movieID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.Schedule
	for rows.Next() {
		var s models.Schedule
		if err := rows.Scan(&s.ID, &s.MovieID, &s.CinemaID, &s.TimeID, &s.LocationID, &s.Date); err != nil {
			return nil, err
		}
		schedules = append(schedules, s)
	}
	return schedules, nil
}

func (or *OrderRepo) GetAvailableSeats(ctx context.Context, scheduleID int) ([]models.Seat, error) {
	rows, err := or.db.Query(ctx, `
		SELECT s.id, s.seat_code
		FROM seats s
		WHERE s.id NOT IN (
			SELECT os.seats_id
			FROM orders o
			JOIN order_seats os ON o.id = os.orders_id
			WHERE o.schedules_id = $1
		)
		ORDER BY s.id
	`, scheduleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seats []models.Seat
	for rows.Next() {
		var seat models.Seat
		if err := rows.Scan(&seat.ID, &seat.SeatCode); err != nil {
			return nil, err
		}
		seats = append(seats, seat)
	}
	return seats, nil
}

func (or *OrderRepo) GetTransactionDetail(ctx context.Context, orderID int) (*models.OrderDetail, error) {
	sql := `
		SELECT o.id, o.qr_code, o.users_id, o.schedules_id, o.payments_id,
		       o.fullname, o.email, o.phone_number, o.created_at, o.updated_at,
		       m.id, m.backdrop_path, m.overview, m.popularity, m.poster_path,
		       m.release_date, m.duration, m.title, m.director_name,
		       c.name as cinema_name, l.name as location, t.time, s.date,
		       pm.name as payment,
		       COALESCE(json_agg(json_build_object('id', se.id, 'seat_code', se.seat_code))
		                FILTER (WHERE se.id IS NOT NULL), '[]') as seats
		FROM orders o
		JOIN schedules s ON o.schedules_id = s.id
		JOIN movies m ON s.movies_id = m.id
		JOIN cinemas c ON s.cinemas_id = c.id
		JOIN locations l ON s.locations_id = l.id
		JOIN times t ON s.times_id = t.id
		JOIN payment_methods pm ON o.payments_id = pm.id
		LEFT JOIN order_seats os ON o.id = os.orders_id
		LEFT JOIN seats se ON se.id = os.seats_id
		WHERE o.id = $1
		GROUP BY o.id, m.id, c.name, l.name, t.time, s.date, pm.name
	`

	var d models.OrderDetail
	var seatsJSON []byte

	err := or.db.QueryRow(ctx, sql, orderID).Scan(
		&d.ID, &d.QRCode, &d.UserID, &d.ScheduleID, &d.PaymentID,
		&d.FullName, &d.Email, &d.Phone, &d.CreatedAt, &d.UpdatedAt,
		&d.Movie.ID, &d.Movie.Backdrop, &d.Movie.Overview, &d.Movie.Popularity,
		&d.Movie.Poster, &d.Movie.ReleaseDate, &d.Movie.Duration,
		&d.Movie.Title, &d.Movie.Director,
		&d.CinemaName, &d.Location, &d.TimeStr, &d.Date,
		&d.PaymentName,
		&seatsJSON,
	)
	if err != nil {
		return nil, err
	}

	_ = json.Unmarshal(seatsJSON, &d.Seats)
	return &d, nil
}

func (or *OrderRepo) GetOrderHistory(ctx context.Context, userID uuid.UUID) ([]models.OrderDetail, error) {
	rows, err := or.db.Query(ctx, `
		SELECT o.id, o.qr_code, o.users_id, o.schedules_id, o.payments_id,
		       o.fullname, o.email, o.phone_number, o.created_at, o.updated_at,
		       m.id, m.backdrop_path, m.overview, m.popularity, m.poster_path,
		       m.release_date, m.duration, m.title, m.director_name,
		       c.name as cinema_name, l.name as location, t.time, s.date,
		       pm.name as payment,
		       COALESCE(json_agg(json_build_object('id', se.id, 'seat_code', se.seat_code))
		                FILTER (WHERE se.id IS NOT NULL), '[]') as seats
		FROM orders o
		JOIN schedules s ON o.schedules_id = s.id
		JOIN movies m ON s.movies_id = m.id
		JOIN cinemas c ON s.cinemas_id = c.id
		JOIN locations l ON s.locations_id = l.id
		JOIN times t ON s.times_id = t.id
		JOIN payment_methods pm ON o.payments_id = pm.id
		LEFT JOIN order_seats os ON o.id = os.orders_id
		LEFT JOIN seats se ON se.id = os.seats_id
		WHERE o.users_id = $1
		GROUP BY o.id, m.id, c.name, l.name, t.time, s.date, pm.name
		ORDER BY o.created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.OrderDetail
	for rows.Next() {
		var d models.OrderDetail
		var seatsJSON []byte

		if err := rows.Scan(
			&d.ID, &d.QRCode, &d.UserID, &d.ScheduleID, &d.PaymentID,
			&d.FullName, &d.Email, &d.Phone, &d.CreatedAt, &d.UpdatedAt,
			&d.Movie.ID, &d.Movie.Backdrop, &d.Movie.Overview, &d.Movie.Popularity,
			&d.Movie.Poster, &d.Movie.ReleaseDate, &d.Movie.Duration,
			&d.Movie.Title, &d.Movie.Director,
			&d.CinemaName, &d.Location, &d.TimeStr, &d.Date,
			&d.PaymentName,
			&seatsJSON,
		); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(seatsJSON, &d.Seats)
		orders = append(orders, d)
	}
	return orders, nil
}

//
// ========================
// Payments
// ========================
//

func (or *OrderRepo) GetPayments(ctx context.Context) ([]models.PaymentMethod, error) {
	rows, err := or.db.Query(ctx, `SELECT id, name FROM payment_methods ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []models.PaymentMethod
	for rows.Next() {
		var p models.PaymentMethod
		if err := rows.Scan(&p.ID, &p.Name); err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}
	return payments, nil
}

//
// ========================
// Cinemas
// ========================
//

func (or *OrderRepo) GetCinemas(ctx context.Context) ([]models.Cinema, error) {
	rows, err := or.db.Query(ctx, `SELECT id, name FROM cinemas ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cinemas []models.Cinema
	for rows.Next() {
		var c models.Cinema
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		cinemas = append(cinemas, c)
	}
	return cinemas, nil
}

//
// ========================
// Locations
// ========================
//

func (or *OrderRepo) GetLocations(ctx context.Context) ([]models.Location, error) {
	rows, err := or.db.Query(ctx, `SELECT id, name FROM locations ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []models.Location
	for rows.Next() {
		var l models.Location
		if err := rows.Scan(&l.ID, &l.Name); err != nil {
			return nil, err
		}
		locations = append(locations, l)
	}
	return locations, nil
}

func (or *OrderRepo) GetTimes(ctx context.Context) ([]models.Time, error) {
	rows, err := or.db.Query(ctx, `SELECT id, time FROM times ORDER BY time`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var times []models.Time
	for rows.Next() {
		var t models.Time
		if err := rows.Scan(&t.ID, &t.Time); err != nil {
			return nil, err
		}
		times = append(times, t)
	}
	return times, nil
}
