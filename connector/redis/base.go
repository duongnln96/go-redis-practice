package connector

import (
	"app/config"
	"log"
	"sync"

	"github.com/go-redis/redis"
)

var (
	redisClient *redis.Client
	once        sync.Once
)

type RedisConnector struct {
	Conn *redis.Client
}

func NewRedisConnector() *RedisConnector {

	once.Do(func() {
		redisClient = redis.NewClient(
			&redis.Options{
				Addr:     config.Addr,
				Password: config.Password,
				DB:       config.DB,
			},
		)

		if _, err := redisClient.Ping().Result(); err != nil {
			log.Fatalf("Connect to Redis fail: %v\n", err)
		}
	})

	return &RedisConnector{
		Conn: redisClient,
	}
}
