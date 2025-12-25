package permission

import (
	"github.com/ChangerzaryX1602/SkillSync/internal/handlers"
	"github.com/ChangerzaryX1602/SkillSync/pkg/domain"
	"github.com/ChangerzaryX1602/SkillSync/pkg/models"

	"github.com/gofiber/fiber/v2"
	helpers "github.com/zercle/gofiber-helpers"
)

type permissionHandler struct {
	service  domain.PermissionService
	resource *handlers.RouterResources
}

func NewPermissionHandler(router fiber.Router, resource *handlers.RouterResources, service domain.PermissionService) {
	handler := &permissionHandler{service: service, resource: resource}
	router.Post("/", resource.ReqAuthPerms(models.PermissionGroupName(models.PermissionGroup, models.Create)), handler.CreatePermission())
	router.Get("/me", resource.ReqAuthPerms(models.PermissionGroupName(models.PermissionGroup, models.Me)), handler.GetMyPermissions())
	router.Get("/:id", resource.ReqAuthPerms(models.PermissionGroupName(models.PermissionGroup, models.Read)), handler.GetPermission())
	router.Get("/", resource.ReqAuthPerms(models.PermissionGroupName(models.PermissionGroup, models.List)), handler.GetPermissions())
	router.Put("/:id", resource.ReqAuthPerms(models.PermissionGroupName(models.PermissionGroup, models.Update)), handler.UpdatePermission())
	router.Delete("/:id", resource.ReqAuthPerms(models.PermissionGroupName(models.PermissionGroup, models.Delete)), handler.DeletePermission())
}

// CreatePermission godoc
// @Summary Create a new permission
// @Description Create a new permission
// @Tags Permission
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param permission body models.Permission true "Permission Data"
// @Router /permissions [post]
func (h *permissionHandler) CreatePermission() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var permission models.Permission
		if err := c.BodyParser(&permission); err != nil {
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
		if errorForm := h.service.CreatePermission(permission); errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusCreated).JSON(helpers.ResponseForm{
			Success: true,
			Data:    permission,
		})
	}
}

// GetPermission godoc
// @Summary Get a permission by ID
// @Description Get a permission by ID
// @Tags Permission
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Permission ID"
// @Router /permissions/{id} [get]
func (h *permissionHandler) GetPermission() fiber.Handler {
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
		permission, errorForm := h.service.GetPermission(uint(id))
		if errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
			Data:    permission,
		})
	}
}

// GetPermissions godoc
// @Summary Get all permissions
// @Description Get all permissions with pagination and search
// @Tags Permission
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Param order_by query string false "Order by field"
// @Param sort_by query string false "Sort order (asc/desc)"
// @Param keyword query string false "Search keyword"
// @Param column query string false "Search column"
// @Router /permissions [get]
func (h *permissionHandler) GetPermissions() fiber.Handler {
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
		permissions, paginated, searched, errorForm := h.service.GetPermissions(pagination, search)
		if errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
			Data:    permissions,
			Result: map[string]interface{}{
				"pagination": paginated,
				"search":     searched,
			},
		})
	}
}

// UpdatePermission godoc
// @Summary Update a permission
// @Description Update a permission by ID
// @Tags Permission
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Permission ID"
// @Param permission body models.Permission true "Permission Data"
// @Router /permissions/{id} [put]
func (h *permissionHandler) UpdatePermission() fiber.Handler {
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
		var permission models.Permission
		if err := c.BodyParser(&permission); err != nil {
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
		if errorForm := h.service.UpdatePermission(uint(id), permission); errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
			Data:    permission,
		})
	}
}

// DeletePermission godoc
// @Summary Delete a permission
// @Description Delete a permission by ID
// @Tags Permission
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Permission ID"
// @Router /permissions/{id} [delete]
func (h *permissionHandler) DeletePermission() fiber.Handler {
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
		if errorForm := h.service.DeletePermission(uint(id)); errorForm != nil {
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

// GetMyPermissions godoc
// @Summary Get my permissions
// @Description Get permissions for the authenticated user
// @Tags Permission
// @Accept json
// @Produce json
// @Security BearerAuth
// @Router /permissions/me [get]
func (h *permissionHandler) GetMyPermissions() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := c.Locals("user_id").(string)
		permissions, err := h.resource.ExtractPerms(userId)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{
					{
						Code:    fiber.StatusInternalServerError,
						Source:  helpers.WhereAmI(),
						Title:   "Internal Server Error",
						Message: err.Error(),
					},
				},
			})
		}
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
			Data:    permissions,
		})
	}
}
