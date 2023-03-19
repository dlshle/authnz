package group

import (
	"github.com/dlshle/authnz/pkg/store"
	pb "github.com/dlshle/authnz/proto"
	"github.com/dlshle/gommon/utils"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"google.golang.org/protobuf/proto"
)

type Store interface {
	Get(id string) (*pb.Group, error)
	Delete(id string) error
	Put(group *pb.Group) (*pb.Group, error)
}

type SQLGroupStore struct {
	PbEntityStore store.PBEntityStore
}

func NewSQLStore(db *sqlx.DB) Store {
	return &SQLGroupStore{PbEntityStore: store.NewSQLPBEntityStore(db, "groups")}
}

func (s *SQLGroupStore) Get(id string) (*pb.Group, error) {
	pbEntity, err := s.PbEntityStore.Get(id)
	if err != nil {
		return nil, err
	}
	group := &pb.Group{}
	err = proto.Unmarshal(pbEntity.Payload, group)
	return group, err
}

func (s *SQLGroupStore) Put(group *pb.Group) (ret *pb.Group, err error) {
	var (
		payload  []byte
		pbEntity *store.PBEntity
	)
	err = utils.ProcessWithErrors(func() error {
		if group.Id == "" {
			newID, err := uuid.NewV4()
			if err != nil {
				return err
			}
			group.Id = newID.String()
		}
		return nil
	}, func() error {
		payload, err = proto.Marshal(group)
		return err
	}, func() error {
		pbEntity, err = s.PbEntityStore.Put(&store.PBEntity{ID: group.Id, Payload: payload})
		return err
	}, func() error {
		group.Id = pbEntity.ID
		return nil
	})
	return group, err
}

func (s *SQLGroupStore) Delete(id string) error {
	return s.PbEntityStore.Delete(id)
}
