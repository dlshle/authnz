package server

import (
	"context"
	"net"

	"github.com/dlshle/authnz/internal/config"
	"github.com/dlshle/authnz/internal/contract"
	"github.com/dlshle/authnz/internal/group"
	"github.com/dlshle/authnz/internal/policy"
	"github.com/dlshle/authnz/internal/subject"
	pb "github.com/dlshle/authnz/proto"
	"github.com/dlshle/gommon/errors"
	"github.com/dlshle/gommon/logging"
	"github.com/gofrs/uuid"
	"google.golang.org/grpc"

	"google.golang.org/grpc/reflection"
)

type server struct {
	logger          logging.Logger
	subjectHandler  *subject.Handler
	groupHandler    *group.Handler
	policyHandler   *policy.Handler
	contractHandler *contract.Handler
	*pb.UnimplementedAuthNZServer
}

func NewGRPCServer(
	subjectHandler *subject.Handler,
	groupHandler *group.Handler,
	policyHandler *policy.Handler,
	contractHandler *contract.Handler,
) pb.AuthNZServer {
	return &server{
		logger:          logging.GlobalLogger.WithPrefix("[GRPCServer]"),
		subjectHandler:  subjectHandler,
		groupHandler:    groupHandler,
		policyHandler:   policyHandler,
		contractHandler: contractHandler,
	}
}

func (s *server) Authorize(ctx context.Context, req *pb.AuthorizeRequest) (*pb.AuthorizeResponse, error) {
	engine := policy.NewEngine()
	groups, err := s.contractHandler.GetGroupsBySubjectID(req.SubjectId)
	if err != nil {
		return nil, errors.Error("failed to get groups by subject due to " + err.Error())
	}
	policy, err := s.policyHandler.GetPolicyByID(req.PolicyId)
	if err != nil {
		return nil, errors.Error("failed to get policy due to " + err.Error())
	}
	verdict, err := engine.Check(policy, group.MergeGroups(groups), nil)
	return &pb.AuthorizeResponse{Verdict: verdict}, err
}

func handleRequest[T any](logger logging.Logger, action string, ctx context.Context, reqStr string, handler func() (T, error)) (T, error) {
	logger.Infof(ctx, "received %s request %s", action, reqStr)
	res, err := handler()
	if err != nil {
		logger.Errorf(ctx, "request %s failed with error %s", reqStr, err.Error())
	}
	return res, err
}

func (s *server) AddSubject(ctx context.Context, req *pb.AddSubjectRequest) (*pb.AddSubjectResponse, error) {
	return s.subjectHandler.AddSubject(ctx, req)
}

func (s *server) GetSubject(ctx context.Context, req *pb.SubjectIDRequest) (*pb.Subject, error) {
	return s.subjectHandler.GetSubjectByID(req.SubjectId)
}

func (s *server) FindSubjectsByUserID(ctx context.Context, req *pb.SubjectsByUserIDRequest) (*pb.SubjectsByUserIDResponse, error) {
	pbSubjects, err := s.subjectHandler.FindSubjectsByUserID(req.UserId)
	if err != nil {
		return nil, err
	}
	return &pb.SubjectsByUserIDResponse{Subjects: pbSubjects}, nil
}

func (s *server) CreateGroupsForSubjects(ctx context.Context, req *pb.CreateGroupForSubjectsRequest) (*pb.CreateGroupForSubjectsResponse, error) {
	return s.subjectHandler.CreateGroupsForSubjects(ctx, req.SubjectIds, req.Attributes)
}

func (s *server) AddSubjectWithAttributes(ctx context.Context, req *pb.AddSubjectWithAttributesRequest) (*pb.AddSubjectWithAttributesResponse, error) {
	return s.subjectHandler.AddSubjectWithAttributes(ctx, req.UserId, req.Attributes)
}

func (s *server) DeleteSubject(ctx context.Context, req *pb.SubjectIDRequest) (*pb.EmptyResponse, error) {
	return s.subjectHandler.DeleteSubject(ctx, req.SubjectId)
}

func (s *server) CreateGroup(ctx context.Context, req *pb.GroupRequest) (*pb.GroupResponse, error) {
	return s.groupHandler.CreateGroup(ctx, req.Group)
}

func (s *server) GetGroup(ctx context.Context, req *pb.GroupByIDRequest) (*pb.GroupResponse, error) {
	group, err := s.groupHandler.GetGroupByID(req.GroupId)
	return &pb.GroupResponse{Group: group}, err
}

func (s *server) GetGroupsBySubjectID(ctx context.Context, req *pb.SubjectIDRequest) (*pb.GroupsResponse, error) {
	groups, err := s.groupHandler.GetGroupsBySubjectID(ctx, req.SubjectId)
	if err != nil {
		return nil, err
	}
	return &pb.GroupsResponse{Groups: groups}, nil
}

func (s *server) UpdateGroup(ctx context.Context, req *pb.GroupRequest) (*pb.GroupResponse, error) {
	return s.groupHandler.UpdateGroup(ctx, req.Group)
}

func (s *server) DuplicateGroup(ctx context.Context, req *pb.GroupByIDRequest) (*pb.GroupResponse, error) {
	return s.groupHandler.DuplicateGroup(ctx, req.GroupId)
}

func (s *server) DeleteGroup(ctx context.Context, req *pb.GroupByIDRequest) (*pb.EmptyResponse, error) {
	return s.groupHandler.DeleteGroup(ctx, req.GroupId)
}

func (s *server) CreatePolicy(ctx context.Context, req *pb.PolicyRequest) (*pb.Policy, error) {
	return s.policyHandler.CreatePolicy(ctx, req.Policy)
}

func (s *server) GetPolicy(ctx context.Context, req *pb.PolicyByIDRequest) (*pb.Policy, error) {
	return s.policyHandler.GetPolicyByID(req.PolicyId)
}

func (s *server) UpdatePolicy(ctx context.Context, req *pb.PolicyRequest) (*pb.Policy, error) {
	return s.policyHandler.UpdatePolicy(ctx, req.Policy)
}

func (s *server) DeletePolicy(ctx context.Context, req *pb.PolicyByIDRequest) (*pb.EmptyResponse, error) {
	return s.policyHandler.DeletePolicy(ctx, req.PolicyId)
}

func (s *server) CreateContract(ctx context.Context, req *pb.ContractRequest) (*pb.ContractResponse, error) {
	return s.contractHandler.CreateContract(ctx, req.Contract)
}

func (s *server) DeleteContract(ctx context.Context, req *pb.DeleteContractRequest) (*pb.EmptyResponse, error) {
	return s.contractHandler.DeleteContract(ctx, req.ContractId)
}

func StartServer(serverCfg config.ServerConfig, server pb.AuthNZServer) error {
	lis, err := net.Listen("tcp", serverCfg.GRPC)
	if err != nil {
		return err
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		tracingID, _ := uuid.NewV4()
		ctx = logging.WrapCtx(ctx, "traceID", tracingID.String())
		logging.GlobalLogger.Infof(ctx, "[%s] received grpc request %v ", info.FullMethod, req)
		resp, err = handler(ctx, req)
		logging.GlobalLogger.Infof(ctx, "[%s] request done with response: {%v} err: %v", info.FullMethod, resp, err)
		return
	}))
	pb.RegisterAuthNZServer(s, server)
	reflection.Register(s)
	logging.GlobalLogger.Infof(context.Background(), "server started on %s", serverCfg.GRPC)
	return s.Serve(lis)
}
