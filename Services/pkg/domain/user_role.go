package domain

import (
	"github.com/ChangerzaryX1602/SkillSync/pkg/models"

	helpers "github.com/zercle/gofiber-helpers"
)

type UserRoleRepository interface {
	Migrate() error
	CreateUserRole(userRole models.UserRole) *helpers.ResponseError
	GetUserRole(id uint) (*models.UserRole, *helpers.ResponseError)
	GetUserRoles(pagination models.Pagination, search models.Search) ([]models.UserRole, *models.Pagination, *models.Search, *helpers.ResponseError)
	GetUserRolesByUserID(userID uint) ([]models.UserRole, *helpers.ResponseError)
	UpdateUserRole(id uint, userRole models.UserRole) *helpers.ResponseError
	DeleteUserRole(id uint) *helpers.ResponseError
	DeleteUserRolesByUserID(userID uint) *helpers.ResponseError
}

type UserRoleService interface {
	CreateUserRole(userRole models.UserRole) []helpers.ResponseError
	GetUserRole(id uint) (*models.UserRole, []helpers.ResponseError)
	GetUserRoles(pagination models.Pagination, search models.Search) ([]models.UserRole, *models.Pagination, *models.Search, []helpers.ResponseError)
	GetUserRolesByUserID(userID uint) ([]models.UserRole, []helpers.ResponseError)
	UpdateUserRole(id uint, userRole models.UserRole) []helpers.ResponseError
	DeleteUserRole(id uint) []helpers.ResponseError
	DeleteUserRolesByUserID(userID uint) []helpers.ResponseError
}
