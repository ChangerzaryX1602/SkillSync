package user_role

import (
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

func (s *userRoleService) CreateUserRole(userRole models.UserRole) []helpers.ResponseError {
	err := s.repository.CreateUserRole(userRole)
	if err != nil {
		return []helpers.ResponseError{*err}
	}
	return nil
}

func (s *userRoleService) GetUserRole(id uint) (*models.UserRole, []helpers.ResponseError) {
	userRole, err := s.repository.GetUserRole(id)
	if err != nil {
		return nil, []helpers.ResponseError{*err}
	}
	return userRole, nil
}

func (s *userRoleService) GetUserRoles(pagination models.Pagination, search models.Search) ([]models.UserRole, *models.Pagination, *models.Search, []helpers.ResponseError) {
	userRoles, paginated, searched, err := s.repository.GetUserRoles(pagination, search)
	if err != nil {
		return nil, nil, nil, []helpers.ResponseError{*err}
	}
	return userRoles, paginated, searched, nil
}

func (s *userRoleService) GetUserRolesByUserID(userID uint) ([]models.UserRole, []helpers.ResponseError) {
	userRoles, err := s.repository.GetUserRolesByUserID(userID)
	if err != nil {
		return nil, []helpers.ResponseError{*err}
	}
	return userRoles, nil
}

func (s *userRoleService) UpdateUserRole(id uint, userRole models.UserRole) []helpers.ResponseError {
	err := s.repository.UpdateUserRole(id, userRole)
	if err != nil {
		return []helpers.ResponseError{*err}
	}
	return nil
}

func (s *userRoleService) DeleteUserRole(id uint) []helpers.ResponseError {
	err := s.repository.DeleteUserRole(id)
	if err != nil {
		return []helpers.ResponseError{*err}
	}
	return nil
}

func (s *userRoleService) DeleteUserRolesByUserID(userID uint) []helpers.ResponseError {
	err := s.repository.DeleteUserRolesByUserID(userID)
	if err != nil {
		return []helpers.ResponseError{*err}
	}
	return nil
}
