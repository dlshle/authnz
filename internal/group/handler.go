package group

import (
	"context"
	"strings"

	"github.com/dlshle/authnz/internal/contract"
	"github.com/dlshle/authnz/pkg/store"
	pb "github.com/dlshle/authnz/proto"
)

type Handler struct {
	store         Store
	contractStore contract.Store
}

func NewHandler(store Store, contractStore contract.Store) *Handler {
	return &Handler{store: store, contractStore: contractStore}
}

func (h *Handler) CreateGroup(ctx context.Context, group *pb.Group) (*pb.GroupResponse, error) {
	newGroup, err := h.store.Put(group)
	if err != nil {
		return nil, err
	}
	return &pb.GroupResponse{Group: newGroup}, nil
}

func (h *Handler) UpdateGroup(ctx context.Context, group *pb.Group) (*pb.GroupResponse, error) {
	updatedGroup, err := h.store.Put(group)
	if err != nil {
		return nil, err
	}
	return &pb.GroupResponse{Group: updatedGroup}, nil
}

func (h *Handler) DuplicateGroup(ctx context.Context, groupID string) (*pb.GroupResponse, error) {
	existingGroup, err := h.store.Get(groupID)
	if err != nil {
		return nil, err
	}
	existingGroup.Id = ""
	newGroup, err := h.store.Put(existingGroup)
	return &pb.GroupResponse{Group: newGroup}, err
}

func (h *Handler) DeleteGroup(ctx context.Context, groupID string) (*pb.EmptyResponse, error) {
	var (
		err error
	)
	h.store.WithTx(func(s store.SQLTransactional) error {
		if _, err = h.store.TxGet(s, groupID); err != nil {
			// check if group exists
			return err
		}
		// delete all contracts by group id
		err = h.contractStore.TxDeleteContractsByGroupID(s, groupID)
		if err != nil && !strings.HasPrefix(err.Error(), "not found") {
			return err
		}
		return h.store.TxDelete(s, groupID)
	})
	return &pb.EmptyResponse{}, err
}

func (h *Handler) GetGroupByID(groupID string) (*pb.Group, error) {
	return h.store.Get(groupID)
}

func (h *Handler) GetGroupsBySubjectID(ctx context.Context, subjectID string) ([]*pb.Group, error) {
	var (
		contracts []contract.Contract
		groups    []*pb.Group
		err       error
	)
	h.store.WithTx(func(tx store.SQLTransactional) error {
		contracts, err = h.contractStore.TxListAllContractsBySubject(tx, subjectID)
		if err != nil {
			return err
		}
		groupIDs := make([]string, len(contracts), len(contracts))
		for i, contract := range contracts {
			groupIDs[i] = contract.GroupID
		}
		groups, err = h.store.TxBulkGet(tx, groupIDs)
		return err
	})
	return groups, err
}
