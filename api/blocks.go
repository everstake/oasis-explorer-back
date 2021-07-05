package api

import (
	"fmt"
	"net/http"
	"oasisTracker/common/apperrors"
	response "oasisTracker/common/http/responce"
	"oasisTracker/common/log"
	"oasisTracker/smodels"

	"go.uber.org/zap"
)

func (api *API) GetBlocksList(w http.ResponseWriter, r *http.Request) {

	params := smodels.NewBlockParams()
	err := api.queryDecoder.Decode(&params, r.URL.Query())
	if err != nil {
		response.JsonError(w, err)
		return
	}

	err = params.Validate()
	if err != nil {
		log.Error("params.Validate", zap.Error(err))
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam, err.Error()))
		return
	}

	blocks, count, err := api.services.GetBlockList(params)
	if err != nil {
		log.Error("GetBlockList api error", zap.Error(err))
		response.JsonError(w, err)
		return
	}

	w.Header().Set(TotalCountHeader, fmt.Sprint(count))
	Json(w, blocks)
}
