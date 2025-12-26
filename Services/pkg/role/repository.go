package role

import (
	"context"

	"github.com/ChangerzaryX1602/SkillSync/pkg/domain"
	"github.com/ChangerzaryX1602/SkillSync/pkg/models"
	"github.com/ChangerzaryX1602/SkillSync/pkg/utils"

	"github.com/gofiber/fiber/v2"
	helpers "github.com/zercle/gofiber-helpers"
	"gorm.io/gorm"
)

type roleRepository struct {
	resources models.Resources
}

func NewRoleRepository(resources models.Resources) domain.RoleRepository {
	return &roleRepository{resources: resources}
}

func (r *roleRepository) Migrate() error {
	return r.resources.MainDbConn.AutoMigrate(&models.Role{})
}

func (r *roleRepository) CreateRole(ctx context.Context, role models.Role) *helpers.ResponseError {
	if err := r.resources.MainDbConn.WithContext(ctx).Create(&role).Error; err != nil {
		return &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}
	return nil
}

func (r *roleRepository) GetRole(ctx context.Context, id uint) (*models.Role, *helpers.ResponseError) {
	var role models.Role
	if err := r.resources.MainDbConn.WithContext(ctx).Where("id = ?", id).First(&role).Error; err != nil {
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
	return &role, nil
}

func (r *roleRepository) GetRoles(ctx context.Context, pagination models.Pagination, search models.Search) ([]models.Role, *models.Pagination, *models.Search, *helpers.ResponseError) {
	var roles []models.Role
	db := r.resources.MainDbConn.WithContext(ctx).Model(&models.Role{})
	db = utils.ApplySearch(ctx, r.resources.FastHTTPClient, db, search, false)
	db = utils.ApplyPagination(db, &pagination, models.Role{})
	if err := db.Find(&roles).Error; err != nil {
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
	return roles, &pagination, &search, nil
}

func (r *roleRepository) UpdateRole(ctx context.Context, id uint, role models.Role) *helpers.ResponseError {
	if err := r.resources.MainDbConn.WithContext(ctx).Where("id = ?", id).First(&models.Role{}).Error; err != nil {
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
	if err := r.resources.MainDbConn.WithContext(ctx).Model(&models.Role{}).Where("id = ?", id).Updates(role).Error; err != nil {
		return &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}
	return nil
}

func (r *roleRepository) DeleteRole(ctx context.Context, id uint) *helpers.ResponseError {
	if err := r.resources.MainDbConn.WithContext(ctx).Where("id = ?", id).First(&models.Role{}).Error; err != nil {
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
	if err := r.resources.MainDbConn.WithContext(ctx).Delete(&models.Role{}, id).Error; err != nil {
		return &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}
	return nil
}
