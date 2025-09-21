package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Darari17/be-tickitz-full/internal/dtos"
	"github.com/Darari17/be-tickitz-full/internal/models"
	"github.com/Darari17/be-tickitz-full/internal/repositories"
	"github.com/Darari17/be-tickitz-full/internal/utils"
	"github.com/gin-gonic/gin"
)

type AdminController struct {
	adminRepository *repositories.AdminRepository
}

func NewAdminController(ar *repositories.AdminRepository) *AdminController {
	return &AdminController{
		adminRepository: ar,
	}
}

// CreateMovie godoc
// @Summary Create new movie
// @Description Add a new movie with genres, casts, schedules, poster & backdrop
// @Tags Admin - Movies
// @Accept multipart/form-data
// @Produce json
// @Param title formData string true "Movie Title"
// @Param overview formData string false "Movie Overview"
// @Param director_name formData string false "Director Name"
// @Param duration formData int false "Duration in minutes"
// @Param release_date formData string false "Release date in format YYYY-MM-DD"
// @Param popularity formData number false "Popularity"
// @Param poster formData file false "Poster image"
// @Param backdrop formData file false "Backdrop image"
// @Param genres formData []int false "Genre IDs"
// @Param casts formData []int false "Cast IDs"
// @Param schedules formData string false "Schedules JSON [{cinema_id, location_id, date, time_ids}]"
// @Success 201 {object} dtos.Response
// @Failure 400 {object} dtos.Response
// @Failure 500 {object} dtos.Response
// @Router /admin/movies [post]
func (ac *AdminController) CreateMovie(c *gin.Context) {
	var body dtos.CreateMovieRequest
	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid request data",
		})
		return
	}

	movie := &models.Movie{
		Title:       body.Title,
		Overview:    body.Overview,
		Director:    body.Director,
		Duration:    body.Duration,
		ReleaseDate: body.ReleaseDate,
		Popularity:  body.Popularity,
	}

	if body.Poster != nil {
		path := utils.SaveImage(c, body.Poster, "poster")
		if path == "" {
			return
		}
		movie.Poster = path
	}
	if body.Backdrop != nil {
		path := utils.SaveImage(c, body.Backdrop, "backdrop")
		if path == "" {
			return
		}
		movie.Backdrop = path
	}

	var schedules []map[string]interface{}
	for _, s := range body.Schedules {
		date, _ := time.Parse("2006-01-02", s.Date)
		schedules = append(schedules, map[string]interface{}{
			"date":        date,
			"cinema_id":   s.CinemaID,
			"location_id": s.LocationID,
			"time_ids":    s.TimeIDs,
		})
	}

	created, err := ac.adminRepository.CreateMovie(c, movie, body.Genres, body.Casts, schedules)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, dtos.Response{
		Code:    http.StatusCreated,
		Success: true,
		Message: "Movie created successfully",
		Data:    created,
	})
}

// GetMovies godoc
// @Summary Get list of movies
// @Description Retrieve all movies
// @Tags Admin - Movies
// @Produce json
// @Success 200 {object} dtos.Response
// @Failure 500 {object} dtos.Response
// @Router /admin/movies [get]
func (ac *AdminController) GetMovies(c *gin.Context) {
	movies, err := ac.adminRepository.GetMovies(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to fetch movies",
		})
		return
	}

	c.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    movies,
	})
}

// GetMovieByID godoc
// @Summary Get movie detail
// @Description Retrieve a movie by ID
// @Tags Admin - Movies
// @Produce json
// @Param id path int true "Movie ID"
// @Success 200 {object} dtos.Response
// @Failure 404 {object} dtos.Response
// @Router /admin/movies/{id} [get]
func (ac *AdminController) GetMovieByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	movie, err := ac.adminRepository.GetMovieByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, dtos.Response{
			Code:    http.StatusNotFound,
			Success: false,
			Message: "Movie not found",
		})
		return
	}

	c.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    movie,
	})
}

// UpdateMovie godoc
// @Summary Update movie
// @Description Update movie data including genres, casts, and images
// @Tags Admin - Movies
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Movie ID"
// @Param title formData string false "Movie Title"
// @Param overview formData string false "Movie Overview"
// @Param director_name formData string false "Director Name"
// @Param duration formData int false "Duration in minutes"
// @Param release_date formData string false "Release date in format YYYY-MM-DD"
// @Param popularity formData number false "Popularity"
// @Param poster formData file false "Poster image"
// @Param backdrop formData file false "Backdrop image"
// @Param genres formData []int false "Genre IDs"
// @Param casts formData []int false "Cast IDs"
// @Success 200 {object} dtos.Response
// @Failure 400 {object} dtos.Response
// @Failure 404 {object} dtos.Response
// @Failure 500 {object} dtos.Response
// @Router /admin/movies/{id} [patch]
func (ac *AdminController) UpdateMovie(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var body dtos.UpdateMovieRequest

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid request data",
		})
		return
	}

	_, err := ac.adminRepository.GetMovieByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, dtos.Response{
			Code:    http.StatusNotFound,
			Success: false,
			Message: "Movie not found",
		})
		return
	}

	update := make(map[string]interface{})
	if body.Title != nil {
		update["title"] = *body.Title
	}
	if body.Overview != nil {
		update["overview"] = *body.Overview
	}
	if body.Director != nil {
		update["director_name"] = *body.Director
	}
	if body.Duration != nil {
		update["duration"] = *body.Duration
	}
	if body.ReleaseDate != nil {
		update["release_date"] = *body.ReleaseDate
	}
	if body.Popularity != nil {
		update["popularity"] = *body.Popularity
	}
	if body.Poster != nil {
		path := utils.SaveImage(c, body.Poster, "poster")
		if path == "" {
			return
		}
		update["poster_path"] = path
	}
	if body.Backdrop != nil {
		path := utils.SaveImage(c, body.Backdrop, "backdrop")
		if path == "" {
			return
		}
		update["backdrop_path"] = path
	}

	if err := ac.adminRepository.UpdateMovie(c, id, update, body.Genres, body.Casts); err != nil {
		c.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: err.Error(),
		})
		return
	}

	updated, _ := ac.adminRepository.GetMovieByID(c, id)
	c.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Message: "Movie updated successfully",
		Data:    updated,
	})
}

// DeleteMovie godoc
// @Summary Delete movie
// @Description Soft delete a movie by ID
// @Tags Admin - Movies
// @Produce json
// @Param id path int true "Movie ID"
// @Success 200 {object} dtos.Response
// @Failure 500 {object} dtos.Response
// @Router /admin/movies/{id} [delete]
func (ac *AdminController) DeleteMovie(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := ac.adminRepository.SoftDeleteMovie(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to delete movie",
		})
		return
	}

	c.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Message: "Movie deleted successfully",
	})
}
