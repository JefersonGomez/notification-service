package redisclient

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func ConectarRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_URL"),
	})

	ctx := context.Background()

	_, err := RDB.Ping(ctx).Result()
	if err != nil {
		panic("No se pudo conectar a Redis: " + err.Error())
	}

	fmt.Println("Redis conectado")
}
