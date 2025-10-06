package example

import (
	"log"
	"time"

	"github.com/hinha/baselith/persistence"
)

// Example usage of the persistence layer
func ExampleUsage() {
	// Create a MySQL configuration using the builder pattern
	mysqlConfig, err := persistence.NewDBConfigBuilder().
		Driver("mysql").
		Host("localhost").
		Port(3306).
		Database("myapp").
		Username("user").
		Password("password").
		MaxIdleConns(10).
		MaxOpenConns(100).
		ConnMaxLifetime(time.Hour).
		Build()
	if err != nil {
		log.Fatal("Failed to build MySQL config:", err)
	}

	// Create a PostgreSQL configuration using the builder pattern
	postgresConfig, err := persistence.NewDBConfigBuilder().
		Driver("postgres").
		Host("localhost").
		Port(5432).
		Database("myapp").
		Username("user").
		Password("password").
		SSLMode("disable").
		MaxIdleConns(10).
		MaxOpenConns(100).
		ConnMaxLifetime(time.Hour).
		AddParam("TimeZone", "UTC").
		Build()
	if err != nil {
		log.Fatal("Failed to build PostgreSQL config:", err)
	}

	// Create connector factory
	factory := persistence.NewConnectorFactory()

	// Connect to MySQL
	_, err = factory.Connect(mysqlConfig)
	if err != nil {
		log.Printf("Failed to connect to MySQL: %v", err)
	} else {
		log.Println("Successfully connected to MySQL")
	}

	// Connect to PostgreSQL
	_, err = factory.Connect(postgresConfig)
	if err != nil {
		log.Printf("Failed to connect to PostgreSQL: %v", err)
	} else {
		log.Println("Successfully connected to PostgreSQL")
	}

	// Alternatively, create specific connectors
	mysqlConnector := persistence.NewMySQLConnector(mysqlConfig)
	db1, err := mysqlConnector.Connect()
	if err != nil {
		log.Printf("Failed to connect to MySQL with specific connector: %v", err)
	} else {
		log.Println("Successfully connected to MySQL using specific connector")
		_ = db1 // Use the variable to avoid "declared and not used" error
	}

	postgresConnector := persistence.NewPostgresConnector(postgresConfig)
	db2, err := postgresConnector.Connect()
	if err != nil {
		log.Printf("Failed to connect to PostgreSQL with specific connector: %v", err)
	} else {
		log.Println("Successfully connected to PostgreSQL using specific connector")
		_ = db2 // Use the variable to avoid "declared and not used" error
	}
}
