package redis

import (
	"os"
	"strconv"

	"github.com/gofiber/storage/redis/v3"
)

var redisClient *redis.Storage

func Connect() {
	var redisPort, _ = strconv.Atoi(os.Getenv("REDIS_PORT"))
	var redisDB, _ = strconv.Atoi(os.Getenv("REDIS_DB"))

	redisClient = redis.New(redis.Config{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     redisPort,
		Database: redisDB,
		// redis cluster setup
		// Addrs:    strings.Split(strings.TrimSpace(os.Getenv("REDIS_CLUSTER")), ","),
	})
}
