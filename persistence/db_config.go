package persistence

import (
	"fmt"
	"time"
)

// DBConfig represents the configuration for database connections
type DBConfig struct {
	Driver          string
	Host            string
	Port            int
	Database        string
	Username        string
	Password        string
	SSLMode         string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	Params          map[string]string
	Schema          string // for postgres only
}

// DBConfigBuilder implements the builder pattern for DBConfig
type DBConfigBuilder struct {
	config DBConfig
}

// NewDBConfigBuilder creates a new instance of DBConfigBuilder
func NewDBConfigBuilder() *DBConfigBuilder {
	return &DBConfigBuilder{
		config: DBConfig{
			MaxIdleConns:    10,
			MaxOpenConns:    100,
			ConnMaxLifetime: time.Hour,
		},
	}
}

// Driver sets the database driver (mysql or postgres)
func (b *DBConfigBuilder) Driver(driver string) *DBConfigBuilder {
	b.config.Driver = driver
	return b
}

// Host sets the database host
func (b *DBConfigBuilder) Host(host string) *DBConfigBuilder {
	b.config.Host = host
	return b
}

// Port sets the database port
func (b *DBConfigBuilder) Port(port int) *DBConfigBuilder {
	b.config.Port = port
	return b
}

// Database sets the database name
func (b *DBConfigBuilder) Database(database string) *DBConfigBuilder {
	b.config.Database = database
	return b
}

// Username sets the database username
func (b *DBConfigBuilder) Username(username string) *DBConfigBuilder {
	b.config.Username = username
	return b
}

// Password sets the database password
func (b *DBConfigBuilder) Password(password string) *DBConfigBuilder {
	b.config.Password = password
	return b
}

// SSLMode sets the SSL mode for PostgreSQL
func (b *DBConfigBuilder) SSLMode(sslMode string) *DBConfigBuilder {
	b.config.SSLMode = sslMode
	return b
}

// MaxIdleConns sets the maximum number of idle connections
func (b *DBConfigBuilder) MaxIdleConns(maxIdleConns int) *DBConfigBuilder {
	b.config.MaxIdleConns = maxIdleConns
	return b
}

// MaxOpenConns sets the maximum number of open connections
func (b *DBConfigBuilder) MaxOpenConns(maxOpenConns int) *DBConfigBuilder {
	b.config.MaxOpenConns = maxOpenConns
	return b
}

// ConnMaxLifetime sets the maximum lifetime of a connection
func (b *DBConfigBuilder) ConnMaxLifetime(connMaxLifetime time.Duration) *DBConfigBuilder {
	b.config.ConnMaxLifetime = connMaxLifetime
	return b
}

// AddParam adds a custom parameter to the connection string
func (b *DBConfigBuilder) AddParam(key, value string) *DBConfigBuilder {
	if b.config.Params == nil {
		b.config.Params = make(map[string]string)
	}
	b.config.Params[key] = value
	return b
}

func (b *DBConfigBuilder) Schema(schema string) *DBConfigBuilder {
	b.config.Schema = schema
	return b
}

// Build returns the final DBConfig
func (b *DBConfigBuilder) Build() (*DBConfig, error) {
	if b.config.Driver == "" {
		return nil, fmt.Errorf("driver is required")
	}
	if b.config.Host == "" {
		return nil, fmt.Errorf("host is required")
	}
	if b.config.Port == 0 {
		return nil, fmt.Errorf("port is required")
	}
	if b.config.Database == "" {
		return nil, fmt.Errorf("database name is required")
	}
	if b.config.Schema == "" {
		b.config.Schema = "public"
	}

	return &b.config, nil
}
