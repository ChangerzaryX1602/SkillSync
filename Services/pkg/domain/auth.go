package domain

import (
	"context"
	"time"

	"github.com/ChangerzaryX1602/SkillSync/pkg/models"
	"github.com/golang-jwt/jwt/v4"

	helpers "github.com/zercle/gofiber-helpers"
)

type AuthRepository interface {
	SignToken(ctx context.Context, user models.User, host string, ttl time.Duration) (string, *helpers.ResponseError)
	SaveRefreshToken(ctx context.Context, userID uint, token string, ttl time.Duration) *helpers.ResponseError
	GetRefreshToken(ctx context.Context, userID uint) (string, *helpers.ResponseError)
	DeleteRefreshToken(ctx context.Context, userID uint) *helpers.ResponseError
	ParseToken(ctx context.Context, tokenString string) (*jwt.Token, *helpers.ResponseError)
}

type AuthService interface {
	Register(ctx context.Context, user models.User) []helpers.ResponseError
	Login(ctx context.Context, user models.User, host string) (*string, *string, []helpers.ResponseError)
	RefreshToken(ctx context.Context, refreshToken string) (*string, *string, []helpers.ResponseError)
	GetUserByID(ctx context.Context, userId uint) (*models.User, []helpers.ResponseError)
}
