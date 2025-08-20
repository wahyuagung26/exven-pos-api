package subscription_plans

import (
	"github.com/exven/pos-system/modules/subscription_plans/handlers"
	"github.com/exven/pos-system/modules/subscription_plans/persistence"
	"github.com/exven/pos-system/modules/subscription_plans/services"
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
	m.container.RegisterSingleton("subscription_plans.repository", func() interface{} {
		return persistence.NewSubscriptionPlanRepository(m.db)
	})

	// Register services
	m.container.RegisterSingleton("subscription_plans.service", func() interface{} {
		repo := persistence.NewSubscriptionPlanRepository(m.db)
		return services.NewSubscriptionPlanService(repo)
	})

	// Register handlers
	m.container.RegisterSingleton("subscription_plans.handler", func() interface{} {
		repo := persistence.NewSubscriptionPlanRepository(m.db)
		service := services.NewSubscriptionPlanService(repo)
		return handlers.NewSubscriptionPlanHandler(service)
	})
}

func (m *Module) GetHandler() *handlers.SubscriptionPlanHandler {
	repo := persistence.NewSubscriptionPlanRepository(m.db)
	service := services.NewSubscriptionPlanService(repo)
	return handlers.NewSubscriptionPlanHandler(service)
}