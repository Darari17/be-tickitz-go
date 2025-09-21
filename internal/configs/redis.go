package configs

import (
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

func InitRedis() (*redis.Client, error) {
	rdbHost := os.Getenv("RDBHOST")
	rdbPort := os.Getenv("RDBPORT")
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", rdbHost, rdbPort),
	})

	return rdb, nil
}
