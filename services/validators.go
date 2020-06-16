package services

import (
	"oasisTracker/services/render"
	"oasisTracker/smodels"
)

func (s *ServiceFacade) GetValidatorList(listParams smodels.ValidatorParams) ([]smodels.Validator, error) {
	blk, err := s.dao.GetLastBlock()
	if err != nil {
		return nil, err
	}

	resp, err := s.dao.GetValidatorsList(listParams)
	if err != nil {
		return nil, err
	}

	for i := range resp {
		//Availability Score
		availabilityScore := resp[i].SignaturesCount
		if resp[i].BlocksCount > 0 {
			//Temp without proposed stat
			availabilityPercent := float64(resp[i].BlocksCount) / float64(resp[i].BlocksCount)
			availabilityScore += uint64(availabilityPercent * float64(blk.Height-resp[i].StartBlockLevel))
		}

		if !resp[i].IsActive {
			resp[i].Status = smodels.StatusInActive
			continue
		}
		resp[i].Status = smodels.StatusActive
	}

	return render.ValidatorsList(resp), nil
}
