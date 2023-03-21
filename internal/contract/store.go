package contract

import (
	"github.com/dlshle/authnz/pkg/store"
	pb "github.com/dlshle/authnz/proto"
	"github.com/dlshle/gommon/errors"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"google.golang.org/protobuf/proto"
)

type Store interface {
	AddNewContract(subjectID, groupID string) (*pb.Contract, error)
	TxAddNewContract(tx store.SQLTransactional, subjectID, groupID string) (*pb.Contract, error)
	DeleteContractByContractID(contractID string) error
	DeleteContract(subjectID, groupID string) error
	TxDeleteContract(tx store.SQLTransactional, contractID string) error
	TxDeleteContractsByGroupID(tx store.SQLTransactional, groupID string) error
	ListAllContractsBySubject(subjectID string) ([]Contract, error)
	ListGroupsBySubjectID(subjectID string) ([]*pb.Group, error)
	TxListAllContractsBySubject(tx store.SQLTransactional, subjectID string) ([]Contract, error)
	TxListContractsByGroupID(tx store.SQLTransactional, groupID string) ([]Contract, error)
}

type contractStore struct {
	db *sqlx.DB
}

func NewContractStore(db *sqlx.DB) Store {
	return &contractStore{db: db}
}

func (s *contractStore) AddNewContract(subjectID, groupID string) (*pb.Contract, error) {
	return s.TxAddNewContract(s.db, subjectID, groupID)
}

func (s *contractStore) TxAddNewContract(tx store.SQLTransactional, subjectID, groupID string) (*pb.Contract, error) {
	contracts, err := s.ListAllContractsBySubject(subjectID)
	if err != nil {
		return nil, err
	}
	for _, c := range contracts {
		if c.GroupID == groupID {
			return nil, errors.Error("contract by " + subjectID + ":" + groupID + " already exists")
		}
	}
	contractID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	res, err := tx.Exec("INSERT INTO contracts (id, subject_id, group_id) VALUES ($1, $2, $3)", contractID.String(), subjectID, groupID)
	if err != nil {
		return nil, err
	}
	newContract := &pb.Contract{Id: contractID.String(), SubjectId: subjectID, GroupId: groupID}
	return newContract, store.CheckErrorForRowsAffected(res, "no record is inserted for "+subjectID+":"+groupID)
}

func (s *contractStore) DeleteContract(subjectID, groupID string) error {
	res, err := s.db.Exec("DELETE FROM contracts WHERE subject_id = $1 AND group_id = $2", subjectID, groupID)
	if err != nil {
		return err
	}
	return store.CheckErrorForRowsAffected(res, "no record is found for "+subjectID+":"+groupID)
}

func (s *contractStore) DeleteContractByContractID(contractID string) error {
	return s.TxDeleteContract(s.db, contractID)
}

func (s *contractStore) TxDeleteContractsByGroupID(tx store.SQLTransactional, groupID string) error {
	res, err := tx.Exec("DELETE FROM contracts WHERE group_id = $1", groupID)
	if err != nil {
		return err
	}
	return store.CheckErrorForRowsAffected(res, "not found: no group is found by group_id "+groupID)
}

func (s *contractStore) TxDeleteContract(tx store.SQLTransactional, contractID string) error {
	res, err := tx.Exec("DELETE FROM contracts WHERE id = $1", contractID)
	if err != nil {
		return err
	}
	return store.CheckErrorForRowsAffected(res, "no record is found for contractID "+contractID)
}

func (s *contractStore) ListGroupsBySubjectID(subjectID string) ([]*pb.Group, error) {
	pbEntities := []store.PBEntity{}
	err := s.db.Select(&pbEntities, "SELECT * FROM groups WHERE id IN (SELECT group_id FROM contracts WHERE subject_id = $1)", subjectID)
	if err != nil {
		return nil, err
	}
	groups := make([]*pb.Group, len(pbEntities), len(pbEntities))
	for i, pbEntity := range pbEntities {
		pbGroup := &pb.Group{}
		err = proto.Unmarshal(pbEntity.Payload, pbGroup)
		if err != nil {
			return nil, err
		}
		groups[i] = pbGroup
	}
	return groups, nil
}

func (s *contractStore) ListAllContractsBySubject(subjectID string) ([]Contract, error) {
	return s.TxListAllContractsBySubject(s.db, subjectID)
}

func (s *contractStore) TxListAllContractsBySubject(tx store.SQLTransactional, subjectID string) ([]Contract, error) {
	contracts := []Contract{}
	err := tx.Select(&contracts, "SELECT * FROM contracts WHERE subject_id = $1", subjectID)
	return contracts, err
}

func (s *contractStore) TxListContractsByGroupID(tx store.SQLTransactional, groupID string) ([]Contract, error) {
	contracts := []Contract{}
	err := tx.Select(&contracts, "SELECT * FROM contracts WHERE group_id = $1", groupID)
	return contracts, err
}
