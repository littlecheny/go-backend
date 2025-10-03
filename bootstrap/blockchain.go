package bootstrap

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/littlecheny/go-backend/domain"
	"github.com/littlecheny/go-backend/service"
)

func InitBlockchain(ctx context.Context, env *Env, app *Application) error {

}

func selectDefaultNetwork(env *Env) (network string) {
	network = env.DefaultNetwork
	if network == "" {
		network = "mainnet"
	}
	return network
}

func resolveNetworkRPC(env *Env, network string) (rpcURL string, isMock bool){
	switch network {
	case "mainnet":
		rpcURL = env.EthereumMainnetRPC
	case "sepolia":
		rpcURL = env.EthereumSepoliaRPC
	case "goerli":
		rpcURL = env.EthereumGoerliRPC
	default:
		isMock = true
		rpcURL = "http://localhost:8545"
	}
	return rpcURL, isMock
}

func connectEthereum(rpcURL string) (svc domain.EthereumService, err error) {
	//使用以太坊服务建立客户端连接（利用service/ethereum_service.go)
	ethClient, err := ethclient.DialContext( rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum node: %v", err)
	}
	svc = service.NewEthereumService(ethClient)
	return svc, nil
}

func healthCheckEthereum(ctx context.Context, svc service.EthereumService) (err error) {
	// 检查以太坊节点是否可访问
	_, err = svc.Client().BlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("failed to get block number: %v", err)
	}
	return nil
}

func regi
