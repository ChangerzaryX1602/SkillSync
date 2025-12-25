package domain

import (
	"github.com/ChangerzaryX1602/SkillSync/pkg/models"

	helpers "github.com/zercle/gofiber-helpers"
)

type UserRepository interface {
	Migrate() error
	CreateUser(user models.User) *helpers.ResponseError
	GetUser(id uint) (*models.User, *helpers.ResponseError)
	GetUsers(pagination models.Pagination, search models.Search) ([]models.User, *models.Pagination, *models.Search, *helpers.ResponseError)
	UpdateUser(id uint, user models.User) *helpers.ResponseError
	DeleteUser(id uint) *helpers.ResponseError
	GetUserByEmail(email string) (*models.User, *helpers.ResponseError)
}
type UserService interface {
	CreateUser(user models.User) []helpers.ResponseError
	GetUser(id uint) (*models.User, []helpers.ResponseError)
	GetUsers(pagination models.Pagination, search models.Search) ([]models.User, *models.Pagination, *models.Search, []helpers.ResponseError)
	UpdateUser(id uint, user models.User) []helpers.ResponseError
	DeleteUser(id uint) []helpers.ResponseError
	GetUserByEmail(email string) (*models.User, []helpers.ResponseError)
}
