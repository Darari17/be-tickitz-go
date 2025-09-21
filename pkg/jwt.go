package pkg

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

func NewJWTClaims(u uuid.UUID, e string, r string) *Claims {
	return &Claims{
		UserID: u,
		Email:  e,
		Role:   r,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
			Issuer:    os.Getenv("JWT_ISSUER"),
		},
	}
}

func (c *Claims) GenerateToken() (string, error) {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return "", errors.New("no secret keys found")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString([]byte(secretKey))
}

func (c *Claims) VerifyToken(token string) error {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return errors.New("no secret keys found")
	}

	parsedToken, err := jwt.ParseWithClaims(token, c, func(t *jwt.Token) (any, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return err
	}

	if !parsedToken.Valid {
		return jwt.ErrTokenExpired
	}

	iss, err := parsedToken.Claims.GetIssuer()
	if err != nil {
		return err
	}

	if iss != os.Getenv("JWT_ISSUER") {
		return jwt.ErrTokenInvalidIssuer
	}

	return nil
}
