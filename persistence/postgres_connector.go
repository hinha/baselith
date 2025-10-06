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
	db *gorm.DB
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
	if pc.config.Password != "" {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
			pc.config.Host,
			pc.config.Username,
			pc.config.Password,
			pc.config.Database,
			pc.config.Port,
			pc.config.SSLMode,
		)
	} else {
		// no required password
		dsn = fmt.Sprintf("host=%s user=%s dbname=%s port=%d sslmode=%s",
			pc.config.Host,
			pc.config.Username,
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
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %v", err)
	}

	pc.db = db
	return db, nil
}

// Close closes the database connection
func (pc *PostgresConnector) Close() error {
	if pc.db != nil {
		sqlDB, err := pc.db.DB()
		if err != nil {
			return fmt.Errorf("failed to get underlying *sql.DB: %v", err)
		}
		return sqlDB.Close()
	}
	return nil
}
