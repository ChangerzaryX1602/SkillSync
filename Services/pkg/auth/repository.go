package auth

import (
	"fmt"
	"time"

	"github.com/ChangerzaryX1602/SkillSync/pkg/domain"
	"github.com/ChangerzaryX1602/SkillSync/pkg/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	helpers "github.com/zercle/gofiber-helpers"
)

type authRepository struct {
	resources models.Resources
}

func NewAuthRepository(resources models.Resources) domain.AuthRepository {
	return &authRepository{resources: resources}
}
func (r *authRepository) SignToken(user models.User, host string, ttl time.Duration) (string, *helpers.ResponseError) {
	token := jwt.NewWithClaims(r.resources.JwtResources.JwtSigningMethod, &jwt.RegisteredClaims{})
	claims := token.Claims.(*jwt.RegisteredClaims)
	claims.Subject = fmt.Sprintf("%d", user.ID)
	claims.Issuer = host
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(ttl))
	signToken, err := token.SignedString(r.resources.JwtResources.JwtSignKey)
	if err != nil {
		return "", &helpers.ResponseError{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Failed to sign json web token",
		}
	}
	return signToken, nil
}

func (r *authRepository) SaveRefreshToken(userID uint, token string, ttl time.Duration) *helpers.ResponseError {
	key := fmt.Sprintf("refresh:%d", userID)
	if err := r.resources.RedisStorage.Set(key, []byte(token), ttl); err != nil {
		return &helpers.ResponseError{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Failed to save refresh token",
		}
	}
	return nil
}

func (r *authRepository) GetRefreshToken(userID uint) (string, *helpers.ResponseError) {
	key := fmt.Sprintf("refresh:%d", userID)
	tokenBytes, err := r.resources.RedisStorage.Get(key)
	if err != nil {
		return "", &helpers.ResponseError{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Failed to get refresh token",
		}
	}
	if tokenBytes == nil {
		return "", &helpers.ResponseError{
			Message: "Refresh token not found",
			Code:    fiber.StatusUnauthorized,
			Source:  helpers.WhereAmI(),
			Title:   "Unauthorized",
		}
	}
	return string(tokenBytes), nil
}

func (r *authRepository) DeleteRefreshToken(userID uint) *helpers.ResponseError {
	key := fmt.Sprintf("refresh:%d", userID)
	if err := r.resources.RedisStorage.Delete(key); err != nil {
		return &helpers.ResponseError{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Failed to delete refresh token",
		}
	}
	return nil
}

func (r *authRepository) ParseToken(tokenString string) (*jwt.Token, *helpers.ResponseError) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, r.resources.JwtResources.JwtKeyfunc)
	if err != nil {
		return nil, &helpers.ResponseError{
			Message: err.Error(),
			Code:    fiber.StatusUnauthorized,
			Source:  helpers.WhereAmI(),
			Title:   "Invalid token",
		}
	}
	return token, nil
}
