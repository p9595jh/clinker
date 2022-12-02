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
	rpcclient "layer/internal/infrastructure/rpc/client"
	"layer/internal/infrastructure/rpc/proto/clink"
)

type TransactionService struct {
	rpcUrl         string
	pk             *ecdsa.PrivateKey
	chainId        *big.Int
	ca             common.Address
	ethClient      *ethclient.Client
	clinkerErc721  *erc721clinker.ClinkerERC721
	clinkRpcClient *rpcclient.ClinkRpcClient
}

func NewTransactionService(clinkRpcClient *rpcclient.ClinkRpcClient) *TransactionService {
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
		rpcUrl:         rpcUrl,
		pk:             pk,
		chainId:        chainId,
		ca:             ca,
		ethClient:      client,
		clinkerErc721:  clinker,
		clinkRpcClient: clinkRpcClient,
	}
}

func (*TransactionService) name() string {
	return "TransactionService"
}

func (s *TransactionService) refresh() {
	if _, err := s.ethClient.BlockNumber(context.Background()); errors.Is(err, syscall.ECONNRESET) {
		if s.ethClient, err = ethclient.Dial(s.rpcUrl); err != nil {
			logger.Error(s.name()).E(err).W()
		} else if s.clinkerErc721, err = erc721clinker.NewClinkerERC721(s.ca, s.ethClient); err != nil {
			logger.Error(s.name()).E(err).W()
		}
	}
}

func (s *TransactionService) Create(kind clink.Kind, userAddress common.Address, data string) (string, error) {
	s.refresh()

	// create transaction and return its hash first
	tx, err := s.clinkerErc721.Mint(nil, userAddress, data)
	if err != nil {
		logger.Error(s.name()).E(err).W()
		return "", err
	}

	// then it will be sent
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		signedTx, err := types.SignTx(tx, types.NewEIP155Signer(s.chainId), s.pk)
		if err != nil {
			s.clinkRpcClient.ConfirmError(kind, err)
			logger.Error(s.name()).E(err).W()
			return
		}

		if err = s.ethClient.SendTransaction(ctx, signedTx); err != nil {
			s.clinkRpcClient.ConfirmError(kind, err)
			logger.Error(s.name()).E(err).W()
		} else {
			hash := signedTx.Hash().String()
			s.clinkRpcClient.Confirm(kind, hash)
			logger.Info(s.name()).Wf("[%s] confirmed: %s", clink.Kind_name[int32(kind)], hash)
			// db insertion needed
		}
	}()

	return tx.Hash().String(), nil
}

func (s *TransactionService) Initializer() {
	sink := make(chan *erc721clinker.ClinkerERC721AddressAvailable)
	sub, err := s.clinkerErc721.WatchAddressAvailable(nil, sink, nil)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case err := <-sub.Err():
			logger.Error(s.name(), "Subscription").E(err).W()
		case event := <-sink:
			address := event.Available.String()
			logger.Info(s.name(), "Subscription").Wf("[USER] confirmed: %s", address)
			s.clinkRpcClient.Confirm(clink.Kind_USER, address)
		}
	}
}
