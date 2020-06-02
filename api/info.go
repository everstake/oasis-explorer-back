package api

import (
	"go.uber.org/zap"
	"net/http"
	response "oasisTracker/common/http/responce"
	"oasisTracker/common/log"
)

func (api *API) GetInfo(w http.ResponseWriter, r *http.Request) {

	info, err := api.services.GetInfo()
	if err != nil {
		log.Error("GetInfo err", zap.Error(err))
		response.JsonError(w, err)
		return
	}

	Json(w, info)
}
