package role_permission

import (
	"github.com/ChangerzaryX1602/SkillSync/pkg/domain"
	"github.com/ChangerzaryX1602/SkillSync/pkg/models"

	helpers "github.com/zercle/gofiber-helpers"
)

type rolePermissionService struct {
	repository domain.RolePermissionRepository
}

func NewRolePermissionService(repository domain.RolePermissionRepository) domain.RolePermissionService {
	return &rolePermissionService{repository: repository}
}

func (s *rolePermissionService) CreateRolePermission(rolePermission models.RolePermission) []helpers.ResponseError {
	err := s.repository.CreateRolePermission(rolePermission)
	if err != nil {
		return []helpers.ResponseError{*err}
	}
	return nil
}

func (s *rolePermissionService) GetRolePermission(id uint) (*models.RolePermission, []helpers.ResponseError) {
	rolePermission, err := s.repository.GetRolePermission(id)
	if err != nil {
		return nil, []helpers.ResponseError{*err}
	}
	return rolePermission, nil
}

func (s *rolePermissionService) GetRolePermissions(pagination models.Pagination, search models.Search) ([]models.RolePermission, *models.Pagination, *models.Search, []helpers.ResponseError) {
	rolePermissions, paginated, searched, err := s.repository.GetRolePermissions(pagination, search)
	if err != nil {
		return nil, nil, nil, []helpers.ResponseError{*err}
	}
	return rolePermissions, paginated, searched, nil
}

func (s *rolePermissionService) GetRolePermissionsByRoleID(roleID uint) ([]models.RolePermission, []helpers.ResponseError) {
	rolePermissions, err := s.repository.GetRolePermissionsByRoleID(roleID)
	if err != nil {
		return nil, []helpers.ResponseError{*err}
	}
	return rolePermissions, nil
}

func (s *rolePermissionService) UpdateRolePermission(id uint, rolePermission models.RolePermission) []helpers.ResponseError {
	err := s.repository.UpdateRolePermission(id, rolePermission)
	if err != nil {
		return []helpers.ResponseError{*err}
	}
	return nil
}

func (s *rolePermissionService) DeleteRolePermission(id uint) []helpers.ResponseError {
	err := s.repository.DeleteRolePermission(id)
	if err != nil {
		return []helpers.ResponseError{*err}
	}
	return nil
}

func (s *rolePermissionService) DeleteRolePermissionsByRoleID(roleID uint) []helpers.ResponseError {
	err := s.repository.DeleteRolePermissionsByRoleID(roleID)
	if err != nil {
		return []helpers.ResponseError{*err}
	}
	return nil
}
