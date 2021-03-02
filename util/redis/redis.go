package redis

import (
	"context"
	"fsbm/conf"
	"github.com/go-redis/redis/v8"
	"time"
)

var redisClient *redis.Client

func init() {
	redisCfg := conf.GlobalConfig.Redis
	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisCfg.Addr,
		Password: redisCfg.Password,
		DB:       redisCfg.DB,
	})
}

func GetWithRetry(ctx context.Context, key string) (res string, err error) {
	for i := 0; i < 5; i++ {
		res, err = redisClient.Get(ctx, key).Result()
		if err == nil {
			return
		}
		if err == redis.Nil {
			return "", nil
		}
	}
	return
}

func SetWithRetry(ctx context.Context, key string, value interface{}, expiration time.Duration) (err error) {
	for i := 0; i < 5; i++ {
		_, err = redisClient.Set(ctx, key, value, expiration).Result()
		if err == nil {
			return
		}
	}
	return
}

func Del(ctx context.Context, key string) {
	redisClient.Del(ctx, key)
	return
}
