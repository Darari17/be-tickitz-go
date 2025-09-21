package dtos

import (
	"mime/multipart"

	"github.com/google/uuid"
)

type UserRequest struct {
	Email    string `form:"email" json:"email" binding:"required,email" example:"user@mail.com"`
	Password string `form:"password" json:"password" binding:"required" example:"Password123"`
}

type ProfileRequest struct {
	FirstName   *string `form:"firstname" json:"firstname" example:"Farid"`
	LastName    *string `form:"lastname" json:"lastname" example:"Darari"`
	PhoneNumber *string `form:"phone_number" json:"phone_number"`
}

type AvatarProfileRequest struct {
	Avatar *multipart.FileHeader `form:"avatar"`
}

type ChangePasswordRequest struct {
	OldPassword string `form:"old_password" json:"old_password" binding:"required"`
	NewPassword string `form:"new_password" json:"new_password" binding:"required"`
}

type UserResponse struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
	Token  string    `json:"token"`
}

type ProfileResponse struct {
	UserID      uuid.UUID `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	FirstName   *string   `json:"firstname" example:"Farid"`
	LastName    *string   `json:"lastname" example:"Darari"`
	PhoneNumber *string   `json:"phone_number" example:"08123456789"`
	Avatar      *string   `json:"avatar" example:"https://example.com/avatar.png"`
	Point       *int      `json:"point" example:"100"`
}
