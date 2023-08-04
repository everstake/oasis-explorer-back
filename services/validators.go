package services

import (
	"fmt"
	"math"
	"oasisTracker/common/apperrors"
	"oasisTracker/dmodels"
	"oasisTracker/services/render"
	"oasisTracker/smodels"
	"time"
)

const getValidatorListEP = "/data/validators"

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

type validatorsRespStr struct {
	arr     []smodels.Validator
	counter uint64
}

func (s *ServiceFacade) GetValidatorList(listParams smodels.ValidatorParams) ([]smodels.Validator, uint64, error) {

	raw, ok, err := s.apiCache.Get(fmt.Sprintf("%s?limit=%d&offset=%d&validator=%s",
		getValidatorListEP, listParams.Limit, listParams.Offset, listParams.ValidatorID))
	if err != nil {
		return nil, 0, err
	}

	if !ok {
		lastBlock, err := s.dao.GetLastBlock()
		if err != nil {
			return nil, 0, err
		}

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
			resp[i].CurrentEpoch = lastBlock.Epoch

			if !resp[i].IsActive {
				resp[i].Status = smodels.StatusInActive
				continue
			}
			resp[i].Status = smodels.StatusActive
		}

		for i := range resp {
			if math.IsNaN(resp[i].DayUptime) {
				resp[i].DayUptime = 0
			}
		}

		info := validatorsRespStr{
			arr:     render.ValidatorsList(resp),
			counter: count,
		}

		err = s.apiCache.Save(fmt.Sprintf("%s?limit=%d&offset=%d&validator=%s",
			getValidatorListEP, listParams.Limit, listParams.Offset, listParams.ValidatorID), info, time.Second*30)
		if err != nil {
			return nil, 0, err
		}

		return render.ValidatorsList(resp), count, nil
	} else {
		info := raw.(validatorsRespStr)
		return info.arr, info.counter, err
	}
}

func (s *ServiceFacade) GetValidatorListNew(listParams smodels.ValidatorParams) ([]smodels.Validator, uint64, error) {

	raw, ok, err := s.apiCache.Get(fmt.Sprintf("%s?limit=%d&offset=%d&validator=%s",
		getValidatorListEP, listParams.Limit, listParams.Offset, listParams.ValidatorID))
	if err != nil {
		return nil, 0, err
	}

	if !ok {
		vl, err := s.dao.GetValidatorsListNew(listParams)
		if err != nil {
			return nil, 0, err
		}

		blockInfo, err := s.pDao.GetBlocksInfo()
		if err != nil {
			return nil, 0, err
		}

		blockDayInfo, err := s.pDao.GetBlocksDayInfo()
		if err != nil {
			return nil, 0, err
		}

		validatorInfo, err := s.pDao.GetValidatorsInfo()
		if err != nil {
			return nil, 0, err
		}

		valMap := make(map[string]dmodels.ValidatorInfoWithDay)
		for _, v := range validatorInfo {
			valMap[v.Address] = v
		}

		for i := range vl {
			vl[i].DayUptime = float64(valMap[vl[i].EntityAddress].DaySigs) / float64(blockDayInfo.TotalDayBlocks)
			vl[i].TotalUptime = float64(valMap[vl[i].EntityAddress].TotalSigs) / float64(blockInfo.LastLvl-s.genesisHeight-1)
			vl[i].CurrentEpoch = blockInfo.Epoch

			if valMap[vl[i].EntityAddress].DaySigs == 0 {
				vl[i].Status = smodels.StatusInActive
				continue
			}
			vl[i].Status = smodels.StatusActive

			if math.IsNaN(vl[i].DayUptime) {
				vl[i].DayUptime = 0
			}

			vl[i].SignedBlocksCount = valMap[vl[i].EntityAddress].TotalSigs
			vl[i].ProposedBlocksCount = valMap[vl[i].EntityAddress].TotalBlocks
		}

		info := validatorsRespStr{
			arr:     render.ValidatorsList(vl),
			counter: uint64(len(validatorInfo)),
		}

		err = s.apiCache.Save(fmt.Sprintf("%s?limit=%d&offset=%d&validator=%s",
			getValidatorListEP, listParams.Limit, listParams.Offset, listParams.ValidatorID), info, time.Second*30)
		if err != nil {
			return nil, 0, err
		}

		return render.ValidatorsList(vl), uint64(len(validatorInfo)), nil
	} else {
		info := raw.(validatorsRespStr)
		return info.arr, info.counter, err
	}
}

func (s *ServiceFacade) GetValidatorStatsChartData(accountID string, params smodels.ChartParams) ([]smodels.ValidatorStats, error) {

	validators, err := s.dao.GetValidatorsListNew(smodels.ValidatorParams{
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

	blocks, err := s.dao.GetBlocksListNew(smodels.BlockParams{
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
		availabilityPercent := float64(signatures) / float64(blocks)
		availabilityScore += uint64(availabilityPercent * float64(currentHeight-nodeRegisterBlock))
	}

	return availabilityScore
}
