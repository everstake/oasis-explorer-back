package api

import (
	"go.uber.org/zap"
	"net/http"
	"oasisTracker/common/apperrors"
	response "oasisTracker/common/http/responce"
	"oasisTracker/common/log"
	"oasisTracker/smodels"
)

func (api *API) GetTransactionsVolume(w http.ResponseWriter, r *http.Request) {

	params := smodels.ChartParams{}
	err := api.queryDecoder.Decode(&params, r.URL.Query())
	if err != nil {
		log.Error("err", zap.Error(err))
		response.JsonError(w, err)
		return
	}

	err = params.Validate()
	if err != nil {
		log.Error("params.Validate", zap.Error(err))
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam, err.Error()))
	}

	data, err := api.services.GetChartData(params)
	if err != nil {
		log.Error("GetChartData api error", zap.Error(err))
		response.JsonError(w, err)
		return
	}

	Json(w, data)
}
