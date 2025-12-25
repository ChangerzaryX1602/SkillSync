package handlers

import (
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

// RouterResources DB handler
type RouterResources struct {
	JwtKeyfunc jwt.Keyfunc
	MainDbConn *gorm.DB
}

// NewRouterResources returns a new DBHandler
func NewRouterResources(jwtKeyfunc jwt.Keyfunc, mainDbConn *gorm.DB) *RouterResources {
	return &RouterResources{
		JwtKeyfunc: jwtKeyfunc,
		MainDbConn: mainDbConn,
	}
}
