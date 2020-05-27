package services

import (
	"oasisTracker/services/render"
	"oasisTracker/smodels"
)

func (s *ServiceFacade) GetTransactionsList(params smodels.TransactionsParams) ([]smodels.Transaction, error) {

	txs, err := s.dao.GetTransactionsList(params)
	if err != nil {
		return nil, err
	}

	return render.Transactions(txs), nil
}
