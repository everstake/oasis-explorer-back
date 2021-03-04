package api

import (
	"net/http"
	"net/url"
	"oasisTracker/common/apperrors"
	response "oasisTracker/common/http/responce"
	"oasisTracker/common/log"
	"oasisTracker/smodels"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func (api *API) GetTransactionsVolume(w http.ResponseWriter, r *http.Request) {

	params := smodels.ChartParams{}
	err := api.queryDecoder.Decode(&params, r.URL.Query())
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadRequest))
		return
	}

	err = params.Validate()
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam, err.Error()))
		return
	}

	data, err := api.services.GetChartData(params)
	if err != nil {
		log.Error("GetChartData api error", zap.Error(err))
		response.JsonError(w, err)
		return
	}

	Json(w, data)
}

func (api *API) GetEscrowRatio(w http.ResponseWriter, r *http.Request) {

	params := smodels.ChartParams{}
	err := api.queryDecoder.Decode(&params, r.URL.Query())
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadRequest))
		return
	}

	err = params.Validate()
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam, err.Error()))
		return
	}

	data, err := api.services.GetEscrowRatioChartData(params)
	if err != nil {
		log.Error("GetChartData api error", zap.Error(err))
		response.JsonError(w, err)
		return
	}

	Json(w, data)
}

func (api *API) GetTotalAccountsCountChart(w http.ResponseWriter, r *http.Request) {

	params := smodels.ChartParams{}
	err := api.queryDecoder.Decode(&params, r.URL.Query())
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadRequest))
		return
	}

	err = params.Validate()
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam, err.Error()))
		return
	}

	data, err := api.services.GetTotalAccountsCountChartData(params)
	if err != nil {
		log.Error("GetChartData api error", zap.Error(err))
		response.JsonError(w, err)
		return
	}

	Json(w, data)
}

func (api *API) GetAvgBlockTimeChart(w http.ResponseWriter, r *http.Request) {

	params := smodels.ChartParams{}
	err := api.queryDecoder.Decode(&params, r.URL.Query())
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadRequest))
		return
	}

	err = params.Validate()
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam, err.Error()))
		return
	}

	data, err := api.services.GetAvgBlockTimeChartData(params)
	if err != nil {
		log.Error("GetChartData api error", zap.Error(err))
		response.JsonError(w, err)
		return
	}

	Json(w, data)
}

func (api *API) GetFeeVolumeChart(w http.ResponseWriter, r *http.Request) {

	params := smodels.ChartParams{}
	err := api.queryDecoder.Decode(&params, r.URL.Query())
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadRequest))
		return
	}

	err = params.Validate()
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam, err.Error()))
		return
	}

	data, err := api.services.GetFeeVolumeChartData(params)
	if err != nil {
		log.Error("GetChartData api error", zap.Error(err))
		response.JsonError(w, err)
		return
	}

	Json(w, data)
}

func (api *API) GetOperationsCountChart(w http.ResponseWriter, r *http.Request) {

	params := smodels.ChartParams{}
	err := api.queryDecoder.Decode(&params, r.URL.Query())
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadRequest))
		return
	}

	err = params.Validate()
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam, err.Error()))
		return
	}

	data, err := api.services.GetOperationsCountChartData(params)
	if err != nil {
		log.Error("GetChartData api error", zap.Error(err))
		response.JsonError(w, err)
		return
	}

	Json(w, data)
}

func (api *API) GetReclaimAmountChart(w http.ResponseWriter, r *http.Request) {

	params := smodels.ChartParams{}
	err := api.queryDecoder.Decode(&params, r.URL.Query())
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadRequest))
		return
	}

	err = params.Validate()
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam, err.Error()))
		return
	}

	data, err := api.services.GetReclaimAmountChartData(params)
	if err != nil {
		log.Error("GetChartData api error", zap.Error(err))
		response.JsonError(w, err)
		return
	}

	Json(w, data)
}

func (api *API) GetTopEscrowRatioChart(w http.ResponseWriter, r *http.Request) {

	params := smodels.CommonParams{
		Limit: 10,
	}
	err := api.queryDecoder.Decode(&params, r.URL.Query())
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadRequest))
		return
	}

	err = params.Validate()
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam, err.Error()))
		return
	}

	data, err := api.services.GetTopEscrowRatioChart(params)
	if err != nil {
		log.Error("GetChartData api error", zap.Error(err))
		response.JsonError(w, err)
		return
	}

	Json(w, data)
}

func (api *API) GetValidatorStats(w http.ResponseWriter, r *http.Request) {
	urlAcc, ok := mux.Vars(r)["account_id"]
	if !ok || urlAcc == "" {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam, "account_id"))
		return
	}

	account, err := url.QueryUnescape(urlAcc)
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam, "account_id"))
		return
	}

	params := smodels.ChartParams{}
	err = api.queryDecoder.Decode(&params, r.URL.Query())
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

	data, err := api.services.GetValidatorStatsChartData(account, params)
	if err != nil {
		log.Error("GetValidatorStatsChartData api error", zap.Error(err))
		response.JsonError(w, err)
		return
	}

	Json(w, data)
}

func (api *API) GetBalanceChart(w http.ResponseWriter, r *http.Request) {
	urlAcc, ok := mux.Vars(r)["account_id"]
	if !ok || urlAcc == "" {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam, "account_id"))
		return
	}

	account, err := url.QueryUnescape(urlAcc)
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam, "account_id"))
		return
	}

	params := smodels.ChartParams{}
	err = api.queryDecoder.Decode(&params, r.URL.Query())
	if err != nil {
		log.Error("err", zap.Error(err))
		response.JsonError(w, err)
		return
	}

	err = params.Validate()
	if err != nil {
		log.Error("params.Validate", zap.Error(err))
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam, err.Error()))
		return
	}

	data, err := api.services.GetBalanceChartData(account, params)
	if err != nil {
		log.Error("GetBalanceChartData api error", zap.Error(err))
		response.JsonError(w, err)
		return
	}

	Json(w, data)
}
