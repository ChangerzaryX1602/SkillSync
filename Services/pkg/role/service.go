package role

import (
	"github.com/ChangerzaryX1602/SkillSync/pkg/domain"
	"github.com/ChangerzaryX1602/SkillSync/pkg/models"

	helpers "github.com/zercle/gofiber-helpers"
)

type roleService struct {
	repository domain.RoleRepository
}

func NewRoleService(repository domain.RoleRepository) domain.RoleService {
	return &roleService{repository: repository}
}

func (s *roleService) CreateRole(role models.Role) []helpers.ResponseError {
	err := s.repository.CreateRole(role)
	if err != nil {
		return []helpers.ResponseError{*err}
	}
	return nil
}

func (s *roleService) GetRole(id uint) (*models.Role, []helpers.ResponseError) {
	role, err := s.repository.GetRole(id)
	if err != nil {
		return nil, []helpers.ResponseError{*err}
	}
	return role, nil
}

func (s *roleService) GetRoles(pagination models.Pagination, search models.Search) ([]models.Role, *models.Pagination, *models.Search, []helpers.ResponseError) {
	roles, paginated, searched, err := s.repository.GetRoles(pagination, search)
	if err != nil {
		return nil, nil, nil, []helpers.ResponseError{*err}
	}
	return roles, paginated, searched, nil
}

func (s *roleService) UpdateRole(id uint, role models.Role) []helpers.ResponseError {
	err := s.repository.UpdateRole(id, role)
	if err != nil {
		return []helpers.ResponseError{*err}
	}
	return nil
}

func (s *roleService) DeleteRole(id uint) []helpers.ResponseError {
	err := s.repository.DeleteRole(id)
	if err != nil {
		return []helpers.ResponseError{*err}
	}
	return nil
}
