package handlers

import (
	"github.com/gofiber/storage/redis"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

// RouterResources DB handler
type RouterResources struct {
	JwtKeyfunc   jwt.Keyfunc
	MainDbConn   *gorm.DB
	RedisStorage *redis.Storage
}

// NewRouterResources returns a new DBHandler
func NewRouterResources(jwtKeyfunc jwt.Keyfunc, mainDbConn *gorm.DB, redisStorage *redis.Storage) *RouterResources {
	return &RouterResources{
		JwtKeyfunc:   jwtKeyfunc,
		MainDbConn:   mainDbConn,
		RedisStorage: redisStorage,
	}
}
