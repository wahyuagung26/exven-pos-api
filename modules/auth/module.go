package auth

import (
	"github.com/exven/pos-system/internal/config"
	"github.com/exven/pos-system/modules/auth/domain"
	"github.com/exven/pos-system/modules/auth/persistence"
	"github.com/exven/pos-system/modules/auth/services"
	"github.com/exven/pos-system/shared/container"
	"github.com/exven/pos-system/shared/infrastructure/cache"
	"github.com/exven/pos-system/shared/infrastructure/messaging"
	"gorm.io/gorm"
)

type Module struct {
	container container.Container
	db        *gorm.DB
	redis     *cache.RedisClient
	eventBus  messaging.EventBus
	jwtConfig config.JWTConfig
}

func NewModule(
	container container.Container,
	db *gorm.DB,
	redis *cache.RedisClient,
	eventBus messaging.EventBus,
	jwtConfig config.JWTConfig,
) *Module {
	return &Module{
		container: container,
		db:        db,
		redis:     redis,
		eventBus:  eventBus,
		jwtConfig: jwtConfig,
	}
}

func (m *Module) Register() {
	m.container.RegisterSingleton("auth.userRepository", func() interface{} {
		return persistence.NewUserRepository(m.db)
	})

	m.container.RegisterSingleton("auth.sessionRepository", func() interface{} {
		return persistence.NewSessionRepository()
	})

	m.container.RegisterSingleton("auth.tokenService", func() interface{} {
		return services.NewTokenService(
			m.jwtConfig.Secret,
			m.jwtConfig.ExpiryHours,
			m.jwtConfig.RefreshExpiryDays,
		)
	})

	m.container.RegisterSingleton("auth.passwordService", func() interface{} {
		return services.NewPasswordService()
	})

	m.container.RegisterSingleton("auth.service", func() interface{} {
		userRepo := m.container.MustGet("auth.userRepository").(domain.UserRepository)
		sessionRepo := m.container.MustGet("auth.sessionRepository").(domain.SessionRepository)
		tokenService := m.container.MustGet("auth.tokenService").(domain.TokenService)
		passwordService := m.container.MustGet("auth.passwordService").(domain.PasswordService)

		return services.NewAuthService(
			userRepo,
			sessionRepo,
			tokenService,
			passwordService,
			m.eventBus,
		)
	})
}
