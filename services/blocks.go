package services

import (
	"fmt"
	"oasisTracker/services/render"
	"oasisTracker/smodels"
	"time"
)

const getBlocksListEP = "/data/blocks"

type blocksRespStr struct {
	arr     []smodels.Block
	counter uint64
}

func (s *ServiceFacade) GetBlockList(params smodels.BlockParams) ([]smodels.Block, uint64, error) {

	raw, ok, err := s.apiCache.Get(fmt.Sprintf("%s?limit=%d&offset=%d&from=%d&to=%d&proposer=%v&id=%v&lvl=%v",
		getBlocksListEP, params.Limit, params.Offset, params.From, params.To, params.Proposer, params.BlockID,
		params.BlockLevel))
	if err != nil {
		return nil, 0, err
	}

	if !ok {
		count, err := s.dao.BlocksCount(params)
		if err != nil {
			return nil, 0, err
		}

		blocks, err := s.dao.GetBlocksList(params)
		if err != nil {
			return nil, 0, err
		}

		rendered := render.Blocks(blocks)
		info := blocksRespStr{
			arr:     rendered,
			counter: count,
		}

		err = s.apiCache.Save(fmt.Sprintf("%s?limit=%d&offset=%d&from=%d&to=%d&proposer=%v&id=%v&lvl=%v",
			getBlocksListEP, params.Limit, params.Offset, params.From, params.To, params.Proposer, params.BlockID,
			params.BlockLevel), info, time.Second*6)
		if err != nil {
			return nil, 0, err
		}

		return rendered, count, nil
	} else {
		info := raw.(blocksRespStr)
		return info.arr, info.counter, err
	}
}
