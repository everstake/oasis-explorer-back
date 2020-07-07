package services

import (
	"oasisTracker/services/render"
	"oasisTracker/smodels"
)

func (s *ServiceFacade) GetChartData(params smodels.ChartParams) ([]smodels.ChartData, error) {

	data, err := s.dao.GetChartsData(params)
	if err != nil {
		return nil, err
	}

	return render.ChartData(data), nil
}

func (s *ServiceFacade) GetEscrowRatioChartData(params smodels.ChartParams) ([]smodels.ChartData, error) {

	data, err := s.dao.GetEscrowRatioChartData(params)
	if err != nil {
		return nil, err
	}

	return render.ChartData(data), nil
}

func (s *ServiceFacade) GetTotalAccountsCountChartData(params smodels.ChartParams) ([]smodels.ChartData, error) {

	data, err := s.dao.GetTotalAccountsCountChartData(params)
	if err != nil {
		return nil, err
	}

	return render.ChartData(data), nil
}

func (s *ServiceFacade) GetAvgBlockTimeChartData(params smodels.ChartParams) ([]smodels.ChartData, error) {

	data, err := s.dao.GetAvgBlockTimeChartData(params)
	if err != nil {
		return nil, err
	}

	return render.ChartData(data), nil
}

func (s *ServiceFacade) GetFeeVolumeChartData(params smodels.ChartParams) ([]smodels.ChartData, error) {

	data, err := s.dao.GetFeeVolumeChartData(params)
	if err != nil {
		return nil, err
	}

	return render.ChartData(data), nil
}

func (s *ServiceFacade) GetOperationsCountChartData(params smodels.ChartParams) ([]smodels.ChartData, error) {

	data, err := s.dao.GetOperationsCountChartData(params)
	if err != nil {
		return nil, err
	}

	return render.ChartData(data), nil
}

func (s *ServiceFacade) GetReclaimAmountChartData(params smodels.ChartParams) ([]smodels.ChartData, error) {

	data, err := s.dao.GetReclaimAmountChartData(params)
	if err != nil {
		return nil, err
	}

	return render.ChartData(data), nil
}

func (s *ServiceFacade) GetTopEscrowRatioChart(params smodels.CommonParams) (resp []smodels.TopEscrowRatioChart, err error) {
	totalBalance, err := s.dao.GetLastDayTotalBalance()
	if err != nil {
		return resp, err
	}

	topAccounts, err := s.dao.GetTopEscrowAccounts(params.Limit)
	if err != nil {
		return resp, err
	}

	for i := range topAccounts {

		resp = append(resp, smodels.TopEscrowRatioChart{
			AccountID:   topAccounts[i].Account,
			AccountName: topAccounts[i].AccountName,
			EsrowRatio:  float64(topAccounts[i].EscrowBalanceActive) / float64(totalBalance.EscrowBalanceActive),
		})

	}

	return resp, nil
}

func (s *ServiceFacade) GetBalanceChartData(accountID string, params smodels.ChartParams) ([]smodels.BalanceChartData, error) {

	data, err := s.dao.GetBalanceChartData(accountID, params)
	if err != nil {
		return nil, err
	}

	return render.BalanceChartData(data), nil
}
