package user

import (
	"strconv"

	"github.com/ChangerzaryX1602/SkillSync/internal/handlers"
	"github.com/ChangerzaryX1602/SkillSync/pkg/domain"
	"github.com/ChangerzaryX1602/SkillSync/pkg/models"

	"github.com/gofiber/fiber/v2"
	helpers "github.com/zercle/gofiber-helpers"
)

type userHandler struct {
	service domain.UserService
}

func NewUserHandler(router fiber.Router, resource *handlers.RouterResources, service domain.UserService) {
	handler := &userHandler{service: service}
	router.Post("/", resource.ReqAuthPerms(models.PermissionGroupName(models.UserGroup, models.Create)), handler.CreateUser())
	router.Get("/me", resource.ReqAuthPerms(models.PermissionGroupName(models.UserGroup, models.Me)), handler.GetMe())
	router.Get("/:id", resource.ReqAuthPerms(models.PermissionGroupName(models.UserGroup, models.Read)), handler.GetUser())
	router.Get("/", resource.ReqAuthPerms(models.PermissionGroupName(models.UserGroup, models.List)), handler.GetUsers())
	router.Patch("/:id", resource.ReqAuthPerms(models.PermissionGroupName(models.UserGroup, models.Update)), handler.UpdateUser())
	router.Delete("/:id", resource.ReqAuthPerms(models.PermissionGroupName(models.UserGroup, models.Delete)), handler.DeleteUser())
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body models.User true "User Data"
// @Router /users [post]
func (h *userHandler) CreateUser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var user models.User
		if err := c.BodyParser(&user); err != nil {
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
		if errorForm := h.service.CreateUser(c.UserContext(), user); errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusCreated).JSON(helpers.ResponseForm{
			Success: true,
			Data:    user,
		})
	}
}

// GetUser godoc
// @Summary Get a user by ID
// @Description Get a user by ID
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Router /users/{id} [get]
func (h *userHandler) GetUser() fiber.Handler {
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
		user, errorForm := h.service.GetUser(c.UserContext(), uint(id))
		if errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
			Data:    user,
		})
	}
}

// GetUsers godoc
// @Summary Get all users
// @Description Get all users with pagination and search
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Param order_by query string false "Order by field"
// @Param sort_by query string false "Sort order (asc/desc)"
// @Param keyword query string false "Search keyword"
// @Param column query string false "Search column"
// @Router /users [get]
func (h *userHandler) GetUsers() fiber.Handler {
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
		users, paginated, searched, errorForm := h.service.GetUsers(c.UserContext(), pagination, search)
		if errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
			Data:    users,
			Result: map[string]interface{}{
				"pagination": paginated,
				"search":     searched,
			},
		})
	}
}

// UpdateUser godoc
// @Summary Update a user by ID
// @Description Update a user by ID
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param user body models.User true "User Data"
// @Router /users/{id} [patch]
func (h *userHandler) UpdateUser() fiber.Handler {
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
		var user models.User
		if err := c.BodyParser(&user); err != nil {
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
		if errorForm := h.service.UpdateUser(c.UserContext(), uint(id), user); errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
			Data:    nil,
		})
	}
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user by ID
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Router /users/{id} [delete]
func (h *userHandler) DeleteUser() fiber.Handler {
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
		if errorForm := h.service.DeleteUser(c.UserContext(), uint(id)); errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
			Data:    nil,
		})
	}
}

// GetMe godoc
// @Summary Get current user
// @Description Get current user
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Router /users/me [get]
func (h *userHandler) GetMe() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := c.Locals("user_id").(string)
		userIdInt, err := strconv.Atoi(userId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{
					{
						Code:    fiber.StatusBadRequest,
						Source:  helpers.WhereAmI(),
						Title:   "Invalid User ID",
						Message: err.Error(),
					},
				},
			})
		}
		user, errorForm := h.service.GetUser(c.UserContext(), uint(userIdInt))
		if errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
			Data:    user,
		})
	}
}
