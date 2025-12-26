package auth

import (
	"context"
	"strconv"
	"time"

	"github.com/ChangerzaryX1602/SkillSync/pkg/domain"
	"github.com/ChangerzaryX1602/SkillSync/pkg/models"
	"github.com/ChangerzaryX1602/SkillSync/pkg/utils"

	"github.com/golang-jwt/jwt/v4"
	helpers "github.com/zercle/gofiber-helpers"
)

type authService struct {
	repository domain.AuthRepository
	service    domain.UserService
}

func NewAuthService(repository domain.AuthRepository, service domain.UserService) domain.AuthService {
	return &authService{repository: repository, service: service}
}
func (s *authService) Login(ctx context.Context, user models.User, host string) (*string, *string, []helpers.ResponseError) {
	if !utils.IsValidEmail(user.Email) {
		return nil, nil, []helpers.ResponseError{
			{
				Code:    400,
				Source:  helpers.WhereAmI(),
				Title:   "Bad Request",
				Message: "Invalid email format",
			},
		}
	}
	userE, errorFormList := s.service.GetUserByEmail(ctx, user.Email)
	if errorFormList != nil {
		return nil, nil, errorFormList
	}
	checked := models.CheckPasswordHash(user.PasswordTemp, userE.Password)
	if !checked {
		return nil, nil, []helpers.ResponseError{
			{
				Code:    401,
				Source:  helpers.WhereAmI(),
				Title:   "Unauthorized",
				Message: "Invalid credentials",
			},
		}
	}
	token, errorForm := s.repository.SignToken(ctx, *userE, host, time.Minute*15)
	if errorForm != nil {
		return nil, nil, []helpers.ResponseError{*errorForm}
	}

	refreshToken, errorForm := s.repository.SignToken(ctx, *userE, host, time.Hour*24*7)
	if errorForm != nil {
		return nil, nil, []helpers.ResponseError{*errorForm}
	}

	if errorForm := s.repository.SaveRefreshToken(ctx, userE.ID, refreshToken, time.Hour*24*7); errorForm != nil {
		return nil, nil, []helpers.ResponseError{*errorForm}
	}

	return &token, &refreshToken, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*string, *string, []helpers.ResponseError) {
	token, errorForm := s.repository.ParseToken(ctx, refreshToken)
	if errorForm != nil {
		return nil, nil, []helpers.ResponseError{*errorForm}
	}

	if !token.Valid {
		return nil, nil, []helpers.ResponseError{
			{
				Code:    401,
				Source:  helpers.WhereAmI(),
				Title:   "Unauthorized",
				Message: "Invalid token",
			},
		}
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, nil, []helpers.ResponseError{
			{
				Code:    401,
				Source:  helpers.WhereAmI(),
				Title:   "Unauthorized",
				Message: "Invalid token claims",
			},
		}
	}

	userId, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		return nil, nil, []helpers.ResponseError{
			{
				Code:    401,
				Source:  helpers.WhereAmI(),
				Title:   "Unauthorized",
				Message: "Invalid user ID in token",
			},
		}
	}

	storedToken, errorForm := s.repository.GetRefreshToken(ctx, uint(userId))
	if errorForm != nil {
		return nil, nil, []helpers.ResponseError{*errorForm}
	}

	if storedToken != refreshToken {
		return nil, nil, []helpers.ResponseError{
			{
				Code:    401,
				Source:  helpers.WhereAmI(),
				Title:   "Unauthorized",
				Message: "Refresh token mismatch (Reuse detected?)",
			},
		}
	}

	user, errorFormList := s.service.GetUser(ctx, uint(userId))
	if errorFormList != nil {
		return nil, nil, errorFormList
	}

	newAccessToken, errorForm := s.repository.SignToken(ctx, *user, claims.Issuer, time.Minute*15)
	if errorForm != nil {
		return nil, nil, []helpers.ResponseError{*errorForm}
	}

	newRefreshToken, errorForm := s.repository.SignToken(ctx, *user, claims.Issuer, time.Hour*24*7)
	if errorForm != nil {
		return nil, nil, []helpers.ResponseError{*errorForm}
	}

	if errorForm := s.repository.SaveRefreshToken(ctx, user.ID, newRefreshToken, time.Hour*24*7); errorForm != nil {
		return nil, nil, []helpers.ResponseError{*errorForm}
	}

	return &newAccessToken, &newRefreshToken, nil
}
func (s *authService) Register(ctx context.Context, user models.User) []helpers.ResponseError {
	if !utils.IsValidEmail(user.Email) {
		return []helpers.ResponseError{
			{
				Code:    400,
				Source:  helpers.WhereAmI(),
				Title:   "Bad Request",
				Message: "Invalid email format",
			},
		}
	}
	err := s.service.CreateUser(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (s *authService) GetUserByID(ctx context.Context, userId uint) (*models.User, []helpers.ResponseError) {
	user, err := s.service.GetUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	user.Password = ""
	user.PasswordTemp = ""
	return user, nil
}
