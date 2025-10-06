package baselith

// DBInterface adalah interface minimal yang dibutuhkan untuk testing
type DBInterface interface {
	Exec(query string, values ...interface{}) DBResult
}

// DBResult adalah interface untuk hasil eksekusi SQL
type DBResult interface {
	Error() error
	RowsAffected() int64
}
