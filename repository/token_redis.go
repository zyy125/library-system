package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenRdb struct {
	rdb *redis.Client
} 

var Rdb *TokenRdb

const (
	RefreshTokenPrefix = "refresh_token:"  // Refresh Token 前缀
	BlacklistPrefix    = "blacklist:"      // 黑名单前缀
	RefreshTokenTTL    = 7 * 24 * time. Hour // 7天
)

func NewRedis(rdb *redis.Client) {
	Rdb = &TokenRdb{
		rdb: rdb,
	}
}

func (r *TokenRdb)StoreRefreshToken(ctx context.Context, userID uint64, tokenID, tokenString string) error {
	key := fmt.Sprintf("%s%d:%s", RefreshTokenPrefix, userID, tokenID)
	return r.rdb.Set(ctx, key, tokenString, RefreshTokenTTL).Err()
}

func (r *TokenRdb)GetRefreshToken(ctx context.Context, userID uint64, tokenID string) (string, error) {
	key := fmt.Sprintf("%s%d:%s", RefreshTokenPrefix, userID, tokenID)
	return r.rdb.Get(ctx, key).Result()
}

// DeleteRefreshToken 删除 Refresh Token
func (r *TokenRdb)DeleteRefreshToken(ctx context.Context, userID uint64, tokenID string) error {
	key := fmt. Sprintf("%s%d:%s", RefreshTokenPrefix, userID, tokenID)
	return r.rdb.Del(ctx, key).Err()
}

// DeleteAllUserRefreshTokens 删除用户的所有 Refresh Token（踢出所有设备）
func (r *TokenRdb)DeleteAllUserRefreshTokens(ctx context.Context, userID uint64) error {
	pattern := fmt.Sprintf("%s%d:*", RefreshTokenPrefix, userID)
	iter := r.rdb.Scan(ctx, 0, pattern, 0).Iterator()
	
	for iter.Next(ctx) {
		if err := r.rdb.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	
	return iter.Err()
}

// AddToBlacklist 将 Token 加入黑名单
func (r *TokenRdb)AddToBlacklist(ctx context.Context, tokenID string, ttl time.Duration) error {
	key := fmt.Sprintf("%s%s", BlacklistPrefix, tokenID)
	return r.rdb.Set(ctx, key, "1", ttl).Err()
}

// IsInBlacklist 检查 Token 是否在黑名单
func (r *TokenRdb)IsInBlacklist(ctx context.Context, tokenID string) (bool, error) {
	key := fmt.Sprintf("%s%s", BlacklistPrefix, tokenID)
	result, err := r.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}