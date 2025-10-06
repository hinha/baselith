package persistence

import (
	"fmt"

	"gorm.io/gorm"
)

// ConnectorFactory provides a factory method to create database connectors based on the driver
type ConnectorFactory struct{}

// NewConnectorFactory creates a new instance of ConnectorFactory
func NewConnectorFactory() *ConnectorFactory {
	return &ConnectorFactory{}
}

// CreateConnector creates a new database connector based on the driver
func (cf *ConnectorFactory) CreateConnector(config *DBConfig) (DBConnector, error) {
	switch config.Driver {
	case "mysql":
		return NewMySQLConnector(config), nil
	case "postgres", "postgresql":
		return NewPostgresConnector(config), nil
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", config.Driver)
	}
}

// Connect creates and returns a GORM database connection based on the provided configuration
func (cf *ConnectorFactory) Connect(config *DBConfig) (*gorm.DB, error) {
	connector, err := cf.CreateConnector(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connector: %w", err)
	}

	return connector.Connect()
}
