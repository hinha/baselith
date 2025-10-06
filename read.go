package baselith

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func parseXML(path string) (*xmlMigrations, string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, "", err
	}
	var doc xmlMigrations
	if err := xml.Unmarshal(b, &doc); err != nil {
		return nil, "", err
	}
	SortChangelogsByID(doc.Items)
	base := filepath.Dir(path)
	return &doc, base, nil
}

func readSQL(base, file string, relative bool) (string, error) {
	p := file
	if relative {
		p = filepath.Join(base, file)
	}
	b, err := os.ReadFile(p)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// SortChangelogsByID sorts a slice of xmlChangelog by their numeric ID prefix (e.g., "000", "001", ...)
func SortChangelogsByID(changelogs []xmlChangelog) {
	sort.SliceStable(changelogs, func(i, j int) bool {
		getPrefix := func(id string) int {
			parts := strings.SplitN(id, "_", 2)
			if len(parts) > 0 {
				if n, err := strconv.Atoi(parts[0]); err == nil {
					return n
				}
			}
			return 0 // fallback if malformed
		}
		return getPrefix(changelogs[i].ID) < getPrefix(changelogs[j].ID)
	})
}
