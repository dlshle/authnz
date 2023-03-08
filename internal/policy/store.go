package policy

import (
	pb "github.com/dlshle/authnz/proto"
)

type PolicyStore interface {
	Get(id string) (*pb.Policy, error)
	Delete(id string) (*pb.Policy, error)
	Update(id string, policy *pb.Policy) (*pb.Policy, error)
}
