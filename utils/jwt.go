package utils

import (
    "log"
    "os"
    "time"
    "errors"
    "github.com/google/uuid"
    "github.com/golang-jwt/jwt/v5"
)

var accessTokenSecret []byte
var refreshTokenSecret []byte

func init() {
    accessTokenSecret := os.Getenv("JWT_ACCESS_SECRET")
    if accessTokenSecret == "" {
        log.Fatal("JWT_ACCESS_SECRET 环境变量未设置")
    }

    refreshTokenSecre := os.Getenv("JWT_REFRESH_SECRET")
    if refreshTokenSecre == "" {
        log.Fatal("JWT_REFRESH_SECRET 环境变量未设置")
    }
}

// AccessTokenClaims Access Token 的 Claims
type AccessTokenClaims struct {
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	TokenID  string `json:"token_id"` // 用于黑名单
	jwt.RegisteredClaims
}

// RefreshTokenClaims Refresh Token 的 Claims
type RefreshTokenClaims struct {
	UserID  uint64 `json:"user_id"`
	TokenID string `json:"token_id"` // 唯一标识
	jwt. RegisteredClaims
}
// ========== Token 生成 ==========

// GenerateTokenPair 生成 Access Token 和 Refresh Token
func GenerateTokenPair(userID uint64, username, role string) (accessToken, refreshToken, tokenID string, err error) {
	// 生成唯一的 TokenID
	tokenID = generateTokenID()
	if err != nil {
		return "", "", "", err
	}

	// 生成 Access Token
	accessToken, err = GenerateAccessToken(userID, username, role, tokenID)
	if err != nil {
		return "", "", "", err
	}

	// 生成 Refresh Token
	refreshToken, err = GenerateRefreshToken(userID, tokenID)
	if err != nil {
		return "", "", "", err
	}

	return accessToken, refreshToken, tokenID, nil
}

// GenerateAccessToken 生成 Access Token
func GenerateAccessToken(userID uint64, username, role, tokenID string) (string, error) {
	now := time.Now()
	claims := AccessTokenClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		TokenID:  tokenID,
		RegisteredClaims: jwt. RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)), // 24小时
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(accessTokenSecret)
}

// GenerateRefreshToken 生成 Refresh Token
func GenerateRefreshToken(userID uint64, tokenID string) (string, error) {
	now := time.Now()
	claims := RefreshTokenClaims{
		UserID:  userID,
		TokenID: tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(7 * 24 * time. Hour)), // 7天
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(refreshTokenSecret)
}

// ========== Token 验证 ==========

// ValidateAccessToken 验证 Access Token
func ValidateAccessToken(tokenString string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(token *jwt. Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("无效的签名方法")
		}
		return accessTokenSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*AccessTokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的 Token")
}

// ValidateRefreshToken 验证 Refresh Token
func ValidateRefreshToken(tokenString string) (*RefreshTokenClaims, error) {
	token, err := jwt. ParseWithClaims(tokenString, &RefreshTokenClaims{}, func(token *jwt. Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("无效的签名方法")
		}
		return refreshTokenSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token. Claims.(*RefreshTokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors. New("无效的 Refresh Token")
}

func generateTokenID() string {
    return uuid.NewString()
}

// GetRemainingTime 获取 Token 剩余时间
func GetRemainingTime(expiresAt time.Time) time.Duration {
	return time.Until(expiresAt)
}
