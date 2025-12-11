package database

import (
	"context"
	"library-system/config"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

// InitRedis 初始化 Redis 连接
func InitRedis() (*redis.Client, error) {
	cfg := config.GetRedisConfig()
	RedisClient = redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10, // 连接池大小
		MinIdleConns: 5,  // 最小空闲连接数
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := RedisClient.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("Redis 连接失败: %v", err)
	}

	log.Println("Redis 已连接")
	return RedisClient, nil
}

func CloseRedis() {
	if RedisClient != nil {
		if err := RedisClient.Close(); err != nil {
			log. Printf("Redis 关闭失败: %v", err)
		}
	}
}

