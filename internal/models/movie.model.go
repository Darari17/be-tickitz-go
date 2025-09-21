package models

import (
	"time"
)

type Movie struct {
	ID          int        `db:"id" json:"id"`
	Backdrop    string     `db:"backdrop_path" json:"backdrop_path"`
	Overview    string     `db:"overview" json:"overview"`
	Popularity  float64    `db:"popularity" json:"popularity"`
	Poster      string     `db:"poster_path" json:"poster_path"`
	ReleaseDate time.Time  `db:"release_date" json:"release_date"`
	Duration    int        `db:"duration" json:"duration"`
	Title       string     `db:"title" json:"title"`
	Director    string     `db:"director_name" json:"director_name"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deleted_at"`
	Genres      []Genre    `db:"-" json:"genres"`
	Casts       []Cast     `db:"-" json:"casts"`
}

type Cast struct {
	ID   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type MovieCast struct {
	ID      int `db:"id" json:"id"`
	MovieID int `db:"movies_id" json:"movies_id"`
	CastID  int `db:"casts_id" json:"casts_id"`
}

type Genre struct {
	ID   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type MovieGenre struct {
	ID      int `db:"id" json:"id"`
	MovieID int `db:"movies_id" json:"movies_id"`
	GenreID int `db:"genres_id" json:"genres_id"`
}
