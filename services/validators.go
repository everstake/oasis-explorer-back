package services

import (
	"fmt"
	"oasisTracker/common/apperrors"
	"oasisTracker/services/render"
	"oasisTracker/smodels"
)

func (s *ServiceFacade) GetValidatorInfo(accountID string) (val smodels.Validator, err error) {
	validators, err := s.GetValidatorList(smodels.ValidatorParams{
		CommonParams: smodels.CommonParams{Limit: 1},
		ValidatorID:  accountID,
	})

	if err != nil {
		return val, err
	}

	if len(validators) == 0 {
		return val, apperrors.New(apperrors.ErrNotFound, "account_id")
	}

	return validators[0], nil
}

func (s *ServiceFacade) GetPublicValidatorsSearchList() (list []smodels.ValidatorEntity, err error) {

	val, err := s.dao.PublicValidatorsSearchList()
	if err != nil {
		return nil, err
	}

	return render.PublicValidatorSearch(val), nil
}

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
		ValidatorID:  accountID,
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
		stats[i].AvailabilityScore = calcAvailabilityScore(stats[i].BlocksCount, stats[i].SignaturesCount, validators[0].StartBlockLevel, stats[i].LastBlock)
	}

	return render.ValidatorStatList(stats), nil
}

func (s *ServiceFacade) GetValidatorDelegators(validatorID string, params smodels.CommonParams) ([]smodels.Delegator, error) {

	delegators, err := s.dao.GetValidatorDelegators(validatorID, params)
	if err != nil {
		return nil, err
	}

	return render.DelegatorList(delegators), nil
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
