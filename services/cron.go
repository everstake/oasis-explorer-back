package services

import (
	"github.com/oasisprotocol/oasis-core/go/common/grpc"
	stakingAPI "github.com/oasisprotocol/oasis-core/go/staking/api"
	"github.com/roylee0704/gron"
	"go.uber.org/zap"
	grpcCommon "google.golang.org/grpc"
	"oasisTracker/common/log"
	"oasisTracker/conf"
	"oasisTracker/dao"
	"oasisTracker/services/metaregistry"
	"time"
)

func AddToCron(cron *gron.Cron, cfg conf.Config, dao *dao.DaoImpl) {

	if cfg.Cron.ParseValidatorsRegisterInterval > 0 {
		dur := time.Duration(cfg.Cron.ParseValidatorsRegisterInterval) * time.Minute
		log.Info("Sheduling counter saver every", zap.Duration("dur", dur))
		cron.AddFunc(gron.Every(dur), func() {
			log.Info("Start")
			grpcConn, err := grpc.Dial(cfg.Scanner.NodeConfig, grpcCommon.WithInsecure())
			if err != nil {
				log.Error("grpc.Dial failed:", zap.Error(err))
				return
			}

			defer grpcConn.Close()

			err = metaregistry.UpdatePublicValidators(dao, stakingAPI.NewStakingClient(grpcConn))
			if err != nil {
				log.Error("public validators update saver failed:", zap.Error(err))
				return
			}

		})
	} else {
		log.Info("no sheduling counter due to missing ParseValidatorRegisterInterval in config")
	}

}
