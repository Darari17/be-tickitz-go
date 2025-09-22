package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Darari17/be-tickitz-full/internal/dtos"
	"github.com/Darari17/be-tickitz-full/internal/repositories"
	"github.com/gin-gonic/gin"
)

type MovieController struct {
	movieRepository *repositories.MovieRepository
}

func NewMovieController(mr *repositories.MovieRepository) *MovieController {
	return &MovieController{
		movieRepository: mr,
	}
}

// GetUpcomingMovies godoc
// @Summary Get upcoming movies
// @Description Retrieve a paginated list of upcoming movies
// @Tags Movies
// @Produce json
// @Param page query int false "Page number"
// @Success 200 {object} dtos.Response
// @Failure 500 {object} dtos.ErrResponse
// @Router /movies/upcoming [get]
func (mc *MovieController) GetUpcomingMovies(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	movies, _, err := mc.movieRepository.GetUpcomingMovies(c.Request.Context(), page)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to fetch upcoming movies",
		})
		return
	}

	c.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Message: "Get upcoming movies successfully",
		Data:    movies,
	})
}

// GetPopularMovies godoc
// @Summary Get popular movies
// @Description Retrieve a paginated list of popular movies
// @Tags Movies
// @Produce json
// @Param page query int false "Page number"
// @Success 200 {object} dtos.Response
// @Failure 500 {object} dtos.ErrResponse
// @Router /movies/popular [get]
func (mc *MovieController) GetPopularMovies(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	movies, _, err := mc.movieRepository.GetPopularMovies(c.Request.Context(), page)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to fetch popular movies",
		})
		return
	}

	c.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Message: "Get popular movies successfully",
		Data:    movies,
	})
}

// GetAllMovies godoc
// @Summary Get all movies
// @Description Retrieve a paginated list of movies with optional search and genre filter
// @Tags Movies
// @Produce json
// @Param page query int false "Page number"
// @Param search query string false "Search by title"
// @Param genre query int false "Filter by genre ID"
// @Success 200 {object} dtos.Response
// @Failure 500 {object} dtos.ErrResponse
// @Router /movies [get]
func (mc *MovieController) GetAllMovies(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	search := c.DefaultQuery("search", "")
	genreStr := c.DefaultQuery("genre", "")
	genreID := 0
	if genreStr != "" {
		if g, err := strconv.Atoi(genreStr); err == nil {
			genreID = g
		}
	}

	movies, total, err := mc.movieRepository.GetMovies(c.Request.Context(), page, search, genreID)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to fetch movies",
		})
		return
	}

	const limit = 12
	meta := map[string]interface{}{
		"page":        page,
		"limit":       limit,
		"total":       total,
		"total_pages": (total + limit - 1) / limit,
	}

	c.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Message: "Get movies successfully",
		Data: map[string]interface{}{
			"movies": movies,
			"meta":   meta,
		},
	})
}

// GetMovieDetail godoc
// @Summary Get movie detail
// @Description Retrieve detailed information about a movie by ID
// @Tags Movies
// @Produce json
// @Param id path int true "Movie ID"
// @Success 200 {object} dtos.Response
// @Failure 500 {object} dtos.ErrResponse
// @Router /movies/{id} [get]
func (mh *MovieController) GetMovieDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid movie ID",
		})
		return
	}

	movie, err := mh.movieRepository.GetMovieDetails(c.Request.Context(), id)
	if err != nil {
		log.Println(err.Error())
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
		Message: "Get movie details successfully",
		Data:    movie,
	})
}

// GetGenres godoc
// @Summary Get all genres
// @Description Retrieve list of all available genres
// @Tags Movies
// @Produce json
// @Success 200 {object} dtos.Response
// @Failure 500 {object} dtos.ErrResponse
// @Router /movies/genres [get]
func (mc *MovieController) GetGenres(c *gin.Context) {
	genres, err := mc.movieRepository.GetAllGenres(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to fetch genres",
		})
		return
	}

	c.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Message: "Get genres successfully",
		Data:    genres,
	})
}
