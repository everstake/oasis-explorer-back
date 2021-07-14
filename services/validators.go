package services

import (
	"fmt"
	"oasisTracker/common/apperrors"
	"oasisTracker/services/render"
	"oasisTracker/smodels"
)

func (s *ServiceFacade) GetValidatorInfo(accountID string) (val smodels.Validator, err error) {
	validators, _, err := s.GetValidatorList(smodels.ValidatorParams{
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

func (s *ServiceFacade) GetValidatorList(listParams smodels.ValidatorParams) ([]smodels.Validator, uint64, error) {
	count, err := s.dao.GetValidatorsCount(listParams)
	if err != nil {
		return nil, 0, err
	}

	resp, err := s.dao.GetValidatorsList(listParams)
	if err != nil {
		return nil, 0, err
	}

	for i := range resp {

		resp[i].DayUptime = float64(resp[i].DaySignedBlocks) / float64(resp[i].DayBlocksCount)
		resp[i].TotalUptime = float64(resp[i].SignedBlocksCount) / float64(resp[i].LastBlockLevel-s.genesisHeight-1)

		if !resp[i].IsActive {
			resp[i].Status = smodels.StatusInActive
			continue
		}
		resp[i].Status = smodels.StatusActive
	}

	return render.ValidatorsList(resp), count, nil
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

func (s *ServiceFacade) GetValidatorBlocks(validatorID string, params smodels.CommonParams) ([]smodels.Block, error) {
	entity, err := s.dao.GetAccountValidatorInfo(validatorID)
	if err != nil {
		return nil, err
	}

	blocks, err := s.dao.GetBlocksList(smodels.BlockParams{
		CommonParams: params,
		Proposer:     []string{entity.GetEntity().ConsensusAddress},
	})
	if err != nil {
		return nil, err
	}

	return render.Blocks(blocks), nil
}

func (s *ServiceFacade) GetValidatorRewards(validatorID string, params smodels.CommonParams) ([]smodels.Reward, error) {
	rewards, err := s.dao.GetValidatorRewards(validatorID, params)
	if err != nil {
		return nil, err
	}

	return render.Rewards(rewards), nil
}

func (s *ServiceFacade) GetValidatorRewardsStat(validatorID string) (stat smodels.RewardStat, err error) {
	rewardsStat, err := s.dao.GetValidatorRewardsStat(validatorID)
	if err != nil {
		return stat, err
	}

	return render.RewardStat(rewardsStat), nil
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
