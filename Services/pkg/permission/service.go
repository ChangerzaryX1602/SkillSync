package permission

import (
	"context"

	"github.com/ChangerzaryX1602/SkillSync/pkg/domain"
	"github.com/ChangerzaryX1602/SkillSync/pkg/models"

	helpers "github.com/zercle/gofiber-helpers"
)

type permissionService struct {
	repository domain.PermissionRepository
}

func NewPermissionService(repository domain.PermissionRepository) domain.PermissionService {
	return &permissionService{repository: repository}
}

func (s *permissionService) CreatePermission(ctx context.Context, permission models.Permission) []helpers.ResponseError {
	err := s.repository.CreatePermission(ctx, permission)
	if err != nil {
		return []helpers.ResponseError{*err}
	}
	return nil
}

func (s *permissionService) GetPermission(ctx context.Context, id uint) (*models.Permission, []helpers.ResponseError) {
	permission, err := s.repository.GetPermission(ctx, id)
	if err != nil {
		return nil, []helpers.ResponseError{*err}
	}
	return permission, nil
}

func (s *permissionService) GetPermissions(ctx context.Context, pagination models.Pagination, search models.Search) ([]models.Permission, *models.Pagination, *models.Search, []helpers.ResponseError) {
	permissions, paginated, searched, err := s.repository.GetPermissions(ctx, pagination, search)
	if err != nil {
		return nil, nil, nil, []helpers.ResponseError{*err}
	}
	return permissions, paginated, searched, nil
}

func (s *permissionService) UpdatePermission(ctx context.Context, id uint, permission models.Permission) []helpers.ResponseError {
	err := s.repository.UpdatePermission(ctx, id, permission)
	if err != nil {
		return []helpers.ResponseError{*err}
	}
	return nil
}

func (s *permissionService) DeletePermission(ctx context.Context, id uint) []helpers.ResponseError {
	err := s.repository.DeletePermission(ctx, id)
	if err != nil {
		return []helpers.ResponseError{*err}
	}
	return nil
}
