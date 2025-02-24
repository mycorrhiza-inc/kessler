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
	"kessler/logger"
	"kessler/rag"
	"kessler/search"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type AccessTokenData struct {
	AccessToken string `json:"access_token"`
}

// CORS middleware function
func corsDomainMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		domain := os.Getenv("DOMAIN")
		if domain != "" {
			domain = "https://" + domain
			// Set CORS headers
			w.Header().Set("Access-Control-Allow-Origin", domain) // or specify allowed origin
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*") // or specify allowed origin
		}
		// Putting these in here temporarially, it seems that our cursed wrapping multiple
		// timeouts inside multiple routers broke mux's traditonal cors method handling.
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Origin")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

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
	logger.Init(os.Getenv("GO_ENV"))
	defer logger.Sync()

	log := logger.GetLogger("main")

	log.Debug("starting application",
		zap.String("env", os.Getenv("GO_ENV")),
		zap.Int("port", 4041))

	log.Info("connecting to database")
	if err := database.Init(30); err != nil {
		log.Fatal("unable to connect to connect to the database ", zap.Error(err))
	}
	log.Info("database connection successiful")

	defer database.ConnPool.Close()
	// log.Info("initializing memecached")
	// if err := cache.InitMemcached(); err != nil {
	// 	log.Fatal("unable to connect to memcached", zap.Error(err))
	// }
	// log.Info("cache initialized")

	log.Info("registering api routes")
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

	// Commenting out temporarially, it seems that our cursed wrapping multiple
	// timeouts inside multiple routers broke it.
	// r.Use(mux.CORSMethodMiddleware(r))
	r.Use(corsDomainMiddleware)
	log.Info("routes registered")

	server := &http.Server{
		Addr:    ":4041",
		Handler: r,
		// Set longer timeouts at server level to allow for admin operations
		ReadTimeout:  adminTimeout,
		WriteTimeout: adminTimeout,
	}
	log.Info("server started",
		zap.String("address", ":4041"))
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Webserver Failed: %s", zap.Error(err))
	}
}
