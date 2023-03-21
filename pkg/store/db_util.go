package store

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

func WithSQLXTx(db *sqlx.DB, cb func(SQLTransactional) error) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	err = cb(tx)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func MakeInQueryClause(elements []string) string {
	var clauseBuilder strings.Builder
	clauseBuilder.WriteRune('(')
	l := len(elements)
	for i, ele := range elements {
		clauseBuilder.WriteRune('\'')
		clauseBuilder.WriteString(ele)
		clauseBuilder.WriteRune('\'')
		if i < l-1 {
			clauseBuilder.WriteRune(',')
		}
	}
	clauseBuilder.WriteRune(')')
	return clauseBuilder.String()
}
