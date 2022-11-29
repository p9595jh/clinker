package service

import (
	"context"
	"errors"
	"layer/internal/port/proto/clink"

	"google.golang.org/grpc"
)

type ClinkerClientService struct {
	clink.ClinkClient
}

func (s *ClinkerClientService) Confirm(ctx context.Context, in *clink.ConfirmRequest, opts ...grpc.CallOption) (*clink.ConfirmResponse, error) {
	if in.Error != "" {
		return &clink.ConfirmResponse{Kind: in.Kind, Id: in.Id, Confirmed: false, Error: in.Error}, errors.New(in.Error)
	} else {
		return &clink.ConfirmResponse{Kind: in.Kind, Id: in.Id, Confirmed: true}, nil
	}
}
