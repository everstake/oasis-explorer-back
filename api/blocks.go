package api

import (
	"go.uber.org/zap"
	"net/http"
	response "oasisTracker/common/http/responce"
	"oasisTracker/common/log"
	"oasisTracker/smodels"
)

func (api *API) GetBlocksList(w http.ResponseWriter, r *http.Request) {

	params := smodels.BlockParams{}
	err := api.queryDecoder.Decode(&params, r.URL.Query())
	if err != nil {
		response.JsonError(w, err)
		return
	}

	blocks, err := api.services.GetBlockList(params)
	if err != nil {
		log.Error("GetBlockList api error", zap.Error(err))
		response.JsonError(w, err)
		return
	}

	Json(w, blocks)
}
