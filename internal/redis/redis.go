package redis

import (
	"errors"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var CacheConn *redis.Client

var (
	ErrRedisAddrAndPassIsRequiredEnvVariables = errors.New("redis addr and password is required env variables")
)

func EstablishRedisConnection() (*redis.Client, error) {
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPass := os.Getenv("REDIS_PASS")

	if redisAddr == "" || redisPass == "" {
		return nil, ErrRedisAddrAndPassIsRequiredEnvVariables
	}

	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		redisDB = 1
	}

	redisProtocol, err := strconv.Atoi(os.Getenv("REDIS_PROTOCOL"))
	if err != nil {
		redisProtocol = 2
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPass,
		DB:       redisDB,
		Protocol: redisProtocol,
	})

	return rdb, nil
}
