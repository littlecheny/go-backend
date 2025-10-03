package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionTransaction = "transactions"
)

// TransactionStatus 交易状态
type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusConfirmed TransactionStatus = "confirmed"
	TransactionStatusFailed    TransactionStatus = "failed"
)

// TransactionType 交易类型
type TransactionType string

const (
	TransactionTypeSend    TransactionType = "send"
	TransactionTypeReceive TransactionType = "receive"
)

// Transaction 交易模型
type Transaction struct {
	ID              primitive.ObjectID `bson:"_id" json:"id"`
	UserID          primitive.ObjectID `bson:"user_id" json:"user_id"`
	WalletID        primitive.ObjectID `bson:"wallet_id" json:"wallet_id"`
	Hash            string             `bson:"hash" json:"hash"`
	From            string             `bson:"from" json:"from"`
	To              string             `bson:"to" json:"to"`
	Value           string             `bson:"value" json:"value"`           // Wei格式的金额
	GasPrice        string             `bson:"gas_price" json:"gas_price"`   // Wei格式的Gas价格
	GasLimit        uint64             `bson:"gas_limit" json:"gas_limit"`
	GasUsed         uint64             `bson:"gas_used" json:"gas_used"`
	Nonce           uint64             `bson:"nonce" json:"nonce"`
	Status          TransactionStatus  `bson:"status" json:"status"`
	Type            TransactionType    `bson:"type" json:"type"`
	Network         string             `bson:"network" json:"network"`
	BlockNumber     uint64             `bson:"block_number" json:"block_number"`
	BlockHash       string             `bson:"block_hash" json:"block_hash"`
	TransactionFee  string             `bson:"transaction_fee" json:"transaction_fee"` // 实际交易费用
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
	ConfirmedAt     *time.Time         `bson:"confirmed_at" json:"confirmed_at,omitempty"`
}

// TransactionSendRequest 发送交易请求
type TransactionSendRequest struct {
	WalletID string `json:"wallet_id" binding:"required"`
	To       string `json:"to" binding:"required"`
	Amount   string `json:"amount" binding:"required"` // ETH格式的金额
	Password string `json:"password" binding:"required"`
	GasPrice string `json:"gas_price,omitempty"` // 可选，自动估算
}

// TransactionResponse 交易响应
type TransactionResponse struct {
	ID             primitive.ObjectID `json:"id"`
	Hash           string             `json:"hash"`
	From           string             `json:"from"`
	To             string             `json:"to"`
	Value          string             `json:"value"`          // ETH格式的金额
	GasPrice       string             `json:"gas_price"`      // Gwei格式
	GasLimit       uint64             `json:"gas_limit"`
	GasUsed        uint64             `json:"gas_used"`
	Status         TransactionStatus  `json:"status"`
	Type           TransactionType    `json:"type"`
	Network        string             `json:"network"`
	BlockNumber    uint64             `json:"block_number"`
	TransactionFee string             `json:"transaction_fee"` // ETH格式
	CreatedAt      time.Time          `json:"created_at"`
	ConfirmedAt    *time.Time         `json:"confirmed_at,omitempty"`
}

// TransactionRepository 交易仓库接口
type TransactionRepository interface {
	Create(c context.Context, transaction *Transaction) error
	GetByID(c context.Context, id string) (Transaction, error)
	GetByHash(c context.Context, hash string) (Transaction, error)
	GetByUserID(c context.Context, userID string, limit, offset int) ([]Transaction, error)
	GetByWalletID(c context.Context, walletID string, limit, offset int) ([]Transaction, error)
	Update(c context.Context, transaction *Transaction) error
	UpdateStatus(c context.Context, hash string, status TransactionStatus) error
	GetPendingTransactions(c context.Context) ([]Transaction, error)
}

// TransactionUsecase 交易用例接口
type TransactionUsecase interface {
	SendTransaction(c context.Context, userID string, req *TransactionSendRequest) (*TransactionResponse, error)
	GetTransactions(c context.Context, userID string, limit, offset int) ([]TransactionResponse, error)
	GetTransaction(c context.Context, userID, transactionID string) (*TransactionResponse, error)
	GetTransactionByHash(c context.Context, hash string) (*TransactionResponse, error)
	UpdateTransactionStatus(c context.Context, hash string, status TransactionStatus) error
	EstimateGas(c context.Context, from, to, value string) (uint64, error)
	GetGasPrice(c context.Context, network string) (string, error)
}