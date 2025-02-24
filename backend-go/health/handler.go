package health

import (
	"kessler/search"
	"net/http"

	"github.com/gorilla/mux"
)

func DefineHealthRoutes(r *mux.Router) {
	r.HandleFunc(
		"/complete-check",
		CompleteHealthCheckHandler,
	).Methods(http.MethodPost)
	r.HandleFunc(
		"/complete-check",
		MinimalHealthCheckHandler,
	).Methods(http.MethodPost)
}

func CompleteHealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	err := MinimalHealthCheck(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Passed Health Check"))
}

func MinimalHealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	err := MinimalHealthCheck(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Passed Health Check"))
}

func MinimalHealthCheck(r *http.Request) error {
	q := database.GetTx()
	ctx := r.Context()
	_, err := q.HealthCheck(ctx)
	if err != nil {
		return err
	}
	_, err = search.SearchQuickwit(search.ExampleSearchRequest)
	if err != nil {
		return err
	}

	return err
}
