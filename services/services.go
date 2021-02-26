package services

import (
	"github.com/oasisprotocol/oasis-core/go/common/grpc"
	"github.com/oasisprotocol/oasis-core/go/staking/api"
	"github.com/patrickmn/go-cache"
	grpcCommon "google.golang.org/grpc"
	"oasisTracker/conf"
	"oasisTracker/dao"
	"oasisTracker/services/cmc"
	"oasisTracker/smodels"
	"time"
)

type (
	Service interface {
		GetInfo() (smodels.Info, error)
		GetBlockList(params smodels.BlockParams) ([]smodels.Block, error)
		GetTransactionsList(params smodels.TransactionsParams) ([]smodels.Transaction, error)
		GetAccountInfo(accountID string) (smodels.Account, error)
		GetAccountList(listParams smodels.AccountListParams) ([]smodels.AccountList, error)
		GetValidatorInfo(string) (smodels.Validator, error)
		GetValidatorList(listParams smodels.ValidatorParams) ([]smodels.Validator, error)
		GetPublicValidatorsSearchList() ([]smodels.ValidatorEntity, error)
		GetValidatorDelegators(validatorID string, params smodels.CommonParams) ([]smodels.Delegator, error)
		GetValidatorBlocks(validatorID string, params smodels.CommonParams) ([]smodels.Block, error)
		GetValidatorRewards(validatorID string, params smodels.CommonParams) ([]smodels.Reward, error)
		GetValidatorRewardsStat(validatorID string) (stat smodels.RewardStat, err error)

		GetChartData(params smodels.ChartParams) ([]smodels.ChartData, error)
		GetBalanceChartData(accountID string, params smodels.ChartParams) ([]smodels.BalanceChartData, error)
		GetEscrowRatioChartData(params smodels.ChartParams) ([]smodels.ChartData, error)
		GetValidatorStatsChartData(accountID string, params smodels.ChartParams) ([]smodels.ValidatorStats, error)
		GetTotalAccountsCountChartData(params smodels.ChartParams) ([]smodels.ChartData, error)
		GetAvgBlockTimeChartData(params smodels.ChartParams) ([]smodels.ChartData, error)
		GetFeeVolumeChartData(params smodels.ChartParams) ([]smodels.ChartData, error)
		GetOperationsCountChartData(params smodels.ChartParams) ([]smodels.ChartData, error)
		GetReclaimAmountChartData(params smodels.ChartParams) ([]smodels.ChartData, error)
		GetTopEscrowRatioChart(params smodels.CommonParams) (resp []smodels.TopEscrowRatioChart, err error)
	}

	ServiceFacade struct {
		cfg                conf.Config
		dao                dao.ServiceDAO
		nodeAPI            api.Backend
		cache              *cache.Cache
		marketDataProvider cmc.MarketDataProvider
	}
)

const (
	topEscrowCacheKey = "top_escrow_percent"
	cacheTTL          = 1 * time.Minute
)

func NewService(cfg conf.Config, dao dao.ServiceDAO) *ServiceFacade {
	grpcConn, err := grpc.Dial(cfg.Scanner.NodeConfig, grpcCommon.WithInsecure())
	if err != nil {
		return nil
	}

	sAPI := api.NewStakingClient(grpcConn)

	return &ServiceFacade{
		cfg:                cfg,
		dao:                dao,
		nodeAPI:            sAPI,
		cache:              cache.New(cacheTTL, cacheTTL),
		marketDataProvider: cmc.NewCoinGecko(),
	}
}
