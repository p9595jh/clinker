package service

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"layer/common/abi/erc721clinker"
	"layer/common/config"
	"layer/common/logger"
	"layer/common/util"
	"layer/internal/infrastructure/database/entity"
	"layer/internal/infrastructure/database/repository"
	rpcclient "layer/internal/infrastructure/rpc/client"
	"layer/internal/infrastructure/rpc/proto/clink"
)

type TransactionService struct {
	rpcUrl          string
	pk              *ecdsa.PrivateKey
	chainId         *big.Int
	ca              common.Address
	ethClient       *ethclient.Client
	clinkerErc721   *erc721clinker.ClinkerERC721
	clinkRpcClient  *rpcclient.ClinkRpcClient
	clinkRepository *repository.ClinkRepository
	caRepository    *repository.CaRepository
	txnRepository   *repository.TxnRepository
}

func NewTransactionService(
	clinkRpcClient *rpcclient.ClinkRpcClient,
	clinkRepository *repository.ClinkRepository,
	caRepository *repository.CaRepository,
	txnRepository *repository.TxnRepository,
) *TransactionService {
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

	ca, err := caRepository.FindOne(
		bson.M{"name": "clinker"},
		options.FindOne().SetSort(bson.D{
			{Key: "timestamp", Value: -1},
		}),
	)
	if err != nil {
		panic(err)
	}

	clinker, err := erc721clinker.NewClinkerERC721(ca.Address, client)
	if err != nil {
		panic(err)
	}
	logger.Info("TransactionService").D("ClinkerERC721Address", ca.Address).W()

	return &TransactionService{
		rpcUrl:          rpcUrl,
		pk:              pk,
		chainId:         chainId,
		ca:              ca.Address,
		ethClient:       client,
		clinkerErc721:   clinker,
		clinkRpcClient:  clinkRpcClient,
		clinkRepository: clinkRepository,
		caRepository:    caRepository,
		txnRepository:   txnRepository,
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
		ctx := context.TODO()

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

			// insert clink
			go func() {
				id, err := s.clinkRepository.Insert(&entity.Clink{
					TxHash:    hash,
					Address:   userAddress.String(),
					Timestamp: util.Now(),
				})
				if err != nil {
					logger.Error(s.name()).E(err).W()
				} else {
					logger.Info(s.name()).Wf("clink %v inserted", id)
				}
			}()

			// insert txn
			go func() {
				s.refresh()

				// limit 5s
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
				defer cancel()

				tx, _, err := s.txFromHash(ctx, hash)
				if err != nil {
					logger.Error(s.name()).E(err).W()
					return
				}

				gasPrice, gasUsed, txFee, err := s.gasFromTx(ctx, tx)
				if err != nil {
					logger.Error(s.name()).E(err).W()
					return
				}

				id, err := s.txnRepository.Insert(&entity.Txn{
					Hash:      hash,
					GasPrice:  gasPrice,
					GasUsed:   gasUsed,
					Fee:       txFee,
					Timestamp: util.Now(),
				})
				if err != nil {
					logger.Error(s.name()).E(err).W()
				} else {
					logger.Info(s.name()).Wf("txn %v inserted", id)
				}
			}()
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

func (s *TransactionService) txFromHash(ctx context.Context, hash string) (*types.Transaction, bool, error) {
	return s.ethClient.TransactionByHash(ctx, common.HexToHash(hash))
}

func (s *TransactionService) gasFromTx(ctx context.Context, tx *types.Transaction) (gasPrice, gasUsed, txFee *big.Int, err error) {
	receipt, err := s.ethClient.TransactionReceipt(ctx, tx.Hash())
	if err != nil {
		return
	}

	block, err := s.ethClient.BlockByNumber(ctx, receipt.BlockNumber)
	if err != nil {
		return
	}

	baseFee := block.BaseFee()
	if baseFee == nil {
		baseFee = tx.GasPrice()
	}

	gasPrice = new(big.Int).Add(baseFee, tx.EffectiveGasTipValue(baseFee))
	gasUsed = big.NewInt(int64(receipt.GasUsed))
	txFee = new(big.Int).Mul(gasPrice, gasUsed)
	return
}
