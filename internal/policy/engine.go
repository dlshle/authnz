package policy

import (
	pb "github.com/dlshle/authnz/proto"
)

type Engine interface {
	Check(policy *pb.Policy, group *pb.Group, ctx *pb.ContextProperty) pb.Verdict
}
