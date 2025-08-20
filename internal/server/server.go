package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/exven/pos-system/internal/config"
	"github.com/exven/pos-system/modules/auth/domain"
	"github.com/exven/pos-system/modules/auth/handlers"
	"github.com/exven/pos-system/shared/container"
	"github.com/exven/pos-system/shared/middleware"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
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

	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
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

}

func (s *Server) healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "healthy",
		"service": s.config.App.Name,
		"env":     s.config.App.Env,
	})
}
