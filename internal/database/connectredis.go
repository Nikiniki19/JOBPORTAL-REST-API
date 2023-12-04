package database

import (
	"job-portal-api/config"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func RedisConnection() *redis.Client {
	cfg := config.GetConfig()
	PASS, _ := strconv.Atoi(cfg.RedisConfig.Db)
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisConfig.Address,
		Password: cfg.RedisConfig.Password, // no password set
		DB:       PASS,                     // use default DB
	})
	return rdb
}
