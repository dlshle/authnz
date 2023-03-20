package store

import "database/sql"

type SQLTransactional interface {
	Select(dest interface{}, query string, args ...interface{}) error
	Exec(query string, args ...any) (sql.Result, error)
}
