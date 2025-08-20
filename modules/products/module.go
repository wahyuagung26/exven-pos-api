package products

import (
	"github.com/exven/pos-system/modules/products/handlers"
	"github.com/exven/pos-system/modules/products/persistence"
	"github.com/exven/pos-system/modules/products/services"
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
	m.container.RegisterSingleton("products.categoryRepository", func() interface{} {
		return persistence.NewProductCategoryRepository(m.db)
	})

	// Register services
	m.container.RegisterSingleton("products.categoryService", func() interface{} {
		repo := persistence.NewProductCategoryRepository(m.db)
		return services.NewProductCategoryService(repo)
	})

	// Register handlers
	m.container.RegisterSingleton("products.handler", func() interface{} {
		categoryRepo := persistence.NewProductCategoryRepository(m.db)
		categoryService := services.NewProductCategoryService(categoryRepo)
		return handlers.NewProductHandler(categoryService)
	})
}

func (m *Module) GetHandler() *handlers.ProductHandler {
	categoryRepo := persistence.NewProductCategoryRepository(m.db)
	categoryService := services.NewProductCategoryService(categoryRepo)
	return handlers.NewProductHandler(categoryService)
}