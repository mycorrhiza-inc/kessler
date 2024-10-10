package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/blake2b"

	"github.com/mycorrhiza-inc/kessler/backend-go/gen/dbstore"
	"github.com/mycorrhiza-inc/kessler/backend-go/rag"
	"github.com/mycorrhiza-inc/kessler/backend-go/search"
)

func PgPoolConfig() *pgxpool.Config {
	const defaultMaxConns = int32(10)
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

	dbConfig.BeforeAcquire = func(ctx context.Context, c *pgx.Conn) bool {
		log.Println("Before acquiring the connection pool to the database!!")
		return true
	}

	dbConfig.AfterRelease = func(c *pgx.Conn) bool {
		log.Println("After releasing the connection pool to the database!!")
		return true
	}

	dbConfig.BeforeClose = func(c *pgx.Conn) {
		log.Println("Closed the connection pool to the database!!")
	}

	return dbConfig
}

type UserValidation struct {
	validated bool
	userID    string
}

func makeTokenValidator(dbtx_val dbstore.DBTX) func(r *http.Request) UserValidation {
	return_func := func(r *http.Request) UserValidation {
		token := r.Header.Get("Authorization")
		if token == "" {
			return UserValidation{
				validated: true,
				userID:    "anonomous",
			}
		}
		if strings.HasPrefix(token, "Bearer thaum_") {
			q := *dbstore.New(dbtx_val)
			const trim = len("Bearer thaum_")
			// Replacing this with PBKDF2 or something would be more secure, but it should matter since every API key can be gaurenteed to have at least 128/256 bits of strength.
			hash := blake2b.Sum256([]byte(token[trim:]))
			encodedHash := base64.StdEncoding.EncodeToString(hash[:])
			fmt.Println(encodedHash)
			ctx := r.Context()
			result, err := q.CheckIfThaumaturgyAPIKeyExists(ctx, encodedHash)
			if result.KeyBlake3Hash == encodedHash && err != nil {
				return UserValidation{userID: "thaumaturgy", validated: true}
			}
			return UserValidation{validated: false}
		}
		return UserValidation{validated: false}
	}
	return return_func
}

func makeAuthMiddleware(dbtx_val dbstore.DBTX) func(http.Handler) http.Handler {
	tokenValidator := makeTokenValidator(dbtx_val)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userInfo := tokenValidator(r)
			if userInfo.validated {
				r.Header.Set("Authorization", fmt.Sprintf("Authenticated %s", userInfo.userID))
				next.ServeHTTP(w, r)

			} else {
				http.Error(w, "Forbidden", http.StatusForbidden)
			}
		})
	}
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

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Create database connection
	connPool, err := pgxpool.NewWithConfig(context.Background(), PgPoolConfig())
	if err != nil {
		log.Fatal("Error while creating connection to the database!!")
	}

	// connection, err := connPool.Acquire(context.Background())
	// if err != nil {
	// 	log.Fatal("Error while acquiring connection from the database pool!!")
	// }
	//
	// err = connection.Ping(context.Background())
	// if err != nil {
	// 	log.Fatal("Could not ping database")
	// }
	//
	// fmt.Println("Connected to the database!!")
	// connection.Release()

	defer connPool.Close()

	mux := mux.NewRouter()
	misc_s := mux.PathPrefix("/api/v2").Subrouter()
	misc_s.HandleFunc("/search", search.HandleSearchRequest)
	misc_s.HandleFunc("/rag/basic_chat", rag.HandleBasicChatRequest)
	misc_s.HandleFunc("/rag/chat", rag.HandleRagChatRequest)
	const timeout = time.Second * 10

	muxWithMiddlewares := http.TimeoutHandler(mux, timeout, "Timeout!")
	handler := corsMiddleware(muxWithMiddlewares)

	server := &http.Server{
		Addr:         ":4041",
		Handler:      handler,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	}

	log.Println("Starting server on :4041")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Webserver Failed: %s", err)
	}
}
