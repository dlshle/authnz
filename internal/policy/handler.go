package policy

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

func (h *Handler) CreatePolicy(ctx context.Context, req *pb.PolicyRequest) (*pb.Policy, error) {
	policy, err := h.store.Put(req.Policy)
	if err != nil {
		return nil, err
	}
	return policy, nil
}

func (h *Handler) UpdatePolicy(ctx context.Context, req *pb.PolicyRequest) (*pb.Policy, error) {
	updatedPolicy, err := h.store.Put(req.Policy)
	if err != nil {
		return nil, err
	}
	return updatedPolicy, nil
}

func (h *Handler) DeletePolicy(ctx context.Context, req *pb.PolicyByIDRequest) (*pb.EmptyResponse, error) {
	err := h.store.Delete(req.PolicyId)
	return &pb.EmptyResponse{}, err
}

func (h *Handler) GetPolicyByID(policyID string) (*pb.Policy, error) {
	return h.store.Get(policyID)
}
