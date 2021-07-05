package dao

import (
	"fmt"
	"oasisTracker/conf"
	"oasisTracker/dao/clickhouse"
	"oasisTracker/dao/postgres"
	"oasisTracker/dmodels"
	"oasisTracker/smodels"
)

type (
	DAO interface {
		TaskDAO
		GetParserDAO() (ParserDAO, error)
	}

	TaskDAO interface {
		CreateTask(task dmodels.Task) error
		GetTasks(bool) (tasks []dmodels.Task, err error)
		GetLastTask(title string) (task dmodels.Task, found bool, err error)
		UpdateTask(task dmodels.Task) error
	}

	ServiceDAO interface {
		GetAccountList(listParams smodels.AccountListParams) (resp []dmodels.AccountList, err error)
		AccountsCount() (count uint64, err error)
		GetAccountTiming(accountID string) (dmodels.AccountTime, error)

		GetLastBlock() (dmodels.Block, error)
		BlocksCount(params smodels.BlockParams) (count uint64, err error)
		GetBlocksList(params smodels.BlockParams) ([]dmodels.RowBlock, error)

		GetTransactionsList(params smodels.TransactionsParams) ([]dmodels.Transaction, error)

		//Charts
		GetChartsData(params smodels.ChartParams) ([]dmodels.ChartData, error)
		GetEscrowRatioChartData(params smodels.ChartParams) (resp []dmodels.ChartData, err error)
		GetBalanceChartData(accountID string, params smodels.ChartParams) (resp []dmodels.BalanceChartData, err error)
		GetTotalAccountsCountChartData(params smodels.ChartParams) (resp []dmodels.ChartData, err error)
		GetAvgBlockTimeChartData(params smodels.ChartParams) (resp []dmodels.ChartData, err error)
		GetFeeVolumeChartData(params smodels.ChartParams) (resp []dmodels.ChartData, err error)
		GetOperationsCountChartData(params smodels.ChartParams) (resp []dmodels.ChartData, err error)
		GetReclaimAmountChartData(params smodels.ChartParams) (resp []dmodels.ChartData, err error)

		GetTopEscrowAccounts(uint64) ([]dmodels.AccountBalance, error)
		GetLastDayTotalBalance() (dmodels.DayBalance, error)

		GetAccountValidatorInfo(accountID string) (resp dmodels.EntityNodesContainer, err error)
		GetEntity(string) (dmodels.EntityRegistryTransaction, error)
		GetEntityActiveDepositorsCount(accountID string) (count uint64, err error)

		GetValidatorsList(params smodels.ValidatorParams) (resp []dmodels.ValidatorView, err error)
		PublicValidatorsSearchList() (resp []dmodels.ValidatorView, err error)
		GetValidatorDayStats(string, smodels.ChartParams) (resp []dmodels.ValidatorStats, err error)
		GetValidatorDelegators(validatorID string, params smodels.CommonParams) ([]dmodels.Delegator, error)
		GetValidatorRewards(validatorID string, params smodels.CommonParams) ([]dmodels.Reward, error)
		GetValidatorRewardsStat(validatorID string) (resp dmodels.RewardsStat, err error)
	}

	ParserDAO interface {
		CreateBlocks(blocks []dmodels.Block) error
		CreateBlockSignatures(sig []dmodels.BlockSignature) error
		CreateAccountBalances(accounts []dmodels.AccountBalance) error
		CreateTransfers(transfers []dmodels.Transaction) error
		CreateRegisterNodeTransactions(txs []dmodels.NodeRegistryTransaction) error
		CreateRegisterEntityTransactions(txs []dmodels.EntityRegistryTransaction) error
		CreateRewards(txs []dmodels.Reward) error
		//To resync from last block
		GetLastBlock() (dmodels.Block, error)
	}

	DaoImpl struct {
		*clickhouse.Clickhouse
		*postgres.DAO
	}
)

func New(cfg conf.Config) (*DaoImpl, error) {
	m, err := postgres.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("postgres.New: %s", err.Error())
	}
	ch, err := clickhouse.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("clickhouse.New: %s", err.Error())
	}
	return &DaoImpl{
		Clickhouse: ch,
		DAO:        m,
	}, nil
}

func (d DaoImpl) GetParserDAO() (ParserDAO, error) {
	return d.Clickhouse, nil
}

func (d DaoImpl) GetServiceDAO() ServiceDAO {
	return d.Clickhouse
}
