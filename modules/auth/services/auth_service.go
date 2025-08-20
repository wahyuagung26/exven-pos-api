package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/exven/pos-system/modules/auth/domain"
	"github.com/exven/pos-system/shared/infrastructure/messaging"
)

type AuthService struct {
	userRepo        domain.UserRepository
	sessionRepo     domain.SessionRepository
	tokenService    domain.TokenService
	passwordService domain.PasswordService
	eventBus        messaging.EventBus
}

func NewAuthService(
	userRepo domain.UserRepository,
	sessionRepo domain.SessionRepository,
	tokenService domain.TokenService,
	passwordService domain.PasswordService,
	eventBus messaging.EventBus,
) *AuthService {
	return &AuthService{
		userRepo:        userRepo,
		sessionRepo:     sessionRepo,
		tokenService:    tokenService,
		passwordService: passwordService,
		eventBus:        eventBus,
	}
}

func (s *AuthService) Login(ctx context.Context, credentials domain.LoginCredentials) (*domain.TokenPair, *domain.User, error) {
	// Try to find user by username first (globally across all tenants)
	user, err := s.userRepo.FindByUsernameGlobal(ctx, credentials.Username)
	if err != nil {
		// If username not found, try by email
		user, err = s.userRepo.FindByEmailGlobal(ctx, credentials.Username)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid credentials")
		}
	}

	if !user.IsActive {
		return nil, nil, fmt.Errorf("user account is inactive")
	}

	if err := s.passwordService.VerifyPassword(user.PasswordHash, credentials.Password); err != nil {
		return nil, nil, fmt.Errorf("invalid credentials")
	}

	accessToken, err := s.tokenService.GenerateAccessToken(user)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(user)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	session := &domain.Session{
		ID:           generateSessionID(),
		UserID:       user.ID,
		TenantID:     user.TenantID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(24 * time.Hour),
		CreatedAt:    time.Now(),
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, nil, fmt.Errorf("failed to create session: %w", err)
	}

	now := time.Now()
	user.LastLoginAt = &now
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, nil, fmt.Errorf("failed to update user: %w", err)
	}

	event := messaging.NewEvent("user.logged_in", user.TenantID, user.ID, map[string]interface{}{
		"username": user.Username,
		"ip":       ctx.Value("ip"),
	})

	fmt.Printf("Publishing login event, eventBus is nil: %t\n", s.eventBus == nil)
	if s.eventBus != nil {
		s.eventBus.Publish(ctx, "auth.login", event)
	} else {
		fmt.Printf("Skipping event publishing as eventBus is nil\n")
	}

	return &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    3600,
	}, user, nil
}

func (s *AuthService) Register(ctx context.Context, req domain.RegisterRequest) (*domain.User, error) {
	hashedPassword, err := s.passwordService.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &domain.User{
		Username:  req.Username,
		Email:     req.Email,
		FullName:  req.FullName,
		Phone:     req.Phone,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	user.PasswordHash = hashedPassword

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	event := messaging.NewEvent("user.registered", user.TenantID, user.ID, map[string]interface{}{
		"username": user.Username,
		"email":    user.Email,
	})
	s.eventBus.Publish(ctx, "auth.register", event)

	return user, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*domain.TokenPair, error) {
	claims, err := s.tokenService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if !user.IsActive {
		return nil, fmt.Errorf("user account is inactive")
	}

	newAccessToken, err := s.tokenService.GenerateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, err := s.tokenService.GenerateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &domain.TokenPair{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    3600,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, userID uint64) error {
	if err := s.sessionRepo.DeleteByUserID(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete sessions: %w", err)
	}

	event := messaging.NewEvent("user.logged_out", 0, userID, map[string]interface{}{})
	s.eventBus.Publish(ctx, "auth.logout", event)

	return nil
}

func (s *AuthService) ValidateToken(ctx context.Context, token string) (*domain.User, error) {
	claims, err := s.tokenService.ValidateAccessToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if !user.IsActive {
		return nil, fmt.Errorf("user account is inactive")
	}

	return user, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, userID uint64, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	if err := s.passwordService.VerifyPassword(user.PasswordHash, oldPassword); err != nil {
		return fmt.Errorf("invalid old password")
	}

	hashedPassword, err := s.passwordService.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.PasswordHash = hashedPassword
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if err := s.sessionRepo.DeleteByUserID(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete sessions: %w", err)
	}

	return nil
}

func (s *AuthService) ResetPassword(ctx context.Context, email string) error {
	user, err := s.userRepo.FindByEmailGlobal(ctx, email)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	event := messaging.NewEvent("password.reset_requested", user.TenantID, user.ID, map[string]interface{}{
		"email": user.Email,
	})
	s.eventBus.Publish(ctx, "auth.password_reset", event)

	return nil
}

func (s *AuthService) VerifyEmail(ctx context.Context, token string) error {
	return errors.New("not implemented")
}

func generateSessionID() string {
	return fmt.Sprintf("sess_%d", time.Now().UnixNano())
}

type PasswordHash interface {
	Hash(password string) (string, error)
	Verify(hashedPassword, password string) error
}
