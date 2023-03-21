package migration

import "github.com/jmoiron/sqlx"

var v1 = `
CREATE TABLE IF NOT EXISTS subjects (
    id uuid,
	user_id varchar(255),
	PRIMARY KEY ( id )
);

CREATE TABLE IF NOT EXISTS groups (
	id uuid,
	payload bytea,
	PRIMARY KEY ( id )
);

CREATE TABLE IF NOT EXISTS policies (
	id uuid,
	payload bytea,
	PRIMARY KEY ( id )
);

CREATE TABLE IF NOT EXISTS contracts (
	id uuid,
	subject_id uuid,
	group_id uuid,
	PRIMARY KEY ( id )
);
`

var migration_scripts = []string{v1}

func ExecMigration(db *sqlx.DB) error {
	for _, migration := range migration_scripts {
		_, err := db.Exec(migration)
		if err != nil {
			return err
		}
	}
	return nil
}
