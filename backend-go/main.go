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
	"github.com/jackc/pgx"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/blake2b"

	"github.com/mycorrhiza-inc/kessler/backend-go/gen/dbstore"
	"github.com/mycorrhiza-inc/kessler/backend-go/rag"
	"github.com/mycorrhiza-inc/kessler/backend-go/search"
)

type UserValidation struct {
	validated bool
	userID    string
}

func makeTokenValidator(q *dbstore.Queries) func(r *http.Request) UserValidation {
	return_func := func(r *http.Request) UserValidation {
		token := r.Header.Get("Authorization")
		if token == "" {
			return UserValidation{
				validated: true,
				userID:    "anonomous",
			}
		}
		if strings.HasPrefix(token, "Bearer thaum_") {
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

func makeAuthMiddleware(q *dbstore.Queries) func(http.Handler) http.Handler {
	tokenValidator := makeTokenValidator(q)
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

	ctx := context.Background()

	// conn, err := pgx.Connect(ctx, "user=pqgotest dbname=pqgotest sslmode=verify-full")
	conn, err := pgx.Connect(ctx, pgConnString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return
	}
	defer conn.Close(ctx)
	queries := dbstore.New(conn)

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
