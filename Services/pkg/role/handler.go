package role

import (
	"github.com/ChangerzaryX1602/SkillSync/internal/handlers"
	"github.com/ChangerzaryX1602/SkillSync/pkg/domain"
	"github.com/ChangerzaryX1602/SkillSync/pkg/models"

	"github.com/gofiber/fiber/v2"
	helpers "github.com/zercle/gofiber-helpers"
)

type roleHandler struct {
	service domain.RoleService
}

func NewRoleHandler(router fiber.Router, resource *handlers.RouterResources, service domain.RoleService) {
	handler := &roleHandler{service: service}
	router.Post("/", resource.ReqAuthPerms(models.PermissionGroupName(models.RoleGroup, models.Create)), handler.CreateRole())
	router.Get("/:id", resource.ReqAuthPerms(models.PermissionGroupName(models.RoleGroup, models.Read)), handler.GetRole())
	router.Get("/", resource.ReqAuthPerms(models.PermissionGroupName(models.RoleGroup, models.List)), handler.GetRoles())
	router.Put("/:id", resource.ReqAuthPerms(models.PermissionGroupName(models.RoleGroup, models.Update)), handler.UpdateRole())
	router.Delete("/:id", resource.ReqAuthPerms(models.PermissionGroupName(models.RoleGroup, models.Delete)), handler.DeleteRole())
}

// CreateRole godoc
// @Summary Create a new role
// @Description Create a new role
// @Tags Role
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param role body models.Role true "Role Data"
// @Router /roles [post]
func (h *roleHandler) CreateRole() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var role models.Role
		if err := c.BodyParser(&role); err != nil {
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
		if errorForm := h.service.CreateRole(c.UserContext(), role); errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusCreated).JSON(helpers.ResponseForm{
			Success: true,
			Data:    role,
		})
	}
}

// GetRole godoc
// @Summary Get a role by ID
// @Description Get a role by ID
// @Tags Role
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Router /roles/{id} [get]
func (h *roleHandler) GetRole() fiber.Handler {
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
		role, errorForm := h.service.GetRole(c.UserContext(), uint(id))
		if errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
			Data:    role,
		})
	}
}

// GetRoles godoc
// @Summary Get all roles
// @Description Get all roles with pagination and search
// @Tags Role
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Param order_by query string false "Order by field"
// @Param sort_by query string false "Sort order (asc/desc)"
// @Param keyword query string false "Search keyword"
// @Param column query string false "Search column"
// @Router /roles [get]
func (h *roleHandler) GetRoles() fiber.Handler {
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
		roles, paginated, searched, errorForm := h.service.GetRoles(c.UserContext(), pagination, search)
		if errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
			Data:    roles,
			Result: map[string]interface{}{
				"pagination": paginated,
				"search":     searched,
			},
		})
	}
}

// UpdateRole godoc
// @Summary Update a role
// @Description Update a role by ID
// @Tags Role
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Param role body models.Role true "Role Data"
// @Router /roles/{id} [put]
func (h *roleHandler) UpdateRole() fiber.Handler {
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
		var role models.Role
		if err := c.BodyParser(&role); err != nil {
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
		if errorForm := h.service.UpdateRole(c.UserContext(), uint(id), role); errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
			Data:    role,
		})
	}
}

// DeleteRole godoc
// @Summary Delete a role
// @Description Delete a role by ID
// @Tags Role
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Router /roles/{id} [delete]
func (h *roleHandler) DeleteRole() fiber.Handler {
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
		if errorForm := h.service.DeleteRole(c.UserContext(), uint(id)); errorForm != nil {
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
