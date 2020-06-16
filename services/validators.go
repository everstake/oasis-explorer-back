package services

import (
	"oasisTracker/services/render"
	"oasisTracker/smodels"
)

func (s *ServiceFacade) GetValidatorList(listParams smodels.ValidatorParams) ([]smodels.Validator, error) {

	resp, err := s.dao.GetValidatorsList(listParams)
	if err != nil {
		return nil, err
	}

	for i := range resp {
		if !resp[i].IsActive {
			resp[i].Status = smodels.StatusInActive
			continue
		}
		resp[i].Status = smodels.StatusActive
	}

	return render.ValidatorsList(resp), nil
}
