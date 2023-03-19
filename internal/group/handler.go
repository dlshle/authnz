package group

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

func (h *Handler) CreateGroup(ctx context.Context, req *pb.GroupRequest) (*pb.GroupResponse, error) {
	newGroup, err := h.store.Put(req.Group)
	if err != nil {
		return nil, err
	}
	return &pb.GroupResponse{Group: newGroup}, nil
}

func (h *Handler) UpdateGroup(ctx context.Context, req *pb.GroupRequest) (*pb.GroupResponse, error) {
	updatedGroup, err := h.store.Put(req.Group)
	if err != nil {
		return nil, err
	}
	return &pb.GroupResponse{Group: updatedGroup}, nil
}

func (h *Handler) DeleteGroup(ctx context.Context, req *pb.GroupByIDRequest) (*pb.EmptyResponse, error) {
	err := h.store.Delete(req.GroupId)
	return &pb.EmptyResponse{}, err
}

func (h *Handler) GetGroupByID(groupID string) (*pb.Group, error) {
	return h.store.Get(groupID)
}
