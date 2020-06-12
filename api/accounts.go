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

	accs, err := api.services.GetAccountList(params)
	if err != nil {
		log.Error("GetAccountList api error", zap.Error(err))
		response.JsonError(w, err)
		return
	}

	Json(w, accs)
}
