package domain

import (
	"context"

	"github.com/ChangerzaryX1602/SkillSync/pkg/models"

	helpers "github.com/zercle/gofiber-helpers"
)

type UserRoleRepository interface {
	Migrate() error
	CreateUserRole(ctx context.Context, userRole models.UserRole) *helpers.ResponseError
	GetUserRole(ctx context.Context, id uint) (*models.UserRole, *helpers.ResponseError)
	GetUserRoles(ctx context.Context, pagination models.Pagination, search models.Search) ([]models.UserRole, *models.Pagination, *models.Search, *helpers.ResponseError)
	GetUserRolesByUserID(ctx context.Context, userID uint) ([]models.UserRole, *helpers.ResponseError)
	UpdateUserRole(ctx context.Context, id uint, userRole models.UserRole) *helpers.ResponseError
	DeleteUserRole(ctx context.Context, id uint) *helpers.ResponseError
	DeleteUserRolesByUserID(ctx context.Context, userID uint) *helpers.ResponseError
}

type UserRoleService interface {
	CreateUserRole(ctx context.Context, userRole models.UserRole) []helpers.ResponseError
	GetUserRole(ctx context.Context, id uint) (*models.UserRole, []helpers.ResponseError)
	GetUserRoles(ctx context.Context, pagination models.Pagination, search models.Search) ([]models.UserRole, *models.Pagination, *models.Search, []helpers.ResponseError)
	GetUserRolesByUserID(ctx context.Context, userID uint) ([]models.UserRole, []helpers.ResponseError)
	UpdateUserRole(ctx context.Context, id uint, userRole models.UserRole) []helpers.ResponseError
	DeleteUserRole(ctx context.Context, id uint) []helpers.ResponseError
	DeleteUserRolesByUserID(ctx context.Context, userID uint) []helpers.ResponseError
}
