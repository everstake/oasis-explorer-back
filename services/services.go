package services

import (
	"google.golang.org/grpc/credentials/google"
	"oasisTracker/conf"
	"oasisTracker/dao"
	"oasisTracker/services/cmc"
	"oasisTracker/smodels"
	"time"

	"github.com/oasisprotocol/oasis-core/go/common/grpc"
	"github.com/oasisprotocol/oasis-core/go/staking/api"
	"github.com/patrickmn/go-cache"
	grpcCommon "google.golang.org/grpc"
)

type (
	Service interface {
		GetInfo() (smodels.Info, error)
		GetBlockList(params smodels.BlockParams) ([]smodels.Block, uint64, error)
		GetTransactionsList(params smodels.TransactionsParams) ([]smodels.Transaction, uint64, error)
		GetAccountInfo(accountID string) (smodels.Account, error)
		GetAccountList(listParams smodels.AccountListParams) ([]smodels.AccountList, uint64, error)
		GetAccountRewards(accountID string, params smodels.CommonParams) ([]smodels.Reward, error)
		GetAccountRewardsStat(validatorID string) (stat smodels.RewardStat, err error)

		GetValidatorInfo(string) (smodels.Validator, error)
		GetValidatorList(listParams smodels.ValidatorParams) ([]smodels.Validator, uint64, error)
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
		genesisHeight      uint64
	}
)

const (
	topEscrowCacheKey = "top_escrow_percent"
	cacheTTL          = 1 * time.Minute
)

func NewService(cfg conf.Config, dao dao.ServiceDAO, genStartBlock uint64) *ServiceFacade {
	credentials := google.NewDefaultCredentials().TransportCredentials()
	grpcConn, err := grpc.Dial(cfg.Scanner.NodeConfig, grpcCommon.WithTransportCredentials(credentials))
	if err != nil {
		return nil
	}

	sAPI := api.NewStakingClient(grpcConn)

	return &ServiceFacade{
		cfg:                cfg,
		dao:                dao,
		nodeAPI:            sAPI,
		genesisHeight:      genStartBlock,
		cache:              cache.New(cacheTTL, cacheTTL),
		marketDataProvider: cmc.NewCoinGecko(),
	}
}
