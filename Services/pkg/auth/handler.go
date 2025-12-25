package auth

import (
	"github.com/ChangerzaryX1602/SkillSync/internal/handlers"
	"github.com/ChangerzaryX1602/SkillSync/pkg/domain"
	"github.com/ChangerzaryX1602/SkillSync/pkg/models"

	"github.com/gofiber/fiber/v2"
	helpers "github.com/zercle/gofiber-helpers"
)

type authHandler struct {
	service         domain.AuthService
	routerResources *handlers.RouterResources
}

func NewAuthHandler(router fiber.Router, routerResources *handlers.RouterResources, service domain.AuthService) {
	handler := &authHandler{service: service, routerResources: routerResources}
	router.Post("/login", handler.Login())
	router.Post("/register", handler.Register())
	router.Post("/refresh", handler.RefreshToken())
}

// Login godoc
// @Summary Login to the application
// @Description Login with username and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body models.User true "User Login Credentials"
// @Router /auth/login [post]
func (h *authHandler) Login() fiber.Handler {
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
		host := c.Hostname()
		token, refreshToken, errorForm := h.service.Login(user, host)
		if errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success":       true,
			"token":         token,
			"refresh_token": refreshToken,
		})
	}
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Refresh access token using refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param refresh_token body map[string]string true "Refresh Token"
// @Router /auth/refresh [post]
func (h *authHandler) RefreshToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body map[string]string
		if err := c.BodyParser(&body); err != nil {
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

		refreshTokenGiven, ok := body["refresh_token"]
		if !ok {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{
					{
						Code:    fiber.StatusBadRequest,
						Source:  helpers.WhereAmI(),
						Title:   "Invalid Request",
						Message: "refresh_token is required",
					},
				},
			})
		}

		token, refreshToken, errorForm := h.service.RefreshToken(refreshTokenGiven)
		if errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success":       true,
			"token":         token,
			"refresh_token": refreshToken,
		})
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body models.User true "User Registration Data"
// @Router /auth/register [post]
func (h *authHandler) Register() fiber.Handler {
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
		errorForm := h.service.Register(user)
		if errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusCreated).JSON(helpers.ResponseForm{
			Success: true,
		})
	}
}
