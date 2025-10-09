package services

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/littlecheny/go-backend/domain"
)

type ethereumSepoliaService struct {
	client *ethclient.Client
}

func NewEthereumSepoliaService(rpcURL string) (domain.EthereumService, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}
	return &ethereumSepoliaService{client: client}, nil
}

func (s *ethereumSepoliaService) GetTransactionStatus(ctx context.Context, hash string) (string, error) {
	receipt, err := s.client.TransactionReceipt(ctx, common.HexToHash(hash))
	if err != nil {
		return "", err
	}
	switch receipt.Status {
	case 0:
		return "failed", nil
	case 1:
		return "success", nil
	default:
		return "unknown", nil
	}
}

func (s *ethereumSepoliaService) GetBalance(ctx context.Context, address string) (string, error) {
	balance, err := s.client.BalanceAt(ctx, common.HexToAddress(address), nil)
	if err != nil {
		return "", err
	}
	return balance.String(), nil
}

func (s *ethereumSepoliaService) GetLatestBlockNumber(ctx context.Context) (uint64, error) {
	blockNumber, err := s.client.BlockNumber(ctx)
	if err != nil {
		return 0, err
	}
	return blockNumber, nil
}
