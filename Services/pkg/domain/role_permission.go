package domain

import (
	"github.com/ChangerzaryX1602/SkillSync/pkg/models"

	helpers "github.com/zercle/gofiber-helpers"
)

type RolePermissionRepository interface {
	Migrate() error
	CreateRolePermission(rolePermission models.RolePermission) *helpers.ResponseError
	GetRolePermission(id uint) (*models.RolePermission, *helpers.ResponseError)
	GetRolePermissions(pagination models.Pagination, search models.Search) ([]models.RolePermission, *models.Pagination, *models.Search, *helpers.ResponseError)
	GetRolePermissionsByRoleID(roleID uint) ([]models.RolePermission, *helpers.ResponseError)
	UpdateRolePermission(id uint, rolePermission models.RolePermission) *helpers.ResponseError
	DeleteRolePermission(id uint) *helpers.ResponseError
	DeleteRolePermissionsByRoleID(roleID uint) *helpers.ResponseError
}

type RolePermissionService interface {
	CreateRolePermission(rolePermission models.RolePermission) []helpers.ResponseError
	GetRolePermission(id uint) (*models.RolePermission, []helpers.ResponseError)
	GetRolePermissions(pagination models.Pagination, search models.Search) ([]models.RolePermission, *models.Pagination, *models.Search, []helpers.ResponseError)
	GetRolePermissionsByRoleID(roleID uint) ([]models.RolePermission, []helpers.ResponseError)
	UpdateRolePermission(id uint, rolePermission models.RolePermission) []helpers.ResponseError
	DeleteRolePermission(id uint) []helpers.ResponseError
	DeleteRolePermissionsByRoleID(roleID uint) []helpers.ResponseError
}
