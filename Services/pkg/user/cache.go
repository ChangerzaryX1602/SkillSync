package user

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
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
	c.InvalidateAllLists()
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
func (c *userCache) InvalidateAllLists() {
	if c.store == nil {
		return
	}
	pattern := fmt.Sprintf("%s:*", models.PkgUserGetUsers)
	conn := c.store.Conn()
	ctx := context.Background()
	keys, err := conn.Keys(ctx, pattern).Result()
	if err != nil {
		return
	}
	if len(keys) == 0 {
		return
	}
	const numWorkers = 100
	keyChan := make(chan string, len(keys))
	var wg sync.WaitGroup
	workerCount := numWorkers
	if len(keys) < numWorkers {
		workerCount = len(keys)
	}
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for key := range keyChan {
				_ = c.store.Delete(key)
			}
		}()
	}
	for _, key := range keys {
		keyChan <- key
	}
	close(keyChan)
	wg.Wait()
}
func (c *userCache) GetList(pagination models.Pagination, search models.Search) ([]models.User, *models.Pagination, *models.Search, bool) {
	if c.store == nil {
		return nil, nil, nil, false
	}
	key := fmt.Sprintf("%s:%s:%s", models.PkgUserGetUsers, pagination.GetPaginationString(), search.GetSearchString())
	fmt.Println("key get", key)
	bytes, err := c.store.Get(key)
	if err != nil || len(bytes) == 0 {
		return nil, nil, nil, false
	}

	var usersCache UsersCache
	if err := json.Unmarshal(bytes, &usersCache); err != nil {
		fmt.Println("err", err)
		return nil, nil, nil, false
	}
	return usersCache.Users, &usersCache.Pagination, &usersCache.Search, true
}

func (c *userCache) SetList(pagination models.Pagination, search models.Search, users []models.User) {
	if c.store == nil {
		return
	}
	go func(us []models.User, p models.Pagination, s models.Search) {
		usersCache := UsersCache{
			Users:      us,
			Pagination: p,
			Search:     s,
		}
		key := fmt.Sprintf("%s:%s:%s", models.PkgUserGetUsers, p.GetPaginationString(), s.GetSearchString())
		if bytes, err := json.Marshal(usersCache); err == nil {
			fmt.Println("key set", key)
			err = c.store.Set(key, bytes, utils.RandomJitter(cacheTTLList))
			if err != nil {
				fmt.Println(err)
			}
		}
	}(users, pagination, search)
}

type UsersCache struct {
	Users []models.User `json:"users"`
	models.Pagination
	models.Search
}
