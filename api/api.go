package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"oasisTracker/common/apperrors"
	"oasisTracker/common/log"
	"oasisTracker/conf"
	"oasisTracker/services"
	"time"
)

const (
	gracefulTimeout  = time.Second * 10
	actionsAPIPrefix = ""
)

type (
	API struct {
		router       *mux.Router
		server       *http.Server
		cfg          conf.Config
		services     services.Service
		queryDecoder *schema.Decoder
	}

	// Route stores an API route data
	Route struct {
		Path       string
		Method     string
		Func       func(http.ResponseWriter, *http.Request)
		Middleware []negroni.HandlerFunc
	}
)

func NewAPI(cfg conf.Config, s services.Service) *API {
	queryDecoder := schema.NewDecoder()
	queryDecoder.IgnoreUnknownKeys(true)
	api := &API{
		cfg:          cfg,
		services:     s,
		queryDecoder: queryDecoder,
	}
	api.initialize()
	return api
}

// Run starts the http server and binds the handlers.
func (api *API) Run() error {
	return api.startServe()
}

func (api *API) Stop() error {
	ctx, _ := context.WithTimeout(context.Background(), gracefulTimeout)
	return api.server.Shutdown(ctx)
}

func (api *API) Title() string {
	return "API"
}

func (api *API) initialize(handlerArr ...negroni.Handler) {
	api.router = mux.NewRouter().UseEncodedPath()

	wrapper := negroni.New()

	for _, handler := range handlerArr {
		wrapper.Use(handler)
	}

	api.router.
		PathPrefix("/static").
		Handler(http.StripPrefix("/static", http.FileServer(http.Dir("./resources/static"))))

	wrapper.Use(cors.New(cors.Options{
		AllowedOrigins:   api.cfg.API.CORSAllowedOrigins,
		AllowCredentials: true,
		AllowedMethods:   []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "User-Env"},
	}))

	//public routes
	HandleActions(api.router, wrapper, actionsAPIPrefix, []*Route{
		{Path: "/", Method: http.MethodGet, Func: api.Index},
		{Path: "/health", Method: http.MethodGet, Func: api.Health},
		{Path: "/api", Method: http.MethodGet, Func: api.GetSwaggerAPI},
		{Path: "/metrics_config", Method: http.MethodGet, Func: api.GetMetricsConfig},
		{Path: "/data/info", Method: http.MethodGet, Func: api.GetInfo},
		{Path: "/data/accounts", Method: http.MethodGet, Func: api.GetAccountList},
		{Path: "/data/accounts/{account_id}", Method: http.MethodGet, Func: api.GetAccountInfo},
		{Path: "/data/validators", Method: http.MethodGet, Func: api.GetValidatorsList},
		{Path: "/data/blocks", Method: http.MethodGet, Func: api.GetBlocksList},
		{Path: "/data/transactions", Method: http.MethodGet, Func: api.GetTransactionsList},
		{Path: "/chart/transactions_volume", Method: http.MethodGet, Func: api.GetTransactionsVolume},
		{Path: "/chart/escrow_ratio", Method: http.MethodGet, Func: api.GetEscrowRatio},
		{Path: "/chart/validator_stat/{account_id}", Method: http.MethodGet, Func: api.GetValidatorStats},
	})

	api.server = &http.Server{Addr: fmt.Sprintf(":%d", api.cfg.API.ListenOnPort), Handler: api.router}
}

func (api *API) startServe() error {
	log.Info("Start listening server on port", zap.Uint64("port", api.cfg.API.ListenOnPort))
	err := api.server.ListenAndServe()
	if err == http.ErrServerClosed {
		log.Warn("API server was closed")
		return nil
	}
	if err != nil {
		return fmt.Errorf("cannot run API service: %s", err.Error())
	}
	return nil
}

// HandleActions is used to handle all given routes
func HandleActions(router *mux.Router, wrapper *negroni.Negroni, prefix string, routes []*Route) {
	for _, r := range routes {
		w := wrapper.With()
		for _, m := range r.Middleware {
			w.Use(m)
		}

		w.Use(negroni.Wrap(http.HandlerFunc(r.Func)))
		router.Handle(prefix+r.Path, w).Methods(r.Method, "OPTIONS")
	}
}

// Json writes to ResponseWriter a single JSON-object
func Json(w http.ResponseWriter, data interface{}) {
	js, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// JsonError writes to ResponseWriter error
func JsonError(w http.ResponseWriter, err error) {
	var e *apperrors.Error
	var ok bool

	if e, ok = err.(*apperrors.Error); !ok {
		e = apperrors.FromError(err)
	}

	js, _ := json.Marshal(e.ToMap())
	w.WriteHeader(e.GetHttpCode())
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (api *API) GetSwaggerAPI(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadFile("./resources/templates/swagger.html")
	if err != nil {
		log.Error("GetSwaggerAPI: ReadFile", zap.Error(err))
		return
	}
	_, err = w.Write(body)
	if err != nil {
		log.Error("GetSwaggerAPI: Write", zap.Error(err))
		return
	}
}
