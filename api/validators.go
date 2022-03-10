package api

import (
	"fmt"
	"net/http"
	"net/url"
	"oasisTracker/common/apperrors"
	response "oasisTracker/common/http/responce"
	"oasisTracker/common/log"
	"oasisTracker/smodels"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
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

	validators, count, err := api.services.GetValidatorList(params)
	if err != nil {
		log.Error("GetValidatorList api error", zap.Error(err))
		response.JsonError(w, err)
		return
	}

	for i := range validators {
		validators[i].DayUptime = 0
	}

	w.Header().Set(TotalCountHeader, fmt.Sprint(count))
	Json(w, validators)
}

func (api *API) GetPublicValidatorsSearchList(w http.ResponseWriter, r *http.Request) {
	validators, err := api.services.GetPublicValidatorsSearchList()
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

	validators, err := api.services.GetValidatorInfo(account)
	if err != nil {
		log.Error("GetValidatorList api error", zap.Error(err))
		response.JsonError(w, err)
		return
	}

	Json(w, validators)
}

func (api *API) GetValidatorDelegators(w http.ResponseWriter, r *http.Request) {
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

	params := smodels.CommonParams{}
	err = api.queryDecoder.Decode(&params, r.URL.Query())
	if err != nil {
		response.JsonError(w, err)
		return
	}

	delegators, err := api.services.GetValidatorDelegators(account, params)
	if err != nil {
		log.Error("GetValidatorList api error", zap.Error(err))
		response.JsonError(w, err)
		return
	}

	Json(w, delegators)
}

func (api *API) GetValidatorBlocks(w http.ResponseWriter, r *http.Request) {
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

	params := smodels.CommonParams{}
	err = api.queryDecoder.Decode(&params, r.URL.Query())
	if err != nil {
		response.JsonError(w, err)
		return
	}

	blocks, err := api.services.GetValidatorBlocks(account, params)
	if err != nil {
		log.Error("GetValidatorBlocks api error", zap.Error(err))
		response.JsonError(w, err)
		return
	}

	Json(w, blocks)
}

func (api *API) GetValidatorRewards(w http.ResponseWriter, r *http.Request) {
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

	params := smodels.CommonParams{}
	err = api.queryDecoder.Decode(&params, r.URL.Query())
	if err != nil {
		response.JsonError(w, err)
		return
	}

	rewards, err := api.services.GetValidatorRewards(account, params)
	if err != nil {
		log.Error("GetValidatorRewards api error", zap.Error(err))
		response.JsonError(w, err)
		return
	}

	Json(w, rewards)
}

func (api *API) GetValidatorRewardsStat(w http.ResponseWriter, r *http.Request) {
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

	stat, err := api.services.GetValidatorRewardsStat(account)
	if err != nil {
		log.Error("GetValidatorRewards api error", zap.Error(err))
		response.JsonError(w, err)
		return
	}

	Json(w, stat)
}
