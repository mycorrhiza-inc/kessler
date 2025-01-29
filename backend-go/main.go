package main

import (
	"context"
	"kessler/admin"
	"kessler/autocomplete"
	"kessler/crud"
	"kessler/health"
	"kessler/jobs"
	"kessler/rag"
	"kessler/search"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func PgPoolConfig() *pgxpool.Config {
	const defaultMaxConns = int32(30)
	const defaultMinConns = int32(0)
	const defaultMaxConnLifetime = time.Hour
	const defaultMaxConnIdleTime = time.Minute * 30
	const defaultHealthCheckPeriod = time.Minute
	const defaultConnectTimeout = time.Second * 5

	// Your own Database URL
	DATABASE_URL := os.Getenv("DATABASE_CONNECTION_STRING")

	dbConfig, err := pgxpool.ParseConfig(DATABASE_URL)
	if err != nil {
		log.Fatal("Failed to create a config, error: ", err)
	}

	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout
	// Removed to clean up logging in golang
	dbConfig.BeforeAcquire = func(ctx context.Context, c *pgx.Conn) bool {
		// log.Println("Before acquiring the connection pool to the database!!")
		return true
	}

	dbConfig.AfterRelease = func(c *pgx.Conn) bool {
		// log.Println("After releasing the connection pool to the database!!")
		return true
	}

	dbConfig.BeforeClose = func(c *pgx.Conn) {
		// log.Println("Closed the connection pool to the database!!")
	}

	return dbConfig
}

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

type RouteGroup struct {
	router         *mux.Router
	timeoutHandler http.Handler
	prefixes       []string
}

// NewRouteGroup creates a new RouteGroup with specified timeout and middleware
func NewRouteGroup(timeout time.Duration, timeoutMsg string, middleware ...mux.MiddlewareFunc) *RouteGroup {
	r := mux.NewRouter()
	for _, m := range middleware {
		r.Use(m)
	}
	return &RouteGroup{
		router:         r,
		timeoutHandler: http.TimeoutHandler(r, timeout, timeoutMsg),
		prefixes:       []string{},
	}
}

// HandlePrefix registers a new path prefix and sets up routes using the provided function
func (rg *RouteGroup) HandlePrefix(prefix string, setup func(*mux.Router)) {
	sub := rg.router.PathPrefix(prefix).Subrouter()
	setup(sub)
	rg.prefixes = append(rg.prefixes, prefix)
}

// HandleFunc registers a direct route without path prefixing
func (rg *RouteGroup) HandleFunc(path string, handler http.HandlerFunc) {
	rg.router.HandleFunc(path, handler)
}

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }
	// Create database connection
	var connPool *pgxpool.Pool
	connPool, err := pgxpool.NewWithConfig(context.Background(), PgPoolConfig())
	if err != nil {
		log.Fatal("Error while creating connection to the database!!")
	}

	defer connPool.Close()

	// mux_route := mux.NewRouter()
	// static.HandleStaticGenerationRouting(mux, connPool)
	const timeout = time.Second * 20
	const adminTimeout = time.Minute * 10
	dbMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "db", connPool)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	// Initialize route groups
	adminGroup := NewRouteGroup(adminTimeout, "Admin Timeout!", dbMiddleware)
	regularGroup := NewRouteGroup(timeout, "Timeout!", dbMiddleware)

	// Setup admin routes
	adminGroup.HandlePrefix("/v2/admin/", admin.DefineAdminRoutes)
	adminGroup.HandlePrefix("/v2/jobs", jobs.DefineJobRoutes)
	// Add more admin prefixes as needed

	// Setup regular routes
	regularGroup.HandlePrefix("/v2/public", crud.DefineCrudRoutes)
	regularGroup.HandlePrefix("/v2/health", health.DefineHealthRoutes)
	regularGroup.HandlePrefix("/v2/search", search.DefineSearchRoutes)
	regularGroup.HandlePrefix("/v2/autocomplete", autocomplete.DefineAutocompleteRoutes)

	regularGroup.HandleFunc("/v2/version_hash", HandleVersionHash)
	regularGroup.HandleFunc("/v2/rag/basic_chat", rag.HandleBasicChatRequest)
	regularGroup.HandleFunc("/v2/rag/chat", rag.HandleRagChatRequest)

	// Final handler automatically routes based on registered prefixes
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		currentPath := r.URL.Path
		for _, prefix := range adminGroup.prefixes {
			if strings.HasPrefix(currentPath, prefix) {
				adminGroup.timeoutHandler.ServeHTTP(w, r)
				return
			}
		}
		regularGroup.timeoutHandler.ServeHTTP(w, r)
	})

	handler := corsMiddleware(finalHandler)

	server := &http.Server{
		Addr:    ":4041",
		Handler: handler,
		// Set longer timeouts at server level to allow for admin operations
		ReadTimeout:  adminTimeout,
		WriteTimeout: adminTimeout,
	}

	log.Println("Starting server on :4041")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Webserver Failed: %s", err)
	}
}
