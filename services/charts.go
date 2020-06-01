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
