package api

import (
	"go.uber.org/zap"
	"net/http"
	response "oasisTracker/common/http/responce"
	"oasisTracker/common/log"
	"oasisTracker/smodels"
)

func (api *API) GetValidatorsList(w http.ResponseWriter, r *http.Request) {

	params := smodels.ValidatorParams{}
	err := api.queryDecoder.Decode(&params, r.URL.Query())
	if err != nil {
		response.JsonError(w, err)
		return
	}

	validators, err := api.services.GetValidatorList(params)
	if err != nil {
		log.Error("GetValidatorList api error", zap.Error(err))
		response.JsonError(w, err)
		return
	}

	Json(w, validators)
}
