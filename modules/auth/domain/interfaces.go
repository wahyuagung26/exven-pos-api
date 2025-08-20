package domain

import (
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uint64) error
	FindByID(ctx context.Context, id uint64) (*User, error)
	FindByUsername(ctx context.Context, tenantID uint64, username string) (*User, error)
	FindByEmail(ctx context.Context, tenantID uint64, email string) (*User, error)
	FindByUsernameGlobal(ctx context.Context, username string) (*User, error)
	FindByEmailGlobal(ctx context.Context, email string) (*User, error)
	FindAll(ctx context.Context, tenantID uint64, limit, offset int) ([]*User, error)
	Count(ctx context.Context, tenantID uint64) (int64, error)
}

type SessionRepository interface {
	Create(ctx context.Context, session *Session) error
	Delete(ctx context.Context, sessionID string) error
	FindByID(ctx context.Context, sessionID string) (*Session, error)
	FindByUserID(ctx context.Context, userID uint64) ([]*Session, error)
	DeleteByUserID(ctx context.Context, userID uint64) error
	DeleteExpired(ctx context.Context) error
}

type AuthService interface {
	Login(ctx context.Context, credentials LoginCredentials) (*TokenPair, *User, error)
	Register(ctx context.Context, req RegisterRequest) (*User, error)
	RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error)
	Logout(ctx context.Context, userID uint64) error
	ValidateToken(ctx context.Context, token string) (*User, error)
	ChangePassword(ctx context.Context, userID uint64, oldPassword, newPassword string) error
	ResetPassword(ctx context.Context, email string) error
	VerifyEmail(ctx context.Context, token string) error
}

type TokenService interface {
	GenerateAccessToken(user *User) (string, error)
	GenerateRefreshToken(user *User) (string, error)
	ValidateAccessToken(token string) (*TokenClaims, error)
	ValidateRefreshToken(token string) (*TokenClaims, error)
}

type PasswordService interface {
	HashPassword(password string) (string, error)
	VerifyPassword(hashedPassword, password string) error
}