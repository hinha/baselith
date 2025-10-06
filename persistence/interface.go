package persistence

import "gorm.io/gorm"

// DBConnector defines the interface for database connections
type DBConnector interface {
	Connect() (*gorm.DB, error)
	Close() error
	GetConfig() *DBConfig
}

// Factory defines the interface for creating database connectors
type Factory interface {
	CreateConnector(config *DBConfig) (DBConnector, error)
}
