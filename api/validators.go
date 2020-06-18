package api

import (
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"oasisTracker/common/apperrors"
	response "oasisTracker/common/http/responce"
	"oasisTracker/common/log"
	"oasisTracker/smodels"
)

func (api *API) GetValidatorsList(w http.ResponseWriter, r *http.Request) {

	params := smodels.NewValidatorListParams()
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

	validators, err := api.services.GetValidatorList(params)
	if err != nil {
		log.Error("GetValidatorList api error", zap.Error(err))
		response.JsonError(w, err)
		return
	}

	Json(w, validators)
}

func (api *API) GetValidatorInfo(w http.ResponseWriter, r *http.Request) {
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

	validators, err := api.services.GetValidatorList(smodels.ValidatorParams{
		CommonParams: smodels.CommonParams{Limit: 1},
		ValidatorID:  account,
	})

	if err != nil {
		log.Error("GetValidatorList api error", zap.Error(err))
		response.JsonError(w, err)
		return
	}

	Json(w, validators)
}
