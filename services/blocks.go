package services

import (
	"oasisTracker/services/render"
	"oasisTracker/smodels"
)

func (s *ServiceFacade) GetBlockList(params smodels.BlockParams) ([]smodels.Block, error) {

	blocks, err := s.dao.GetBlocksList(params)
	if err != nil {
		return nil, err
	}

	return render.Blocks(blocks), nil
}
