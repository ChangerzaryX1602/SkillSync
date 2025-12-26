package domain

import (
	"context"

	"github.com/ChangerzaryX1602/SkillSync/pkg/models"

	helpers "github.com/zercle/gofiber-helpers"
)

type PermissionRepository interface {
	Migrate() error
	CreatePermission(ctx context.Context, permission models.Permission) *helpers.ResponseError
	GetPermission(ctx context.Context, id uint) (*models.Permission, *helpers.ResponseError)
	GetPermissions(ctx context.Context, pagination models.Pagination, search models.Search) ([]models.Permission, *models.Pagination, *models.Search, *helpers.ResponseError)
	UpdatePermission(ctx context.Context, id uint, permission models.Permission) *helpers.ResponseError
	DeletePermission(ctx context.Context, id uint) *helpers.ResponseError
}

type PermissionService interface {
	CreatePermission(ctx context.Context, permission models.Permission) []helpers.ResponseError
	GetPermission(ctx context.Context, id uint) (*models.Permission, []helpers.ResponseError)
	GetPermissions(ctx context.Context, pagination models.Pagination, search models.Search) ([]models.Permission, *models.Pagination, *models.Search, []helpers.ResponseError)
	UpdatePermission(ctx context.Context, id uint, permission models.Permission) []helpers.ResponseError
	DeletePermission(ctx context.Context, id uint) []helpers.ResponseError
}
