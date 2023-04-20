package services

import "C"
import (
	"fmt"
	"github.com/oasisprotocol/oasis-core/go/common/grpc"
	stakingAPI "github.com/oasisprotocol/oasis-core/go/staking/api"
	"github.com/roylee0704/gron"
	"go.uber.org/zap"
	grpcCommon "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/google"
	"oasisTracker/common/log"
	"oasisTracker/common/modules"
	"oasisTracker/conf"
	"oasisTracker/dao"
	"oasisTracker/services/metaregistry"
	"oasisTracker/services/scanners"
	"time"
)

func (s *ServiceFacade) AddToCron(cron *gron.Cron, cfg conf.Config, dao *dao.DaoImpl) {

	if cfg.Cron.ParseValidatorsRegisterInterval > 0 {
		dur := time.Duration(cfg.Cron.ParseValidatorsRegisterInterval) * time.Minute
		log.Info("Scheduling counter saver every", zap.Duration("dur", dur))
		cron.AddFunc(gron.Every(dur), func() {
			log.Info("Start")
			credentials := google.NewDefaultCredentials().TransportCredentials()
			grpcConn, err := grpc.Dial(cfg.Scanner.NodeConfig, grpcCommon.WithTransportCredentials(credentials))
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
		log.Info("no scheduling counter due to missing ParseValidatorRegisterInterval in config")
	}

	dur := time.Minute * 10
	log.Info("Scheduling delay checker every", zap.Duration("dur", dur))
	cron.AddFunc(gron.Every(dur), func() {
		log.Info("Start")
		err := s.CheckDelay()
		if err != nil {
			log.Error("delay checker failed:", zap.Error(err))
			return
		}

	})
}

func (s *ServiceFacade) CheckDelay() error {
	block, err := s.dao.GetLastBlock()
	if err != nil {
		return fmt.Errorf("dao.GetLastBlock: %v", err)
	}

	if block.CreatedAt.Before(time.Now().Add(-time.Minute * 15)) {
		nw, err := scanners.NewWatcher(s.cfg, s.D)
		if err != nil {
			log.Fatal("Watcher.New", zap.Error(err))
		}
		//nm := scanners.NewManager(s.cfg, s.D)

		modules.Stop(s.Modules[2:])
		s.Modules = s.Modules[:2]

		s.Modules = append(s.Modules, nw)
		modules.Run(s.Modules[2:])
	}

	return nil
}
