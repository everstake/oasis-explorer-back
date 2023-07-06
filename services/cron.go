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
	"oasisTracker/smodels"
	"time"
)

func (s *ServiceFacade) AddToCron(cron *gron.Cron, cfg conf.Config, dao *dao.DaoImpl) {

	if cfg.Cron.ParseValidatorsRegisterInterval > 0 {
		dur := time.Duration(cfg.Cron.ParseValidatorsRegisterInterval) * time.Minute
		log.Info("Scheduling counter saver every", zap.Duration("dur", dur))
		cron.AddFunc(gron.Every(dur), func() {
			log.Info("Start ParseValidatorsRegisterInterval")
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

			log.Info("End ParseValidatorsRegisterInterval")
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

	// todo delete some func
	//err = s.MigrateBlocks()
	//if err != nil {
	//	log.Error("MigrateBlocks failed:", zap.Error(err))
	//	return
	//}
	//
	//err = s.MigrateValidators()
	//if err != nil {
	//	log.Error("MigrateValidators failed:", zap.Error(err))
	//	return
	//}

	err = s.SyncBlocksStats()
	if err != nil {
		log.Error("SyncBlocksStats failed:", zap.Error(err))
		return
	}
}

func (s *ServiceFacade) CheckDelay() error {
	block, err := s.dao.GetLastBlock()
	if err != nil {
		return fmt.Errorf("dao.GetLastBlock: %v", err)
	}

	if block.CreatedAt.Before(time.Now().Add(-time.Minute * 7)) {
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

func (s *ServiceFacade) MigrateValidators() error {
	log.Info("MigrateValidators start")
	count, err := s.dao.GetValidatorsCount(smodels.ValidatorParams{})
	if err != nil {
		return fmt.Errorf("dao.GetValidatorsCount: %v", err)
	}

	for i := uint64(0); i < count; {
		validators, err := s.dao.GetValidatorsList(smodels.ValidatorParams{
			CommonParams: smodels.CommonParams{
				Limit:  200,
				Offset: i,
			},
		})
		if err != nil {
			return fmt.Errorf("dao.GetValidatorsList: %v", err)
		}

		err = s.pDao.MigrateValidatorsInfo(validators)
		if err != nil {
			return fmt.Errorf("pDao.MigrateValidatorsInfo: %v", err)
		}

		i += 200
		time.Sleep(time.Second)
	}
	log.Info("MigrateValidators done")

	return nil
}

func (s *ServiceFacade) MigrateBlocks() error {
	log.Info("MigrateBlocks start")

	bCount, err := s.dao.BlocksCount(smodels.BlockParams{})
	if err != nil {
		return fmt.Errorf("dao.BlocksCount: %v", err)
	}

	log.Info(fmt.Sprintf("blocks count: %d", bCount))
	limit := uint64(10000)
	for {
		offset, err := s.pDao.GetBlocksMigrationOffset()
		if err != nil {
			return fmt.Errorf("pDao.GetBlocksMigrationOffset: %v", err)
		}

		if offset >= bCount {
			break
		}

		log.Info(fmt.Sprintf("offset = %d", offset))
		blocks, err := s.dao.GetBlocksList(smodels.BlockParams{
			CommonParams: smodels.CommonParams{
				Limit:  limit,
				Offset: offset,
			},
		})
		if err != nil {
			return fmt.Errorf("dao.GetBlocksList: %v", err)
		}

		if len(blocks) > 1 {
			log.Info(fmt.Sprintf("inserting from %d to %d", blocks[0].Height, blocks[len(blocks)-1].Height))

			err = s.D.CreateBlocks(blocks)
			if err != nil {
				return fmt.Errorf("D.CreateBlocks: %v", err)
			}

			err = s.D.SaveBlocks(blocks)
			if err != nil {
				return fmt.Errorf("D.SaveBlocks: %v", err)
			}

			offset += uint64(len(blocks))

			log.Info(fmt.Sprintf("save offset: %d", offset))
			err = s.pDao.UpdateBlocksMigrationOffset(offset)
			if err != nil {
				return fmt.Errorf("pDao.UpdateBlocksMigrationOffset: %v", err)
			}
		} else {
			log.Info(fmt.Sprintf("Len blocks: %d", len(blocks)))
		}
	}
	log.Info("MigrateBlocks end")

	return nil
}
