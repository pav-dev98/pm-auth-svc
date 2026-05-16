package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secret           string
	expiration       time.Duration
	refreshExpiration time.Duration
}

func NewJWTService(secret string, expiration, refreshExpiration time.Duration) *JWTService {
	return &JWTService{
		secret:           secret,
		expiration:       expiration,
		refreshExpiration: refreshExpiration,
	}
}

func (j *JWTService) GenerateToken(ID uint, email string) (string, error) {
	claims := jwt.MapClaims{
		"email":  email,
		"userID": ID,
		"exp":    time.Now().Add(j.expiration).Unix(),
		"iat":    time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}

func (j *JWTService) GenerateRefreshToken(ID uint) (string, error) {
	claims := jwt.MapClaims{
		"userID": ID,
		"type":   "refresh",
		"exp":    time.Now().Add(j.refreshExpiration).Unix(),
		"iat":    time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}