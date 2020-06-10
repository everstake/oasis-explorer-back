package api

import (
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"oasisTracker/common/apperrors"
	response "oasisTracker/common/http/responce"
	"oasisTracker/common/log"
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
