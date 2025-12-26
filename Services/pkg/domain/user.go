package domain

import (
	"context"

	"github.com/ChangerzaryX1602/SkillSync/pkg/models"

	helpers "github.com/zercle/gofiber-helpers"
)

type UserRepository interface {
	Migrate() error
	CreateUser(ctx context.Context, user models.User) *helpers.ResponseError
	GetUser(ctx context.Context, id uint) (*models.User, *helpers.ResponseError)
	GetUsers(ctx context.Context, pagination models.Pagination, search models.Search) ([]models.User, *models.Pagination, *models.Search, *helpers.ResponseError)
	UpdateUser(ctx context.Context, id uint, user models.User) *helpers.ResponseError
	DeleteUser(ctx context.Context, id uint) *helpers.ResponseError
	GetUserByEmail(ctx context.Context, email string) (*models.User, *helpers.ResponseError)
}
type UserService interface {
	CreateUser(ctx context.Context, user models.User) []helpers.ResponseError
	GetUser(ctx context.Context, id uint) (*models.User, []helpers.ResponseError)
	GetUsers(ctx context.Context, pagination models.Pagination, search models.Search) ([]models.User, *models.Pagination, *models.Search, []helpers.ResponseError)
	UpdateUser(ctx context.Context, id uint, user models.User) []helpers.ResponseError
	DeleteUser(ctx context.Context, id uint) []helpers.ResponseError
	GetUserByEmail(ctx context.Context, email string) (*models.User, []helpers.ResponseError)
}
