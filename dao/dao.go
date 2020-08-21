package dao

import (
	"fmt"
	"oasisTracker/conf"
	"oasisTracker/dao/clickhouse"
	"oasisTracker/dao/mysql"
	"oasisTracker/dmodels"
	"oasisTracker/smodels"
)

type (
	DAO interface {
		MySql
		GetParserDAO() (ParserDAO, error)
	}
	MySql interface {
		CreateTask(task dmodels.Task) error
		GetTasks(bool) (tasks []dmodels.Task, err error)
		GetLastTask() (task dmodels.Task, found bool, err error)
		UpdateTask(task dmodels.Task) error
	}

	ServiceDAO interface {
		GetAccountList(listParams smodels.AccountListParams) (resp []dmodels.AccountList, err error)
		GetAccountTiming(accountID string) (dmodels.AccountTime, error)

		GetLastBlock() (dmodels.Block, error)
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
		GetEntityActiveDepositorsCount(accountID string) (count uint64, err error)

		GetValidatorsList(params smodels.ValidatorParams) (resp []dmodels.Validator, err error)
		PublicValidatorsSearchList() (resp []dmodels.Validator, err error)
		GetValidatorDayStats(string, smodels.ChartParams) (resp []dmodels.ValidatorStats, err error)
		GetValidatorDelegators(validatorID string, params smodels.CommonParams) ([]dmodels.Delegator, error)
	}

	ParserDAO interface {
		CreateBlocks(blocks []dmodels.Block) error
		CreateBlockSignatures(sig []dmodels.BlockSignature) error
		CreateAccountBalances(accounts []dmodels.AccountBalance) error
		CreateTransfers(transfers []dmodels.Transaction) error
		CreateRegisterNodeTransactions(txs []dmodels.NodeRegistryTransaction) error
		CreateRegisterEntityTransactions(txs []dmodels.EntityRegistryTransaction) error
		CreateRewards(txs []dmodels.Reward) error
	}

	daoImpl struct {
		*clickhouse.Clickhouse
		*mysql.MysqlDAO
	}
)

func New(cfg conf.Config) (*daoImpl, error) {
	m, err := mysql.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("mysql.New: %s", err.Error())
	}
	ch, err := clickhouse.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("clickhouse.New: %s", err.Error())
	}
	return &daoImpl{
		Clickhouse: ch,
		MysqlDAO:   m,
	}, nil
}

func (d daoImpl) GetParserDAO() (ParserDAO, error) {
	return d.Clickhouse, nil
}

func (d daoImpl) GetServiceDAO() ServiceDAO {
	return d.Clickhouse
}
