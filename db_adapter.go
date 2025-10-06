package baselith

import "gorm.io/gorm"

// DBAdapter adapts gorm.DB to a simpler interface for testing
type DBAdapter struct {
	*gorm.DB
}

// NewDBAdapter creates a new adapter for gorm.DB
func NewDBAdapter(db *gorm.DB) *DBAdapter {
	return &DBAdapter{DB: db}
}

// dbResult is the internal implementation of DBResult
type dbResult struct {
	err          error
	rowsAffected int64
}

func (r *dbResult) Error() error {
	return r.err
}

func (r *dbResult) RowsAffected() int64 {
	return r.rowsAffected
}

// Exec executes SQL query
func (a *DBAdapter) Exec(query string, values ...interface{}) DBResult {
	result := a.DB.Exec(query, values...)
	return &dbResult{
		err:          result.Error,
		rowsAffected: result.RowsAffected,
	}
}
