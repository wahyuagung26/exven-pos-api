package services

import (
	"fmt"
	"time"

	"github.com/exven/pos-system/modules/auth/domain"
	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	jwtSecret          string
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

func NewTokenService(jwtSecret string, accessExpiryHours, refreshExpiryDays int) *TokenService {
	return &TokenService{
		jwtSecret:          jwtSecret,
		accessTokenExpiry:  time.Duration(accessExpiryHours) * time.Hour,
		refreshTokenExpiry: time.Duration(refreshExpiryDays) * 24 * time.Hour,
	}
}

func (s *TokenService) GenerateAccessToken(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":   user.ID,
		"tenant_id": user.TenantID,
		"username":  user.Username,
		"role_id":   user.RoleID,
		"exp":       time.Now().Add(s.accessTokenExpiry).Unix(),
		"iat":       time.Now().Unix(),
		"type":      "access",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *TokenService) GenerateRefreshToken(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":   user.ID,
		"tenant_id": user.TenantID,
		"exp":       time.Now().Add(s.refreshTokenExpiry).Unix(),
		"iat":       time.Now().Unix(),
		"type":      "refresh",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *TokenService) ValidateAccessToken(tokenString string) (*domain.TokenClaims, error) {
	return s.validateToken(tokenString, "access")
}

func (s *TokenService) ValidateRefreshToken(tokenString string) (*domain.TokenClaims, error) {
	return s.validateToken(tokenString, "refresh")
}

func (s *TokenService) validateToken(tokenString string, tokenType string) (*domain.TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	if claims["type"] != tokenType {
		return nil, fmt.Errorf("invalid token type")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid user_id in token")
	}

	tenantID, ok := claims["tenant_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid tenant_id in token")
	}

	username, ok := claims["username"].(string)
	if !ok && tokenType == "access" {
		return nil, fmt.Errorf("invalid username in token")
	}

	roleID := float64(0)
	if tokenType == "access" {
		roleID, ok = claims["role_id"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid role_id in token")
		}
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid exp in token")
	}

	return &domain.TokenClaims{
		UserID:    uint64(userID),
		TenantID:  uint64(tenantID),
		Username:  username,
		RoleID:    uint64(roleID),
		ExpiresAt: time.Unix(int64(exp), 0),
	}, nil
}