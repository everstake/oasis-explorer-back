package api

import (
	"net/http"
	"oasisTracker/conf"
)

func (api *API) Index(w http.ResponseWriter, r *http.Request) {
	Json(w, map[string]string{
		"service": conf.Service,
	})
}

func (api *API) Health(w http.ResponseWriter, r *http.Request) {
	Json(w, map[string]bool{
		"status": true,
	})
}
