package rpcserver

import (
	"context"
	"encoding/json"
	"fmt"
	"layer/common/logger"
	"layer/internal/domain/service"
	"layer/internal/infrastructure/rpc/proto/clink"
	"net"

	"github.com/ethereum/go-ethereum/common"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

type ClinkRpcServer struct {
	clink.ClinkServer
	port               int
	transactionService *service.TransactionService
}

func NewClinkRpcServer(port int, transactionService *service.TransactionService) *ClinkRpcServer {
	return &ClinkRpcServer{
		port:               port,
		transactionService: transactionService,
	}
}

func (*ClinkRpcServer) name() string {
	return "ClinkRpcServer"
}

func (s *ClinkRpcServer) VestigeCreate(ctx context.Context, req *clink.VestigeCreateRequest) (*clink.TxHashResponse, error) {
	address := common.HexToAddress(req.Address)
	b, err := json.Marshal(req)
	if err != nil {
		logger.Error(s.name()).E(err).W()
		return &clink.TxHashResponse{Kind: clink.Kind_VESTIGE, Error: err.Error()}, err
	}

	hash, err := s.transactionService.Create(clink.Kind_VESTIGE, address, string(b))
	if err != nil {
		logger.Error(s.name()).E(err).W()
		return &clink.TxHashResponse{Kind: clink.Kind_VESTIGE, Error: err.Error()}, err
	}
	return &clink.TxHashResponse{Kind: clink.Kind_VESTIGE, TxHash: hash}, nil
}

func (s *ClinkRpcServer) AppraisalCreate(ctx context.Context, req *clink.AppraisalCreateRequest) (*clink.TxHashResponse, error) {
	address := common.HexToAddress(req.Address)
	b, err := json.Marshal(req)
	if err != nil {
		logger.Error(s.name()).E(err).W()
		return &clink.TxHashResponse{Kind: clink.Kind_APPRAISAL, Error: err.Error()}, err
	}

	hash, err := s.transactionService.Create(clink.Kind_APPRAISAL, address, string(b))
	if err != nil {
		logger.Error(s.name()).E(err).W()
		return &clink.TxHashResponse{Kind: clink.Kind_APPRAISAL, Error: err.Error()}, err
	}
	return &clink.TxHashResponse{Kind: clink.Kind_APPRAISAL, TxHash: hash}, nil
}

func (s *ClinkRpcServer) logger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	logger.Info(s.name(), info.FullMethod).D("request", req)
	return handler(ctx, req)
}

func (s *ClinkRpcServer) Listen() <-chan error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(s.logger)),
	)
	clink.RegisterClinkServer(grpcServer, s)

	ch := make(chan error)
	go func() { ch <- grpcServer.Serve(lis) }()
	return ch
}
