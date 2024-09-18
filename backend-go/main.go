package main

import (
	"net/http"
	"time"

	"log"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/mycorrhizainc/kessler/backend/search"
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

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mux := mux.NewRouter()
	mux.HandleFunc("/api/v2/search", search.HandleSearchRequest)

	muxWithMiddlewares := http.TimeoutHandler(mux, time.Second*3, "Timeout!")
	handler := corsMiddleware(muxWithMiddlewares)

	server := &http.Server{
		Addr:         ":4041",
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Starting server on :4041")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Webserver Failed: %s", err)
	}

}
