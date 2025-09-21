package utils

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func GetRedis(c context.Context, rdb *redis.Client, key string, dest any) (bool, error) {
	cmd := rdb.Get(c, key)
	if cmd.Err() != nil {
		if cmd.Err() == redis.Nil {
			log.Printf("Key %s does not exist\n", key)
			return false, nil
		}

		log.Println("Redis Error\nCause:", cmd.Err().Error())
		return false, nil
	}

	data, err := cmd.Bytes()
	if err != nil {
		log.Println("Redis Bytes Error\nCause:", err.Error())
		return false, nil
	}

	if err := json.Unmarshal(data, dest); err != nil {
		log.Println("Redis Unmarshal Error\nCause:", err.Error())
		return false, nil
	}

	return true, nil
}

func SetRedis(c context.Context, rdb *redis.Client, key string, value any, time time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		log.Println("Redis Marshal Error\nCause:", err.Error())
		return nil
	}

	if err := rdb.Set(c, key, data, time).Err(); err != nil {
		log.Println("Redis Set Error\nCause:", err.Error())
		return nil
	}

	return nil
}
