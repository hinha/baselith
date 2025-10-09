package baselith

import (
	"fmt"
	"log"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/hinha/baselith/persistence"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

func doAction(m *gormigrate.Gormigrate, sub, toID string) error {
	switch sub {
	case "down":
		if toID != "" {
			// RollbackTo: will revert ALL migrations above the target; the target ID itself will NOT be rolled back
			return m.RollbackTo(toID)
		}
		return m.RollbackLast()
	case "to":
		if toID == "" {
			return fmt.Errorf("--to <ID> required for 'to'")
		}
		return m.MigrateTo(toID)
	case "redo":
		// gormigrate does not have RedoLast(); implement it manually
		if err := m.RollbackLast(); err != nil {
			return err
		}
		return m.Migrate()
	default: // "up"
		return m.Migrate()
	}
}

func Run(cmd *cobra.Command, _ []string) {
	fmt.Fprintln(cmd.OutOrStdout(), "Welcome to Baselith! Use --help for more information.")
	switch {
	case ConfigYaml && ConfigPath == "":
		yamlCfg, err := ReadConfigYAML()
		if err != nil {
			log.Fatal("Failed to read YAML config:", err)
			return
		}
		Driver = yamlCfg.Driver
		Host = yamlCfg.Host
		Port = yamlCfg.Port
		Dbname = yamlCfg.Dbname
		User = yamlCfg.User
		Password = yamlCfg.Password
		Schema = yamlCfg.Schema
	}
	if Driver == "" || Host == "" || Port == 0 || Dbname == "" || User == "" {
		log.Fatal("Database connection parameters are required when not using a config file.")
	}

	doc, baseDir, err := loadMigrationsXML(Folder.JoinPath("migrations.xml"))
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Printf("Base directory: %s\n", baseDir)

	config, err := persistence.NewDBConfigBuilder().Driver(Driver).Host(Host).Port(Port).Database(Dbname).
		Username(User).
		Password(Password).
		MaxIdleConns(10).
		MaxOpenConns(100).
		Schema(Schema).
		ConnMaxLifetime(time.Hour).Build()
	if err != nil {
		log.Fatal("Failed to build config:", err)
		return
	}
	Schema = config.Schema

	factory := persistence.NewConnectorFactory()
	connect, err := factory.CreateConnector(config)
	if err != nil {
		log.Fatal("Failed to create connector:", err)
		return
	}

	db, err := connect.Connect()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		return
	}

	sql, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
		return
	}
	defer sql.Close()

	dbAdapter := NewDBAdapter(db)
	if err := migrationTable(config.Driver, dbAdapter); err != nil {
		log.Fatal("Failed to migration table:", err)
		return
	}

	// load migrations from XML
	txMigs, notxMigs, metasTx, metasNoTx, err := readMigrationsXML(doc, baseDir)
	if err != nil {
		log.Fatal(err)
		return
	}

	// command dispatcher
	switch Sub {
	case "status":
		if err := cmdStatus(db, append(txMigs, notxMigs...)); err != nil {
			log.Fatal(err)
		}
	case "history":
		if err := cmdHistory(db); err != nil {
			log.Fatal(err)
		}
	case "up", "down", "to", "redo":
		if err := runMutations(db, Sub, ToID, txMigs, notxMigs, metasTx, metasNoTx, Schema); err != nil {
			log.Fatal(err)
		}
	}
}

func migrationTable(driver string, db DBInterface) error {
	var createTableQuery string
	var alters []string

	switch driver {
	case "postgres", "postgresql":
		createTableQuery = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.schema_migrations (
	id varchar(255) PRIMARY KEY,
	applied_at timestamptz NOT NULL DEFAULT now(),
	author      varchar(128) NOT NULL DEFAULT 'unknown',
	labels      varchar(255) NOT NULL DEFAULT 'unknown',
	kind        varchar(32)
)`, Schema)
		alters = []string{
			fmt.Sprintf(`ALTER TABLE %s.schema_migrations ADD COLUMN IF NOT EXISTS applied_at timestamptz NOT NULL DEFAULT now()`, Schema),
			fmt.Sprintf(`ALTER TABLE %s.schema_migrations ADD COLUMN IF NOT EXISTS author varchar(128) NOT NULL DEFAULT 'unknown'`, Schema),
			fmt.Sprintf(`ALTER TABLE %s.schema_migrations ADD COLUMN IF NOT EXISTS labels varchar(255) NOT NULL DEFAULT 'unknown'`, Schema),
			fmt.Sprintf(`ALTER TABLE %s.schema_migrations ADD COLUMN IF NOT EXISTS kind varchar(32)`, Schema),
			fmt.Sprintf(`ALTER TABLE %s.schema_migrations ADD COLUMN IF NOT EXISTS transactional boolean NOT NULL DEFAULT true`, Schema),
		}
	case "mysql":
		createTableQuery = `CREATE TABLE IF NOT EXISTS schema_migrations (
	id varchar(255) PRIMARY KEY,
	applied_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
	author      varchar(128) NOT NULL DEFAULT 'unknown',
	labels      varchar(255) NOT NULL DEFAULT 'unknown',
	kind        varchar(32),
	transactional boolean NOT NULL DEFAULT true
)`
		alters = []string{
			`ALTER TABLE schema_migrations ADD COLUMN IF NOT EXISTS applied_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP`,
			`ALTER TABLE schema_migrations ADD COLUMN IF NOT EXISTS author varchar(128) NOT NULL DEFAULT 'unknown'`,
			`ALTER TABLE schema_migrations ADD COLUMN IF NOT EXISTS labels varchar(255) NOT NULL DEFAULT 'unknown'`,
			`ALTER TABLE schema_migrations ADD COLUMN IF NOT EXISTS kind varchar(32)`,
			`ALTER TABLE schema_migrations ADD COLUMN IF NOT EXISTS transactional boolean NOT NULL DEFAULT true`,
		}
	default:
		return fmt.Errorf("unsupported driver: %s", driver)
	}

	log.Printf("Creating migration table with query: %s", createTableQuery)
	result := db.Exec(createTableQuery)
	if err := result.Error(); err != nil {
		return fmt.Errorf("failed to create migration table: %w", err)
	}

	for _, q := range alters {
		log.Printf("Altering table with query: %s", q)
		result = db.Exec(q)
		if err := result.Error(); err != nil {
			return fmt.Errorf("failed to alter table: %w", err)
		}
	}
	return nil
}

func runMutations(
	db *gorm.DB,
	sub, toID string,
	txMigs, notxMigs []*gormigrate.Migration,
	metasTx, metasNoTx map[string]Meta,
	schema string,
) error {
	dbAdapter := NewDBAdapter(db)

	// LOCK session for batch transactional
	// Use database-specific locking mechanism
	var result DBResult
	if Driver == "postgres" || Driver == "postgresql" {
		exec := dbAdapter.Exec(`SELECT pg_advisory_lock( hashtext('gormigrate:xml:tx') )`)
		if err := exec.Error(); err != nil {
			return err
		}
		result = exec
		defer dbAdapter.Exec(`SELECT pg_advisory_unlock( hashtext('gormigrate:xml:tx') )`)
	} else if Driver == "mysql" {
		exec := dbAdapter.Exec(`SELECT GET_LOCK('gormigrate:xml:tx', 10)`)
		if err := exec.Error(); err != nil {
			return err
		}
		result = exec
		defer dbAdapter.Exec(`SELECT RELEASE_LOCK('gormigrate:xml:tx')`)
	}

	// Runner transactional
	mtx := gormigrate.New(db, &gormigrate.Options{
		TableName:      schema + ".schema_migrations",
		IDColumnName:   "id",
		IDColumnSize:   255,
		UseTransaction: true,
	}, txMigs)

	if err := doAction(mtx, sub, toID); err != nil {
		return err
	}

	if err := syncMetadata(dbAdapter, schema, metasTx); err != nil {
		return err
	}

	// NON-transactional batch use lock session a side for race condition
	if len(notxMigs) > 0 {
		result = dbAdapter.Exec(`SELECT pg_advisory_lock( hashtext('gormigrate:xml:notx') )`)
		if err := result.Error(); err != nil {
			return err
		}
		defer dbAdapter.Exec(`SELECT pg_advisory_unlock( hashtext('gormigrate:xml:notx') )`)

		mntx := gormigrate.New(db, &gormigrate.Options{
			TableName:      schema + ".schema_migrations",
			IDColumnName:   "id",
			IDColumnSize:   255,
			UseTransaction: false,
		}, notxMigs)

		if err := doAction(mntx, sub, toID); err != nil {
			return err
		}

		if err := syncMetadata(dbAdapter, schema, metasNoTx); err != nil {
			return err
		}
	}
	return nil
}
