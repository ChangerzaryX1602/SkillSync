package role_permission

import (
	"github.com/ChangerzaryX1602/SkillSync/pkg/domain"
	"github.com/ChangerzaryX1602/SkillSync/pkg/models"
	"github.com/ChangerzaryX1602/SkillSync/pkg/utils"

	"github.com/gofiber/fiber/v2"
	helpers "github.com/zercle/gofiber-helpers"
	"gorm.io/gorm"
)

type rolePermissionRepository struct {
	resources models.Resources
}

func NewRolePermissionRepository(resources models.Resources) domain.RolePermissionRepository {
	return &rolePermissionRepository{resources: resources}
}

func (r *rolePermissionRepository) Migrate() error {
	return r.resources.MainDbConn.AutoMigrate(&models.RolePermission{})
}

func (r *rolePermissionRepository) CreateRolePermission(rolePermission models.RolePermission) *helpers.ResponseError {
	if err := r.resources.MainDbConn.Create(&rolePermission).Error; err != nil {
		return &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}
	return nil
}

func (r *rolePermissionRepository) GetRolePermission(id uint) (*models.RolePermission, *helpers.ResponseError) {
	var rolePermission models.RolePermission
	if err := r.resources.MainDbConn.Where("id = ?", id).First(&rolePermission).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &helpers.ResponseError{
				Code:    fiber.StatusNotFound,
				Source:  helpers.WhereAmI(),
				Title:   "Not Found",
				Message: err.Error(),
			}
		}
		return nil, &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}
	return &rolePermission, nil
}

func (r *rolePermissionRepository) GetRolePermissions(pagination models.Pagination, search models.Search) ([]models.RolePermission, *models.Pagination, *models.Search, *helpers.ResponseError) {
	var rolePermissions []models.RolePermission
	db := r.resources.MainDbConn.Model(&models.RolePermission{})
	db = utils.ApplySearch(db, search, false)
	db = utils.ApplyPagination(db, &pagination, models.RolePermission{})
	if err := db.Find(&rolePermissions).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil, nil, &helpers.ResponseError{
				Code:    fiber.StatusNotFound,
				Source:  helpers.WhereAmI(),
				Title:   "Not Found",
				Message: err.Error(),
			}
		}
		return nil, nil, nil, &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}
	return rolePermissions, &pagination, &search, nil
}

func (r *rolePermissionRepository) GetRolePermissionsByRoleID(roleID uint) ([]models.RolePermission, *helpers.ResponseError) {
	var rolePermissions []models.RolePermission
	if err := r.resources.MainDbConn.Where("role_id = ?", roleID).Find(&rolePermissions).Error; err != nil {
		return nil, &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}
	return rolePermissions, nil
}

func (r *rolePermissionRepository) UpdateRolePermission(id uint, rolePermission models.RolePermission) *helpers.ResponseError {
	if err := r.resources.MainDbConn.Where("id = ?", id).First(&models.RolePermission{}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &helpers.ResponseError{
				Code:    fiber.StatusNotFound,
				Source:  helpers.WhereAmI(),
				Title:   "Not Found",
				Message: err.Error(),
			}
		}
		return &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}
	if err := r.resources.MainDbConn.Model(&models.RolePermission{}).Where("id = ?", id).Updates(rolePermission).Error; err != nil {
		return &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}
	return nil
}

func (r *rolePermissionRepository) DeleteRolePermission(id uint) *helpers.ResponseError {
	if err := r.resources.MainDbConn.Where("id = ?", id).First(&models.RolePermission{}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &helpers.ResponseError{
				Code:    fiber.StatusNotFound,
				Source:  helpers.WhereAmI(),
				Title:   "Not Found",
				Message: err.Error(),
			}
		}
		return &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}
	if err := r.resources.MainDbConn.Delete(&models.RolePermission{}, id).Error; err != nil {
		return &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}
	return nil
}

func (r *rolePermissionRepository) DeleteRolePermissionsByRoleID(roleID uint) *helpers.ResponseError {
	if err := r.resources.MainDbConn.Where("role_id = ?", roleID).Delete(&models.RolePermission{}).Error; err != nil {
		return &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}
	return nil
}
