package utils

import (
	"errors"

	"github.com/Darari17/be-tickitz-full/internal/models"
	"github.com/Darari17/be-tickitz-full/pkg"
	"github.com/gin-gonic/gin"
)

var (
	ErrClaimsNotFound   = errors.New("claims not found in context, token might be missing")
	ErrInvalidClaimsFmt = errors.New("invalid claims format")
)

func GetUser(c *gin.Context) (*models.UserContext, error) {
	claims, exists := c.Get("claims")
	if !exists {
		return nil, ErrClaimsNotFound
	}

	userClaims, ok := claims.(*pkg.Claims)
	if !ok || userClaims == nil {
		return nil, ErrInvalidClaimsFmt
	}

	userCtx := &models.UserContext{
		ID:    userClaims.UserID,
		Email: userClaims.Email,
		Role:  userClaims.Role,
	}

	return userCtx, nil
}
