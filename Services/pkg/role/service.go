package role

import (
	"context"

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

func (s *roleService) CreateRole(ctx context.Context, role models.Role) []helpers.ResponseError {
	err := s.repository.CreateRole(ctx, role)
	if err != nil {
		return []helpers.ResponseError{*err}
	}
	return nil
}

func (s *roleService) GetRole(ctx context.Context, id uint) (*models.Role, []helpers.ResponseError) {
	role, err := s.repository.GetRole(ctx, id)
	if err != nil {
		return nil, []helpers.ResponseError{*err}
	}
	return role, nil
}

func (s *roleService) GetRoles(ctx context.Context, pagination models.Pagination, search models.Search) ([]models.Role, *models.Pagination, *models.Search, []helpers.ResponseError) {
	roles, paginated, searched, err := s.repository.GetRoles(ctx, pagination, search)
	if err != nil {
		return nil, nil, nil, []helpers.ResponseError{*err}
	}
	return roles, paginated, searched, nil
}

func (s *roleService) UpdateRole(ctx context.Context, id uint, role models.Role) []helpers.ResponseError {
	err := s.repository.UpdateRole(ctx, id, role)
	if err != nil {
		return []helpers.ResponseError{*err}
	}
	return nil
}

func (s *roleService) DeleteRole(ctx context.Context, id uint) []helpers.ResponseError {
	err := s.repository.DeleteRole(ctx, id)
	if err != nil {
		return []helpers.ResponseError{*err}
	}
	return nil
}
