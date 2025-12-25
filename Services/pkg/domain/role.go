package domain

import (
	"github.com/ChangerzaryX1602/SkillSync/pkg/models"

	helpers "github.com/zercle/gofiber-helpers"
)

type RoleRepository interface {
	Migrate() error
	CreateRole(role models.Role) *helpers.ResponseError
	GetRole(id uint) (*models.Role, *helpers.ResponseError)
	GetRoles(pagination models.Pagination, search models.Search) ([]models.Role, *models.Pagination, *models.Search, *helpers.ResponseError)
	UpdateRole(id uint, role models.Role) *helpers.ResponseError
	DeleteRole(id uint) *helpers.ResponseError
}

type RoleService interface {
	CreateRole(role models.Role) []helpers.ResponseError
	GetRole(id uint) (*models.Role, []helpers.ResponseError)
	GetRoles(pagination models.Pagination, search models.Search) ([]models.Role, *models.Pagination, *models.Search, []helpers.ResponseError)
	UpdateRole(id uint, role models.Role) []helpers.ResponseError
	DeleteRole(id uint) []helpers.ResponseError
}
