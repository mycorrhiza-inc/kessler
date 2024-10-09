package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"github.com/mycorrhiza-inc/kessler/backend-go/rag"
	"github.com/mycorrhiza-inc/kessler/backend-go/search"
)

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
