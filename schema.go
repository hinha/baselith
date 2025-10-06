package baselith

import (
	"fmt"
	"log"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

const (
	sqlPostgresSchema = `SELECT * FROM %s.schema_migrations ORDER BY applied_at`
	sqlMysqlSchema    = `SELECT * FROM schema_migrations ORDER BY applied_at DESC`
)

type Meta struct {
	Author        string
	Labels        string
	Kind          string // "struct" | "sql" | "index" | etc
	Transactional bool
}

type xmlMigrations struct {
	Items []xmlChangelog `xml:"changeLog"`
}

type xmlChangelog struct {
	ID            string      `xml:"id,attr"`
	Kind          string      `xml:"kind,attr"` // "struct" | "sql"
	Author        string      `xml:"author,attr"`
	Labels        string      `xml:"labels,attr"`
	Transactional *bool       `xml:"transactional,attr"` // default: true
	Table         *xmlTable   `xml:"table"`
	IncludeUp     *xmlInclude `xml:"include"`
	IncludeDown   *xmlInclude `xml:"includeDown"`
}

type xmlTable struct {
	Name string `xml:"name,attr"` // e.g. "public.m_roles"
}

type xmlInclude struct {
	File string `xml:"file,attr"`
	Rel  string `xml:"relativeToChangelogFile,attr"` // "true"/"false"
}

type migRow struct {
	ID        string
	AppliedAt time.Time
}

func upsertMeta(db DBInterface, schema, id string, meta Meta) error {
	result := db.Exec(
		fmt.Sprintf(`UPDATE %s.schema_migrations
		             SET author = ?, labels = ?, kind = ?, transactional = ?
		             WHERE id = ?`, schema),
		meta.Author, meta.Labels, meta.Kind, meta.Transactional, id,
	)
	if err := result.Error(); err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("metadata not updated: id %q not found in %s.schema_migrations", id, schema)
	}
	return nil
}

// syncMetadata updates the metadata for multiple migrations in the schema_migrations table.
func syncMetadata(db DBInterface, schema string, metas map[string]Meta) error {
	for id, meta := range metas {
		if err := upsertMeta(db, schema, id, meta); err != nil {
			return fmt.Errorf("failed to sync metadata for %q: %w", id, err)
		}
	}
	return nil
}

func cmdHistory(db *gorm.DB) error {
	var rows []migRow
	if Driver == "postgres" || Driver == "postgresql" {
		if err := db.Raw(fmt.Sprintf(sqlPostgresSchema, Schema)).
			Scan(&rows).Error; err != nil {
			return err
		}
	} else if Driver == "mysql" {
		if err := db.Raw(sqlMysqlSchema).
			Scan(&rows).Error; err != nil {
			return err
		}
	}

	log.Println("== Migration History ==")
	for _, r := range rows {
		fmt.Printf("%s\t%s\n", r.AppliedAt.Format(time.RFC3339), r.ID)
	}
	return nil
}

func cmdStatus(db *gorm.DB, all []*gormigrate.Migration) error {
	var rows []migRow
	if Driver == "postgres" || Driver == "postgresql" {
		if err := db.Raw(fmt.Sprintf(sqlPostgresSchema, Schema)).
			Scan(&rows).Error; err != nil {
			return err
		}
	} else if Driver == "mysql" {
		if err := db.Raw(sqlMysqlSchema).
			Scan(&rows).Error; err != nil {
			return err
		}
	}
	applied := map[string]time.Time{}
	for _, r := range rows {
		applied[r.ID] = r.AppliedAt
	}

	log.Println("== Migration IsActive ==")
	for _, gm := range all {
		if t, ok := applied[gm.ID]; ok {
			log.Printf("✓ %s\t(%s)\n", gm.ID, t.Format(time.RFC3339))
		} else {
			log.Printf("• %s\t(PENDING)\n", gm.ID)
		}
	}

	// optional: info drift (ada di DB tapi tidak ada di XML)
	for id := range applied {
		found := false
		for _, gm := range all {
			if gm.ID == id {
				found = true
				break
			}
		}
		if !found {
			log.Printf("! drift: applied but missing in XML -> %s\n", id)
		}
	}
	return nil
}
