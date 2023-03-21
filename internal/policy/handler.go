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

func (h *Handler) CreatePolicy(ctx context.Context, policy *pb.Policy) (*pb.Policy, error) {
	policy, err := h.store.Put(policy)
	if err != nil {
		return nil, err
	}
	return policy, nil
}

func (h *Handler) UpdatePolicy(ctx context.Context, policy *pb.Policy) (*pb.Policy, error) {
	updatedPolicy, err := h.store.Put(policy)
	if err != nil {
		return nil, err
	}
	return updatedPolicy, nil
}

func (h *Handler) DeletePolicy(ctx context.Context, policyID string) (*pb.EmptyResponse, error) {
	err := h.store.Delete(policyID)
	return &pb.EmptyResponse{}, err
}

func (h *Handler) GetPolicyByID(policyID string) (*pb.Policy, error) {
	return h.store.Get(policyID)
}
