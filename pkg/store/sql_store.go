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
	return s.TxGet(s.Db, id)
}

func (s *SQLPBEntityStore) TxGet(tx SQLTransactional, id string) (*PBEntity, error) {
	entities := []PBEntity{}
	err := tx.Select(&entities, "SELECT * FROM "+s.tableName+" WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	if len(entities) == 0 {
		return nil, errors.Error("no record found for " + id)
	}
	return &entities[0], err
}

func (s *SQLPBEntityStore) TxBulkGet(tx SQLTransactional, ids []string) ([]*PBEntity, error) {
	entities := []PBEntity{}
	if len(ids) == 0 {
		return []*PBEntity{}, nil
	}
	sql := "SELECT * FROM " + s.tableName + " WHERE id IN " + MakeInQueryClause(ids)
	logging.GlobalLogger.Infof(context.Background(), "query sql: %s", sql)
	err := tx.Select(&entities, sql)
	if err != nil {
		return nil, err
	}
	pbEntities := make([]*PBEntity, len(entities), len(entities))
	for i, entity := range entities {
		pbEntities[i] = &entity
	}
	return pbEntities, err
}

func (s *SQLPBEntityStore) Delete(id string) error {
	return s.TxDelete(s.Db, id)
}

func (s *SQLPBEntityStore) TxDelete(tx SQLTransactional, id string) error {
	res, err := tx.Exec("DELETE FROM "+s.tableName+" WHERE id = $1", id)
	if err != nil {
		return err
	}
	return CheckErrorForRowsAffected(res, id+" is not found")
}

func (s *SQLPBEntityStore) Put(entity *PBEntity) (*PBEntity, error) {
	return s.TxPut(s.Db, entity)
}

func (s *SQLPBEntityStore) TxPut(tx SQLTransactional, entity *PBEntity) (*PBEntity, error) {
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
		result, err = s.execUpsert(tx, entity)
	} else {
		result, err = s.execUpsert(tx, entity)
	}
	if err != nil {
		logging.GlobalLogger.Infof(context.Background(), "error: %v", err)
		return entity, err
	}
	err = CheckErrorForRowsAffected(result, "no row is affected")
	return entity, err

}

func (s *SQLPBEntityStore) WithTx(cb func(SQLTransactional) error) error {
	return WithSQLXTx(s.Db, cb)
}

func (s *SQLPBEntityStore) execUpsert(tx SQLTransactional, entity *PBEntity) (sql.Result, error) {
	return tx.Exec("INSERT INTO "+s.tableName+" (id, payload) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET id = $1, payload = $2", entity.ID, entity.Payload)
}

func (s *SQLPBEntityStore) execUpdate(tx SQLTransactional, entity *PBEntity) (sql.Result, error) {
	return tx.Exec("UPDATE "+s.tableName+" SET payload = $2 WHERE id = $1", entity.ID, entity.Payload)
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
