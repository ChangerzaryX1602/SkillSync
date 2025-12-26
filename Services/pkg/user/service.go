package user

import (
	"context"

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
func (s *userService) CreateUser(ctx context.Context, user models.User) []helpers.ResponseError {
	err := s.repository.CreateUser(ctx, user)
	if err != nil {
		return []helpers.ResponseError{*err}
	}
	return nil
}
func (s *userService) GetUser(ctx context.Context, id uint) (*models.User, []helpers.ResponseError) {
	user, err := s.repository.GetUser(ctx, id)
	if err != nil {
		return nil, []helpers.ResponseError{*err}
	}
	return user, nil
}
func (s *userService) GetUsers(ctx context.Context, pagination models.Pagination, search models.Search) ([]models.User, *models.Pagination, *models.Search, []helpers.ResponseError) {
	users, paginated, searched, err := s.repository.GetUsers(ctx, pagination, search)
	if err != nil {
		return nil, nil, nil, []helpers.ResponseError{*err}
	}
	return users, paginated, searched, nil
}
func (s *userService) UpdateUser(ctx context.Context, id uint, user models.User) []helpers.ResponseError {
	err := s.repository.UpdateUser(ctx, id, user)
	if err != nil {
		return []helpers.ResponseError{*err}
	}
	return nil
}
func (s *userService) DeleteUser(ctx context.Context, id uint) []helpers.ResponseError {
	err := s.repository.DeleteUser(ctx, id)
	if err != nil {
		return []helpers.ResponseError{*err}
	}
	return nil
}
func (s *userService) GetUserByEmail(ctx context.Context, email string) (*models.User, []helpers.ResponseError) {
	user, err := s.repository.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, []helpers.ResponseError{*err}
	}
	return user, nil
}
