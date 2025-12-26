package user

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

type userCache struct {
	store *redis.Storage
}

func newUserCache(store *redis.Storage) *userCache {
	return &userCache{store: store}
}

func (c *userCache) Get(id uint) (*models.User, bool) {
	if c.store == nil {
		return nil, false
	}
	key := fmt.Sprintf("%s:%d", models.PkgUserGetUser, id)

	bytes, err := c.store.Get(key)
	if err != nil || len(bytes) == 0 {
		return nil, false
	}

	var user models.User
	if err := json.Unmarshal(bytes, &user); err != nil {
		return nil, false
	}
	return &user, true
}

func (c *userCache) Set(user *models.User) {
	if c.store == nil || user == nil {
		return
	}

	go func(u models.User) {
		key := fmt.Sprintf("%s:%d", models.PkgUserGetUser, u.ID)
		if bytes, err := json.Marshal(u); err == nil {
			_ = c.store.Set(key, bytes, utils.RandomJitter(cacheTTL))
		}
	}(*user)
}

func (c *userCache) Invalidate(id uint, email string) {
	if c.store == nil {
		return
	}
	go func() {
		key := fmt.Sprintf("%s:%d", models.PkgUserGetUser, id)
		_ = c.store.Delete(key)
	}()
	go func() {
		key := fmt.Sprintf("%s:%s", models.PkgUserGetUserGmail, email)
		_ = c.store.Delete(key)
	}()
}

func (c *userCache) GetByEmail(email string) (*models.User, bool) {
	if c.store == nil {
		return nil, false
	}
	key := fmt.Sprintf("%s:%s", models.PkgUserGetUserGmail, email)

	bytes, err := c.store.Get(key)
	if err != nil || len(bytes) == 0 {
		return nil, false
	}

	var user models.User
	if err := json.Unmarshal(bytes, &user); err != nil {
		return nil, false
	}
	return &user, true
}

func (c *userCache) SetByEmail(email string, user *models.User) {
	if c.store == nil || user == nil {
		return
	}
	go func(u models.User) {
		key := fmt.Sprintf("%s:%s", models.PkgUserGetUserGmail, email)
		if bytes, err := json.Marshal(u); err == nil {
			_ = c.store.Set(key, bytes, utils.RandomJitter(cacheTTL))
		}
	}(*user)
}

func (c *userCache) GetList(pagination models.Pagination, search models.Search) ([]models.User, bool) {
	if c.store == nil {
		return nil, false
	}
	key := fmt.Sprintf("%s:%s:%s", models.PkgUserGetUsers, pagination.GetPaginationString(), search.GetSearchString())
	fmt.Println("key get", key)
	bytes, err := c.store.Get(key)
	if err != nil || len(bytes) == 0 {
		return nil, false
	}

	var users []models.User
	if err := json.Unmarshal(bytes, &users); err != nil {
		fmt.Println("err", err)
		return nil, false
	}
	return users, true
}

func (c *userCache) SetList(pagination models.Pagination, search models.Search, users []models.User) {
	if c.store == nil {
		return
	}
	go func(us []models.User, p models.Pagination, s models.Search) {
		key := fmt.Sprintf("%s:%s:%s", models.PkgUserGetUsers, p.GetPaginationString(), s.GetSearchString())
		if bytes, err := json.Marshal(us); err == nil {
			fmt.Println("key set", key)
			err = c.store.Set(key, bytes, utils.RandomJitter(cacheTTLList))
			if err != nil {
				fmt.Println(err)
			}
		}
	}(users, pagination, search)
}
