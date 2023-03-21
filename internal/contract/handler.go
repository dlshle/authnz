package contract

import (
	"context"

	pb "github.com/dlshle/authnz/proto"
)

type Handler struct {
	store Store
}

func NewHandler(store Store) *Handler {
	return &Handler{store: store}
}

func (h *Handler) CreateContract(ctx context.Context, contract *pb.Contract) (*pb.ContractResponse, error) {
	contract, err := h.store.AddNewContract(contract.SubjectId, contract.GroupId)
	return &pb.ContractResponse{Contract: contract}, err
}

func (h *Handler) DeleteContract(ctx context.Context, contractID string) (*pb.EmptyResponse, error) {
	err := h.store.DeleteContractByContractID(contractID)
	return &pb.EmptyResponse{}, err
}

func (h *Handler) GetGroupsBySubjectID(subjectID string) ([]*pb.Group, error) {
	return h.store.ListGroupsBySubjectID(subjectID)
}
