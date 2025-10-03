package bootstrap

import(
	"log"
	"github.com/spf13/viper"
)

type Env struct{
	AppEnv			string `mapstructure:"APP_ENV"`
	ServerAddress   string `mapstructure:"SERVER_ADDRESS"`
	ContextTimeout  int `mapstructure:"CONTEXT_TIMEOUT"`
	DBHost			string `mapstructure:"DB_HOST"`
	DBPort			string `mapstructure:"DB_PORT"`
	DBUser			string `mapstructure:"DB_USER"`
	DBPass			string `mapstructure:"DB_PASS"`
	DBName                 string `mapstructure:"DB_NAME"`
	AccessTokenExpiryHour  int    `mapstructure:"ACCESS_TOKEN_EXPIRY_HOUR"`
	RefreshTokenExpiryHour int    `mapstructure:"REFRESH_TOKEN_EXPIRY_HOUR"`
	AccessTokenSecret      string `mapstructure:"ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret     string `mapstructure:"REFRESH_TOKEN_SECRET"`
	
	// Redis配置
	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPort     string `mapstructure:"REDIS_PORT"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDB       int    `mapstructure:"REDIS_DB"`
	
	// 区块链配置
	EthereumMainnetRPC string `mapstructure:"ETHEREUM_MAINNET_RPC"`
	EthereumSepoliaRPC string `mapstructure:"ETHEREUM_SEPOLIA_RPC"`
	EthereumGoerliRPC  string `mapstructure:"ETHEREUM_GOERLI_RPC"`
	DefaultNetwork     string `mapstructure:"DEFAULT_NETWORK"`
	
	// 加密配置
	WalletEncryptionKey string `mapstructure:"WALLET_ENCRYPTION_KEY"`
}

func NewEnv() *Env {
	env := Env{}
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env: ", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	if env.AppEnv == "development" {
		log.Println("The App is running in development env")
	}

	return &env
}