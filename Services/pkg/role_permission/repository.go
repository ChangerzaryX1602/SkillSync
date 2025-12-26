package role_permission

import (
	"context"

	"github.com/ChangerzaryX1602/SkillSync/pkg/domain"
	"github.com/ChangerzaryX1602/SkillSync/pkg/models"
	txcontext "github.com/ChangerzaryX1602/SkillSync/pkg/tx_context"
	"github.com/ChangerzaryX1602/SkillSync/pkg/utils"

	"github.com/gofiber/fiber/v2"
	helpers "github.com/zercle/gofiber-helpers"
	"gorm.io/gorm"
)

type rolePermissionRepository struct {
	resources models.Resources
	cache     *rolePermissionCache
}

func NewRolePermissionRepository(resources models.Resources) domain.RolePermissionRepository {
	return &rolePermissionRepository{resources: resources, cache: newRolePermissionCache(resources.RedisStorage)}
}

func (r *rolePermissionRepository) Migrate() error {
	return r.resources.MainDbConn.AutoMigrate(&models.RolePermission{})
}

func (r *rolePermissionRepository) CreateRolePermission(ctx context.Context, rolePermission models.RolePermission) *helpers.ResponseError {
	db := r.getDB(ctx)
	if err := db.WithContext(ctx).Create(&rolePermission).Error; err != nil {
		return &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}
	r.cache.InvalidateByRoleID(rolePermission.RoleID)
	return nil
}

func (r *rolePermissionRepository) GetRolePermission(ctx context.Context, id uint) (*models.RolePermission, *helpers.ResponseError) {
	db := r.getDB(ctx)
	var rolePermission models.RolePermission
	if err := db.WithContext(ctx).Where("id = ?", id).First(&rolePermission).Error; err != nil {
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

func (r *rolePermissionRepository) GetRolePermissions(ctx context.Context, pagination models.Pagination, search models.Search) ([]models.RolePermission, *models.Pagination, *models.Search, *helpers.ResponseError) {
	db := r.getDB(ctx)
	var rolePermissions []models.RolePermission
	db = db.WithContext(ctx).Model(&models.RolePermission{})
	db = utils.ApplySearch(ctx, r.resources.FastHTTPClient, db, search, false)
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

func (r *rolePermissionRepository) GetRolePermissionsByRoleID(ctx context.Context, roleID uint) ([]models.RolePermission, *helpers.ResponseError) {
	if rolePermissions, found := r.cache.GetByRoleID(roleID); found {
		return rolePermissions, nil
	}
	db := r.getDB(ctx)
	var rolePermissions []models.RolePermission
	if err := db.WithContext(ctx).Where("role_id = ?", roleID).Find(&rolePermissions).Error; err != nil {
		return nil, &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}
	r.cache.SetByRoleID(roleID, rolePermissions)
	return rolePermissions, nil
}

func (r *rolePermissionRepository) UpdateRolePermission(ctx context.Context, id uint, rolePermission models.RolePermission) *helpers.ResponseError {
	db := r.getDB(ctx)
	if err := db.WithContext(ctx).Where("id = ?", id).First(&models.RolePermission{}).Error; err != nil {
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
	if err := db.WithContext(ctx).Model(&models.RolePermission{}).Where("id = ?", id).Updates(rolePermission).Error; err != nil {
		return &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}
	r.cache.InvalidateByRoleID(rolePermission.RoleID)
	return nil
}

func (r *rolePermissionRepository) DeleteRolePermission(ctx context.Context, id uint) *helpers.ResponseError {
	db := r.getDB(ctx)
	var rolePermission models.RolePermission
	if err := db.WithContext(ctx).Where("id = ?", id).First(&rolePermission).Error; err != nil {
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
	if err := db.WithContext(ctx).Delete(&models.RolePermission{}, id).Error; err != nil {
		return &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}
	r.cache.InvalidateByRoleID(rolePermission.RoleID)
	return nil
}

func (r *rolePermissionRepository) DeleteRolePermissionsByRoleID(ctx context.Context, roleID uint) *helpers.ResponseError {
	db := r.getDB(ctx)
	if err := db.WithContext(ctx).Where("role_id = ?", roleID).Delete(&models.RolePermission{}).Error; err != nil {
		return &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}
	r.cache.InvalidateByRoleID(roleID)
	return nil
}

func (r *rolePermissionRepository) getDB(ctx context.Context) *gorm.DB {
	if tx := txcontext.GetTxFromContext(ctx); tx != nil {
		return tx
	}
	return r.resources.MainDbConn
}
