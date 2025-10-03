package domain

import (
	"context"
	"time"
)

// NetworkConfig 网络配置
type NetworkConfig struct {
	Name     string `json:"name"`
	ChainID  int64  `json:"chain_id"`
	RPC      string `json:"rpc"`
	Explorer string `json:"explorer"`
	Symbol   string `json:"symbol"`
	Decimals int    `json:"decimals"`
}

// BlockchainStatus 区块链状态
type BlockchainStatus struct {
	Network       string    `json:"network"`
	LatestBlock   uint64    `json:"latest_block"`
	GasPrice      string    `json:"gas_price"`      // Gwei
	PeerCount     int       `json:"peer_count"`
	Syncing       bool      `json:"syncing"`
	LastUpdated   time.Time `json:"last_updated"`
}

// BlockInfo 区块信息
type BlockInfo struct {
	Number       uint64    `json:"number"`
	Hash         string    `json:"hash"`
	ParentHash   string    `json:"parent_hash"`
	Timestamp    time.Time `json:"timestamp"`
	Transactions int       `json:"transactions"`
	GasUsed      uint64    `json:"gas_used"`
	GasLimit     uint64    `json:"gas_limit"`
}

// EthereumService 以太坊服务接口
type EthereumService interface {
	// 网络连接
	Connect(rpcURL string) error
	Disconnect() error
	IsConnected() bool
	
	// 账户管理
	CreateAccount() (address, privateKey, mnemonic string, err error)
	ImportAccount(mnemonic string) (address, privateKey string, err error)
	GetBalance(address string) (string, error)
	
	// 交易操作
	SendTransaction(from, to, privateKey, value string, gasPrice *string) (string, error)
	GetTransaction(hash string) (*TransactionResponse, error)
	EstimateGas(from, to, value string) (uint64, error)
	GetGasPrice() (string, error)
	GetNonce(address string) (uint64, error)
	
	// 区块链信息
	GetLatestBlock() (*BlockInfo, error)
	GetBlockByNumber(number uint64) (*BlockInfo, error)
	GetNetworkID() (int64, error)
	
	// 监控
	SubscribeNewHeads(ctx context.Context) (<-chan *BlockInfo, error)
	WatchTransactions(ctx context.Context, addresses []string) (<-chan *TransactionResponse, error)
}

// RedisService Redis服务接口
type RedisService interface {
	// 基础操作
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
	Del(key string) error
	Exists(key string) (bool, error)
	
	// 缓存操作
	SetBalance(address string, balance string, expiration time.Duration) error
	GetBalance(address string) (string, error)
	SetGasPrice(network string, gasPrice string, expiration time.Duration) error
	GetGasPrice(network string) (string, error)
	
	// 交易缓存
	SetTransaction(hash string, tx *TransactionResponse, expiration time.Duration) error
	GetTransaction(hash string) (*TransactionResponse, error)
	
	// 区块缓存
	SetBlock(number uint64, block *BlockInfo, expiration time.Duration) error
	GetBlock(number uint64) (*BlockInfo, error)
}

// BlockchainUsecase 区块链用例接口
type BlockchainUsecase interface {
	GetNetworkStatus(c context.Context, network string) (*BlockchainStatus, error)
	GetLatestBlocks(c context.Context, network string, limit int) ([]BlockInfo, error)
	GetBlock(c context.Context, network string, number uint64) (*BlockInfo, error)
	GetSupportedNetworks(c context.Context) ([]NetworkConfig, error)
	SwitchNetwork(c context.Context, network string) error
}