package services

import (
	"github.com/oasislabs/oasis-core/go/common/grpc"
	"github.com/oasislabs/oasis-core/go/staking/api"
	grpcCommon "google.golang.org/grpc"
	"oasisTracker/conf"
	"oasisTracker/dao"
	"oasisTracker/smodels"
)

type (
	Service interface {
		GetBlockList(params smodels.BlockParams) ([]smodels.Block, error)
		GetTransactionsList(params smodels.TransactionsParams) ([]smodels.Transaction, error)
		GetAccountInfo(accountID string) (smodels.Account, error)
		GetChartData(params smodels.ChartParams) ([]smodels.ChartData, error)
	}

	ServiceFacade struct {
		cfg     conf.Config
		dao     dao.ServiceDAO
		nodeAPI api.Backend
	}
)

func NewService(cfg conf.Config, dao dao.ServiceDAO) *ServiceFacade {
	grpcConn, err := grpc.Dial(cfg.Scanner.NodeConfig, grpcCommon.WithInsecure())
	if err != nil {
		return nil
	}

	sAPI := api.NewStakingClient(grpcConn)

	return &ServiceFacade{
		cfg:     cfg,
		dao:     dao,
		nodeAPI: sAPI,
	}
}
