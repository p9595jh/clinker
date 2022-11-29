package service

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"syscall"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"layer/common/abi/erc721clinker"
	"layer/common/config"
	"layer/common/logger"
	"layer/internal/port/proto/clink"
)

type TransactionService struct {
	rpcUrl               string
	pk                   *ecdsa.PrivateKey
	chainId              *big.Int
	ca                   common.Address
	client               *ethclient.Client
	clinker              *erc721clinker.ClinkerERC721
	clinkerClientService *ClinkerClientService
}

func NewTransactionService(clinkerClientService *ClinkerClientService) *TransactionService {
	rpcUrl := config.V.GetString("eth.url")
	pk, err := crypto.HexToECDSA(config.V.GetString("eth.pk"))
	if err != nil {
		panic(err)
	}

	client, err := ethclient.DialContext(context.Background(), rpcUrl)
	if err != nil {
		panic(err)
	}

	chainId, err := client.NetworkID(context.Background())
	if err != nil {
		panic(err)
	}

	ca := common.HexToAddress(config.V.GetString("eth.clinker.address"))
	clinker, err := erc721clinker.NewClinkerERC721(ca, client)
	if err != nil {
		panic(err)
	}

	return &TransactionService{
		rpcUrl:               rpcUrl,
		pk:                   pk,
		chainId:              chainId,
		ca:                   ca,
		client:               client,
		clinker:              clinker,
		clinkerClientService: clinkerClientService,
	}
}

func (*TransactionService) name() string {
	return "TransactionService"
}

func (s *TransactionService) refresh() {
	if _, err := s.client.BlockNumber(context.Background()); errors.Is(err, syscall.ECONNRESET) {
		if s.client, err = ethclient.Dial(s.rpcUrl); err != nil {
			logger.Error(s.name()).E(err).W()
		} else if s.clinker, err = erc721clinker.NewClinkerERC721(s.ca, s.client); err != nil {
			logger.Error(s.name()).E(err).W()
		}
	}
}

func (s *TransactionService) Create(kind clink.Kind, userAddress common.Address, data string) (string, error) {
	s.refresh()

	// create transaction and return its hash first
	tx, err := s.clinker.Mint(nil, userAddress, data)
	if err != nil {
		logger.Error(s.name()).E(err).W()
		return "", err
	}

	// then it will be sent
	go func() {
		ctxRpc, ctxRpcCancel := context.WithCancel(context.Background())
		defer ctxRpcCancel()

		ctxChain, ctxChainCancel := context.WithCancel(context.Background())
		defer ctxChainCancel()

		signedTx, err := types.SignTx(tx, types.NewEIP155Signer(s.chainId), s.pk)
		if err != nil {
			s.clinkerClientService.Confirm(ctxRpc, &clink.ConfirmRequest{Kind: kind, Error: err.Error()})
			logger.Error(s.name()).E(err).W()
			return
		}

		if err = s.client.SendTransaction(ctxChain, signedTx); err != nil {
			s.clinkerClientService.Confirm(ctxRpc, &clink.ConfirmRequest{Kind: kind, Error: err.Error()})
			logger.Error(s.name()).E(err).W()
		} else {
			s.clinkerClientService.Confirm(ctxRpc, &clink.ConfirmRequest{Kind: kind, Id: signedTx.Hash().String()})
			// db insertion needed
		}
	}()

	return tx.Hash().String(), nil
}
