package baselith

import (
	"fmt"
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
