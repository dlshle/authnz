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

func (h *Handler) CreateContract(ctx context.Context, req *pb.ContractRequest) (*pb.ContractResponse, error) {
	contract, err := h.store.AddNewContract(req.Contract.SubjectId, req.Contract.GroupId)
	return &pb.ContractResponse{Contract: contract}, err
}

func (h *Handler) DeleteContract(ctx context.Context, req *pb.DeleteContractRequest) (*pb.EmptyResponse, error) {
	err := h.store.DeleteContractByContractID(req.ContractId)
	return &pb.EmptyResponse{}, err
}

func (h *Handler) GetGroupsBySubjectID(subjectID string) ([]*pb.Group, error) {
	return h.store.ListGroupsBySubjectID(subjectID)
}
