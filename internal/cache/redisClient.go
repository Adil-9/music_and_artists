package cache

import (
	"groupie-tracker/logger"
	"os"

	"github.com/redis/go-redis/v9"
	"gopkg.in/yaml.v3"
)

var RedisClient *redis.Client

// initializing redis client to send data in cache
func RedisClientInit() {
	redisConfFile, err := os.ReadFile("internal/cache/redisConfig.yaml")
	if err != nil {
		logger.ErrorLog.Fatal("Error reading redis config file:", err)
	}

	var redisClt redis.Options

	if err = yaml.Unmarshal(redisConfFile, &redisClt); err != nil {
		logger.ErrorLog.Fatal("Error redis file unmarshalling:", err)
	}

	RedisClient = redis.NewClient(&redisClt)
	// logger.ErrorLog.Println("check 2")
}
