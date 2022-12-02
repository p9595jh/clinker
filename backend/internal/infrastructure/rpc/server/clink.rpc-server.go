package rpcserver

import (
	"clinker-backend/common/logger"
	"clinker-backend/internal/infrastructure/database/entity"
	"clinker-backend/internal/infrastructure/database/repository"
	"clinker-backend/internal/infrastructure/database/repository/reposh"
	"clinker-backend/internal/infrastructure/rpc/proto/clink"
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
)

type ClinkRpcServer struct {
	clink.ClinkServer
	port                int
	userRepository      repository.UserRepository
	vestigeRepository   repository.VestigeRepository
	appraisalRepository repository.AppraisalRepository
}

func NewClinkRpcServer(
	port int,
	userRepository repository.UserRepository,
	vestigeRepository repository.VestigeRepository,
	appraisalRepository repository.AppraisalRepository,
) *ClinkRpcServer {
	return &ClinkRpcServer{
		port:                port,
		userRepository:      userRepository,
		vestigeRepository:   vestigeRepository,
		appraisalRepository: appraisalRepository,
	}
}

func (*ClinkRpcServer) name() string {
	return "ClinkRpcServer"
}

func (s *ClinkRpcServer) Confirm(ctx context.Context, req *clink.ConfirmRequest) (*clink.ConfirmResponse, error) {
	var err error
	switch req.Kind {
	case clink.Kind_USER:
		err = s.userRepository.Update(
			&reposh.EntityParts[entity.User]{EntityFn: func(e *entity.User) { e.Id = req.Id }},
			&reposh.EntityParts[entity.User]{EntityFn: func(e *entity.User) { e.Confirmed = true }},
		).Error
	case clink.Kind_VESTIGE:
		err = s.vestigeRepository.Update(
			&reposh.EntityParts[entity.Vestige]{EntityFn: func(e *entity.Vestige) { e.TxHash = req.Id }},
			&reposh.EntityParts[entity.Vestige]{EntityFn: func(e *entity.Vestige) { e.Confirmed = true }},
		).Error
	case clink.Kind_APPRAISAL:
		err = s.appraisalRepository.Update(
			&reposh.EntityParts[entity.Appraisal]{EntityFn: func(e *entity.Appraisal) { e.TxHash = req.Id }},
			&reposh.EntityParts[entity.Appraisal]{EntityFn: func(e *entity.Appraisal) { e.Confirmed = true }},
		).Error
	}
	if err != nil {
		return nil, err
	} else {
		return &clink.ConfirmResponse{
			Kind:      req.Kind,
			Id:        req.Id,
			Confirmed: true,
		}, nil
	}
}

func (s *ClinkRpcServer) Listen() <-chan error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	clink.RegisterClinkServer(grpcServer, s)

	ch := make(chan error)
	go func() { ch <- grpcServer.Serve(lis) }()
	logger.Info(s.name()).Wf("Listeneing %s at %d", s.name(), s.port)
	return ch
}
