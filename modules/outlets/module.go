package outlets

import (
	"github.com/exven/pos-system/modules/outlets/handlers"
	"github.com/exven/pos-system/modules/outlets/persistence"
	"github.com/exven/pos-system/modules/outlets/services"
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
	m.container.RegisterSingleton("outlets.outletRepository", func() interface{} {
		return persistence.NewOutletRepository(m.db)
	})

	// Register services
	m.container.RegisterSingleton("outlets.outletService", func() interface{} {
		repo := persistence.NewOutletRepository(m.db)
		return services.NewOutletService(repo)
	})

	// Register handlers
	m.container.RegisterSingleton("outlets.handler", func() interface{} {
		repo := persistence.NewOutletRepository(m.db)
		service := services.NewOutletService(repo)
		return handlers.NewOutletHandler(service)
	})
}

func (m *Module) GetHandler() *handlers.OutletHandler {
	repo := persistence.NewOutletRepository(m.db)
	service := services.NewOutletService(repo)
	return handlers.NewOutletHandler(service)
}