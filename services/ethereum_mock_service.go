package services

import (
	"context"

	"github.com/littlecheny/go-backend/domain"
)

type ethereumMockService struct {
}

func NewEthereumMockService() domain.EthereumService {
	return &ethereumMockService{}
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

/*
type EthereumService interface {
	// 发送交易
	//SendTransaction(ctx context.Context, req *TransactionSendRequest) (string, error)

	// 查询交易状态
	GetTransactionStatus(ctx context.Context, hash string) (string, error)

	// 查询余额
	GetBalance(ctx context.Context, address string) (string, error)

	// 查询最新区块高度
	GetLatestBlockNumber(ctx context.Context) (uint64, error)
}*/
