package main

import (
	// "fmt"
	// "os"
	// "context"
	// "github.com/jackc/pgx/v5"
	"log"
	"net/http"

	"github.com/mycorrhizainc/kessler/backend/search"
)

var (
	addr = "localhost:4041"
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
	//
	// set up db connection
	// conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	// 	os.Exit(1)
	// }
	// // close connection when server exits
	// defer conn.Close(context.Background())

	mux := http.NewServeMux()
	mux.HandleFunc("/search", search.HandleSearchRequest)

	handler := corsMiddleware(mux)

	log.Printf("server is listening at %s", addr)
	log.Fatal(http.ListenAndServe(addr, handler))
}
