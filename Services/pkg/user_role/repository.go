package user_role

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

type userRoleRepository struct {
	resources models.Resources
	cache     *userRoleCache
}

func NewUserRoleRepository(resources models.Resources) domain.UserRoleRepository {
	return &userRoleRepository{resources: resources, cache: newUserRoleCache(resources.RedisStorage)}
}

func (r *userRoleRepository) Migrate() error {
	return r.resources.MainDbConn.AutoMigrate(&models.UserRole{})
}

func (r *userRoleRepository) CreateUserRole(ctx context.Context, userRole models.UserRole) *helpers.ResponseError {
	db := r.getDB(ctx)
	if err := db.WithContext(ctx).Create(&userRole).Error; err != nil {
		return &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}
	r.cache.InvalidateByUserID(userRole.UserID)
	return nil
}

func (r *userRoleRepository) GetUserRole(ctx context.Context, id uint) (*models.UserRole, *helpers.ResponseError) {
	db := r.getDB(ctx)
	var userRole models.UserRole
	if err := db.WithContext(ctx).Where("id = ?", id).First(&userRole).Error; err != nil {
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

func (r *userRoleRepository) GetUserRoles(ctx context.Context, pagination models.Pagination, search models.Search) ([]models.UserRole, *models.Pagination, *models.Search, *helpers.ResponseError) {
	db := r.getDB(ctx)
	var userRoles []models.UserRole
	db = db.WithContext(ctx).Model(&models.UserRole{})
	db = utils.ApplySearch(ctx, r.resources.FastHTTPClient, db, search, false)
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

func (r *userRoleRepository) GetUserRolesByUserID(ctx context.Context, userID uint) ([]models.UserRole, *helpers.ResponseError) {
	if userRoles, found := r.cache.GetByUserID(userID); found {
		return userRoles, nil
	}
	db := r.getDB(ctx)
	var userRoles []models.UserRole
	if err := db.WithContext(ctx).Where("user_id = ?", userID).Find(&userRoles).Error; err != nil {
		return nil, &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}
	r.cache.SetByUserID(userID, userRoles)
	return userRoles, nil
}

func (r *userRoleRepository) UpdateUserRole(ctx context.Context, id uint, userRole models.UserRole) *helpers.ResponseError {
	db := r.getDB(ctx)
	if err := db.WithContext(ctx).Where("id = ?", id).First(&models.UserRole{}).Error; err != nil {
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
	if err := db.WithContext(ctx).Model(&models.UserRole{}).Where("id = ?", id).Updates(userRole).Error; err != nil {
		return &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}
	r.cache.InvalidateByUserID(userRole.UserID)
	return nil
}

func (r *userRoleRepository) DeleteUserRole(ctx context.Context, id uint) *helpers.ResponseError {
	db := r.getDB(ctx)
	var userRole models.UserRole
	if err := db.WithContext(ctx).Where("id = ?", id).First(&userRole).Error; err != nil {
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
	if err := db.WithContext(ctx).Delete(&models.UserRole{}, id).Error; err != nil {
		return &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}
	r.cache.InvalidateByUserID(userRole.UserID)
	return nil
}

func (r *userRoleRepository) DeleteUserRolesByUserID(ctx context.Context, userID uint) *helpers.ResponseError {
	db := r.getDB(ctx)
	if err := db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.UserRole{}).Error; err != nil {
		return &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}
	r.cache.InvalidateByUserID(userID)
	return nil
}

func (r *userRoleRepository) getDB(ctx context.Context) *gorm.DB {
	if tx := txcontext.GetTxFromContext(ctx); tx != nil {
		return tx
	}
	return r.resources.MainDbConn
}
