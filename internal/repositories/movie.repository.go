package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Darari17/be-tickitz-full/internal/models"
	"github.com/Darari17/be-tickitz-full/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type MovieRepository struct {
	db  *pgxpool.Pool
	rdb *redis.Client
}

func NewMovieRepository(db *pgxpool.Pool, rdb *redis.Client) *MovieRepository {
	return &MovieRepository{
		db:  db,
		rdb: rdb,
	}
}

// ======================================================
// Upcoming Movies (Redis full, per page)
// ======================================================
func (mr *MovieRepository) GetUpcomingMovies(c context.Context, page int) ([]models.Movie, int, error) {
	const pageSize = 12
	offset := (page - 1) * pageSize

	rdbKey := fmt.Sprintf("movies:upcoming:page:%d", page)
	var cached []models.Movie
	ok, err := utils.GetRedis(c, mr.rdb, rdbKey, &cached)
	if err == nil && ok {
		var total int
		_ = mr.db.QueryRow(c, "SELECT COUNT(*) FROM movies WHERE release_date > NOW()").Scan(&total)
		return cached, total, nil
	}

	var total int
	if err := mr.db.QueryRow(c, "SELECT COUNT(*) FROM movies WHERE release_date > NOW()").Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT m.id, m.backdrop_path, m.overview, m.popularity, m.poster_path,
		       m.release_date, m.duration, m.title, m.director_name,
		       m.created_at, m.updated_at, m.deleted_at,
		       COALESCE(JSON_AGG(DISTINCT jsonb_build_object('id', g.id, 'name', g.name))
		                FILTER (WHERE g.id IS NOT NULL), '[]') AS genres,
		       COALESCE(JSON_AGG(DISTINCT jsonb_build_object('id', c.id, 'name', c.name))
		                FILTER (WHERE c.id IS NOT NULL), '[]') AS casts
		FROM movies m
		LEFT JOIN movies_genres mg ON m.id = mg.movies_id
		LEFT JOIN genres g ON g.id = mg.genres_id
		LEFT JOIN movies_casts mc ON m.id = mc.movies_id
		LEFT JOIN casts c ON c.id = mc.casts_id
		WHERE m.release_date > NOW()
		GROUP BY m.id
		ORDER BY m.release_date ASC
		LIMIT $1 OFFSET $2;
	`

	rows, err := mr.db.Query(c, query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie
		var genresJSON, castsJSON []byte
		if err := rows.Scan(
			&m.ID, &m.Backdrop, &m.Overview, &m.Popularity, &m.Poster,
			&m.ReleaseDate, &m.Duration, &m.Title, &m.Director,
			&m.CreatedAt, &m.UpdatedAt, &m.DeletedAt,
			&genresJSON, &castsJSON,
		); err != nil {
			return nil, 0, err
		}
		_ = json.Unmarshal(genresJSON, &m.Genres)
		_ = json.Unmarshal(castsJSON, &m.Casts)
		movies = append(movies, m)
	}

	if err := utils.SetRedis(c, mr.rdb, rdbKey, movies, time.Hour*1); err != nil {
		log.Printf("redis set error: %v\n", err)
	}

	return movies, total, nil
}

// ======================================================
// Popular Movies (Redis full, per page)
// ======================================================
func (mr *MovieRepository) GetPopularMovies(c context.Context, page int) ([]models.Movie, int, error) {
	const pageSize = 12
	offset := (page - 1) * pageSize

	rdbKey := fmt.Sprintf("movies:popular:page:%d", page)
	var cached []models.Movie
	ok, err := utils.GetRedis(c, mr.rdb, rdbKey, &cached)
	if err == nil && ok {
		var total int
		_ = mr.db.QueryRow(c, "SELECT COUNT(*) FROM movies").Scan(&total)
		return cached, total, nil
	}

	var total int
	if err := mr.db.QueryRow(c, "SELECT COUNT(*) FROM movies").Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT m.id, m.backdrop_path, m.overview, m.popularity, m.poster_path,
		       m.release_date, m.duration, m.title, m.director_name,
		       m.created_at, m.updated_at, m.deleted_at,
		       COALESCE(JSON_AGG(DISTINCT jsonb_build_object('id', g.id, 'name', g.name))
		                FILTER (WHERE g.id IS NOT NULL), '[]') AS genres,
		       COALESCE(JSON_AGG(DISTINCT jsonb_build_object('id', c.id, 'name', c.name))
		                FILTER (WHERE c.id IS NOT NULL), '[]') AS casts
		FROM movies m
		LEFT JOIN movies_genres mg ON m.id = mg.movies_id
		LEFT JOIN genres g ON g.id = mg.genres_id
		LEFT JOIN movies_casts mc ON m.id = mc.movies_id
		LEFT JOIN casts c ON c.id = mc.casts_id
		GROUP BY m.id
		ORDER BY m.popularity DESC
		LIMIT $1 OFFSET $2;
	`

	rows, err := mr.db.Query(c, query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie
		var genresJSON, castsJSON []byte
		if err := rows.Scan(
			&m.ID, &m.Backdrop, &m.Overview, &m.Popularity, &m.Poster,
			&m.ReleaseDate, &m.Duration, &m.Title, &m.Director,
			&m.CreatedAt, &m.UpdatedAt, &m.DeletedAt,
			&genresJSON, &castsJSON,
		); err != nil {
			return nil, 0, err
		}
		_ = json.Unmarshal(genresJSON, &m.Genres)
		_ = json.Unmarshal(castsJSON, &m.Casts)
		movies = append(movies, m)
	}

	if err := utils.SetRedis(c, mr.rdb, rdbKey, movies, time.Hour*1); err != nil {
		log.Printf("redis set error: %v\n", err)
	}

	return movies, total, nil
}

func (mr *MovieRepository) GetMovies(ctx context.Context, page int, search string, genreID int) ([]models.Movie, int, error) {
	const limit = 12 // 🔥 12 movie per halaman

	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	where := []string{"1=1"}
	args := []interface{}{}

	if strings.TrimSpace(search) != "" {
		args = append(args, "%"+strings.TrimSpace(search)+"%")
		args = append(args, "%"+strings.TrimSpace(search)+"%")
		idx := len(args)
		where = append(where, fmt.Sprintf("(m.title ILIKE $%d OR m.overview ILIKE $%d)", idx-1, idx))
	}

	if genreID > 0 {
		args = append(args, genreID)
		idx := len(args)
		where = append(where, fmt.Sprintf("EXISTS (SELECT 1 FROM movies_genres mg2 WHERE mg2.movies_id = m.id AND mg2.genres_id = $%d)", idx))
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(DISTINCT m.id) FROM movies m WHERE %s", strings.Join(where, " AND "))
	var total int
	if err := mr.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	limitIdx := len(args) + 1
	offsetIdx := len(args) + 2
	args = append(args, limit, offset)

	query := fmt.Sprintf(`
		SELECT
			m.id, m.backdrop_path, m.overview, m.popularity, m.poster_path,
			m.release_date, m.duration, m.title, m.director_name,
			m.created_at, m.updated_at, m.deleted_at,
			COALESCE(
				JSON_AGG(DISTINCT jsonb_build_object('id', g.id, 'name', g.name))
				FILTER (WHERE g.id IS NOT NULL), '[]'
			) AS genres
		FROM movies m
		LEFT JOIN movies_genres mg ON m.id = mg.movies_id
		LEFT JOIN genres g ON g.id = mg.genres_id
		WHERE %s
		GROUP BY m.id
		ORDER BY m.release_date DESC
		LIMIT $%d OFFSET $%d
	`, strings.Join(where, " AND "), limitIdx, offsetIdx)

	rows, err := mr.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie
		var genresJSON []byte
		if err := rows.Scan(
			&m.ID, &m.Backdrop, &m.Overview, &m.Popularity, &m.Poster,
			&m.ReleaseDate, &m.Duration, &m.Title, &m.Director,
			&m.CreatedAt, &m.UpdatedAt, &m.DeletedAt,
			&genresJSON,
		); err != nil {
			return nil, 0, err
		}
		_ = json.Unmarshal(genresJSON, &m.Genres)
		movies = append(movies, m)
	}

	return movies, total, nil
}

// ======================================================
// Movie Detail (tidak pakai Redis)
// ======================================================
func (mr *MovieRepository) GetMovieDetails(c context.Context, id int) (*models.Movie, error) {
	query := `
		SELECT m.id, m.backdrop_path, m.overview, m.popularity, m.poster_path,
		       m.release_date, m.duration, m.title, m.director_name,
		       m.created_at, m.updated_at, m.deleted_at,
		       COALESCE(
		           JSON_AGG(DISTINCT jsonb_build_object('id', g.id, 'name', g.name))
		           FILTER (WHERE g.id IS NOT NULL), '[]'
		       ) AS genres,
		       COALESCE(
		           JSON_AGG(DISTINCT jsonb_build_object('id', c.id, 'name', c.name))
		           FILTER (WHERE c.id IS NOT NULL), '[]'
		       ) AS casts
		FROM movies m
		LEFT JOIN movies_genres mg ON m.id = mg.movies_id
		LEFT JOIN genres g ON g.id = mg.genres_id
		LEFT JOIN movies_casts mc ON m.id = mc.movies_id
		LEFT JOIN casts c ON c.id = mc.casts_id
		WHERE m.id = $1
		GROUP BY m.id
	`

	var m models.Movie
	var genresJSON, castsJSON []byte

	err := mr.db.QueryRow(c, query, id).Scan(
		&m.ID, &m.Backdrop, &m.Overview, &m.Popularity, &m.Poster,
		&m.ReleaseDate, &m.Duration, &m.Title, &m.Director,
		&m.CreatedAt, &m.UpdatedAt, &m.DeletedAt,
		&genresJSON, &castsJSON,
	)
	if err != nil {
		return nil, err
	}

	_ = json.Unmarshal(genresJSON, &m.Genres)
	_ = json.Unmarshal(castsJSON, &m.Casts)

	return &m, nil
}

// GetAllGenres mengambil semua genre yang ada di tabel genres
func (mr *MovieRepository) GetAllGenres(ctx context.Context) ([]models.Genre, error) {
	const redisKey = "genres:all"
	var cached []models.Genre

	// coba ambil dari redis
	ok, err := utils.GetRedis(ctx, mr.rdb, redisKey, &cached)
	if err == nil && ok {
		return cached, nil
	}

	rows, err := mr.db.Query(ctx, `SELECT id, name FROM genres ORDER BY name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var genres []models.Genre
	for rows.Next() {
		var g models.Genre
		if err := rows.Scan(&g.ID, &g.Name); err != nil {
			return nil, err
		}
		genres = append(genres, g)
	}

	// simpan ke redis 1 jam
	if err := utils.SetRedis(ctx, mr.rdb, redisKey, genres, time.Hour); err != nil {
		log.Printf("redis set error: %v\n", err)
	}

	return genres, nil
}
