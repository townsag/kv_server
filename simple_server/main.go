package main

import (
	"encoding/json"
	"net/http"
	"log"
	"log/slog"
	"strconv"
	"os"
	"github.com/townsag/kv_server/kv_store"
	"github.com/townsag/kv_server/simple_server/middleware"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type kvHandler struct {
	store kv_store.Store
}

func newKVHandler(store kv_store.Store) *kvHandler {
	return &kvHandler{
		store: store,
	}
}

func (h *kvHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGet(w, r)
	case http.MethodPost:
		h.handlePost(w, r)
	case http.MethodDelete:
		h.handleDelete(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

type KeyValueRequest struct {
	Key 	string `json:"key"`
	Value 	string `json:"value"`
}

type KeyRequest struct {
	Key		string `json:"key"`
}

func writeJsonResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	var tempEncoder *json.Encoder = json.NewEncoder(w)
	tempEncoder.Encode(data)
}

func (h *kvHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r.Context())
	var key string = r.URL.Query().Get("key")
	if key == "" {
		logger.Warn(
			"failed to serve request",
			"message", "'key' query parameter is required",
			"status", strconv.Itoa(http.StatusBadRequest),
		)
		writeJsonResponse(w, map[string]string{"error":"'key' query parameter is required"}, http.StatusBadRequest)
		return
	}

	// TODO: need to decide on a proper type to be stored in the key value store
	// for now strings
	value, err := h.store.Get(key)
	if err != nil {
		logger.Warn(
			"failed to retrieve key value pair to store",
			"key", key,
			"error", err.Error(),
		)
		writeJsonResponse(w, map[string]string{"error": err.Error()}, http.StatusNotFound)
		return
	}
	logger.Info(
		"successfully retrieved from the store",
		"key", key,
	)
	writeJsonResponse(w, map[string]interface{}{"value":value}, http.StatusOK)
}

func (h *kvHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r.Context())
	// decode the request body
	var kvRequest KeyValueRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&kvRequest); err != nil {
		logger.Warn(
			"received bad json",
			"error", err.Error(),
		)
		writeJsonResponse(w, map[string]string{"error":"invalid json body"}, http.StatusBadRequest)
		return
	}

	if kvRequest.Key == "" || kvRequest.Value == "" {
		logger.Warn(
			"message body is missing one or more fields",
		)
		writeJsonResponse(w, map[string]string{"error":"key and value are required fields"}, http.StatusBadRequest)
		return
	}

	// update the store
	if err := h.store.Set(kvRequest.Key, kvRequest.Value); err != nil {
		logger.Error(
			"failed to write to store",
			"error", err.Error(),
		)
		writeJsonResponse(w, map[string]string{"error":err.Error()}, http.StatusInternalServerError)
		return
	}
	logger.Info(
		"successfully wrote to the store",
		"key", kvRequest.Key,
	)
	writeJsonResponse(w, map[string]string{"message":"set successful"}, http.StatusAccepted)
}


func (h *kvHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r.Context())
	// decode the request body
	var kRequest KeyRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&kRequest); err != nil {
		logger.Warn(
			"revived bad json",
			"error", err.Error(),
		)
		writeJsonResponse(w, map[string]string{"error":"invalid json body"}, http.StatusBadRequest)
		return
	}

	if kRequest.Key == "" {
		logger.Warn(
			"missing field 'key' in request body",
		)
		writeJsonResponse(w, map[string]string{"error":"field 'key' is required"}, http.StatusBadRequest)
		return
	}

	// update the store to reflect the delete operation
	if err := h.store.Delete(kRequest.Key); err != nil {
		logger.Error(
			"failed to delete a value form the store",
			"key", kRequest.Key,
		)
		writeJsonResponse(w, map[string]string{"message":err.Error()}, http.StatusInternalServerError)
		return
	}
	logger.Info(
		"successfully deleted a value from the store",
		"key", kRequest.Key,
	)
	writeJsonResponse(w, map[string]string{"message":"delete successful"}, http.StatusOK)
}

// https://www.alexedwards.net/blog/an-introduction-to-handlers-and-servemuxes-in-go

func main() {
	// initialize a parent logger and create a logging middleware handler
	var logger *slog.Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	var loggingMiddleware func(http.Handler) http.Handler = middleware.NewLoggingMiddleware(logger)

	// initialize a metrics registry, metrics and create a metrics middleware
	var registry *prometheus.Registry = prometheus.NewRegistry()
	var metrics *middleware.Metrics = middleware.NewMetrics(registry)
	var metricsMiddleware func(http.Handler) http.Handler = middleware.NewMetricsMiddleware(metrics)


	var store kv_store.Store = kv_store.NewMemoryStore()
	var applicationHandler *kvHandler = newKVHandler(store)

	var mux *http.ServeMux = http.NewServeMux()
	mux.Handle("/item", applicationHandler)
	mux.Handle("/item/", http.RedirectHandler("/item", http.StatusTemporaryRedirect))
	mux.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{Registry: registry}))

	// add the metrics, request, and logging middleware in that order of execution
	var handler http.Handler = mux
	handler = loggingMiddleware(handler)
	handler = middleware.RequestIdMiddleware(handler)
	handler = metricsMiddleware(handler)
	
	var port string = ":8000"
	log.Printf("starting kv server on port: %s", port)
	err := http.ListenAndServe(port, handler)
	log.Fatal(err.Error())
}