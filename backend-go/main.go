package main

import (
	"context"
	"fmt"
	"kessler/admin"
	"kessler/autocomplete"
	"kessler/crud"
	"kessler/database"
	"kessler/health"
	"kessler/jobs"
	"kessler/rag"
	"kessler/search"
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/log"

	"github.com/gorilla/mux"
)

type UserValidation struct {
	validated bool
	userID    string
}

type AccessTokenData struct {
	AccessToken string `json:"access_token"`
}

// CORS middleware function
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // or specify allowed origin
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Origin")

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

func HandleVersionHash(w http.ResponseWriter, r *http.Request) {
	// Get the version hash from the environment variable
	versionHash := os.Getenv("VERSION_HASH")
	w.Write([]byte(versionHash))
}

func timeoutMiddleware(timeout time.Duration) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			r = r.WithContext(ctx)

			done := make(chan bool)
			go func() {
				next.ServeHTTP(w, r)
				done <- true
			}()

			select {
			case <-ctx.Done():
				w.WriteHeader(http.StatusGatewayTimeout)
				fmt.Fprintf(w, "Request timeout after %v\n", timeout)
				return
			case <-done:
				return
			}
		})
	}
}

func main() {
	// initialize the database connection pool
	database.Init(30)
	defer database.ConnPool.Close()

	r := mux.NewRouter()
	rootRoute := r.PathPrefix("/v2").Subrouter()

	// default subrouter
	const standardTimeout = time.Second * 20
	standardRoute := rootRoute.PathPrefix("").Subrouter()
	standardRoute.Use(timeoutMiddleware(standardTimeout))

	// standard rest
	publicSubroute := standardRoute.PathPrefix("/public").Subrouter()
	crud.DefineCrudRoutes(publicSubroute)

	// heathcheck
	healthSubroute := standardRoute.PathPrefix("/health").Subrouter()
	health.DefineHealthRoutes(healthSubroute)

	// search endpoints
	searchSubroute := standardRoute.PathPrefix("/search").Subrouter()
	search.DefineSearchRoutes(searchSubroute)

	// search autocomplete sugggestions endpoints
	autocompleteSubroute := standardRoute.PathPrefix("/autocomplete").Subrouter()
	autocomplete.DefineAutocompleteRoutes(autocompleteSubroute)

	ragSubroute := standardRoute.PathPrefix("/rag").Subrouter()
	ragSubroute.HandleFunc("/basic_chat", rag.HandleBasicChatRequest)
	ragSubroute.HandleFunc("/chat", rag.HandleRagChatRequest)

	standardRoute.HandleFunc("/version_hash", HandleVersionHash)

	// admin route
	const adminTimeout = time.Minute * 10
	adminRoute := rootRoute.PathPrefix("/admin").Subrouter()
	adminRoute.Use(timeoutMiddleware(adminTimeout))
	admin.DefineAdminRoutes(adminRoute)

	// jobs routes
	jobSubroute := rootRoute.PathPrefix("/jobs").Subrouter()
	jobs.DefineJobRoutes(jobSubroute)

	mux.CORSMethodMiddleware(r)

	server := &http.Server{
		Addr:    ":4041",
		Handler: r,
		// Set longer timeouts at server level to allow for admin operations
		ReadTimeout:  adminTimeout,
		WriteTimeout: adminTimeout,
	}

	log.Info("Starting server on :4041")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Webserver Failed: %s", err)
	}
}
