package user

import (
	"github.com/ChangerzaryX1602/SkillSync/pkg/domain"
	"github.com/ChangerzaryX1602/SkillSync/pkg/models"

	helpers "github.com/zercle/gofiber-helpers"
)

type userService struct {
	repository domain.UserRepository
}

func NewUserService(repository domain.UserRepository) domain.UserService {
	return &userService{repository: repository}
}
func (s *userService) CreateUser(user models.User) []helpers.ResponseError {
	err := s.repository.CreateUser(user)
	if err != nil {
		return []helpers.ResponseError{*err}
	}
	return nil
}
func (s *userService) GetUser(id uint) (*models.User, []helpers.ResponseError) {
	user, err := s.repository.GetUser(id)
	if err != nil {
		return nil, []helpers.ResponseError{*err}
	}
	return user, nil
}
func (s *userService) GetUsers(pagination models.Pagination, search models.Search) ([]models.User, *models.Pagination, *models.Search, []helpers.ResponseError) {
	users, paginated, searched, err := s.repository.GetUsers(pagination, search)
	if err != nil {
		return nil, nil, nil, []helpers.ResponseError{*err}
	}
	return users, paginated, searched, nil
}
func (s *userService) UpdateUser(id uint, user models.User) []helpers.ResponseError {
	err := s.repository.UpdateUser(id, user)
	if err != nil {
		return []helpers.ResponseError{*err}
	}
	return nil
}
func (s *userService) DeleteUser(id uint) []helpers.ResponseError {
	err := s.repository.DeleteUser(id)
	if err != nil {
		return []helpers.ResponseError{*err}
	}
	return nil
}
func (s *userService) GetUserByEmail(email string) (*models.User, []helpers.ResponseError) {
	user, err := s.repository.GetUserByEmail(email)
	if err != nil {
		return nil, []helpers.ResponseError{*err}
	}
	return user, nil
}
