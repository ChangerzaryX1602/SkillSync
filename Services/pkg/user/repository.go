package user

import (
	"fmt"

	"github.com/ChangerzaryX1602/SkillSync/pkg/domain"
	"github.com/ChangerzaryX1602/SkillSync/pkg/models"
	"github.com/ChangerzaryX1602/SkillSync/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/pgvector/pgvector-go"
	helpers "github.com/zercle/gofiber-helpers"
	"gorm.io/gorm"
)

type userRepository struct {
	resources models.Resources
}

func NewUserRepository(resources models.Resources) domain.UserRepository {
	return &userRepository{resources: resources}
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
func (r *userRepository) CreateUser(user models.User) *helpers.ResponseError {
	searchContext := fmt.Sprintf("%s %s", user.Username, user.Email)
	vec, err := utils.GenerateEmbedding(searchContext, r.resources.EmbeddingKey)
	if err != nil {
		fmt.Printf("Error generating embedding: %v\n", err)
	} else {
		user.Embedding = pgvector.NewVector(vec)
	}

	if err := r.resources.MainDbConn.Create(&user).Error; err != nil {
		return &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}
	var existingRole models.Role
	if err := r.resources.MainDbConn.Where("name = ?", utils.User).First(&existingRole).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &helpers.ResponseError{
				Code:    fiber.StatusBadRequest,
				Source:  helpers.WhereAmI(),
				Title:   "Invalid Role",
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
	var existingUserRole models.UserRole
	if err := r.resources.MainDbConn.Where("user_id = ? AND role_id = ?", user.ID, existingRole.ID).First(&existingUserRole).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			userRole := models.UserRole{
				UserID: user.ID,
				RoleID: existingRole.ID,
			}
			if err := r.resources.MainDbConn.Create(&userRole).Error; err != nil {
				return &helpers.ResponseError{
					Code:    fiber.StatusInternalServerError,
					Source:  helpers.WhereAmI(),
					Title:   "Database Error",
					Message: err.Error(),
				}
			}
		} else {
			return &helpers.ResponseError{
				Code:    fiber.StatusInternalServerError,
				Source:  helpers.WhereAmI(),
				Title:   "Database Error",
				Message: err.Error(),
			}
		}
	}
	return nil
}
func (r *userRepository) GetUser(id uint) (*models.User, *helpers.ResponseError) {
	var user models.User
	if err := r.resources.MainDbConn.Where("id = ?", id).First(&user).Error; err != nil {
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
	return &user, nil
}
func (r *userRepository) GetUsers(pagination models.Pagination, search models.Search) ([]models.User, *models.Pagination, *models.Search, *helpers.ResponseError) {
	var users []models.User
	db := r.resources.MainDbConn.Model(&models.User{})
	db = utils.ApplySearch(db, search, true, r.resources.EmbeddingKey)
	if db.Error != nil {
		return nil, nil, nil, &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: db.Error.Error(),
		}
	}
	db = utils.ApplyPagination(db, &pagination, models.User{})
	if err := db.Debug().Find(&users).Error; err != nil {
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
	return users, &pagination, &search, nil
}
func (r *userRepository) UpdateUser(id uint, user models.User) *helpers.ResponseError {
	if err := r.resources.MainDbConn.Where("id = ?", id).First(&models.User{}).Error; err != nil {
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
	if err := r.resources.MainDbConn.Model(&models.User{}).Where("id = ?", id).Updates(user).Error; err != nil {
		return &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}
	return nil
}
func (r *userRepository) DeleteUser(id uint) *helpers.ResponseError {
	if err := r.resources.MainDbConn.Where("id = ?", id).First(&models.User{}).Error; err != nil {
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
	if err := r.resources.MainDbConn.Model(&models.User{}).Where("id = ?", id).Delete(&models.User{}).Error; err != nil {
		return &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}
	return nil
}
func (r *userRepository) GetUserByEmail(email string) (*models.User, *helpers.ResponseError) {
	var user models.User
	if err := r.resources.MainDbConn.Where("email = ?", email).First(&user).Error; err != nil {
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
	return &user, nil
}
