package api

import (
	"fmt"
	"net/http"
	response "oasisTracker/common/http/responce"
	"oasisTracker/common/log"
	"oasisTracker/smodels"

	"go.uber.org/zap"
)

func (api *API) GetTransactionsList(w http.ResponseWriter, r *http.Request) {

	params := smodels.NewTransactionsParams()
	err := api.queryDecoder.Decode(&params, r.URL.Query())
	if err != nil {
		response.JsonError(w, err)
		return
	}

	txs, count, err := api.services.GetTransactionsList(params)
	if err != nil {
		log.Error("GetTransactionsList api error", zap.Error(err))
		response.JsonError(w, err)
		return
	}

	w.Header().Set(TotalCountHeader, fmt.Sprint(count))
	Json(w, txs)
}
