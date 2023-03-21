package api

import (
	"context"

	pb "github.com/dlshle/authnz/proto"
	"google.golang.org/grpc"
)

type client struct {
	conn       *grpc.ClientConn
	grpcClient pb.AuthNZClient
}

func NewAuthNZClient(endopint string) (*client, error) {
	conn, c, err := connect(endopint)
	if err != nil {
		return nil, err
	}
	return &client{conn: conn, grpcClient: c}, nil
}

func (c *client) Authorize(ctx context.Context, subjectID, policyID string) (verdict pb.Verdict, err error) {
	var resp *pb.AuthorizeResponse
	resp, err = c.grpcClient.Authorize(ctx, &pb.AuthorizeRequest{SubjectId: subjectID, PolicyId: policyID})
	verdict = resp.Verdict
	return
}

func (c *client) AddSubject(ctx context.Context, userID string) (*pb.Subject, error) {
	resp, err := c.grpcClient.AddSubject(ctx, &pb.AddSubjectRequest{UserId: userID})
	if err != nil {
		return nil, err
	}
	return resp.Subject, nil
}

func (c *client) FindSubjectsByUserID(ctx context.Context, userID string) ([]*pb.Subject, error) {
	resp, err := c.grpcClient.FindSubjectsByUserID(ctx, &pb.SubjectsByUserIDRequest{UserId: userID})
	if err != nil {
		return nil, err
	}
	return resp.Subjects, nil
}

func (c *client) AddSubjectWithAttributes(ctx context.Context, userID string, attributes map[string]string) (*pb.AddSubjectWithAttributesResponse, error) {
	var pbAttributes []*pb.Attribute
	for k, v := range attributes {
		pbAttributes = append(pbAttributes, &pb.Attribute{Key: k, Value: v})
	}
	return c.grpcClient.AddSubjectWithAttributes(ctx, &pb.AddSubjectWithAttributesRequest{UserId: userID, Attributes: pbAttributes})
}

func (c *client) CreateContract(ctx context.Context, subjectID, groupID string) (*pb.Contract, error) {
	resp, err := c.grpcClient.CreateContract(ctx, &pb.ContractRequest{Contract: &pb.Contract{SubjectId: subjectID, GroupId: groupID}})
	if err != nil {
		return nil, err
	}
	return resp.Contract, nil
}

func (c *client) CreateGroupForSubjects(ctx context.Context, subjectIDs []string, attributes map[string]string) (*pb.CreateGroupForSubjectsResponse, error) {
	var pbAttributes []*pb.Attribute
	for k, v := range attributes {
		pbAttributes = append(pbAttributes, &pb.Attribute{Key: k, Value: v})
	}
	return c.grpcClient.CreateGroupsForSubjects(ctx, &pb.CreateGroupForSubjectsRequest{SubjectIds: subjectIDs, Attributes: pbAttributes})
}

func (c *client) GetGroupsBySubjectID(ctx context.Context, subjectID string) ([]*pb.Group, error) {
	resp, err := c.grpcClient.GetGroupsBySubjectID(ctx, &pb.SubjectIDRequest{SubjectId: subjectID})
	if err != nil {
		return nil, err
	}
	return resp.Groups, nil
}

func (c *client) Close() error {
	return c.conn.Close()
}

func connect(endpoint string) (*grpc.ClientConn, pb.AuthNZClient, error) {
	// conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	c := pb.NewAuthNZClient(conn)
	return conn, c, nil
}
