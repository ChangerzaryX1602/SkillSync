package permission

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

type permissionCache struct {
	store *redis.Storage
}

func newPermissionCache(store *redis.Storage) *permissionCache {
	return &permissionCache{store: store}
}

func (c *permissionCache) Get(id uint) (*models.Permission, bool) {
	if c.store == nil {
		return nil, false
	}
	key := fmt.Sprintf("%s:%d", models.PkgPermissionGetPermission, id)

	bytes, err := c.store.Get(key)
	if err != nil || len(bytes) == 0 {
		return nil, false
	}

	var permission models.Permission
	if err := json.Unmarshal(bytes, &permission); err != nil {
		return nil, false
	}
	return &permission, true
}

func (c *permissionCache) Set(permission *models.Permission) {
	if c.store == nil || permission == nil {
		return
	}

	go func(p models.Permission) {
		key := fmt.Sprintf("%s:%d", models.PkgPermissionGetPermission, p.ID)
		if bytes, err := json.Marshal(p); err == nil {
			_ = c.store.Set(key, bytes, utils.RandomJitter(cacheTTL))
		}
	}(*permission)
}

func (c *permissionCache) Invalidate(id uint) {
	if c.store == nil {
		return
	}
	go func() {
		key := fmt.Sprintf("%s:%d", models.PkgPermissionGetPermission, id)
		_ = c.store.Delete(key)
	}()
}

func (c *permissionCache) GetList(pagination models.Pagination, search models.Search) ([]models.Permission, bool) {
	if c.store == nil {
		return nil, false
	}
	key := fmt.Sprintf("%s:%s:%s", models.PkgPermissionGetPermissions, pagination.GetPaginationString(), search.GetSearchString())

	bytes, err := c.store.Get(key)
	if err != nil || len(bytes) == 0 {
		return nil, false
	}

	var permissions []models.Permission
	if err := json.Unmarshal(bytes, &permissions); err != nil {
		return nil, false
	}
	return permissions, true
}

func (c *permissionCache) SetList(pagination models.Pagination, search models.Search, permissions []models.Permission) {
	if c.store == nil {
		return
	}
	go func(ps []models.Permission) {
		key := fmt.Sprintf("%s:%s:%s", models.PkgPermissionGetPermissions, pagination.GetPaginationString(), search.GetSearchString())
		if bytes, err := json.Marshal(ps); err == nil {
			_ = c.store.Set(key, bytes, utils.RandomJitter(cacheTTLList))
		}
	}(permissions)
}
