package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type MySQLClient struct {
	DB *gorm.DB
}

// NewMySQLClient creates a MySQL client with retry logic.
func NewMySQLClient(dsn string) (*MySQLClient, error) {
	var db *gorm.DB
	var err error

	maxRetries := 10
	retryInterval := 2 * time.Second

	for i := 1; i <= maxRetries; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
			Logger: logger.New(
				log.New(os.Stdout, "\r\n", log.LstdFlags),
				logger.Config{
					SlowThreshold: time.Millisecond * 200,
					LogLevel:      logger.Warn,
					Colorful:      true,
				},
			),
		})

		if err == nil {
			sqlDB, pingErr := db.DB()
			if pingErr == nil && sqlDB.Ping() == nil {
				log.Println("MySQL connected successfully")
				break
			}
			err = fmt.Errorf("ping failed: %v", pingErr)
		}

		log.Printf("MySQL not ready (attempt %d/%d): %v\n", i, maxRetries, err)

		if i < maxRetries {
			time.Sleep(retryInterval)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL after %d attempts: %w", maxRetries, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(30 * time.Minute)

	return &MySQLClient{DB: db}, nil
}
