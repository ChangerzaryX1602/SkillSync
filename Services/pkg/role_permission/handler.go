package role_permission

import (
	"github.com/ChangerzaryX1602/SkillSync/internal/handlers"
	"github.com/ChangerzaryX1602/SkillSync/pkg/domain"
	"github.com/ChangerzaryX1602/SkillSync/pkg/models"

	"github.com/gofiber/fiber/v2"
	helpers "github.com/zercle/gofiber-helpers"
)

type rolePermissionHandler struct {
	service domain.RolePermissionService
}

func NewRolePermissionHandler(router fiber.Router, resource *handlers.RouterResources, service domain.RolePermissionService) {
	handler := &rolePermissionHandler{service: service}
	router.Post("/", resource.ReqAuthPerms(models.PermissionGroupName(models.RolePermissionGroup, models.Create)), handler.CreateRolePermission())
	router.Get("/:id", resource.ReqAuthPerms(models.PermissionGroupName(models.RolePermissionGroup, models.Read)), handler.GetRolePermission())
	router.Get("/", resource.ReqAuthPerms(models.PermissionGroupName(models.RolePermissionGroup, models.List)), handler.GetRolePermissions())
	router.Get("/role/:roleId", resource.ReqAuthPerms(models.PermissionGroupName(models.RolePermissionGroup, models.List)), handler.GetRolePermissionsByRoleID())
	router.Put("/:id", resource.ReqAuthPerms(models.PermissionGroupName(models.RolePermissionGroup, models.Update)), handler.UpdateRolePermission())
	router.Delete("/:id", resource.ReqAuthPerms(models.PermissionGroupName(models.RolePermissionGroup, models.Delete)), handler.DeleteRolePermission())
	router.Delete("/role/:roleId", resource.ReqAuthPerms(models.PermissionGroupName(models.RolePermissionGroup, models.Delete)), handler.DeleteRolePermissionsByRoleID())
}

// CreateRolePermission godoc
// @Summary Create a new role permission
// @Description Create a new role permission
// @Tags RolePermission
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param rolePermission body models.RolePermission true "RolePermission Data"
// @Router /role_permissions [post]
func (h *rolePermissionHandler) CreateRolePermission() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var rolePermission models.RolePermission
		if err := c.BodyParser(&rolePermission); err != nil {
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
		if errorForm := h.service.CreateRolePermission(c.UserContext(), rolePermission); errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusCreated).JSON(helpers.ResponseForm{
			Success: true,
			Data:    rolePermission,
		})
	}
}

// GetRolePermission godoc
// @Summary Get a role permission by ID
// @Description Get a role permission by ID
// @Tags RolePermission
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "RolePermission ID"
// @Router /role_permissions/{id} [get]
func (h *rolePermissionHandler) GetRolePermission() fiber.Handler {
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
		rolePermission, errorForm := h.service.GetRolePermission(c.UserContext(), uint(id))
		if errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
			Data:    rolePermission,
		})
	}
}

// GetRolePermissions godoc
// @Summary Get all role permissions
// @Description Get all role permissions with pagination and search
// @Tags RolePermission
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Param order_by query string false "Order by field"
// @Param sort_by query string false "Sort order (asc/desc)"
// @Param keyword query string false "Search keyword"
// @Param column query string false "Search column"
// @Router /role_permissions [get]
func (h *rolePermissionHandler) GetRolePermissions() fiber.Handler {
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
		rolePermissions, paginated, searched, errorForm := h.service.GetRolePermissions(c.UserContext(), pagination, search)
		if errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
			Data:    rolePermissions,
			Result: map[string]interface{}{
				"pagination": paginated,
				"search":     searched,
			},
		})
	}
}

// GetRolePermissionsByRoleID godoc
// @Summary Get role permissions by Role ID
// @Description Get role permissions by Role ID
// @Tags RolePermission
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param roleId path int true "Role ID"
// @Router /role_permissions/role/{roleId} [get]
func (h *rolePermissionHandler) GetRolePermissionsByRoleID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleId, err := c.ParamsInt("roleId")
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
		rolePermissions, errorForm := h.service.GetRolePermissionsByRoleID(c.UserContext(), uint(roleId))
		if errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
			Data:    rolePermissions,
		})
	}
}

// UpdateRolePermission godoc
// @Summary Update a role permission
// @Description Update a role permission by ID
// @Tags RolePermission
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "RolePermission ID"
// @Param rolePermission body models.RolePermission true "RolePermission Data"
// @Router /role_permissions/{id} [put]
func (h *rolePermissionHandler) UpdateRolePermission() fiber.Handler {
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
		var rolePermission models.RolePermission
		if err := c.BodyParser(&rolePermission); err != nil {
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
		if errorForm := h.service.UpdateRolePermission(c.UserContext(), uint(id), rolePermission); errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
			Data:    rolePermission,
		})
	}
}

// DeleteRolePermission godoc
// @Summary Delete a role permission
// @Description Delete a role permission by ID
// @Tags RolePermission
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "RolePermission ID"
// @Router /role_permissions/{id} [delete]
func (h *rolePermissionHandler) DeleteRolePermission() fiber.Handler {
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
		if errorForm := h.service.DeleteRolePermission(c.UserContext(), uint(id)); errorForm != nil {
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

// DeleteRolePermissionsByRoleID godoc
// @Summary Delete role permissions by Role ID
// @Description Delete role permissions by Role ID
// @Tags RolePermission
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param roleId path int true "Role ID"
// @Router /role_permissions/role/{roleId} [delete]
func (h *rolePermissionHandler) DeleteRolePermissionsByRoleID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleId, err := c.ParamsInt("roleId")
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
		if errorForm := h.service.DeleteRolePermissionsByRoleID(c.UserContext(), uint(roleId)); errorForm != nil {
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
