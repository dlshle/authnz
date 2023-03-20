package subject

import (
	"context"

	"github.com/dlshle/authnz/internal/contract"
	"github.com/dlshle/authnz/internal/group"
	"github.com/dlshle/authnz/pkg/store"
	pb "github.com/dlshle/authnz/proto"
	"github.com/dlshle/gommon/logging"
)

type Handler struct {
	store         Store
	contractStore contract.Store
	groupStore    group.Store
	logger        logging.Logger
}

func NewHandler(store Store, contractStore contract.Store, groupStore group.Store) *Handler {
	return &Handler{store: store,
		contractStore: contractStore,
		groupStore:    groupStore,
		logger:        logging.GlobalLogger.WithPrefix("[SubjectHandler]")}
}

func (h *Handler) AddSubject(ctx context.Context, req *pb.AddSubjectRequest) (*pb.AddSubjectResponse, error) {
	subject, err := h.store.Put(&pb.Subject{
		UserId: req.UserId,
	})
	if err != nil {
		return nil, err
	}
	return &pb.AddSubjectResponse{Subject: subject}, nil
}

func (h *Handler) DeleteSubject(ctx context.Context, subjectID string) (*pb.EmptyResponse, error) {
	err := h.store.WithTX(func(s store.SQLTransactional) error {
		// 1. list all contracts about subject id
		contracts, err := h.contractStore.TxListAllContractsBySubject(s, subjectID)
		if err != nil {
			h.logger.Errorf(ctx, "failed to list all contracts by %s due to %s", subjectID, err.Error())
			return err
		}

		// 2. delete subject
		err = h.store.TxDelete(s, subjectID)
		if err != nil {
			h.logger.Errorf(ctx, "failed to delete subject %s due to %s", subjectID, err.Error())
			return err
		}

		// 3. check if any group from contracts has no reference other than this subject
		for _, contract := range contracts {
			err = h.contractStore.TxDeleteContract(s, contract.ID)
			if err != nil {
				h.logger.Warnf(ctx, "failed to delete contract %s due to %s", contract.ID, err.Error())
			}
			contractsByGroup, err := h.contractStore.TxListContractsByGroupID(s, contract.GroupID)
			if err != nil {
				h.logger.Warnf(ctx, "failed to list contracts by groupID %s due to %s", contract.GroupID, err.Error())
				continue
			}
			if len(contractsByGroup) == 1 {
				h.logger.Infof(ctx, "delete group by id %s due to zombie group", contract.GroupID)
				err = h.groupStore.TxDelete(s, contract.GroupID)
				if err != nil {
					h.logger.Warnf(ctx, "failed to delete zombie group %s due to %s", contract.GroupID, err.Error())
					continue
				}
			}
		}
		return nil
	})

	return &pb.EmptyResponse{}, err
}

func (h *Handler) GetSubjectByID(subjectID string) (*pb.Subject, error) {
	return h.store.Get(subjectID)
}
