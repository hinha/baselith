package baselith

import (
	"fmt"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func loadMigrationsXML(xmlPath string) (*xmlMigrations, string, error) {
	doc, baseDir, err := parseXML(xmlPath)
	if err != nil {
		return nil, "", err
	}

	return doc, baseDir, nil
}

func readMigrationsXML(doc *xmlMigrations, baseDir string) (tx, noTx []*gormigrate.Migration, metasTx, metasNoTx map[string]Meta, err error) {
	metasTx = make(map[string]Meta)
	metasNoTx = make(map[string]Meta)

	// read and validate each migration
	for _, m := range doc.Items {
		if len(m.Author) == 0 {
			return nil, nil, nil, nil, fmt.Errorf("%s: missing <author>", m.ID)
		}
		if len(m.Labels) == 0 {
			return nil, nil, nil, nil, fmt.Errorf("%s: missing <labels>", m.ID)
		}
		if m.Kind == "" {
			return nil, nil, nil, nil, fmt.Errorf("%s: missing <kind>", m.ID)
		}
		if m.ID == "" {
			return nil, nil, nil, nil, fmt.Errorf("%s: missing <id>", m.ID)
		}

		useTx := true
		if m.Transactional != nil {
			useTx = *m.Transactional
		}

		var upFn, downFn func(*gorm.DB) error
		switch m.Kind {
		case "sql":
			if m.IncludeUp == nil {
				return nil, nil, nil, nil, fmt.Errorf("%s: missing <include> up file", m.ID)
			}
			upSQL, err := readSQL(baseDir, m.IncludeUp.File, m.IncludeUp.Rel == "true")
			if err != nil {
				return nil, nil, nil, nil, fmt.Errorf("%s up: %w", m.ID, err)
			}

			var downSQL string
			if m.IncludeDown != nil {
				downSQL, err = readSQL(baseDir, m.IncludeDown.File, m.IncludeDown.Rel == "true")
				if err != nil {
					return nil, nil, nil, nil, fmt.Errorf("%s down: %w", m.ID, err)
				}
			}

			upFn = func(tx *gorm.DB) error { return tx.Exec(upSQL).Error }
			downFn = func(tx *gorm.DB) error {
				if downSQL == "" {
					return fmt.Errorf("no down SQL for %s", m.ID)
				}
				return tx.Exec(downSQL).Error
			}

		default:
			return nil, nil, nil, nil, fmt.Errorf("%s: unsupported type=%s", m.ID, m.Kind)
		}

		meta := Meta{
			Author:        m.Author,
			Labels:        m.Labels,
			Kind:          m.Kind,
			Transactional: useTx,
		}

		gm := &gormigrate.Migration{
			ID:       m.ID,
			Migrate:  upFn,
			Rollback: downFn,
		}

		if useTx {
			tx = append(tx, gm)
			metasTx[m.ID] = meta
		} else {
			noTx = append(noTx, gm)
			metasNoTx[m.ID] = meta
		}
	}
	return tx, noTx, metasTx, metasNoTx, nil
}
