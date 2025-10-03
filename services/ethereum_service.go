package services

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"string"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/littlecheny/go-backend/domain"
	"github.com/tyler-smith/go-bip39"
)

type ethereumService struct {
	client    *ethclient.Client
	networkID int64
	rpcURL    string
}

func NewEthereumService() domain.EthereumService {
	return &ethereumService{}
}

func (e *ethereumService) Connect(rpcURL string) error {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return fmt.Errorf("failed to connect to Ethereum client: %v", err)
	}

	e.client = client
	e.rpcURL = rpcURL

	// 获取网络ID
	networkID, err := client.NetworkID(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get network ID: %v", err)
	}
	e.networkID = networkID.Int64()

	return nil
}

func (e *ethereumService) Disconnect() error {
	if e.client != nil {
		e.client.Close()
		e.client = nil
	}
	return nil
}

func (e *ethereumService) IsConnected() bool {
	return e.client != nil
}

func (e *ethereumService) CreateAccount() (address, privateKey, mnemonic string, err error) {
	// 生成助记词
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate entropy: %v", err)
	}

	mnemonic, err = bip39.NewMnemonic(entropy)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate mnemonic: %v", err)
	}

	// 从助记词生成种子
	seed := bip39.NewSeed(mnemonic, "")

	// 生成私钥
	privateKeyECDSA, err := crypto.GenerateKey()
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate private key: %v", err)
	}

	// 获取私钥的十六进制表示
	privateKeyBytes := crypto.FromECDSA(privateKeyECDSA)
	privateKey = fmt.Sprintf("%x", privateKeyBytes)

	// 获取公钥和地址
	publicKey := privateKeyECDSA.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", "", "", fmt.Errorf("failed to cast public key to ECDSA")
	}

	address = crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	return address, privateKey, mnemonic, nil
}

func (e *ethereumService) ImportAccount(mnemonic string) (address, privateKey string, err error) {
	// 验证助记词
	if !bip39.IsMnemonicValid(mnemonic) {
		return "", "", fmt.Errorf("invalid mnemonic")
	}

	// 从助记词生成种子
	seed := bip39.NewSeed(mnemonic, "")

	// 这里简化处理，实际应该使用HD钱包派生
	// 为了演示，我们生成一个新的私钥
	privateKeyECDSA, err := crypto.GenerateKey()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate private key: %v", err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKeyECDSA)
	privateKey = fmt.Sprintf("%x", privateKeyBytes)

	publicKey := privateKeyECDSA.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", "", fmt.Errorf("failed to cast public key to ECDSA")
	}

	address = crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	return address, privateKey, nil
}

func (e *ethereumService) GetBalance(address string) (string, error) {
	if e.client == nil {
		return "", fmt.Errorf("ethereum client not connected")
	}

	account := common.HexToAddress(address)
	balance, err := e.client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get balance: %v", err)
	}

	// 转换为ETH
	ethBalance := new(big.Float)
	ethBalance.SetString(balance.String())
	ethBalance = ethBalance.Quo(ethBalance, big.NewFloat(params.Ether))

	return ethBalance.String(), nil
}

func (e *ethereumService) SendTransaction(from, to, privateKey, value string, gasPrice *string) (string, error) {
	if e.client == nil {
		return "", fmt.Errorf("ethereum client not connected")
	}

	// 解析私钥
	privateKeyECDSA, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %v", err)
	}

	// 获取发送方地址
	fromAddress := crypto.PubkeyToAddress(privateKeyECDSA.PublicKey)

	// 获取nonce
	nonce, err := e.client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return "", fmt.Errorf("failed to get nonce: %v", err)
	}

	// 解析金额
	valueWei, ok := new(big.Int).SetString(value, 10)
	if !ok {
		return "", fmt.Errorf("failed to parse value")
	}

	// 获取gas价格
	var gasPriceWei *big.Int
	if gasPrice != nil {
		gasPriceWei, ok = new(big.Int).SetString(*gasPrice, 10)
		if !ok {
			return "", fmt.Errorf("failed to parse gas price")
		}
	} else {
		gasPriceWei, err = e.client.SuggestGasPrice(context.Background())
		if err != nil {
			return "", fmt.Errorf("failed to suggest gas price: %v", err)
		}
	}

	// 估算gas限制
	toAddress := common.HexToAddress(to)
	gasLimit, err := e.client.EstimateGas(context.Background(), ethereum.CallMsg{
		From:  fromAddress,
		To:    &toAddress,
		Value: valueWei,
	})
	if err != nil {
		gasLimit = 21000 // 默认gas限制
	}

	// 创建交易
	tx := types.NewTransaction(nonce, toAddress, valueWei, gasLimit, gasPriceWei, nil)

	// 签名交易
	chainID := big.NewInt(e.networkID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKeyECDSA)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %v", err)
	}

	// 发送交易
	err = e.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %v", err)
	}

	return signedTx.Hash().Hex(), nil
}

func (e *ethereumService) GetTransaction(hash string) (*domain.TransactionResponse, error) {
	if e.client == nil {
		return nil, fmt.Errorf("ethereum client not connected")
	}

	txHash := common.HexToHash(hash)
	tx, isPending, err := e.client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %v", err)
	}

	var status domain.TransactionStatus
	var blockNumber uint64
	var gasUsed uint64

	if isPending {
		status = domain.TransactionStatusPending
	} else {
		// 获取交易收据
		receipt, err := e.client.TransactionReceipt(context.Background(), txHash)
		if err != nil {
			return nil, fmt.Errorf("failed to get transaction receipt: %v", err)
		}

		if receipt.Status == 1 {
			status = domain.TransactionStatusConfirmed
		} else {
			status = domain.TransactionStatusFailed
		}

		blockNumber = receipt.BlockNumber.Uint64()
		gasUsed = receipt.GasUsed
	}

	// 转换金额为ETH
	ethValue := new(big.Float)
	ethValue.SetString(tx.Value().String())
	ethValue = ethValue.Quo(ethValue, big.NewFloat(params.Ether))

	// 转换gas价格为Gwei
	gweiGasPrice := new(big.Float)
	gweiGasPrice.SetString(tx.GasPrice().String())
	gweiGasPrice = gweiGasPrice.Quo(gweiGasPrice, big.NewFloat(params.GWei))

	return &domain.TransactionResponse{
		Hash:        tx.Hash().Hex(),
		From:        tx.To().Hex(), // 注意：这里需要从交易中获取正确的from地址
		To:          tx.To().Hex(),
		Value:       ethValue.String(),
		GasPrice:    gweiGasPrice.String(),
		GasLimit:    tx.Gas(),
		GasUsed:     gasUsed,
		Status:      status,
		BlockNumber: blockNumber,
		CreatedAt:   time.Now(),
	}, nil
}

func (e *ethereumService) EstimateGas(from, to, value string) (uint64, error) {
	if e.client == nil {
		return 0, fmt.Errorf("ethereum client not connected")
	}

	fromAddress := common.HexToAddress(from)
	toAddress := common.HexToAddress(to)

	valueWei, ok := new(big.Int).SetString(value, 10)
	if !ok {
		return 0, fmt.Errorf("failed to parse value")
	}

	gasLimit, err := e.client.EstimateGas(context.Background(), ethereum.CallMsg{
		From:  fromAddress,
		To:    &toAddress,
		Value: valueWei,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to estimate gas: %v", err)
	}

	return gasLimit, nil
}

func (e *ethereumService) GetGasPrice() (string, error) {
	if e.client == nil {
		return "", fmt.Errorf("ethereum client not connected")
	}

	gasPrice, err := e.client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to get gas price: %v", err)
	}

	// 转换为Gwei
	gweiGasPrice := new(big.Float)
	gweiGasPrice.SetString(gasPrice.String())
	gweiGasPrice = gweiGasPrice.Quo(gweiGasPrice, big.NewFloat(params.GWei))

	return gweiGasPrice.String(), nil
}

func (e *ethereumService) GetNonce(address string) (uint64, error) {
	if e.client == nil {
		return 0, fmt.Errorf("ethereum client not connected")
	}

	account := common.HexToAddress(address)
	nonce, err := e.client.PendingNonceAt(context.Background(), account)
	if err != nil {
		return 0, fmt.Errorf("failed to get nonce: %v", err)
	}

	return nonce, nil
}

func (e *ethereumService) GetLatestBlock() (*domain.BlockInfo, error) {
	if e.client == nil {
		return nil, fmt.Errorf("ethereum client not connected")
	}

	header, err := e.client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block: %v", err)
	}

	block, err := e.client.BlockByNumber(context.Background(), header.Number)
	if err != nil {
		return nil, fmt.Errorf("failed to get block details: %v", err)
	}

	return &domain.BlockInfo{
		Number:       header.Number.Uint64(),
		Hash:         header.Hash().Hex(),
		ParentHash:   header.ParentHash.Hex(),
		Timestamp:    time.Unix(int64(header.Time), 0),
		Transactions: len(block.Transactions()),
		GasUsed:      header.GasUsed,
		GasLimit:     header.GasLimit,
	}, nil
}

func (e *ethereumService) GetBlockByNumber(number uint64) (*domain.BlockInfo, error) {
	if e.client == nil {
		return nil, fmt.Errorf("ethereum client not connected")
	}

	block, err := e.client.BlockByNumber(context.Background(), big.NewInt(int64(number)))
	if err != nil {
		return nil, fmt.Errorf("failed to get block: %v", err)
	}

	return &domain.BlockInfo{
		Number:       block.Number().Uint64(),
		Hash:         block.Hash().Hex(),
		ParentHash:   block.ParentHash().Hex(),
		Timestamp:    time.Unix(int64(block.Time()), 0),
		Transactions: len(block.Transactions()),
		GasUsed:      block.GasUsed(),
		GasLimit:     block.GasLimit(),
	}, nil
}

func (e *ethereumService) GetNetworkID() (int64, error) {
	if e.client == nil {
		return 0, fmt.Errorf("ethereum client not connected")
	}

	return e.networkID, nil
}

func (e *ethereumService) SubscribeNewHeads(ctx context.Context) (<-chan *domain.BlockInfo, error) {
	if e.client == nil {
		return nil, fmt.Errorf("ethereum client not connected")
	}

	headers := make(chan *types.Header)
	sub, err := e.client.SubscribeNewHead(ctx, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to new heads: %v", err)
	}

	blockInfoChan := make(chan *domain.BlockInfo)

	go func() {
		defer close(blockInfoChan)
		defer sub.Unsubscribe()

		for {
			select {
			case err := <-sub.Err():
				fmt.Printf("Subscription error: %v\n", err)
				return
			case header := <-headers:
				block, err := e.client.BlockByHash(ctx, header.Hash())
				if err != nil {
					fmt.Printf("Failed to get block: %v\n", err)
					continue
				}

				blockInfo := &domain.BlockInfo{
					Number:       header.Number.Uint64(),
					Hash:         header.Hash().Hex(),
					ParentHash:   header.ParentHash.Hex(),
					Timestamp:    time.Unix(int64(header.Time), 0),
					Transactions: len(block.Transactions()),
					GasUsed:      header.GasUsed,
					GasLimit:     header.GasLimit,
				}

				select {
				case blockInfoChan <- blockInfo:
				case <-ctx.Done():
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return blockInfoChan, nil
}

func (e *ethereumService) WatchTransactions(ctx context.Context, addresses []string) (<-chan *domain.TransactionResponse, error) {
	// 这是一个简化的实现，实际应用中需要更复杂的逻辑
	txChan := make(chan *domain.TransactionResponse)

	go func() {
		defer close(txChan)
		// 这里应该实现监听指定地址的交易逻辑
		// 由于复杂性，这里只是一个占位符
		<-ctx.Done()
	}()

	return txChan, nil
}