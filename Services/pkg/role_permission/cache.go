package role_permission

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ChangerzaryX1602/SkillSync/pkg/models"
	"github.com/ChangerzaryX1602/SkillSync/pkg/utils"
	"github.com/gofiber/storage/redis"
)

const (
	cacheTTL     = 15 * time.Minute
	cacheTTLList = 1 * time.Minute
)

type rolePermissionCache struct {
	store *redis.Storage
}

func newRolePermissionCache(store *redis.Storage) *rolePermissionCache {
	return &rolePermissionCache{store: store}
}

func (c *rolePermissionCache) GetByRoleID(roleID uint) ([]models.RolePermission, bool) {
	if c.store == nil {
		return nil, false
	}
	key := fmt.Sprintf("%s:%d", models.PkgRolePermissionGetByRoleID, roleID)

	bytes, err := c.store.Get(key)
	if err != nil || len(bytes) == 0 {
		return nil, false
	}

	var rolePermissions []models.RolePermission
	if err := json.Unmarshal(bytes, &rolePermissions); err != nil {
		return nil, false
	}
	return rolePermissions, true
}

func (c *rolePermissionCache) SetByRoleID(roleID uint, rolePermissions []models.RolePermission) {
	if c.store == nil {
		return
	}
	go func(rps []models.RolePermission) {
		key := fmt.Sprintf("%s:%d", models.PkgRolePermissionGetByRoleID, roleID)
		if bytes, err := json.Marshal(rps); err == nil {
			_ = c.store.Set(key, bytes, utils.RandomJitter(cacheTTLList))
		}
	}(rolePermissions)
}

func (c *rolePermissionCache) InvalidateByRoleID(roleID uint) {
	if c.store == nil {
		return
	}
	go func() {
		key := fmt.Sprintf("%s:%d", models.PkgRolePermissionGetByRoleID, roleID)
		_ = c.store.Delete(key)
	}()
}
