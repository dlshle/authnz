package store

import (
	"context"
	"database/sql"

	"github.com/dlshle/gommon/errors"
	"github.com/dlshle/gommon/logging"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

type SQLPBEntityStore struct {
	Db        *sqlx.DB
	tableName string
}

func Open(connectionURL string, tableName string) (PBEntityStore, error) {
	db, err := sqlx.Open("postgres", connectionURL)
	return &SQLPBEntityStore{
		Db:        db,
		tableName: tableName,
	}, err
}

func NewSQLPBEntityStore(db *sqlx.DB, tableName string) PBEntityStore {
	return &SQLPBEntityStore{Db: db, tableName: tableName}
}

func (s *SQLPBEntityStore) Get(id string) (*PBEntity, error) {
	entities := []PBEntity{}
	err := s.Db.Select(&entities, "SELECT * FROM "+s.tableName+" WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	if len(entities) == 0 {
		return nil, errors.Error("no record found for " + id)
	}
	return &entities[0], err
}

func (s *SQLPBEntityStore) Delete(id string) error {
	res, err := s.Db.Exec("DELETE FROM "+s.tableName+" WHERE id = $1", id)
	if err != nil {
		return err
	}
	return CheckErrorForRowsAffected(res, id+" is not found")
}

func (s *SQLPBEntityStore) Put(entity *PBEntity) (*PBEntity, error) {
	var (
		result sql.Result
		err    error
	)
	if entity.ID == "" {
		// create
		var newID uuid.UUID
		newID, err = uuid.NewV4()
		if err != nil {
			return nil, err
		}
		entity.ID = newID.String()
		result, err = s.execUpsert(entity)
	} else {
		result, err = s.execUpsert(entity)
	}
	if err != nil {
		logging.GlobalLogger.Infof(context.Background(), "error: %v", err)
		return entity, err
	}
	err = CheckErrorForRowsAffected(result, "no row is affected")
	return entity, err
}

func (s *SQLPBEntityStore) execUpsert(entity *PBEntity) (sql.Result, error) {
	return s.Db.Exec("INSERT INTO "+s.tableName+" (id, payload) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET id = $1, payload = $2", entity.ID, entity.Payload)
}

func (s *SQLPBEntityStore) execUpdate(entity *PBEntity) (sql.Result, error) {
	return s.Db.Exec("UPDATE "+s.tableName+" SET payload = $2 WHERE id = $1", entity.ID, entity.Payload)
}

func CheckErrorForRowsAffected(result sql.Result, onNoRowAffectedMsg string) error {
	if result == nil {
		return errors.Error("nil result")
	}
	if rowsAffected, err := result.RowsAffected(); err != nil || rowsAffected == 0 {
		if err != nil {
			return err
		}
		return errors.Error(onNoRowAffectedMsg)
	}
	return nil
}
