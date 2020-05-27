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
		GetParserDAO() (interface{}, error)
	}
	MySql interface {
		CreateTask(task dmodels.Task) error
		GetTasks(bool) (tasks []dmodels.Task, err error)
		GetLastTask() (task dmodels.Task, found bool, err error)
		UpdateTask(task dmodels.Task) error
	}

	ServiceDAO interface {
		GetAccountTiming(accountID string) (dmodels.AccountTime, error)

		GetBlocksList(params smodels.BlockParams) ([]dmodels.RowBlock, error)

		GetTransactionsList(params smodels.TransactionsParams) ([]dmodels.Transaction, error)

		GetChartsData(params smodels.ChartParams) ([]dmodels.ChartData, error)
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

func (d daoImpl) GetParserDAO() (interface{}, error) {
	return d.Clickhouse, nil
}

func (d daoImpl) GetServiceDAO() ServiceDAO {
	return d.Clickhouse
}
