package persistence

// BaseConnector provides common functionality for database connectors
type BaseConnector struct {
	config *DBConfig
}

// NewBaseConnector creates a new base connector with the given config
func NewBaseConnector(config *DBConfig) *BaseConnector {
	return &BaseConnector{
		config: config,
	}
}

// GetConfig returns the database configuration
func (bc *BaseConnector) GetConfig() *DBConfig {
	return bc.config
}
