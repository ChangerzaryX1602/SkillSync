package domain

import (
	"context"

	"github.com/ChangerzaryX1602/SkillSync/pkg/models"

	helpers "github.com/zercle/gofiber-helpers"
)

type RolePermissionRepository interface {
	Migrate() error
	CreateRolePermission(ctx context.Context, rolePermission models.RolePermission) *helpers.ResponseError
	GetRolePermission(ctx context.Context, id uint) (*models.RolePermission, *helpers.ResponseError)
	GetRolePermissions(ctx context.Context, pagination models.Pagination, search models.Search) ([]models.RolePermission, *models.Pagination, *models.Search, *helpers.ResponseError)
	GetRolePermissionsByRoleID(ctx context.Context, roleID uint) ([]models.RolePermission, *helpers.ResponseError)
	UpdateRolePermission(ctx context.Context, id uint, rolePermission models.RolePermission) *helpers.ResponseError
	DeleteRolePermission(ctx context.Context, id uint) *helpers.ResponseError
	DeleteRolePermissionsByRoleID(ctx context.Context, roleID uint) *helpers.ResponseError
}

type RolePermissionService interface {
	CreateRolePermission(ctx context.Context, rolePermission models.RolePermission) []helpers.ResponseError
	GetRolePermission(ctx context.Context, id uint) (*models.RolePermission, []helpers.ResponseError)
	GetRolePermissions(ctx context.Context, pagination models.Pagination, search models.Search) ([]models.RolePermission, *models.Pagination, *models.Search, []helpers.ResponseError)
	GetRolePermissionsByRoleID(ctx context.Context, roleID uint) ([]models.RolePermission, []helpers.ResponseError)
	UpdateRolePermission(ctx context.Context, id uint, rolePermission models.RolePermission) []helpers.ResponseError
	DeleteRolePermission(ctx context.Context, id uint) []helpers.ResponseError
	DeleteRolePermissionsByRoleID(ctx context.Context, roleID uint) []helpers.ResponseError
}
