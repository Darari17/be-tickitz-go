package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Darari17/be-tickitz-full/internal/dtos"
	"github.com/Darari17/be-tickitz-full/internal/models"
	"github.com/Darari17/be-tickitz-full/internal/repositories"
	"github.com/Darari17/be-tickitz-full/internal/utils"
	"github.com/gin-gonic/gin"
)

type AdminController struct {
	adminRepository *repositories.AdminRepo
}

func NewAdminController(ar *repositories.AdminRepo) *AdminController {
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
// @Param popularity formData int false "Popularity"
// @Param poster formData file false "Poster image"
// @Param backdrop formData file false "Backdrop image"
// @Param genres formData []int false "Genre IDs (contoh: [1,2])" collectionFormat(multi)
// @Param casts formData []int false "Cast IDs (contoh: [3,5,7])" collectionFormat(multi)
// @Param schedules formData string false "Schedules JSON [{cinema_id, location_id, date, time_ids}]"
// @Success 201 {object} dtos.Response
// @Failure 400 {object} dtos.Response
// @Failure 500 {object} dtos.Response
// @Router /admin/movies [post]
// @Security BearerAuth
func (ac *AdminController) CreateMovie(c *gin.Context) {
	var body dtos.CreateMovieRequest
	if err := c.ShouldBind(&body); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid request data",
		})
		return
	}

	genresRaw := c.PostFormArray("genres")
	var genreIDs []int
	for _, g := range genresRaw {
		for _, part := range strings.Split(g, ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			id, err := strconv.Atoi(part)
			if err != nil {
				c.JSON(http.StatusBadRequest, dtos.Response{
					Code:    http.StatusBadRequest,
					Success: false,
					Message: "Invalid genre ID",
				})
				return
			}
			genreIDs = append(genreIDs, id)
		}
	}

	castsRaw := c.PostFormArray("casts")
	var castIDs []int
	for _, g := range castsRaw {
		for _, part := range strings.Split(g, ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			id, err := strconv.Atoi(part)
			if err != nil {
				c.JSON(http.StatusBadRequest, dtos.Response{
					Code:    http.StatusBadRequest,
					Success: false,
					Message: "Invalid cast ID",
				})
				return
			}
			castIDs = append(castIDs, id)
		}
	}

	schedulesRaw := c.PostForm("schedules")
	if schedulesRaw != "" {
		var scheduleReqs []dtos.ScheduleRequest
		if err := json.Unmarshal([]byte(schedulesRaw), &scheduleReqs); err != nil {
			c.JSON(http.StatusBadRequest, dtos.Response{
				Code:    http.StatusBadRequest,
				Success: false,
				Message: "Invalid schedules JSON format",
			})
			return
		}
		body.Schedules = scheduleReqs
	}

	if body.Title == "" {
		c.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "title is required",
		})
		return
	}

	date, err := time.Parse("2006-01-02", body.ReleaseDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "title is required",
		})
		return
	}

	movie := &models.Movie{
		Title:       body.Title,
		Overview:    body.Overview,
		Director:    body.Director,
		Duration:    body.Duration,
		ReleaseDate: date,
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

	created, err := ac.adminRepository.CreateMovie(c, movie, genreIDs, castIDs, schedules)
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
// @Security BearerAuth
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
// @Security BearerAuth
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
// @Param genres formData []int false "Genre IDs (contoh: [1,2,3])"
// @Param casts formData []int false "Cast IDs (contoh: [4,6])"
// @Success 200 {object} dtos.Response
// @Failure 400 {object} dtos.Response
// @Failure 404 {object} dtos.Response
// @Failure 500 {object} dtos.Response
// @Router /admin/movies/{id} [patch]
// @Security BearerAuth
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
// @Security BearerAuth
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
