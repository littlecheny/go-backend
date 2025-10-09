package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type EthereumService struct {
	mock.Mock
}

func (m *EthereumService) GetTransactionStatus(ctx context.Context, hash string) (string, error) {
	return "pending", nil
}

func (m *EthereumService) GetBalance(ctx context.Context, address string) (string, error) {
	return "1000000000000000000", nil
}

func (m *EthereumService) GetLatestBlockNumber(ctx context.Context) (uint64, error) {
	return 1000000, nil
}
