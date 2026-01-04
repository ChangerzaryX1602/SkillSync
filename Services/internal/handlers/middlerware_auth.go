package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ChangerzaryX1602/SkillSync/pkg/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/golang-jwt/jwt/v4"
	helpers "github.com/zercle/gofiber-helpers"
	"gorm.io/gorm"
)

var ApiLimiter = limiter.New(limiter.Config{
	Max:        750,
	Expiration: 30 * time.Second,
	KeyGenerator: func(c *fiber.Ctx) string {
		return c.Get(fiber.HeaderXForwardedFor)
	},
	LimitReached: func(c *fiber.Ctx) error {
		return helpers.NewError(http.StatusTooManyRequests, helpers.WhereAmI(), http.StatusText(http.StatusTooManyRequests))
	},
})

func ExtractBearerToken(authHeader string) (token string, err error) {
	authHeader = strings.TrimSpace(authHeader)
	authHeaders := strings.Split(authHeader, " ")
	if len(authHeaders) != 2 || authHeaders[0] != "Bearer" {
		err = helpers.NewError(http.StatusUnauthorized, helpers.WhereAmI(), "Authorization: Bearer token")
		return
	}
	token = authHeaders[1]
	return
}

func hasAllPerms(userPerms []string, required []string) bool {
	set := make(map[string]struct{}, len(userPerms))
	for _, p := range userPerms {
		set[p] = struct{}{}
	}
	for _, r := range required {
		if _, ok := set[r]; !ok {
			return false
		}
	}
	return true
}

func (r *RouterResources) ExtractPerms(userId string) ([]string, error) {
	cacheKey := fmt.Sprintf("%s:%s", models.InternalHandlerMiddlewareAuthPermission, userId)
	if r.RedisStorage != nil {
		if cachedBytes, err := r.RedisStorage.Get(cacheKey); err == nil && len(cachedBytes) > 0 {
			var cachedPerms []string
			if err := json.Unmarshal(cachedBytes, &cachedPerms); err == nil {
				return cachedPerms, nil
			}
		}
	}

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return nil, helpers.NewError(http.StatusBadRequest, helpers.WhereAmI(), "user ID malformed")
	}
	permissions, respErr := r.GetUserPermissions(uint(userIdInt))
	if respErr != nil {
		return nil, fmt.Errorf("failed to get user permissions: %s", respErr.Message)
	}

	go func(c string, p []string) {
		if r.RedisStorage != nil && len(p) > 0 {
			if marshaledBytes, err := json.Marshal(p); err == nil {
				_ = r.RedisStorage.Set(c, marshaledBytes, 30*time.Second)
			}
		}
	}(cacheKey, permissions)

	return permissions, nil
}
func (r *RouterResources) GetUserPermissions(userId uint) ([]string, *helpers.ResponseError) {
	type Result struct {
		Pkg  string
		Name string
	}

	var results []Result

	err := r.MainDbConn.
		Table("user_roles").
		Select("permissions.pkg, permissions.name").
		Joins("JOIN role_permissions ON role_permissions.role_id = user_roles.role_id").
		Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
		Where("user_roles.user_id = ?", userId).
		Scan(&results).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound || len(results) == 0 {
			return nil, &helpers.ResponseError{
				Code:    fiber.StatusNotFound,
				Source:  helpers.WhereAmI(),
				Title:   "Not Found",
				Message: "user has no permissions",
			}
		}
		return nil, &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}

	userPerms := make([]string, 0, len(results))
	for _, r := range results {
		userPerms = append(userPerms, r.Pkg+":"+r.Name)
	}
	return userPerms, nil
}
func (r *RouterResources) ReqAuthPerms(requiredPerms ...string) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		// Try Authorization header first, then fall back to query param
		tokenStr, err := ExtractBearerToken(c.Get(fiber.HeaderAuthorization))
		if err != nil {
			tokenStr = c.Query("token")
			if tokenStr == "" {
				return helpers.NewError(http.StatusUnauthorized, helpers.WhereAmI(), "Authorization token required")
			}
		}
		claims := new(jwt.RegisteredClaims)
		jwtToken, err := jwt.ParseWithClaims(tokenStr, claims, r.JwtKeyfunc)
		if err != nil {
			return helpers.NewError(http.StatusUnauthorized, helpers.WhereAmI(), err.Error())
		}
		if jwtToken == nil || !jwtToken.Valid {
			return helpers.NewError(http.StatusUnauthorized, helpers.WhereAmI(), http.StatusText(http.StatusUnauthorized))
		}
		userPerms, err := r.ExtractPerms(claims.Subject)
		if err != nil {
			return helpers.NewError(http.StatusUnauthorized, helpers.WhereAmI(), err.Error())
		}
		ok := hasAllPerms(userPerms, requiredPerms)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{
					{
						Code:    fiber.StatusForbidden,
						Source:  helpers.WhereAmI(),
						Title:   "Forbidden",
						Message: fmt.Sprintf("need permission(s): %v", requiredPerms),
					},
				},
			})
		}
		c.Locals("user_id", claims.Subject)
		return c.Next()
	}
}
