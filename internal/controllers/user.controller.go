package controllers

import (
	"log"
	"net/http"

	"github.com/Darari17/be-tickitz-full/internal/dtos"
	"github.com/Darari17/be-tickitz-full/internal/models"
	"github.com/Darari17/be-tickitz-full/internal/repositories"
	"github.com/Darari17/be-tickitz-full/internal/utils"
	"github.com/Darari17/be-tickitz-full/pkg"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userRepository *repositories.UserRepository
}

func NewUserController(ur *repositories.UserRepository) *UserController {
	return &UserController{
		userRepository: ur,
	}
}

// Login godoc
// @Summary User login
// @Description Authenticate user with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body dtos.UserRequest true "User login credentials"
// @Success 200 {object} dtos.Response
// @Failure 500 {object} dtos.ErrResponse
// @Router /auth/login [post]
func (uc *UserController) Login(c *gin.Context) {
	body := dtos.UserRequest{}
	if err := c.ShouldBind(&body); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid Request Body",
			Data:    nil,
		})
		return
	}

	user, err := uc.userRepository.GetEmail(c.Request.Context(), body.Email)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Something went wrong",
			Data:    nil,
		})
		return
	}

	if user == nil {
		c.JSON(http.StatusUnauthorized, dtos.Response{
			Code:    http.StatusUnauthorized,
			Success: false,
			Message: "Invalid Email or Password",
			Data:    nil,
		})
		return
	}

	hash := pkg.HashConfig{}
	valid, err := hash.CompareHashAndPassword(body.Password, user.Password)
	if err != nil || !valid {
		log.Println(err.Error())
		c.JSON(http.StatusUnauthorized, dtos.Response{
			Code:    http.StatusUnauthorized,
			Success: false,
			Message: "Invalid Email or Password",
			Data:    nil,
		})
		return
	}

	claim := pkg.NewJWTClaims(user.ID, user.Email, string(user.Role))
	token, err := claim.GenerateToken()
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to Generate Token",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Message: "Login Successfully",
		Data: dtos.UserResponse{
			UserID: claim.UserID,
			Email:  claim.Email,
			Role:   claim.Role,
			Token:  token,
		},
	})
}

// Register godoc
// @Summary User registration
// @Description Register a new user with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body dtos.UserRequest true "User registration request"
// @Success 201 {object} dtos.Response
// @Failure 500 {object} dtos.ErrResponse
// @Router /auth/register [post]
func (uc *UserController) Register(c *gin.Context) {
	body := dtos.UserRequest{}
	if err := c.ShouldBind(&body); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid Request Body",
			Data:    nil,
		})
		return
	}

	hash := pkg.HashConfig{}
	hash.UseRecommended()
	hashed, err := hash.GenHash(body.Password)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to Hash Password",
			Data:    nil,
		})
		return
	}

	user := models.User{
		Email:    body.Email,
		Password: hashed,
	}

	if err := uc.userRepository.InsertUser(c, &user); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusConflict, dtos.Response{
			Code:    http.StatusConflict,
			Success: false,
			Message: "Email Already Exists",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, dtos.Response{
		Code:    http.StatusCreated,
		Success: true,
		Message: "Register Successfully",
		Data:    nil,
	})
}

// GetProfile godoc
// @Summary Get user profile
// @Description Retrieve logged in user profile
// @Tags Profile
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dtos.Response
// @Failure 500 {object} dtos.ErrResponse
// @Router /profile [get]
func (uc *UserController) GetProfile(c *gin.Context) {
	user, err := utils.GetUser(c)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusUnauthorized, dtos.Response{
			Code:    http.StatusUnauthorized,
			Success: false,
			Message: "Unauthorized: " + err.Error(),
			Data:    nil,
		})
		return
	}

	profile, err := uc.userRepository.GetProfile(c.Request.Context(), user.ID)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusNotFound, dtos.Response{
			Code:    http.StatusNotFound,
			Success: false,
			Message: "Profile not found",
			Data:    nil,
		})
		return
	}

	res := dtos.ProfileResponse{
		UserID:      profile.UserID,
		FirstName:   profile.FirstName,
		LastName:    profile.LastName,
		PhoneNumber: profile.PhoneNumber,
		Avatar:      profile.Avatar,
		Point:       profile.Point,
	}

	c.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Message: "Get profile successfully",
		Data:    res,
	})
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update first name, last name, and phone number
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dtos.ProfileRequest true "Profile update request"
// @Success 200 {object} dtos.Response
// @Failure 500 {object} dtos.ErrResponse
// @Router /profile [patch]
func (uc *UserController) UpdateProfile(c *gin.Context) {
	user, err := utils.GetUser(c)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusUnauthorized, dtos.Response{
			Code:    http.StatusUnauthorized,
			Success: false,
			Message: "Unauthorized: " + err.Error(),
			Data:    nil,
		})
		return
	}

	var req dtos.ProfileRequest
	if err := c.ShouldBind(&req); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid request body",
			Data:    nil,
		})
		return
	}

	profile := &models.Profile{
		UserID:      user.ID,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
	}

	if err := uc.userRepository.UpdateProfile(c.Request.Context(), profile); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to update profile",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Message: "Profile updated successfully",
		Data:    nil,
	})
}

// ChangePassword godoc
// @Summary Change user password
// @Description Update old password to new password
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dtos.ChangePasswordRequest true "Change password request"
// @Success 200 {object} dtos.Response
// @Failure 500 {object} dtos.ErrResponse
// @Router /profile/change-password [patch]
func (uc *UserController) ChangePassword(c *gin.Context) {
	user, err := utils.GetUser(c)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusUnauthorized, dtos.Response{
			Code:    http.StatusUnauthorized,
			Success: false,
			Message: "Unauthorized: " + err.Error(),
			Data:    nil,
		})
		return
	}

	var req dtos.ChangePasswordRequest
	if err := c.ShouldBind(&req); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid request body",
			Data:    nil,
		})
		return
	}

	hashedPassword, err := uc.userRepository.VerifyPassword(c.Request.Context(), user.ID, req.OldPassword)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusUnauthorized, dtos.Response{
			Code:    http.StatusUnauthorized,
			Success: false,
			Message: "Invalid old password",
			Data:    nil,
		})
		return
	}

	hashConfig := pkg.NewHashConfig()
	ok, err := hashConfig.CompareHashAndPassword(req.OldPassword, hashedPassword)
	if err != nil || !ok {
		log.Println(err.Error())
		c.JSON(http.StatusUnauthorized, dtos.Response{
			Code:    http.StatusUnauthorized,
			Success: false,
			Message: "Invalid old password",
			Data:    nil,
		})
		return
	}

	hashConfig.UseRecommended()
	newHashed, err := hashConfig.GenHash(req.NewPassword)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to hash new password",
			Data:    nil,
		})
		return
	}

	if err := uc.userRepository.UpdatePassword(c.Request.Context(), user.ID, newHashed); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to update password",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Message: "Password changed successfully",
		Data:    nil,
	})
}

// ChangeAvatar godoc
// @Summary Change user avatar
// @Description Upload a new avatar image for user profile
// @Tags Profile
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param avatar formData file true "Avatar image"
// @Success 200 {object} dtos.Response
// @Failure 500 {object} dtos.ErrResponse
// @Router /profile/change-avatar [patch]
func (uc *UserController) ChangeAvatar(c *gin.Context) {
	user, err := utils.GetUser(c)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusUnauthorized, dtos.Response{
			Code:    http.StatusUnauthorized,
			Success: false,
			Message: "Unauthorized: " + err.Error(),
			Data:    nil,
		})
		return
	}

	var req dtos.AvatarProfileRequest
	if err := c.ShouldBind(&req); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, dtos.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid request body",
			Data:    nil,
		})
		return
	}

	if req.Avatar == nil {
		c.JSON(http.StatusOK, dtos.Response{
			Code:    http.StatusOK,
			Success: true,
			Message: "No avatar uploaded, nothing to update",
			Data:    nil,
		})
		return
	}

	filename := utils.SaveImage(c, req.Avatar, "avatars")

	if err := uc.userRepository.UpdateAvatar(c.Request.Context(), user.ID, filename); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, dtos.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to update avatar",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, dtos.Response{
		Code:    http.StatusOK,
		Success: true,
		Message: "Avatar updated successfully",
		Data:    nil,
	})
}
