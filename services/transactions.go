package services

import (
	"oasisTracker/services/render"
	"oasisTracker/smodels"
)

func (s *ServiceFacade) GetTransactionsList(params smodels.TransactionsParams) ([]smodels.Transaction, uint64, error) {
	count, err := s.dao.GetTransactionsCount(params)
	if err != nil {
		return nil, 0, err
	}

	txs, err := s.dao.GetTransactionsList(params)
	if err != nil {
		return nil, 0, err
	}

	return render.Transactions(txs), count, nil
}
