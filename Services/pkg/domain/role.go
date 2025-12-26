package domain

import (
	"context"

	"github.com/ChangerzaryX1602/SkillSync/pkg/models"

	helpers "github.com/zercle/gofiber-helpers"
)

type RoleRepository interface {
	Migrate() error
	CreateRole(ctx context.Context, role models.Role) *helpers.ResponseError
	GetRole(ctx context.Context, id uint) (*models.Role, *helpers.ResponseError)
	GetRoles(ctx context.Context, pagination models.Pagination, search models.Search) ([]models.Role, *models.Pagination, *models.Search, *helpers.ResponseError)
	GetRoleByName(ctx context.Context, name string) (*models.Role, *helpers.ResponseError)
	UpdateRole(ctx context.Context, id uint, role models.Role) *helpers.ResponseError
	DeleteRole(ctx context.Context, id uint) *helpers.ResponseError
}

type RoleService interface {
	CreateRole(ctx context.Context, role models.Role) []helpers.ResponseError
	GetRole(ctx context.Context, id uint) (*models.Role, []helpers.ResponseError)
	GetRoles(ctx context.Context, pagination models.Pagination, search models.Search) ([]models.Role, *models.Pagination, *models.Search, []helpers.ResponseError)
	UpdateRole(ctx context.Context, id uint, role models.Role) []helpers.ResponseError
	DeleteRole(ctx context.Context, id uint) []helpers.ResponseError
}
