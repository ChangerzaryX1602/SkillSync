package user_role

import (
	"github.com/ChangerzaryX1602/SkillSync/pkg/domain"
	"github.com/ChangerzaryX1602/SkillSync/pkg/models"
	"github.com/ChangerzaryX1602/SkillSync/pkg/utils"

	"github.com/gofiber/fiber/v2"
	helpers "github.com/zercle/gofiber-helpers"
	"gorm.io/gorm"
)

type userRoleRepository struct {
	resources models.Resources
}

func NewUserRoleRepository(resources models.Resources) domain.UserRoleRepository {
	return &userRoleRepository{resources: resources}
}

func (r *userRoleRepository) Migrate() error {
	return r.resources.MainDbConn.AutoMigrate(&models.UserRole{})
}

func (r *userRoleRepository) CreateUserRole(userRole models.UserRole) *helpers.ResponseError {
	if err := r.resources.MainDbConn.Create(&userRole).Error; err != nil {
		return &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}
	return nil
}

func (r *userRoleRepository) GetUserRole(id uint) (*models.UserRole, *helpers.ResponseError) {
	var userRole models.UserRole
	if err := r.resources.MainDbConn.Where("id = ?", id).First(&userRole).Error; err != nil {
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
	return &userRole, nil
}

func (r *userRoleRepository) GetUserRoles(pagination models.Pagination, search models.Search) ([]models.UserRole, *models.Pagination, *models.Search, *helpers.ResponseError) {
	var userRoles []models.UserRole
	db := r.resources.MainDbConn.Model(&models.UserRole{})
	db = utils.ApplySearch(db, search, false)
	db = utils.ApplyPagination(db, &pagination, models.UserRole{})
	if err := db.Find(&userRoles).Error; err != nil {
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
	return userRoles, &pagination, &search, nil
}

func (r *userRoleRepository) GetUserRolesByUserID(userID uint) ([]models.UserRole, *helpers.ResponseError) {
	var userRoles []models.UserRole
	if err := r.resources.MainDbConn.Where("user_id = ?", userID).Find(&userRoles).Error; err != nil {
		return nil, &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}
	return userRoles, nil
}

func (r *userRoleRepository) UpdateUserRole(id uint, userRole models.UserRole) *helpers.ResponseError {
	if err := r.resources.MainDbConn.Where("id = ?", id).First(&models.UserRole{}).Error; err != nil {
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
	if err := r.resources.MainDbConn.Model(&models.UserRole{}).Where("id = ?", id).Updates(userRole).Error; err != nil {
		return &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}
	return nil
}

func (r *userRoleRepository) DeleteUserRole(id uint) *helpers.ResponseError {
	if err := r.resources.MainDbConn.Where("id = ?", id).First(&models.UserRole{}).Error; err != nil {
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
	if err := r.resources.MainDbConn.Delete(&models.UserRole{}, id).Error; err != nil {
		return &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}
	return nil
}

func (r *userRoleRepository) DeleteUserRolesByUserID(userID uint) *helpers.ResponseError {
	if err := r.resources.MainDbConn.Where("user_id = ?", userID).Delete(&models.UserRole{}).Error; err != nil {
		return &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}
	return nil
}
