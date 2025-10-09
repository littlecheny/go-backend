package bootstrap

import (
	"context"

	"github.com/littlecheny/go-backend/logger"
	"github.com/littlecheny/go-backend/services"
)

func InitBlockchain(ctx context.Context, env *Env, app *Application) error {

	svc, err := services.NewEthereumSepoliaService(env.EthereumSepoliaRPCURL, env.WalletPrivateKeyHex)
	if err != nil {
		return err
	}
	app.EthereumService = svc
	logger.Log.WithContext(ctx).Infof("Network=%s", env.DefaultNetwork)
	//打印结构化日志
	logger.Log.WithContext(ctx).Info("EthereumService initialized")
	return nil
}
