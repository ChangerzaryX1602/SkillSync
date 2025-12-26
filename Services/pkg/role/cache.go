package role

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

type roleCache struct {
	store *redis.Storage
}

func newRoleCache(store *redis.Storage) *roleCache {
	return &roleCache{store: store}
}

func (c *roleCache) Get(id uint) (*models.Role, bool) {
	if c.store == nil {
		return nil, false
	}
	key := fmt.Sprintf("%s:%d", models.PkgRoleGetRole, id)

	bytes, err := c.store.Get(key)
	if err != nil || len(bytes) == 0 {
		return nil, false
	}

	var role models.Role
	if err := json.Unmarshal(bytes, &role); err != nil {
		return nil, false
	}
	return &role, true
}

func (c *roleCache) Set(role *models.Role) {
	if c.store == nil || role == nil {
		return
	}

	go func(r models.Role) {
		key := fmt.Sprintf("%s:%d", models.PkgRoleGetRole, r.ID)
		if bytes, err := json.Marshal(r); err == nil {
			_ = c.store.Set(key, bytes, utils.RandomJitter(cacheTTL))
		}
	}(*role)
}

func (c *roleCache) GetByName(name string) (*models.Role, bool) {
	if c.store == nil {
		return nil, false
	}
	key := fmt.Sprintf("%s:%s", models.PkgRoleGetRoleByName, name)

	bytes, err := c.store.Get(key)
	if err != nil || len(bytes) == 0 {
		return nil, false
	}

	var role models.Role
	if err := json.Unmarshal(bytes, &role); err != nil {
		return nil, false
	}
	return &role, true
}

func (c *roleCache) SetByName(name string, role *models.Role) {
	if c.store == nil || role == nil {
		return
	}

	go func(r models.Role) {
		key := fmt.Sprintf("%s:%s", models.PkgRoleGetRoleByName, name)
		if bytes, err := json.Marshal(r); err == nil {
			_ = c.store.Set(key, bytes, utils.RandomJitter(cacheTTL))
		}
	}(*role)
}

func (c *roleCache) Invalidate(id uint, name string) {
	if c.store == nil {
		return
	}
	go func() {
		key := fmt.Sprintf("%s:%d", models.PkgRoleGetRole, id)
		_ = c.store.Delete(key)
	}()
	go func() {
		key := fmt.Sprintf("%s:%s", models.PkgRoleGetRoleByName, name)
		_ = c.store.Delete(key)
	}()
}

func (c *roleCache) GetList(pagination models.Pagination, search models.Search) ([]models.Role, bool) {
	if c.store == nil {
		return nil, false
	}
	key := fmt.Sprintf("%s:%s:%s", models.PkgRoleGetRoles, pagination.GetPaginationString(), search.GetSearchString())

	bytes, err := c.store.Get(key)
	if err != nil || len(bytes) == 0 {
		return nil, false
	}

	var roles []models.Role
	if err := json.Unmarshal(bytes, &roles); err != nil {
		return nil, false
	}
	return roles, true
}

func (c *roleCache) SetList(pagination models.Pagination, search models.Search, roles []models.Role) {
	if c.store == nil {
		return
	}
	go func(rs []models.Role) {
		key := fmt.Sprintf("%s:%s:%s", models.PkgRoleGetRoles, pagination.GetPaginationString(), search.GetSearchString())
		if bytes, err := json.Marshal(rs); err == nil {
			_ = c.store.Set(key, bytes, utils.RandomJitter(cacheTTLList))
		}
	}(roles)
}
