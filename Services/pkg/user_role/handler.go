package user_role

import (
	"github.com/ChangerzaryX1602/SkillSync/internal/handlers"
	"github.com/ChangerzaryX1602/SkillSync/pkg/domain"
	"github.com/ChangerzaryX1602/SkillSync/pkg/models"

	"github.com/gofiber/fiber/v2"
	helpers "github.com/zercle/gofiber-helpers"
)

type userRoleHandler struct {
	service domain.UserRoleService
}

func NewUserRoleHandler(router fiber.Router, resource *handlers.RouterResources, service domain.UserRoleService) {
	handler := &userRoleHandler{service: service}
	router.Post("/", resource.ReqAuthPerms(models.PermissionGroupName(models.UserRoleGroup, models.Create)), handler.CreateUserRole())
	router.Get("/:id", resource.ReqAuthPerms(models.PermissionGroupName(models.UserRoleGroup, models.Read)), handler.GetUserRole())
	router.Get("/", resource.ReqAuthPerms(models.PermissionGroupName(models.UserRoleGroup, models.List)), handler.GetUserRoles())
	router.Get("/user/:userId", resource.ReqAuthPerms(models.PermissionGroupName(models.UserRoleGroup, models.List)), handler.GetUserRolesByUserID())
	router.Put("/:id", resource.ReqAuthPerms(models.PermissionGroupName(models.UserRoleGroup, models.Update)), handler.UpdateUserRole())
	router.Delete("/:id", resource.ReqAuthPerms(models.PermissionGroupName(models.UserRoleGroup, models.Delete)), handler.DeleteUserRole())
	router.Delete("/user/:userId", resource.ReqAuthPerms(models.PermissionGroupName(models.UserRoleGroup, models.Delete)), handler.DeleteUserRolesByUserID())
}

// CreateUserRole godoc
// @Summary Create a new user role
// @Description Create a new user role
// @Tags UserRole
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userRole body models.UserRole true "UserRole Data"
// @Router /user_roles [post]
func (h *userRoleHandler) CreateUserRole() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var userRole models.UserRole
		if err := c.BodyParser(&userRole); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{
					{
						Code:    fiber.StatusBadRequest,
						Source:  helpers.WhereAmI(),
						Title:   "Invalid Request",
						Message: err.Error(),
					},
				},
			})
		}
		if errorForm := h.service.CreateUserRole(userRole); errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusCreated).JSON(helpers.ResponseForm{
			Success: true,
			Data:    userRole,
		})
	}
}

// GetUserRole godoc
// @Summary Get a user role by ID
// @Description Get a user role by ID
// @Tags UserRole
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "UserRole ID"
// @Router /user_roles/{id} [get]
func (h *userRoleHandler) GetUserRole() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{
					{
						Code:    fiber.StatusBadRequest,
						Source:  helpers.WhereAmI(),
						Title:   "Invalid Request",
						Message: err.Error(),
					},
				},
			})
		}
		userRole, errorForm := h.service.GetUserRole(uint(id))
		if errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
			Data:    userRole,
		})
	}
}

// GetUserRoles godoc
// @Summary Get all user roles
// @Description Get all user roles with pagination and search
// @Tags UserRole
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Param order_by query string false "Order by field"
// @Param sort_by query string false "Sort order (asc/desc)"
// @Param keyword query string false "Search keyword"
// @Param column query string false "Search column"
// @Router /user_roles [get]
func (h *userRoleHandler) GetUserRoles() fiber.Handler {
	return func(c *fiber.Ctx) error {
		pagination := models.Pagination{}
		search := models.Search{}
		if err := c.QueryParser(&pagination); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{
					{
						Code:    fiber.StatusBadRequest,
						Source:  helpers.WhereAmI(),
						Title:   "Invalid Request",
						Message: err.Error(),
					},
				},
			})
		}
		if err := c.QueryParser(&search); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{
					{
						Code:    fiber.StatusBadRequest,
						Source:  helpers.WhereAmI(),
						Title:   "Invalid Request",
						Message: err.Error(),
					},
				},
			})
		}
		userRoles, paginated, searched, errorForm := h.service.GetUserRoles(pagination, search)
		if errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
			Data:    userRoles,
			Result: map[string]interface{}{
				"pagination": paginated,
				"search":     searched,
			},
		})
	}
}

// GetUserRolesByUserID godoc
// @Summary Get user roles by User ID
// @Description Get user roles by User ID
// @Tags UserRole
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userId path int true "User ID"
// @Router /user_roles/user/{userId} [get]
func (h *userRoleHandler) GetUserRolesByUserID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId, err := c.ParamsInt("userId")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{
					{
						Code:    fiber.StatusBadRequest,
						Source:  helpers.WhereAmI(),
						Title:   "Invalid Request",
						Message: err.Error(),
					},
				},
			})
		}
		userRoles, errorForm := h.service.GetUserRolesByUserID(uint(userId))
		if errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
			Data:    userRoles,
		})
	}
}

// UpdateUserRole godoc
// @Summary Update a user role
// @Description Update a user role by ID
// @Tags UserRole
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "UserRole ID"
// @Param userRole body models.UserRole true "UserRole Data"
// @Router /user_roles/{id} [put]
func (h *userRoleHandler) UpdateUserRole() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{
					{
						Code:    fiber.StatusBadRequest,
						Source:  helpers.WhereAmI(),
						Title:   "Invalid Request",
						Message: err.Error(),
					},
				},
			})
		}
		var userRole models.UserRole
		if err := c.BodyParser(&userRole); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{
					{
						Code:    fiber.StatusBadRequest,
						Source:  helpers.WhereAmI(),
						Title:   "Invalid Request",
						Message: err.Error(),
					},
				},
			})
		}
		if errorForm := h.service.UpdateUserRole(uint(id), userRole); errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
			Data:    userRole,
		})
	}
}

// DeleteUserRole godoc
// @Summary Delete a user role
// @Description Delete a user role by ID
// @Tags UserRole
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "UserRole ID"
// @Router /user_roles/{id} [delete]
func (h *userRoleHandler) DeleteUserRole() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{
					{
						Code:    fiber.StatusBadRequest,
						Source:  helpers.WhereAmI(),
						Title:   "Invalid Request",
						Message: err.Error(),
					},
				},
			})
		}
		if errorForm := h.service.DeleteUserRole(uint(id)); errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
		})
	}
}

// DeleteUserRolesByUserID godoc
// @Summary Delete user roles by User ID
// @Description Delete user roles by User ID
// @Tags UserRole
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userId path int true "User ID"
// @Router /user_roles/user/{userId} [delete]
func (h *userRoleHandler) DeleteUserRolesByUserID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId, err := c.ParamsInt("userId")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{
					{
						Code:    fiber.StatusBadRequest,
						Source:  helpers.WhereAmI(),
						Title:   "Invalid Request",
						Message: err.Error(),
					},
				},
			})
		}
		if errorForm := h.service.DeleteUserRolesByUserID(uint(userId)); errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
		})
	}
}
