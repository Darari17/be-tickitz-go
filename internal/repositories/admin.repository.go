package repositories

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Darari17/be-tickitz-full/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminRepository struct {
	db *pgxpool.Pool
}

func NewAdminRepository(db *pgxpool.Pool) *AdminRepository {
	return &AdminRepository{
		db: db,
	}
}

func (r *AdminRepository) GetMovies(ctx context.Context) ([]models.Movie, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, backdrop_path, overview, popularity, poster_path,
		       release_date, duration, title, director_name,
		       created_at, updated_at, deleted_at
		FROM movies
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie
		if err := rows.Scan(
			&m.ID, &m.Backdrop, &m.Overview, &m.Popularity, &m.Poster,
			&m.ReleaseDate, &m.Duration, &m.Title, &m.Director,
			&m.CreatedAt, &m.UpdatedAt, &m.DeletedAt,
		); err != nil {
			return nil, err
		}
		movies = append(movies, m)
	}
	return movies, nil
}

func (r *AdminRepository) GetMovieByID(ctx context.Context, id int) (*models.Movie, error) {
	var m models.Movie
	err := r.db.QueryRow(ctx, `
		SELECT id, backdrop_path, overview, popularity, poster_path,
		       release_date, duration, title, director_name,
		       created_at, updated_at, deleted_at
		FROM movies WHERE id=$1 AND deleted_at IS NULL
	`, id).Scan(
		&m.ID, &m.Backdrop, &m.Overview, &m.Popularity, &m.Poster,
		&m.ReleaseDate, &m.Duration, &m.Title, &m.Director,
		&m.CreatedAt, &m.UpdatedAt, &m.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	gr, _ := r.db.Query(ctx, `
		SELECT g.id, g.name
		FROM genres g
		JOIN movies_genres mg ON g.id = mg.genres_id
		WHERE mg.movies_id=$1
	`, id)
	defer gr.Close()
	for gr.Next() {
		var g models.Genre
		gr.Scan(&g.ID, &g.Name)
		m.Genres = append(m.Genres, g)
	}

	cr, _ := r.db.Query(ctx, `
		SELECT c.id, c.name
		FROM casts c
		JOIN movies_casts mc ON c.id = mc.casts_id
		WHERE mc.movies_id=$1
	`, id)
	defer cr.Close()
	for cr.Next() {
		var c models.Cast
		cr.Scan(&c.ID, &c.Name)
		m.Casts = append(m.Casts, c)
	}

	return &m, nil
}

func (r *AdminRepository) SoftDeleteMovie(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `UPDATE movies SET deleted_at=NOW() WHERE id=$1 AND deleted_at IS NULL`, id)
	return err
}

func (r *AdminRepository) CreateMovie(ctx context.Context, movie *models.Movie, genreIDs, castIDs []int, schedules []map[string]interface{}) (*models.Movie, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	q := `
		INSERT INTO movies (backdrop_path, overview, popularity, poster_path, release_date, duration, title, director_name, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,NOW())
		RETURNING id, created_at
	`
	err = tx.QueryRow(ctx, q,
		movie.Backdrop, movie.Overview, movie.Popularity,
		movie.Poster, movie.ReleaseDate, movie.Duration,
		movie.Title, movie.Director,
	).Scan(&movie.ID, &movie.CreatedAt)
	if err != nil {
		return nil, err
	}

	for _, gid := range genreIDs {
		if _, err := tx.Exec(ctx, `INSERT INTO movies_genres (movies_id, genres_id) VALUES ($1,$2)`, movie.ID, gid); err != nil {
			return nil, err
		}
	}

	for _, cid := range castIDs {
		if _, err := tx.Exec(ctx, `INSERT INTO movies_casts (movies_id, casts_id) VALUES ($1,$2)`, movie.ID, cid); err != nil {
			return nil, err
		}
	}

	for _, s := range schedules {
		date := s["date"].(time.Time)
		cinemaID := s["cinema_id"].(int)
		locationID := s["location_id"].(int)
		timeIDs := s["time_ids"].([]int)

		for _, tid := range timeIDs {
			_, err := tx.Exec(ctx, `
				INSERT INTO schedules (movies_id, cinemas_id, locations_id, times_id, date)
				VALUES ($1,$2,$3,$4,$5)
			`, movie.ID, cinemaID, locationID, tid, date)
			if err != nil {
				return nil, err
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return r.GetMovieByID(ctx, movie.ID)
}

func (r *AdminRepository) UpdateMovie(ctx context.Context, id int, update map[string]interface{}, genreIDs, castIDs []int) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if len(update) > 0 {
		var setClauses []string
		args := []interface{}{}
		i := 1
		for k, v := range update {
			setClauses = append(setClauses, fmt.Sprintf("%s=$%d", k, i))
			args = append(args, v)
			i++
		}
		setClauses = append(setClauses, "updated_at=NOW()")
		query := fmt.Sprintf("UPDATE movies SET %s WHERE id=$%d", strings.Join(setClauses, ","), i)
		args = append(args, id)

		if _, err := tx.Exec(ctx, query, args...); err != nil {
			return err
		}
	} else {
		if _, err := tx.Exec(ctx, `UPDATE movies SET updated_at=NOW() WHERE id=$1`, id); err != nil {
			return err
		}
	}

	if len(genreIDs) > 0 {
		if _, err := tx.Exec(ctx, `DELETE FROM movies_genres WHERE movies_id=$1`, id); err != nil {
			return err
		}
		for _, gid := range genreIDs {
			if _, err := tx.Exec(ctx, `INSERT INTO movies_genres (movies_id, genres_id) VALUES ($1,$2)`, id, gid); err != nil {
				return err
			}
		}
	}

	if len(castIDs) > 0 {
		if _, err := tx.Exec(ctx, `DELETE FROM movies_casts WHERE movies_id=$1`, id); err != nil {
			return err
		}
		for _, cid := range castIDs {
			if _, err := tx.Exec(ctx, `INSERT INTO movies_casts (movies_id, casts_id) VALUES ($1,$2)`, id, cid); err != nil {
				return err
			}
		}
	}

	return tx.Commit(ctx)
}
