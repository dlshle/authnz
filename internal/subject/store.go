package subject

import (
	"github.com/dlshle/authnz/pkg/store"
	pb "github.com/dlshle/authnz/proto"
	"github.com/dlshle/gommon/utils"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"google.golang.org/protobuf/proto"
)

type Store interface {
	Get(id string) (*pb.Subject, error)
	Delete(id string) error
	Put(subject *pb.Subject) (*pb.Subject, error)
}

type SQLSubjectStore struct {
	pbEntityStore store.PBEntityStore
}

func NewSQLStore(db *sqlx.DB) Store {
	return &SQLSubjectStore{pbEntityStore: store.NewSQLPBEntityStore(db, "subjects")}
}

func (s *SQLSubjectStore) Get(id string) (*pb.Subject, error) {
	pbEntity, err := s.pbEntityStore.Get(id)
	if err != nil {
		return nil, err
	}
	subject := &pb.Subject{}
	err = proto.Unmarshal(pbEntity.Payload, subject)
	return subject, err
}

func (s *SQLSubjectStore) Put(subject *pb.Subject) (ret *pb.Subject, err error) {
	var (
		payload  []byte
		pbEntity *store.PBEntity
	)
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
		payload, err = proto.Marshal(subject)
		return err
	}, func() error {
		pbEntity, err = s.pbEntityStore.Put(&store.PBEntity{ID: subject.Id, Payload: payload})
		return err
	}, func() error {
		subject.Id = pbEntity.ID
		return nil
	})
	return subject, err
}

func (s *SQLSubjectStore) Delete(id string) error {
	return s.pbEntityStore.Delete(id)
}
