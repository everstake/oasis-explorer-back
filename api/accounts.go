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

func (api *API) GetAccountInfo(w http.ResponseWriter, r *http.Request) {

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

	acc, err := api.services.GetAccountInfo(account)
	if err != nil {
		log.Error("GetAccountInfo api error", zap.Error(err))
		response.JsonError(w, err)
		return
	}

	Json(w, acc)
}

func (api *API) GetAccountList(w http.ResponseWriter, r *http.Request) {

	params := smodels.NewAccountListParams()
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

	accs, count, err := api.services.GetAccountList(params)
	if err != nil {
		log.Error("GetAccountList api error", zap.Error(err))
		response.JsonError(w, err)
		return
	}

	w.Header().Set(TotalCountHeader, fmt.Sprint(count))
	Json(w, accs)
}

func (api *API) GetAccountRewards(w http.ResponseWriter, r *http.Request) {
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

	rewards, err := api.services.GetAccountRewards(account, params)
	if err != nil {
		log.Error("GetValidatorRewards api error", zap.Error(err))
		response.JsonError(w, err)
		return
	}

	Json(w, rewards)
}

func (api *API) GetAccountRewardsStat(w http.ResponseWriter, r *http.Request) {
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

	stat, err := api.services.GetAccountRewardsStat(account)
	if err != nil {
		log.Error("GetValidatorRewards api error", zap.Error(err))
		response.JsonError(w, err)
		return
	}

	Json(w, stat)
}
