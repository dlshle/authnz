package server

import (
	"context"

	pb "github.com/dlshle/authnz/proto"
	"google.golang.org/grpc"
)

type server struct {
	*pb.UnimplementedAuthNZServer
}

func (s *server) Authorize(context.Context, *pb.AuthorizeRequest) (*pb.AuthorizeResponse, error) {
	return &pb.AuthorizeResponse{}, nil
}

func StartServer(port int) {
	s := grpc.NewServer()
	pb.RegisterAuthNZServer(s, &server{})
}
