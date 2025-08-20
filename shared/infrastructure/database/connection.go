package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Host               string
	Port               int
	User               string
	Password           string
	DBName             string
	SSLMode            string
	MaxConnections     int
	MaxIdleConnections int
	ConnMaxLifetime    time.Duration
}

var db *gorm.DB

func Initialize(config Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		config.Host, config.User, config.Password, config.DBName, config.Port, config.SSLMode)

	gormConfig := &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Info),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	}

	var err error
	db, err = gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxOpenConns(config.MaxConnections)
	sqlDB.SetMaxIdleConns(config.MaxIdleConnections)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established successfully")
	return db, nil
}

func GetDB() *gorm.DB {
	if db == nil {
		panic("database not initialized")
	}
	return db
}

func Close() error {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

func WithTenant(tenantID uint64) *gorm.DB {
	return db.Where("tenant_id = ?", tenantID)
}

func Transaction(fn func(*gorm.DB) error) error {
	return db.Transaction(fn)
}
