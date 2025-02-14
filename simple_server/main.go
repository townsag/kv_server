package main

import (
	"encoding/json"
	"net/http"
	"log"
	"log/slog"
	"os"
	"example.com/kv_store"
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
	var key string = r.URL.Query().Get("key")
	if key == "" {
		writeJsonResponse(w, map[string]string{"error":"'key' query parameter is required"}, http.StatusBadRequest)
		return
	}

	value, err := h.store.Get(key)
	if err != nil {
		writeJsonResponse(w, map[string]string{"error": err.Error()}, http.StatusNotFound)
		return
	}
	writeJsonResponse(w, map[string]interface{}{"value":value}, http.StatusOK)
}

func (h *kvHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	// decode the request body
	var kvRequest KeyValueRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&kvRequest); err != nil {
		writeJsonResponse(w, map[string]string{"error":"invalid json body"}, http.StatusBadRequest)
		return
	}

	if kvRequest.Key == "" || kvRequest.Value == "" {
		writeJsonResponse(w, map[string]string{"error":"key and value are required fields"}, http.StatusBadRequest)
		return
	}

	// update the store
	if err := h.store.Set(kvRequest.Key, kvRequest.Value); err != nil {
		writeJsonResponse(w, map[string]string{"error":err.Error()}, http.StatusInternalServerError)
		return
	}
	writeJsonResponse(w, map[string]string{"message":"set successful"}, http.StatusAccepted)
}


func (h *kvHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	// decode the request body
	var kRequest KeyRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&kRequest); err != nil {
		writeJsonResponse(w, map[string]string{"error":"invalid json body"}, http.StatusBadRequest)
		return
	}

	if kRequest.Key == "" {
		writeJsonResponse(w, map[string]string{"error":"field 'key' is required"}, http.StatusBadRequest)
		return
	}

	// update the store to reflect the delete operation
	if err := h.store.Delete(kRequest.Key); err != nil {
		writeJsonResponse(w, map[string]string{"message":err.Error()}, http.StatusInternalServerError)
		return
	}
	writeJsonResponse(w, map[string]string{"message":"delete successful"}, http.StatusOK)
}

func main() {
	var store kv_store.Store = kv_store.NewMemoryStore()
	var logger *slog.Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	var loggingMiddleware func(http.Handler) http.Handler = NewLoggingMiddleware(logger)
	var handler *kvHandler = newKVHandler(store)

	var mux *http.ServeMux = http.NewServeMux()
	mux.Handle("/item", handler)
	mux.Handle("/item/", http.RedirectHandler("/item", http.StatusTemporaryRedirect))

	var port string = ":8000"
	log.Printf("starting kv server on port: %s", port)
	err := http.ListenAndServe(port, RequestIdMiddleware(loggingMiddleware(mux)))
	log.Fatal(err.Error())
}