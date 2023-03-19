package group

import (
	pb "github.com/dlshle/authnz/proto"
)

type Group struct {
	ID         string
	Attributes map[string]string
}

func FromPB(pbGroup *pb.Group) Group {
	return Group{
		ID:         pbGroup.GetId(),
		Attributes: attributesToMap(pbGroup.GetAttributes()),
	}
}

func attributesToMap(attributes []*pb.Attribute) map[string]string {
	attributeMap := make(map[string]string)
	for _, attribute := range attributes {
		attributeMap[attribute.GetKey()] = attribute.GetValue()
	}
	return attributeMap
}
