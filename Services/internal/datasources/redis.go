package datasources

import (
	"context"
	"fmt"
	"runtime"

	"github.com/gofiber/storage/redis"
)

type RedisConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	DbName   int
}

func ConnectToRedis(config RedisConfig) (storage *redis.Storage, err error) {
	// Initialize redis storage
	storage = redis.New(redis.Config{
		Host:      config.Host,
		Port:      config.Port,
		Username:  config.Username,
		Password:  config.Password,
		Database:  config.DbName,
		Reset:     false,
		TLSConfig: nil,
		PoolSize:  10 * runtime.GOMAXPROCS(0),
	})

	// test connection
	if err = storage.Conn().Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return
}
