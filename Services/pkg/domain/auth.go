package domain

import (
	"time"

	"github.com/ChangerzaryX1602/SkillSync/pkg/models"
	"github.com/golang-jwt/jwt/v4"

	helpers "github.com/zercle/gofiber-helpers"
)

type AuthRepository interface {
	SignToken(user models.User, host string, ttl time.Duration) (string, *helpers.ResponseError)
	SaveRefreshToken(userID uint, token string, ttl time.Duration) *helpers.ResponseError
	GetRefreshToken(userID uint) (string, *helpers.ResponseError)
	DeleteRefreshToken(userID uint) *helpers.ResponseError
	ParseToken(tokenString string) (*jwt.Token, *helpers.ResponseError)
}

type AuthService interface {
	Register(user models.User) []helpers.ResponseError
	Login(user models.User, host string) (*string, *string, []helpers.ResponseError)
	RefreshToken(refreshToken string) (*string, *string, []helpers.ResponseError)
	GetUserByID(userId uint) (*models.User, []helpers.ResponseError)
}
