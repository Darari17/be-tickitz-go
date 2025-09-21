package models

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type User struct {
	ID        uuid.UUID  `db:"id"`
	Email     string     `db:"email"`
	Password  string     `db:"password"`
	Role      Role       `db:"role"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	Profile   Profile    `db:"-"`
}

type Profile struct {
	UserID      uuid.UUID  `db:"user_id" json:"user_id"`
	FirstName   *string    `db:"firstname" json:"firstname,omitempty"`
	LastName    *string    `db:"lastname" json:"lastname,omitempty"`
	PhoneNumber *string    `db:"phone_number" json:"phone_number,omitempty"`
	Avatar      *string    `db:"avatar" json:"avatar,omitempty"`
	Point       *int       `db:"point" json:"point,omitempty"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at" json:"updated_at,omitempty"`
}

type UserContext struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Role  string    `json:"role"`
}
