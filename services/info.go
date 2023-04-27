package services

import (
	"go.uber.org/zap"
	"oasisTracker/common/log"
	"oasisTracker/smodels"
	"time"
)

const usdCurrency = "usd"
const getInfoEP = "/data/info"

func (s *ServiceFacade) GetInfo() (info smodels.Info, err error) {

	raw, ok, err := s.apiCache.Get(getInfoEP)
	if err != nil {
		return info, err
	}

	if !ok {
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

		marketInfo, err := s.marketDataProvider.GetOasisMarketData(usdCurrency, s.cfg.CMC.Key)
		if err != nil {
			return info, err
		}

		info = smodels.Info{
			Height:      block.Height,
			TopEscrow:   ratio.(float64),
			Price:       marketInfo.GetPrice(),
			PriceChange: marketInfo.GetPriceChange(),
			MarketCap:   marketInfo.GetMarketCap(),
			Volume:      marketInfo.GetVolume(),
			Supply:      marketInfo.GetSupply(),
		}

		err = s.apiCache.Save(getInfoEP, info, time.Second*5)
		if err != nil {
			return info, err
		}

		return info, nil
	} else {
		return raw.(smodels.Info), err
	}
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
