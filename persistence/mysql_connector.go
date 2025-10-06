package persistence

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// MySQLConnector implements DBConnector for MySQL database
type MySQLConnector struct {
	*BaseConnector
	db *gorm.DB
}

// NewMySQLConnector creates a new MySQL connector with the given config
func NewMySQLConnector(config *DBConfig) *MySQLConnector {
	return &MySQLConnector{
		BaseConnector: NewBaseConnector(config),
	}
}

// Connect establishes a connection to the MySQL database
func (mc *MySQLConnector) Connect() (*gorm.DB, error) {

	var dsn string
	if mc.config.Password != "" {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			mc.config.Username,
			mc.config.Password,
			mc.config.Host,
			mc.config.Port,
			mc.config.Database,
		)
	} else {
		dsn = fmt.Sprintf("%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			mc.config.Username,
			mc.config.Host,
			mc.config.Port,
			mc.config.Database,
		)
	}

	// Add custom parameters if any
	if len(mc.config.Params) > 0 {
		dsn += "&"
		for key, value := range mc.config.Params {
			dsn += fmt.Sprintf("%s=%s&", key, value)
		}
		// Remove the trailing &
		dsn = dsn[:len(dsn)-1]
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL database: %w", err)
	}
	mc.db = db

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(mc.config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(mc.config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(mc.config.ConnMaxLifetime)

	return db, nil
}

// Close closes the database connection
func (mc *MySQLConnector) Close() error {
	if mc.db != nil {
		sqlDB, err := mc.db.DB()
		if err != nil {
			return fmt.Errorf("failed to get underlying *sql.DB: %v", err)
		}
		return sqlDB.Close()
	}
	return nil
}
