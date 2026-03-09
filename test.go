package contracts

import (
	"fmt"
	"log"

	"your_project_path/contracts" // 引入刚才生成的包

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// 1. 连接到 Sepolia 节点 (使用你的 Alchemy/Infura URL)
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/你的API_KEY")
	if err != nil {
		log.Fatal(err)
	}

	// 2. 你的合约地址 (刚才部署拿到的那个)
	contractAddr := common.HexToAddress("0x你的合约地址")

	// 3. 实例化合约对象
	instance, err := contracts.NewGreeter(contractAddr, client)
	if err != nil {
		log.Fatal(err)
	}

	// 4. 调用只读方法：Greet()
	// 注意：只读方法不需要消耗 Gas，不需要签名
	msg, err := instance.Greet(nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("合约返回的问候语: %s\n", msg)
}
