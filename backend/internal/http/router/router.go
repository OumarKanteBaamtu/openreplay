package router

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
	"log"
	"net/http"
	http3 "openreplay/backend/internal/config/http"
	http2 "openreplay/backend/internal/http/services"
	"openreplay/backend/pkg/monitoring"
	"time"
)

type Router struct {
	router          *mux.Router
	cfg             *http3.Config
	services        *http2.ServicesBuilder
	requestSize     syncfloat64.Histogram
	requestDuration syncfloat64.Histogram
	totalRequests   syncfloat64.Counter
}

func NewRouter(cfg *http3.Config, services *http2.ServicesBuilder, metrics *monitoring.Metrics) (*Router, error) {
	switch {
	case cfg == nil:
		return nil, fmt.Errorf("config is empty")
	case services == nil:
		return nil, fmt.Errorf("services is empty")
	case metrics == nil:
		return nil, fmt.Errorf("metrics is empty")
	}
	e := &Router{
		cfg:      cfg,
		services: services,
	}
	e.initMetrics(metrics)
	e.init()
	return e, nil
}

func (e *Router) init() {
	e.router = mux.NewRouter()

	// Root path
	e.router.HandleFunc("/", e.root)

	handlers := map[string]func(http.ResponseWriter, *http.Request){
		"/v1/web/not-started": e.notStartedHandlerWeb,
		"/v1/web/start":       e.startSessionHandlerWeb,
		"/v1/web/i":           e.pushMessagesHandlerWeb,
		"/v1/ios/start":       e.startSessionHandlerIOS,
		"/v1/ios/i":           e.pushMessagesHandlerIOS,
		"/v1/ios/late":        e.pushLateMessagesHandlerIOS,
		"/v1/ios/images":      e.imagesUploadHandlerIOS,
	}
	prefix := "/ingest"

	for path, handler := range handlers {
		e.router.HandleFunc(path, handler).Methods("POST", "OPTIONS")
		e.router.HandleFunc(prefix+path, handler).Methods("POST", "OPTIONS")
	}

	// CORS middleware
	e.router.Use(e.corsMiddleware)
}

func (e *Router) initMetrics(metrics *monitoring.Metrics) {
	var err error
	e.requestSize, err = metrics.RegisterHistogram("requests_body_size")
	if err != nil {
		log.Printf("can't create requests_body_size metric: %s", err)
	}
	e.requestDuration, err = metrics.RegisterHistogram("requests_duration")
	if err != nil {
		log.Printf("can't create requests_duration metric: %s", err)
	}
	e.totalRequests, err = metrics.RegisterCounter("requests_total")
	if err != nil {
		log.Printf("can't create requests_total metric: %s", err)
	}
}

func (e *Router) root(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (e *Router) GetHandler() http.Handler {
	return e.router
}
