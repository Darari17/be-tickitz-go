package dtos

import (
	"mime/multipart"
	"time"
)

type ScheduleRequest struct {
	CinemaID   int    `json:"cinema_id" form:"cinema_id" example:"2"`
	LocationID int    `json:"location_id" form:"location_id" example:"1"`
	Date       string `json:"date" form:"date" example:"2025-12-01"`
	TimeIDs    []int  `json:"time_ids" form:"time_ids" example:"1"`
}

type CreateMovieRequest struct {
	Title       string                `json:"title" form:"title" binding:"required" example:"Spider-Man: Homecoming"`
	Overview    string                `json:"overview" form:"overview" example:"Film tentang pahlawan..."`
	Director    string                `json:"director_name" form:"director_name" example:"Jon Watts"`
	Duration    int                   `json:"duration" form:"duration" example:"135"`
	ReleaseDate time.Time             `json:"release_date" form:"release_date" time_format:"2006-01-02" example:"2025-12-20"`
	Popularity  float64               `json:"popularity" form:"popularity" example:"87.5"`
	Poster      *multipart.FileHeader `form:"poster"`
	Backdrop    *multipart.FileHeader `form:"backdrop"`
	Genres      []int                 `json:"genres" form:"genres" example:"[1,2]"`
	Casts       []int                 `json:"casts" form:"casts" example:"[3,5,7]"`
	Schedules   []ScheduleRequest     `json:"schedules" form:"schedules"`
}

type UpdateMovieRequest struct {
	Title       *string               `json:"title" form:"title" example:"Avengers: Endgame"`
	Overview    *string               `json:"overview" form:"overview" example:"Film aksi penuh petualangan"`
	Director    *string               `json:"director_name" form:"director_name" example:"Anthony Russo, Joe Russo"`
	Duration    *int                  `json:"duration" form:"duration" example:"150"`
	ReleaseDate *time.Time            `json:"release_date" form:"release_date" time_format:"2006-01-02" example:"2025-12-20"`
	Popularity  *float64              `json:"popularity" form:"popularity" example:"90.2"`
	Poster      *multipart.FileHeader `form:"poster"`
	Backdrop    *multipart.FileHeader `form:"backdrop"`
	Genres      []int                 `json:"genres" form:"genres" example:"[1,2]"`
	Casts       []int                 `json:"casts" form:"casts" example:"[3,5,7]"`
	Schedules   []ScheduleRequest     `json:"schedules" form:"schedules"`
}
