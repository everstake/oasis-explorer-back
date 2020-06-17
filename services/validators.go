package services

import (
	"fmt"
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

		resp[i].AvailabilityScore = calcAvailabilityScore(resp[i].BlocksCount, resp[i].SignaturesCount, resp[i].StartBlockLevel, blk.Height)

		if !resp[i].IsActive {
			resp[i].Status = smodels.StatusInActive
			continue
		}
		resp[i].Status = smodels.StatusActive
	}

	return render.ValidatorsList(resp), nil
}

func (s *ServiceFacade) GetValidatorStatsChartData(accountID string, params smodels.ChartParams) ([]smodels.ValidatorStats, error) {

	validators, err := s.dao.GetValidatorsList(smodels.ValidatorParams{
		CommonParams: smodels.CommonParams{Limit: 1},
		ValidatorID:  "PhDiz71pnE2XeMpfZzcvbpDRZZkM4Bw0iZFcr3LtB9Q=",
	})

	if err != nil {
		return nil, err
	}

	if len(validators) == 0 {
		return nil, fmt.Errorf("Not found")
	}

	stats, err := s.dao.GetValidatorDayStats(validators[0].ConsensusAddress, params)
	if err != nil {
		return nil, err
	}

	for i := range stats {
		stats[i].AvailabilityScore = calcAvailabilityScore(0, 0, validators[0].StartBlockLevel, 0)
	}

	return render.ValidatorStatList(stats), nil
}

func calcAvailabilityScore(blocks, signatures, nodeRegisterBlock, currentHeight uint64) uint64 {

	availabilityScore := signatures
	if blocks > 0 {
		//Temp without proposed stat
		availabilityPercent := float64(blocks) / float64(blocks)
		availabilityScore += uint64(availabilityPercent * float64(currentHeight-nodeRegisterBlock))
	}

	return availabilityScore
}
