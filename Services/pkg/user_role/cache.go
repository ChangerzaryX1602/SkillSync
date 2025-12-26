package user_role

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

type userRoleCache struct {
	store *redis.Storage
}

func newUserRoleCache(store *redis.Storage) *userRoleCache {
	return &userRoleCache{store: store}
}

func (c *userRoleCache) GetByUserID(userID uint) ([]models.UserRole, bool) {
	if c.store == nil {
		return nil, false
	}
	key := fmt.Sprintf("%s:%d", models.PkgUserRoleGetByUserID, userID)

	bytes, err := c.store.Get(key)
	if err != nil || len(bytes) == 0 {
		return nil, false
	}

	var userRoles []models.UserRole
	if err := json.Unmarshal(bytes, &userRoles); err != nil {
		return nil, false
	}
	return userRoles, true
}

func (c *userRoleCache) SetByUserID(userID uint, userRoles []models.UserRole) {
	if c.store == nil {
		return
	}
	go func(urs []models.UserRole) {
		key := fmt.Sprintf("%s:%d", models.PkgUserRoleGetByUserID, userID)
		if bytes, err := json.Marshal(urs); err == nil {
			_ = c.store.Set(key, bytes, utils.RandomJitter(cacheTTLList))
		}
	}(userRoles)
}

func (c *userRoleCache) InvalidateByUserID(userID uint) {
	if c.store == nil {
		return
	}
	go func() {
		key := fmt.Sprintf("%s:%d", models.PkgUserRoleGetByUserID, userID)
		_ = c.store.Delete(key)
	}()
}
