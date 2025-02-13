package main

import (
	"encoding/json"
	"net/http"
	"log"
	"example.com/kv_store"
)

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

func getHandler(s kv_store.Store) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// check that the request verb is get
		if r.Method != http.MethodGet {
			writeJsonResponse(w, map[string]string{"error":"invalid request method"}, http.StatusMethodNotAllowed)
			return
		}

		var key string = r.URL.Query().Get("key")
		if key == "" {
			writeJsonResponse(w, map[string]string{"error":"'key' query parameter is required"}, http.StatusBadRequest)
			return
		}

		value, err := s.Get(key)
		if err != nil {
			writeJsonResponse(w, map[string]string{"error": err.Error()}, http.StatusNotFound)
			return
		}
		writeJsonResponse(w, map[string]interface{}{"value":value}, http.StatusOK)
	}
}

func setHandler(s kv_store.Store) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// check the request verb is post
		if r.Method != http.MethodPost {
			writeJsonResponse(w, map[string]string{"error":"invalid request method"}, http.StatusMethodNotAllowed)
			return
		}

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
		if err := s.Set(kvRequest.Key, kvRequest.Value); err != nil {
			writeJsonResponse(w, map[string]string{"error":err.Error()}, http.StatusInternalServerError)
			return
		}
		writeJsonResponse(w, map[string]string{"message":"set successful"}, http.StatusAccepted)
	}
}

func deleteHandler(s kv_store.Store) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// check that the request method is post
		if r.Method != http.MethodPost {
			writeJsonResponse(w, map[string]string{"error":"invalid request method"}, http.StatusMethodNotAllowed)
			return
		}

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
		if err := s.Delete(kRequest.Key); err != nil {
			writeJsonResponse(w, map[string]string{"message":err.Error()}, http.StatusInternalServerError)
			return
		}
		writeJsonResponse(w, map[string]string{"message":"delete successful"}, http.StatusOK)
	}
}

func main() {
	var store kv_store.Store = kv_store.NewMemoryStore()

	// var mux *http.ServeMux = http.NewServeMux()

	http.HandleFunc("/GET", getHandler(store))
	http.HandleFunc("/SET", setHandler(store))
	http.HandleFunc("/DELETE", deleteHandler(store))

	var port string = ":8000"
	log.Printf("starting kv server on port: %s", port)
	err := http.ListenAndServe(port, nil)
	log.Fatal(err.Error())
}