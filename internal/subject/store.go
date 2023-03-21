package subject

import (
	"context"

	"github.com/dlshle/authnz/pkg/store"
	pb "github.com/dlshle/authnz/proto"
	"github.com/dlshle/gommon/logging"
	"github.com/dlshle/gommon/utils"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

type Store interface {
	Get(id string) (*pb.Subject, error)
	TxGet(tx store.SQLTransactional, id string) (*pb.Subject, error)
	TxBulkGet(tx store.SQLTransactional, ids []string) ([]*pb.Subject, error)
	Delete(id string) error
	TxDelete(tx store.SQLTransactional, id string) error
	Put(subject *pb.Subject) (*pb.Subject, error)
	TxPut(tx store.SQLTransactional, subject *pb.Subject) (ret *pb.Subject, err error)
	WithTX(cb func(store.SQLTransactional) error) error
}

type SQLSubjectStore struct {
	db *sqlx.DB
}

func NewSQLStore(db *sqlx.DB) Store {
	return &SQLSubjectStore{db: db}
}

func (s *SQLSubjectStore) Get(id string) (*pb.Subject, error) {
	return s.TxGet(s.db, id)
}

func (s *SQLSubjectStore) TxGet(tx store.SQLTransactional, id string) (*pb.Subject, error) {
	subject := &Subject{}
	err := tx.Select(subject, "SELECT * FROM subjects WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &pb.Subject{Id: subject.ID, UserId: subject.UserID}, nil
}

func (s *SQLSubjectStore) TxBulkGet(tx store.SQLTransactional, ids []string) ([]*pb.Subject, error) {
	subjects := []Subject{}
	if len(ids) == 0 {
		return []*pb.Subject{}, nil
	}
	sql := "SELECT * FROM subjects WHERE id in " + store.MakeInQueryClause(ids)
	logging.GlobalLogger.Infof(context.Background(), "tx bulk get query: %s", sql)
	err := tx.Select(&subjects, sql)
	if err != nil {
		return nil, err
	}
	pbSubjects := make([]*pb.Subject, len(subjects), len(subjects))
	for i, subject := range subjects {
		pbSubjects[i] = &pb.Subject{Id: subject.ID, UserId: subject.ID}
	}
	return pbSubjects, nil
}

func (s *SQLSubjectStore) Put(subject *pb.Subject) (ret *pb.Subject, err error) {
	return s.TxPut(s.db, subject)
}

func (s *SQLSubjectStore) TxPut(tx store.SQLTransactional, subject *pb.Subject) (ret *pb.Subject, err error) {
	err = utils.ProcessWithErrors(func() error {
		if subject.Id == "" {
			newID, err := uuid.NewV4()
			if err != nil {
				return err
			}
			subject.Id = newID.String()
		}
		return nil
	}, func() error {
		res, err := tx.Exec("INSERT INTO subjects (id, user_id) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET user_id = $2", subject.Id, subject.UserId)
		if err != nil {
			return err
		}
		return store.CheckErrorForRowsAffected(res, "subject "+subject.Id+":"+subject.UserId+" is not inserted")
	})
	return subject, err
}

func (s *SQLSubjectStore) Delete(id string) error {
	return s.TxDelete(s.db, id)
}

func (s *SQLSubjectStore) TxDelete(tx store.SQLTransactional, id string) error {
	res, err := tx.Exec("DELETE FROM subjects WHERE id = $1", id)
	if err != nil {
		return err
	}
	return store.CheckErrorForRowsAffected(res, "subject not found for id "+id)
}

func (s *SQLSubjectStore) WithTX(cb func(store.SQLTransactional) error) error {
	return store.WithSQLXTx(s.db, cb)
}
