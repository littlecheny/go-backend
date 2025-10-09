package domain

import (
	"context"
)

type TransactionSendRequest struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Value    string `json:"value"`
	GasPrice string `json:"gas_price"`
	GasLimit string `json:"gas_limit"`
}

type Transaction struct {
	Hash   string `json:"hash"`
	Status string `json:"status"`
}

type EthereumService interface {
	// 发送交易
	//SendTransaction(ctx context.Context, req *TransactionSendRequest) (string, error)

	// 查询交易状态
	GetTransactionStatus(ctx context.Context, hash string) (string, error)

	// 查询余额
	GetBalance(ctx context.Context, address string) (string, error)

	// 查询最新区块高度
	GetLatestBlockNumber(ctx context.Context) (uint64, error)
}
