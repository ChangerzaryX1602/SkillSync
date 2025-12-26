package user_role

import (
	"context"

	"github.com/ChangerzaryX1602/SkillSync/pkg/domain"
	"github.com/ChangerzaryX1602/SkillSync/pkg/models"

	helpers "github.com/zercle/gofiber-helpers"
)

type userRoleService struct {
	repository domain.UserRoleRepository
}

func NewUserRoleService(repository domain.UserRoleRepository) domain.UserRoleService {
	return &userRoleService{repository: repository}
}

func (s *userRoleService) CreateUserRole(ctx context.Context, userRole models.UserRole) []helpers.ResponseError {
	err := s.repository.CreateUserRole(ctx, userRole)
	if err != nil {
		return []helpers.ResponseError{*err}
	}
	return nil
}

func (s *userRoleService) GetUserRole(ctx context.Context, id uint) (*models.UserRole, []helpers.ResponseError) {
	userRole, err := s.repository.GetUserRole(ctx, id)
	if err != nil {
		return nil, []helpers.ResponseError{*err}
	}
	return userRole, nil
}

func (s *userRoleService) GetUserRoles(ctx context.Context, pagination models.Pagination, search models.Search) ([]models.UserRole, *models.Pagination, *models.Search, []helpers.ResponseError) {
	userRoles, paginated, searched, err := s.repository.GetUserRoles(ctx, pagination, search)
	if err != nil {
		return nil, nil, nil, []helpers.ResponseError{*err}
	}
	return userRoles, paginated, searched, nil
}

func (s *userRoleService) GetUserRolesByUserID(ctx context.Context, userID uint) ([]models.UserRole, []helpers.ResponseError) {
	userRoles, err := s.repository.GetUserRolesByUserID(ctx, userID)
	if err != nil {
		return nil, []helpers.ResponseError{*err}
	}
	return userRoles, nil
}

func (s *userRoleService) UpdateUserRole(ctx context.Context, id uint, userRole models.UserRole) []helpers.ResponseError {
	err := s.repository.UpdateUserRole(ctx, id, userRole)
	if err != nil {
		return []helpers.ResponseError{*err}
	}
	return nil
}

func (s *userRoleService) DeleteUserRole(ctx context.Context, id uint) []helpers.ResponseError {
	err := s.repository.DeleteUserRole(ctx, id)
	if err != nil {
		return []helpers.ResponseError{*err}
	}
	return nil
}

func (s *userRoleService) DeleteUserRolesByUserID(ctx context.Context, userID uint) []helpers.ResponseError {
	err := s.repository.DeleteUserRolesByUserID(ctx, userID)
	if err != nil {
		return []helpers.ResponseError{*err}
	}
	return nil
}
