package permission

import (
	"context"

	"github.com/ChangerzaryX1602/SkillSync/pkg/domain"
	"github.com/ChangerzaryX1602/SkillSync/pkg/models"
	"github.com/ChangerzaryX1602/SkillSync/pkg/utils"

	"github.com/gofiber/fiber/v2"
	helpers "github.com/zercle/gofiber-helpers"
	"gorm.io/gorm"
)

type permissionRepository struct {
	resources models.Resources
}

func NewPermissionRepository(resources models.Resources) domain.PermissionRepository {
	return &permissionRepository{resources: resources}
}

func (r *permissionRepository) Migrate() error {
	return r.resources.MainDbConn.AutoMigrate(&models.Permission{})
}

func (r *permissionRepository) CreatePermission(ctx context.Context, permission models.Permission) *helpers.ResponseError {
	if err := r.resources.MainDbConn.WithContext(ctx).Create(&permission).Error; err != nil {
		return &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}
	return nil
}

func (r *permissionRepository) GetPermission(ctx context.Context, id uint) (*models.Permission, *helpers.ResponseError) {
	var permission models.Permission
	if err := r.resources.MainDbConn.WithContext(ctx).Where("id = ?", id).First(&permission).Error; err != nil {
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
	return &permission, nil
}

func (r *permissionRepository) GetPermissions(ctx context.Context, pagination models.Pagination, search models.Search) ([]models.Permission, *models.Pagination, *models.Search, *helpers.ResponseError) {
	var permissions []models.Permission
	db := r.resources.MainDbConn.WithContext(ctx).Model(&models.Permission{})
	db = utils.ApplySearch(ctx, r.resources.FastHTTPClient, db, search, false)
	db = utils.ApplyPagination(db, &pagination, models.Permission{})
	if err := db.Find(&permissions).Error; err != nil {
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
	return permissions, &pagination, &search, nil
}

func (r *permissionRepository) UpdatePermission(ctx context.Context, id uint, permission models.Permission) *helpers.ResponseError {
	if err := r.resources.MainDbConn.WithContext(ctx).Where("id = ?", id).First(&models.Permission{}).Error; err != nil {
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
	if err := r.resources.MainDbConn.WithContext(ctx).Model(&models.Permission{}).Where("id = ?", id).Updates(permission).Error; err != nil {
		return &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}
	return nil
}

func (r *permissionRepository) DeletePermission(ctx context.Context, id uint) *helpers.ResponseError {
	if err := r.resources.MainDbConn.WithContext(ctx).Where("id = ?", id).First(&models.Permission{}).Error; err != nil {
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
	if err := r.resources.MainDbConn.WithContext(ctx).Delete(&models.Permission{}, id).Error; err != nil {
		return &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}
	return nil
}
