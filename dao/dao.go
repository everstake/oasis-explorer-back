package dao

import (
	"fmt"
	"oasisTracker/conf"
	"oasisTracker/dao/clickhouse"
	"oasisTracker/dao/mysql"
	"oasisTracker/dmodels"
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

	daoImpl struct {
		*clickhouse.Clickhouse
		*mysql.MysqlDAO
	}
)

func New(cfg conf.Config) (DAO, error) {
	m, err := mysql.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("mysql.New: %s", err.Error())
	}
	ch, err := clickhouse.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("clickhouse.New: %s", err.Error())
	}
	return daoImpl{
		Clickhouse: ch,
		MysqlDAO:   m,
	}, nil
}

func (d daoImpl) GetParserDAO() (interface{}, error) {
	return d.Clickhouse.GetChain(), nil
}
