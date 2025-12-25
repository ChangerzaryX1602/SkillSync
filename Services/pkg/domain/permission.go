package domain

import (
	"github.com/ChangerzaryX1602/SkillSync/pkg/models"

	helpers "github.com/zercle/gofiber-helpers"
)

type PermissionRepository interface {
	Migrate() error
	CreatePermission(permission models.Permission) *helpers.ResponseError
	GetPermission(id uint) (*models.Permission, *helpers.ResponseError)
	GetPermissions(pagination models.Pagination, search models.Search) ([]models.Permission, *models.Pagination, *models.Search, *helpers.ResponseError)
	UpdatePermission(id uint, permission models.Permission) *helpers.ResponseError
	DeletePermission(id uint) *helpers.ResponseError
}

type PermissionService interface {
	CreatePermission(permission models.Permission) []helpers.ResponseError
	GetPermission(id uint) (*models.Permission, []helpers.ResponseError)
	GetPermissions(pagination models.Pagination, search models.Search) ([]models.Permission, *models.Pagination, *models.Search, []helpers.ResponseError)
	UpdatePermission(id uint, permission models.Permission) []helpers.ResponseError
	DeletePermission(id uint) []helpers.ResponseError
}
