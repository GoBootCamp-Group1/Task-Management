package cache

import (
	"github.com/GoBootCamp-Group1/Task-Management/config"
	"github.com/redis/go-redis/v9"
	"strconv"
)

func NewRedisConnection(redisConfig config.Redis) (*redis.Client, error) {
	db, err := strconv.Atoi(redisConfig.DB)

	if err != nil {
		return nil, err
	}

	client := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Address,
		Password: redisConfig.Password,
		DB:       db,
	})

	return client, nil
}
