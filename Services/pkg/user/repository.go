package user

import (
	"context"

	"github.com/ChangerzaryX1602/SkillSync/pkg/domain"
	"github.com/ChangerzaryX1602/SkillSync/pkg/models"
	txcontext "github.com/ChangerzaryX1602/SkillSync/pkg/tx_context"
	"github.com/ChangerzaryX1602/SkillSync/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/pgvector/pgvector-go"
	helpers "github.com/zercle/gofiber-helpers"
	"gorm.io/gorm"
)

type userRepository struct {
	resources models.Resources
	cache     *userCache
}

func NewUserRepository(resources models.Resources) domain.UserRepository {
	return &userRepository{resources: resources, cache: newUserCache(resources.RedisStorage)}
}
func (r *userRepository) Migrate() error {
	if err := r.resources.MainDbConn.AutoMigrate(&models.User{}); err != nil {
		return err
	}
	if err := r.resources.MainDbConn.Exec(`
    CREATE INDEX IF NOT EXISTS idx_users_embedding 
    ON users 
    USING hnsw (embedding vector_cosine_ops)
`).Error; err != nil {
		return err
	}
	return nil
}
func (r *userRepository) CreateUser(ctx context.Context, user models.User) *helpers.ResponseError {
	db := r.getDB(ctx)
	vec, err := utils.GenerateEmbeddingByOllama(ctx, r.resources.FastHTTPClient, user.GenerateSearchContext())
	if err == nil {
		user.Embedding = pgvector.NewVector(vec)
	}

	if err := db.WithContext(ctx).Create(&user).Error; err != nil {
		return &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}
	return nil
}
func (r *userRepository) GetUser(ctx context.Context, id uint) (*models.User, *helpers.ResponseError) {
	if userCache, ok := r.cache.Get(id); ok {
		return userCache, nil
	}
	var user models.User
	db := r.getDB(ctx)
	if err := db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
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
	r.cache.Set(&user)
	return &user, nil
}
func (r *userRepository) GetUsers(ctx context.Context, pagination models.Pagination, search models.Search) ([]models.User, *models.Pagination, *models.Search, *helpers.ResponseError) {
	if users, found := r.cache.GetList(pagination, search); found {
		return users, &pagination, &search, nil
	}
	var users []models.User
	db := r.getDB(ctx)
	db = db.WithContext(ctx).Model(&models.User{})
	db = utils.ApplySearch(ctx, r.resources.FastHTTPClient, db, search, true)
	if db.Error != nil {
		return nil, nil, nil, &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: db.Error.Error(),
		}
	}
	db = utils.ApplyPagination(db, &pagination, models.User{})
	if err := db.Find(&users).Error; err != nil {
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
	r.cache.SetList(pagination, search, users)
	return users, &pagination, &search, nil
}
func (r *userRepository) UpdateUser(ctx context.Context, id uint, user models.User) *helpers.ResponseError {
	db := r.getDB(ctx)
	vec, err := utils.GenerateEmbeddingByOllama(ctx, r.resources.FastHTTPClient, user.GenerateSearchContext())
	if err == nil {
		user.Embedding = pgvector.NewVector(vec)
	}
	if err := db.WithContext(ctx).Where("id = ?", id).First(&models.User{}).Error; err != nil {
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
	if err := db.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Updates(user).Error; err != nil {
		return &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}
	r.cache.Invalidate(id)
	return nil
}
func (r *userRepository) DeleteUser(ctx context.Context, id uint) *helpers.ResponseError {
	db := r.getDB(ctx)
	if err := db.WithContext(ctx).Where("id = ?", id).First(&models.User{}).Error; err != nil {
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
	if err := db.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Delete(&models.User{}).Error; err != nil {
		return &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}
	r.cache.Invalidate(id)
	return nil
}
func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, *helpers.ResponseError) {
	if user, found := r.cache.GetByEmail(email); found {
		return user, nil
	}
	db := r.getDB(ctx)
	var user models.User
	if err := db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
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
	r.cache.SetByEmail(email, &user)
	return &user, nil
}
func (r *userRepository) getDB(ctx context.Context) *gorm.DB {
	if tx := txcontext.GetTxFromContext(ctx); tx != nil {
		return tx
	}
	return r.resources.MainDbConn
}
