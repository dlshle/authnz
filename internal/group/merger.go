package group

import (
	pb "github.com/dlshle/authnz/proto"
)

func MergeGroups(groups []*pb.Group) *pb.Group {
	finalGroup := &pb.Group{}
	for _, group := range groups {
		finalGroup.Attributes = append(finalGroup.Attributes, group.Attributes...)
	}
	return finalGroup
}
