package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/exven/pos-system/internal/config"
	"github.com/exven/pos-system/internal/server"
	"github.com/exven/pos-system/modules/auth"
	"github.com/exven/pos-system/shared/container"
	"github.com/exven/pos-system/shared/infrastructure/cache"
	"github.com/exven/pos-system/shared/infrastructure/database"
	"github.com/exven/pos-system/shared/infrastructure/messaging"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db, err := database.Initialize(database.Config{
		Host:               cfg.Database.Host,
		Port:               cfg.Database.Port,
		User:               cfg.Database.User,
		Password:           cfg.Database.Password,
		DBName:             cfg.Database.DBName,
		SSLMode:            cfg.Database.SSLMode,
		MaxConnections:     cfg.Database.MaxConnections,
		MaxIdleConnections: cfg.Database.MaxIdleConnections,
		ConnMaxLifetime:    cfg.Database.ConnMaxLifetime,
	})
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	redisClient, err := cache.NewRedisClient(cache.Config{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		PoolSize: cfg.Redis.PoolSize,
	})
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}
	defer redisClient.Close()

	// RabbitMQ temporarily disabled
	// eventBus, err := messaging.NewRabbitMQEventBus(
	//		cfg.RabbitMQ.URL,
	//		cfg.RabbitMQ.Exchange,
	//		cfg.RabbitMQ.QueuePrefix,
	// )
	// if err != nil {
	//		log.Fatalf("Failed to initialize RabbitMQ: %v", err)
	// }
	// defer eventBus.Close()
	var eventBus messaging.EventBus = nil
	log.Println("RabbitMQ connection disabled")

	di := container.New()

	registerSharedServices(di, cfg, db, redisClient, eventBus)

	authModule := auth.NewModule(di, db, redisClient, eventBus, cfg.JWT)
	authModule.Register()

	srv := server.New(cfg, di)
	log.Println("Server instance created successfully")
	log.Println("Auth module registered successfully")

	go func() {
		log.Printf("Starting server on port %d", cfg.App.Port)
		log.Printf("Server binding to address: :%d", cfg.App.Port)
		log.Printf("Environment: %s", cfg.App.Env)
		log.Printf("CORS Origins: %v", cfg.CORS.AllowedOrigins)
		
		if err := srv.Start(fmt.Sprintf(":%d", cfg.App.Port)); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Gracefully shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server shutdown complete")
}

func registerSharedServices(di *container.DIContainer, cfg *config.Config, db *gorm.DB, redisClient *cache.RedisClient, eventBus messaging.EventBus) {
	di.RegisterSingleton("config", func() *config.Config {
		return cfg
	})

	di.RegisterSingleton("db", func() *gorm.DB {
		return db
	})

	di.RegisterSingleton("redis", func() *cache.RedisClient {
		return redisClient
	})

	di.RegisterSingleton("eventBus", func() messaging.EventBus {
		return eventBus
	})
}