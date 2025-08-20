package customers

import (
	"github.com/exven/pos-system/modules/customers/handlers"
	"github.com/exven/pos-system/modules/customers/persistence"
	"github.com/exven/pos-system/modules/customers/services"
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
	m.container.RegisterSingleton("customers.customerRepository", func() interface{} {
		return persistence.NewCustomerRepository(m.db)
	})

	// Register services
	m.container.RegisterSingleton("customers.customerService", func() interface{} {
		repo := persistence.NewCustomerRepository(m.db)
		return services.NewCustomerService(repo)
	})

	// Register handlers
	m.container.RegisterSingleton("customers.handler", func() interface{} {
		repo := persistence.NewCustomerRepository(m.db)
		service := services.NewCustomerService(repo)
		return handlers.NewCustomerHandler(service)
	})
}

func (m *Module) GetHandler() *handlers.CustomerHandler {
	repo := persistence.NewCustomerRepository(m.db)
	service := services.NewCustomerService(repo)
	return handlers.NewCustomerHandler(service)
}