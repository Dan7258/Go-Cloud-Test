package models

import (
	"context"
	"fmt"
	"time"
)

func SetDataInRedis(key string, value []byte, timeLive time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := RDB.Set(ctx, fmt.Sprintf("%s", key), value, timeLive).Err()
	return err
}

func GetDataFromRedis(key string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return RDB.Get(ctx, key).Bytes()
}

func GetAllKeysFromRedis() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return RDB.Keys(ctx, "*").Result()
}
