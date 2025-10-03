package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionWallet = "wallets"
	CollectionWalletPrivateData = "wallet_private_data"
)

// WalletType 钱包类型
type WalletType string

const (
	WalletTypeHD       WalletType = "hd"       // HD钱包
	WalletTypeImported WalletType = "imported" // 导入钱包
	WalletTypeWatchOnly WalletType = "watch_only" // 只读钱包
)

// WalletStatus 钱包状态
type WalletStatus string

const (
	WalletStatusActive   WalletStatus = "active"   // 活跃
	WalletStatusInactive WalletStatus = "inactive" // 非活跃
	WalletStatusDeleted  WalletStatus = "deleted"  // 已删除
)

// Wallet 钱包模型
type Wallet struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`
	Name        string             `bson:"name" json:"name"`
	Address     string             `bson:"address" json:"address"`
	Type        WalletType         `bson:"type" json:"type"`
	Status      WalletStatus       `bson:"status" json:"status"`
	Network     string             `bson:"network" json:"network"`     // mainnet, sepolia, goerli
	Balance     string             `bson:"balance" json:"balance"`     // 余额 (ETH)
	BalanceUSD  string             `bson:"balance_usd" json:"balance_usd"` // USD余额
	IsDefault   bool               `bson:"is_default" json:"is_default"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// WalletPrivateData 钱包私有数据（加密存储）
type WalletPrivateData struct {
	ID              primitive.ObjectID `bson:"_id" json:"id"`
	WalletID        primitive.ObjectID `bson:"wallet_id" json:"wallet_id"`
	EncryptedKey    string             `bson:"encrypted_key" json:"-"`       // 加密的私钥
	EncryptedMnemonic string           `bson:"encrypted_mnemonic" json:"-"`  // 加密的助记词
	KeyDerivationPath string           `bson:"key_derivation_path" json:"-"` // HD钱包派生路径
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
}

// WalletCreateRequest 创建钱包请求
type WalletCreateRequest struct {
	Name     string     `json:"name" binding:"required"`
	Type     WalletType `json:"type" binding:"required"`
	Network  string     `json:"network" binding:"required"`
	Password string     `json:"password" binding:"required,min=8"`
	Mnemonic string     `json:"mnemonic,omitempty"` // 导入钱包时需要
}

// WalletImportRequest 导入钱包请求
type WalletImportRequest struct {
	Name     string `json:"name" binding:"required"`
	Network  string `json:"network" binding:"required"`
	Mnemonic string `json:"mnemonic" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

// WalletResponse 钱包响应
type WalletResponse struct {
	ID         primitive.ObjectID `json:"id"`
	Name       string             `json:"name"`
	Address    string             `json:"address"`
	Type       WalletType         `json:"type"`
	Status     WalletStatus       `json:"status"`
	Network    string             `json:"network"`
	Balance    string             `json:"balance"`
	BalanceUSD string             `json:"balance_usd"`
	IsDefault  bool               `json:"is_default"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
}

// WalletBalanceResponse 钱包余额响应
type WalletBalanceResponse struct {
	Address    string    `json:"address"`
	Balance    string    `json:"balance"`     // ETH
	BalanceUSD string    `json:"balance_usd"` // USD
	Network    string    `json:"network"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// WalletListResponse 钱包列表响应
type WalletListResponse struct {
	Wallets    []WalletResponse `json:"wallets"`
	TotalCount int              `json:"total_count"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
}

// WalletStatsResponse 钱包统计响应
type WalletStatsResponse struct {
	TotalWallets    int            `json:"total_wallets"`
	TotalBalance    string         `json:"total_balance"`     // ETH
	TotalBalanceUSD string         `json:"total_balance_usd"` // USD
	Networks        []NetworkStats `json:"networks"`
}

// NetworkStats 网络统计
type NetworkStats struct {
	Network         string `json:"network"`
	WalletCount     int    `json:"wallet_count"`
	TotalBalance    string `json:"total_balance"`
	TotalBalanceUSD string `json:"total_balance_usd"`
}

// WalletRepository 钱包仓库接口
type WalletRepository interface {
	// 基础CRUD
	Create(ctx context.Context, wallet *Wallet) error
	GetByID(ctx context.Context, id string) (*Wallet, error)
	GetByAddress(ctx context.Context, address string) (*Wallet, error)
	GetByUserID(ctx context.Context, userID string, page, pageSize int) ([]Wallet, int, error)
	Update(ctx context.Context, wallet *Wallet) error
	Delete(ctx context.Context, id string) error
	
	// 私有数据操作
	CreatePrivateData(ctx context.Context, data *WalletPrivateData) error
	GetPrivateData(ctx context.Context, walletID string) (*WalletPrivateData, error)
	UpdatePrivateData(ctx context.Context, data *WalletPrivateData) error
	DeletePrivateData(ctx context.Context, walletID string) error
	
	// 查询操作
	GetDefaultWallet(ctx context.Context, userID string, network string) (*Wallet, error)
	SetDefaultWallet(ctx context.Context, userID string, walletID string) error
	GetWalletsByNetwork(ctx context.Context, userID string, network string) ([]Wallet, error)
	UpdateBalance(ctx context.Context, walletID string, balance string, balanceUSD string) error
	
	// 统计操作
	GetUserWalletCount(ctx context.Context, userID string) (int, error)
	GetTotalBalance(ctx context.Context, userID string, network string) (string, error)
}

// WalletUsecase 钱包用例接口
type WalletUsecase interface {
	// 钱包管理
	CreateWallet(ctx context.Context, userID string, req *WalletCreateRequest) (*WalletResponse, error)
	ImportWallet(ctx context.Context, userID string, req *WalletImportRequest) (*WalletResponse, error)
	GetWallet(ctx context.Context, userID string, walletID string) (*WalletResponse, error)
	GetWallets(ctx context.Context, userID string, page, pageSize int) (*WalletListResponse, error)
	UpdateWallet(ctx context.Context, userID string, walletID string, name string) (*WalletResponse, error)
	DeleteWallet(ctx context.Context, userID string, walletID string) error
	
	// 余额管理
	GetBalance(ctx context.Context, userID string, walletID string) (*WalletBalanceResponse, error)
	RefreshBalance(ctx context.Context, userID string, walletID string) (*WalletBalanceResponse, error)
	RefreshAllBalances(ctx context.Context, userID string) error
	
	// 默认钱包
	SetDefaultWallet(ctx context.Context, userID string, walletID string) error
	GetDefaultWallet(ctx context.Context, userID string, network string) (*WalletResponse, error)
	
	// 网络管理
	GetWalletsByNetwork(ctx context.Context, userID string, network string) ([]WalletResponse, error)
	SwitchNetwork(ctx context.Context, userID string, walletID string, network string) error
	
	// 导出功能
	ExportPrivateKey(ctx context.Context, userID string, walletID string, password string) (string, error)
	ExportMnemonic(ctx context.Context, userID string, walletID string, password string) (string, error)
	
	// 统计信息
	GetWalletStats(ctx context.Context, userID string) (*WalletStatsResponse, error)
}

// CryptoService 加密服务接口
type CryptoService interface {
	// 加密/解密
	Encrypt(data string, key string) (string, error)
	Decrypt(encryptedData string, key string) (string, error)
	
	// 密钥派生
	DeriveKey(password string, salt string) (string, error)
	GenerateSalt() (string, error)
	
	// 哈希
	Hash(data string) string
	VerifyHash(data string, hash string) bool
}