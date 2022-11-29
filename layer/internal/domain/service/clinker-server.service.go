package service

import (
	"context"
	"encoding/json"
	"layer/common/logger"
	"layer/internal/port/proto/clink"

	"github.com/ethereum/go-ethereum/common"
)

type ClinkerServerService struct {
	clink.ClinkServer
	transactionService *TransactionService
}

func NewClinkerServerService(transactionService *TransactionService) *ClinkerServerService {
	return &ClinkerServerService{
		transactionService: transactionService,
	}
}

func (*ClinkerServerService) name() string {
	return "ClinkerServerService"
}

func (s *ClinkerServerService) VestigeCreate(ctx context.Context, req *clink.VestigeCreateRequest) (*clink.TxHashResponse, error) {
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

func (s *ClinkerServerService) AppraisalCreate(ctx context.Context, req *clink.AppraisalCreateRequest) (*clink.TxHashResponse, error) {
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
