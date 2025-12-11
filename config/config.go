package config

import (
	"os"
)

type Config struct {
	DSN string
}

type RedisConfig struct {
	Addr     string // Redis 地址
	Password string // Redis 密码
	DB       int    // Redis 数据库编号
}

func Load() *Config {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == ""{
		dsn = "root:@tcp(127.0.0.1:3306)/librarySystem?charset=utf8mb4&parseTime=True&loc=Local"
	}

	return &Config{
		DSN: dsn,
	}
}

func GetRedisConfig() *RedisConfig {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379" // 默认地址
	}

	password := os.Getenv("REDIS_PASSWORD")
	// password 默认为空

	return &RedisConfig{
		Addr:     addr,
		Password: password,
		DB:       0, // 默认使用 0 号数据库
	}
}