package roles

import (
	"github.com/exven/pos-system/modules/roles/handlers"
	"github.com/exven/pos-system/modules/roles/persistence"
	"github.com/exven/pos-system/modules/roles/services"
	"github.com/exven/pos-system/shared/container"
	"github.com/exven/pos-system/shared/infrastructure/messaging"
	"gorm.io/gorm"
)

type Module struct {
	container container.Container
	db        *gorm.DB
	eventBus  messaging.EventBus
}

func NewModule(
	container container.Container,
	db *gorm.DB,
	eventBus messaging.EventBus,
) *Module {
	return &Module{
		container: container,
		db:        db,
		eventBus:  eventBus,
	}
}

func (m *Module) Register() {
	// Register repositories
	m.container.RegisterSingleton("roles.repository", func() interface{} {
		return persistence.NewRoleRepository(m.db)
	})

	// Register services
	m.container.RegisterSingleton("roles.service", func() interface{} {
		repo := persistence.NewRoleRepository(m.db)
		return services.NewRoleService(repo)
	})

	// Register handlers
	m.container.RegisterSingleton("roles.handler", func() interface{} {
		repo := persistence.NewRoleRepository(m.db)
		service := services.NewRoleService(repo)
		return handlers.NewRoleHandler(service)
	})
}

func (m *Module) GetHandler() *handlers.RoleHandler {
	repo := persistence.NewRoleRepository(m.db)
	service := services.NewRoleService(repo)
	return handlers.NewRoleHandler(service)
}