package services

import (
	"context"

	"github.com/littlecheny/go-backend/domain"
)

type ethereumMockService struct {}

func NewEthereumMockService() domain.EthereumService {
	return &ethereumMockService{}
}

func (m *ethereumMockService) SendTransaction(ctx context.Context, req *domain.TransactionSendRequest) (string, error) {
	return "0xmocktxhash", nil
}

func (m *ethereumMockService) GetTransactionStatus(ctx context.Context, hash string) (string, error) {
	return "pending", nil
}

func (m *ethereumMockService) GetBalance(ctx context.Context, address string) (string, error) {
	return "1000000000000000000", nil
}

func (m *ethereumMockService) GetLatestBlockNumber(ctx context.Context) (uint64, error) {
	return 1000000, nil
}
