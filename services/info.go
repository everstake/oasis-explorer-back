package services

import (
	"go.uber.org/zap"
	"oasisTracker/common/log"
	"oasisTracker/smodels"
)

func (s *ServiceFacade) GetInfo() (info smodels.Info, err error) {

	block, err := s.dao.GetLastBlock()
	if err != nil {
		return info, err
	}

	ratio, isFound := s.cache.Get(topEscrowCacheKey)
	if !isFound {
		ratio, err = s.getTopEscrowRatio()
		if err != nil {
			log.Error("failed to get staking ratio: ", zap.Error(err))
		}
		s.cache.Set(topEscrowCacheKey, ratio, cacheTTL)
	}

	return smodels.Info{
		Height:    block.Height,
		TopEscrow: ratio.(float64),
	}, nil
}

func (s *ServiceFacade) getTopEscrowRatio() (float64, error) {
	totalBalance, err := s.dao.GetLastDayTotalBalance()
	if err != nil {
		return 0, err
	}

	topAccounts, err := s.dao.GetTopEscrowAccounts(20)
	if err != nil {
		return 0, err
	}

	var topAccountsEscrow uint64
	for i := range topAccounts {
		topAccountsEscrow += topAccounts[i].EscrowBalanceActive
	}

	return float64(topAccountsEscrow) / float64(totalBalance.EscrowBalanceActive) * 100, nil
}
