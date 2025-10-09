package services

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/littlecheny/go-backend/domain"
)

type ethereumSepoliaService struct {
	client         *ethclient.Client
	privateKeyHex  string
}

func NewEthereumSepoliaService(rpcURL string, privateKeyHex string) (domain.EthereumService, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}
	return &ethereumSepoliaService{client: client, privateKeyHex: privateKeyHex}, nil
}

func (s *ethereumSepoliaService) SendTransaction(ctx context.Context, req *domain.TransactionSendRequest) (string, error) {
	if s.privateKeyHex == "" {
		return "", errors.New("private key not configured")
	}

	// 解析参数
	to := common.HexToAddress(req.To)
	value, ok := new(big.Int).SetString(req.Value, 10)
	if !ok {
		return "", ErrInvalidValue
	}
	gasLimit, ok := new(big.Int).SetString(req.GasLimit, 10)
	if !ok {
		return "", ErrInvalidGasLimit
	}
	gasPrice, ok := new(big.Int).SetString(req.GasPrice, 10)
	if !ok {
		return "", ErrInvalidGasPrice
	}

	// 私钥与 from 地址
	privKey, err := crypto.HexToECDSA(s.privateKeyHex)
	if err != nil {
		return "", err
	}
	from := crypto.PubkeyToAddress(privKey.PublicKey)

	// 获取 nonce
	nonce, err := s.client.PendingNonceAt(ctx, from)
	if err != nil {
		return "", err
	}

	// 构造交易
	tx := types.NewTransaction(nonce, to, value, gasLimit.Uint64(), gasPrice, nil)

	chainID, err := s.client.NetworkID(ctx)
	if err != nil {
		return "", err
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privKey)
	if err != nil {
		return "", err
	}

	// 广播交易
	if err := s.client.SendTransaction(ctx, signedTx); err != nil {
		return "", err
	}

	return signedTx.Hash().Hex(), nil
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

// 简化错误声明（可迁移至独立 errors 包）
var (
	ErrInvalidValue    = errors.New("invalid value")
	ErrInvalidGasLimit = errors.New("invalid gas limit")
	ErrInvalidGasPrice = errors.New("invalid gas price")
)
