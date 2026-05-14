package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	redisClient *redis.Client
}

func NewRedisCache(redisURL string) *RedisCache {
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	return &RedisCache{redisClient: redisClient}
}

func (c *RedisCache) SetBlacklistedToken(ctx context.Context, tokenHash string, expirationDuration time.Duration) error {
	return c.redisClient.Set(ctx, "blacklist:"+tokenHash, true, expirationDuration).Err()
}

func (c *RedisCache) IsTokenBlacklisted(ctx context.Context, tokenHash string) (bool, error) {
	result, err := c.redisClient.Get(ctx, "blacklist:"+tokenHash).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return result == "1" || result == "true", nil
}

func (c *RedisCache) SetUserSession(ctx context.Context, userID string, sessionData string, expirationDuration time.Duration) error {
	return c.redisClient.Set(ctx, "session:"+userID, sessionData, expirationDuration).Err()
}

func (c *RedisCache) GetUserSession(ctx context.Context, userID string) (string, error) {
	result, err := c.redisClient.Get(ctx, "session:"+userID).Result()
	if err == redis.Nil {
		return "", nil
	}
	return result, err
}

func (c *RedisCache) DeleteUserSession(ctx context.Context, userID string) error {
	return c.redisClient.Del(ctx, "session:"+userID).Err()
}

func (c *RedisCache) Ping(ctx context.Context) error {
	return c.redisClient.Ping(ctx).Err()
}
