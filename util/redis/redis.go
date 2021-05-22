package redis

import (
	"context"
	"fsbm/conf"
	"github.com/go-redis/redis/v8"
	"strconv"
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

func GetInt64WithRetry(ctx context.Context, key string) (res int64, err error) {
	var val string
	for i := 0; i < 5; i++ {
		val, err = redisClient.Get(ctx, key).Result()
		if err == nil {
			break
		}
		if err == redis.Nil {
			return 0, nil
		}
	}
	res, err = strconv.ParseInt(val, 10, 64)
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

func IncrByWithRetry(ctx context.Context, key string, value int64) (err error) {
	for i := 0; i < 5; i++ {
		_, err = redisClient.IncrBy(ctx, key, value).Result()
		if err == nil {
		}
	}
	return
}

func IncrWithRetry(ctx context.Context, key string) (err error) {
	for i := 0; i < 5; i++ {
		_, err = redisClient.Incr(ctx, key).Result()
		if err == nil {
		}
	}
	return
}
