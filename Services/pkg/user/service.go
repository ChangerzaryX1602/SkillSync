package user

import (
	"context"

	"github.com/ChangerzaryX1602/SkillSync/pkg/domain"
	"github.com/ChangerzaryX1602/SkillSync/pkg/models"
	txcontext "github.com/ChangerzaryX1602/SkillSync/pkg/tx_context"
	"github.com/ChangerzaryX1602/SkillSync/pkg/utils"

	helpers "github.com/zercle/gofiber-helpers"
)

type userService struct {
	repository         domain.UserRepository
	roleRepository     domain.RoleRepository
	userRoleRepository domain.UserRoleRepository
	txManager          *txcontext.TxManager
}

func NewUserService(repository domain.UserRepository, roleRepository domain.RoleRepository, userRoleRepository domain.UserRoleRepository, txManager *txcontext.TxManager) domain.UserService {
	return &userService{
		repository:         repository,
		roleRepository:     roleRepository,
		userRoleRepository: userRoleRepository,
		txManager:          txManager,
	}
}
func (s *userService) CreateUser(ctx context.Context, user models.User) []helpers.ResponseError {
	if !utils.IsValidEmail(user.Email) || user.PasswordTemp == "" {
		return []helpers.ResponseError{
			{
				Code:    400,
				Source:  helpers.WhereAmI(),
				Title:   "Bad Request",
				Message: "Invalid email format or password temp is empty",
			},
		}
	}
	var errCreate []helpers.ResponseError
	err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) []helpers.ResponseError {
		if err := s.repository.CreateUser(txCtx, user); err != nil {
			errCreate = append(errCreate, *err)
			return errCreate
		}
		role, err := s.roleRepository.GetRoleByName(txCtx, string(utils.User))
		if err != nil {
			errCreate = append(errCreate, *err)
			return errCreate
		}

		userRole := models.UserRole{
			UserID: user.ID,
			RoleID: role.ID,
		}
		if err := s.userRoleRepository.CreateUserRole(txCtx, userRole); err != nil {
			errCreate = append(errCreate, *err)
			return errCreate
		}
		return nil
	})
	if err != nil {
		return errCreate
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
	if user.Email != "" {
		if !utils.IsValidEmail(user.Email) {
			return []helpers.ResponseError{
				{
					Code:    400,
					Source:  helpers.WhereAmI(),
					Title:   "Bad Request",
					Message: "Invalid email format",
				},
			}
		}
	}
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
	if !utils.IsValidEmail(email) {
		return nil, []helpers.ResponseError{
			{
				Code:    400,
				Source:  helpers.WhereAmI(),
				Title:   "Bad Request",
				Message: "Invalid email format",
			},
		}
	}
	user, err := s.repository.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, []helpers.ResponseError{*err}
	}
	return user, nil
}
