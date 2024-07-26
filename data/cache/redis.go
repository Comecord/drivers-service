package cache

import (
	"context"
	"drivers-service/config"
	"drivers-service/pkg/logging"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

var logcod = logging.NewLogger(config.GetConfig())
var redisClient *redis.Client

func InitRedis(cfg *config.Config, ctx context.Context) error {
	redisClient = redis.NewClient(&redis.Options{
		Addr:               fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password:           cfg.Redis.Password,
		DB:                 cfg.Redis.Db,
		DialTimeout:        cfg.Redis.DialTimeout,
		ReadTimeout:        cfg.Redis.ReadTimeout,
		WriteTimeout:       cfg.Redis.WriteTimeout,
		PoolSize:           cfg.Redis.PoolSize,
		PoolTimeout:        cfg.Redis.PoolTimeout,
		IdleTimeout:        500 * time.Millisecond,
		IdleCheckFrequency: cfg.Redis.IdleCheckFrequency * time.Millisecond,
	})

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return err
	}
	logcod.Info(logging.Redis, logging.Connection, "Redis cache init", nil)
	return nil
}

func GetRedis() *redis.Client {
	return redisClient
}

func CloseRedis() {
	redisClient.Close()
}

func Set[T any](ctx context.Context, c *redis.Client, key string, value T, duration time.Duration) error {

	v, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.Set(ctx, key, v, duration).Err()
}

func Get[T any](ctx context.Context, c *redis.Client, key string) (T, error) {
	var dest T = *new(T)
	v, err := c.Get(ctx, key).Result()
	if err != nil {
		return dest, err
	}
	err = json.Unmarshal([]byte(v), &dest)
	if err != nil {
		return dest, err
	}
	return dest, nil
}

func HSet[T any](ctx context.Context, c *redis.Client, hashKey string, field string, value T) error {
	v, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.HSet(ctx, hashKey, field, v).Err()
}

func HGet[T any](ctx context.Context, c *redis.Client, hashKey string, field string) (T, error) {
	var dest T = *new(T)
	v, err := c.HGet(ctx, hashKey, field).Result()
	if err != nil {
		return dest, err
	}
	err = json.Unmarshal([]byte(v), &dest)
	if err != nil {
		return dest, err
	}
	return dest, nil
}

func ZAdd[T any](ctx context.Context, c *redis.Client, zsetKey string, score float64, value T) error {
	v, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.ZAdd(ctx, zsetKey, &redis.Z{
		Score:  score,
		Member: v,
	}).Err()
}

func ZRange[T any](ctx context.Context, c *redis.Client, zsetKey string, start, stop int64) ([]T, error) {
	vals, err := c.ZRange(ctx, zsetKey, start, stop).Result()
	if err != nil {
		return nil, err
	}

	var result []T
	for _, v := range vals {
		var dest T
		err = json.Unmarshal([]byte(v), &dest)
		if err != nil {
			return nil, err
		}
		result = append(result, dest)
	}
	return result, nil
}
