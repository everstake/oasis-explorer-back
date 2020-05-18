package services

import (
	"oasisTracker/conf"
	"oasisTracker/dao"
)

type (
	Service interface {
	}

	ServiceFacade struct {
		cfg conf.Config
		dao dao.DAO
	}
)

func NewService(cfg conf.Config, dao dao.DAO) *ServiceFacade {
	return &ServiceFacade{
		cfg: cfg,
		dao: dao,
	}
}
