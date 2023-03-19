package subject

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

func (h *Handler) AddSubject(ctx context.Context, req *pb.AddSubjectRequest) (*pb.AddSubjectResponse, error) {
	subject, err := h.store.Put(&pb.Subject{
		UserId: req.UserId,
	})
	if err != nil {
		return nil, err
	}
	return &pb.AddSubjectResponse{Subject: subject}, nil
}

func (h *Handler) DeleteSubject(ctx context.Context, req *pb.SubjectByIDRequest) (*pb.EmptyResponse, error) {
	err := h.store.Delete(req.SubjectId)
	if err != nil {
		return nil, err
	}
	return &pb.EmptyResponse{}, nil
}

func (h *Handler) GetSubjectByID(subjectID string) (*pb.Subject, error) {
	return h.store.Get(subjectID)
}
