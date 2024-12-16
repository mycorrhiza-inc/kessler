package main

import (
	"context"
	"kessler/admin"
	"kessler/autocomplete"
	"kessler/crud"
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
	"github.com/joho/godotenv"
)

var connPool *pgxpool.Pool

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

func withDBTX(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "db", connPool)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Create database connection
	connPool, err = pgxpool.NewWithConfig(context.Background(), PgPoolConfig())
	if err != nil {
		log.Fatal("Error while creating connection to the database!!")
	}

	defer connPool.Close()

	// mux_route := mux.NewRouter()
	// static.HandleStaticGenerationRouting(mux, connPool)
	const timeout = time.Second * 20
	const adminTimeout = time.Minute * 10

	// Create two separate routers for different timeout requirements
	adminMux := mux.NewRouter()
	adminMux.Use(withDBTX)
	adminRouter := adminMux.PathPrefix("/v2/admin/").Subrouter()
	admin.DefineAdminRoutes(adminRouter)
	adminMux.PathPrefix("/v2/admin/").Handler(adminRouter)
	adminMux.PathPrefix("/v2/crud/").Handler(adminRouter)
	adminWithTimeout := http.TimeoutHandler(adminMux, adminTimeout, "Admin Timeout!")

	// Regular routes with standard timeout
	regularMux := mux.NewRouter()
	regularMux.Use(withDBTX)
	public_subrouter := regularMux.PathPrefix("/v2/public").Subrouter()
	crud.DefineCrudRoutes(public_subrouter)
	regularMux.PathPrefix("/v2/public/").Handler(public_subrouter)
	regularMux.HandleFunc("/v2/search", search.HandleSearchRequest)
	regularMux.HandleFunc("/v2/rag/basic_chat", rag.HandleBasicChatRequest)
	regularMux.HandleFunc("/v2/rag/chat", rag.HandleRagChatRequest)
	regularMux.HandleFunc("/v2/recent_updates", search.HandleRecentUpdatesRequest)
	autocomplete_subrouter := regularMux.PathPrefix("/v2/autocomplete").Subrouter()
	autocomplete.DefineAutocompleteRoutes(autocomplete_subrouter)
	regularWithTimeout := http.TimeoutHandler(regularMux, timeout, "Timeout!")

	// Combine both routers
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/v2/admin/") {
			adminWithTimeout.ServeHTTP(w, r)
		} else {
			regularWithTimeout.ServeHTTP(w, r)
		}
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
