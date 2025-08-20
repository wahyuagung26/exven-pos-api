package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/exven/pos-system/internal/config"
	"github.com/exven/pos-system/modules/auth/domain"
	"github.com/exven/pos-system/modules/auth/handlers"
	"github.com/exven/pos-system/modules/outlets"
	"github.com/exven/pos-system/modules/products"
	"github.com/exven/pos-system/shared/container"
	"github.com/exven/pos-system/shared/middleware"
	"github.com/exven/pos-system/shared/validator"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

type Server struct {
	echo      *echo.Echo
	config    *config.Config
	container container.Container
}

func New(cfg *config.Config, container container.Container) *Server {
	e := echo.New()

	e.HideBanner = true
	e.HidePort = true

	// Set custom validator
	e.Validator = validator.New()

	// Enable detailed request logging
	e.Use(echoMiddleware.LoggerWithConfig(echoMiddleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}, error=${error}, latency=${latency_human}\n",
	}))

	// Enable panic recovery with detailed stack traces
	e.Use(echoMiddleware.RecoverWithConfig(echoMiddleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
		LogLevel:  1,       // DEBUG level
	}))

	e.Use(echoMiddleware.RequestID())

	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins:     cfg.CORS.AllowedOrigins,
		AllowMethods:     cfg.CORS.AllowedMethods,
		AllowHeaders:     cfg.CORS.AllowedHeaders,
		AllowCredentials: true,
	}))

	// Rate limiting disabled for now - can be enabled with external tools like nginx
	// e.Use(echoMiddleware.RateLimiterWithConfig(...))

	e.Use(echoMiddleware.BodyLimit(fmt.Sprintf("%d", cfg.FileUpload.MaxSize)))

	e.Use(echoMiddleware.Secure())

	return &Server{
		echo:      e,
		config:    cfg,
		container: container,
	}
}

func (s *Server) Start(address string) error {
	s.registerRoutes()
	return s.echo.Start(address)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.echo.Shutdown(ctx)
}

func (s *Server) registerRoutes() {
	s.echo.GET("/health", s.healthCheck)

	api := s.echo.Group("/api/v1")

	authService := s.container.MustGet("auth.service").(domain.AuthService)
	authHandler := handlers.NewAuthHandler(authService)
	authHandler.RegisterRoutes(api)

	protected := api.Group("")
	protected.Use(middleware.JWTAuth(s.config.JWT.Secret))
	protected.Use(middleware.TenantContext())

	// Get the products module and register its routes
	db := s.container.MustGet("db").(*gorm.DB)
	productsModule := products.NewModule(s.container, db, nil)
	productHandler := productsModule.GetHandler()
	productHandler.RegisterRoutes(protected)

	// Get the outlets module and register its routes
	outletsModule := outlets.NewModule(s.container, db, nil)
	outletHandler := outletsModule.GetHandler()
	outletHandler.RegisterRoutes(protected)

}

func (s *Server) healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "healthy",
		"service": s.config.App.Name,
		"env":     s.config.App.Env,
	})
}
