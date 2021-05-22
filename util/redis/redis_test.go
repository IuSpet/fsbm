package redis

import (
	"context"
	"fmt"
	"fsbm/conf"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetWithRetry(t *testing.T) {
	key := "test_key"
	value := "test_value"
	ctx := context.Background()
	res, err := GetWithRetry(ctx, key)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, value, res)
}

func TestSetWithRetry(t *testing.T) {
	key := "test_key"
	value := "test_value"
	ctx := context.Background()
	err := SetWithRetry(ctx, key, value, 30*time.Minute)
	if err != nil {
		panic(err)
	}
}

func TestIncrByWithRetry(t *testing.T) {
	key := "test_incrby"
	ctx := context.Background()
	conf.Init()
	err := IncrByWithRetry(ctx, key, 25)
	if err != nil {
		panic(err)
	}
	res, err := GetInt64WithRetry(ctx, key)
	if err != nil {
		panic(err)
	}
	_ = res
}

func TestRawRedis(t *testing.T) {
	key := "test_incrby"
	ctx := context.Background()
	conf.Init()
	res, err := redisClient.IncrBy(ctx, key, 25).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
	value, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(value)
}
