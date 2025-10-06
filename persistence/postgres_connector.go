package persistence

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// PostgresConnector implements DBConnector for PostgreSQL database
type PostgresConnector struct {
	*BaseConnector
}

// NewPostgresConnector creates a new PostgreSQL connector with the given config
func NewPostgresConnector(config *DBConfig) *PostgresConnector {
	return &PostgresConnector{
		BaseConnector: NewBaseConnector(config),
	}
}

// Connect establishes a connection to the PostgreSQL database
func (pc *PostgresConnector) Connect() (*gorm.DB, error) {
	// Build connection string
	var dsn string
	// no required password
	dsn = fmt.Sprintf("host=%s user=%s dbname=%s port=%d sslmode=%s",
		pc.config.Host,
		pc.config.Username,
		pc.config.Database,
		pc.config.Port,
		pc.config.SSLMode,
	)
	if pc.config.Password != "" {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
			pc.config.Host,
			pc.config.Username,
			pc.config.Password,
			pc.config.Database,
			pc.config.Port,
			pc.config.SSLMode,
		)
	}

	// Add custom parameters if any
	for key, value := range pc.config.Params {
		dsn += fmt.Sprintf(" %s=%s", key, value)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(pc.config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(pc.config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(pc.config.ConnMaxLifetime)

	return db, nil
}
