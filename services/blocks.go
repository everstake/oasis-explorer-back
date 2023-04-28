package services

import (
	"oasisTracker/services/render"
	"oasisTracker/smodels"
)

func (s *ServiceFacade) GetBlockList(params smodels.BlockParams) ([]smodels.Block, uint64, error) {

	count, err := s.dao.BlocksCount(params)
	if err != nil {
		return nil, 0, err
	}

	blocks, err := s.dao.GetBlocksList(params)
	if err != nil {
		return nil, 0, err
	}

	return render.Blocks(blocks), count, nil
}
