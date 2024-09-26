package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/mycorrhizainc/kessler/backend-go/rag"
	"github.com/mycorrhizainc/kessler/backend-go/search"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	// This seems wrong, and the lsp is complaining, but format on save always adds it
)

// path swagger is customable
// path (/*any) is required for load the html page own by swagger
// http://localhost:8810/nexsoft/doc/api/swagger/index.html
func SwaggerRouting(router *mux.Router) {
	prefix := "/api/v2"
	router.PathPrefix(prefix).Handler(httpSwagger.Handler(
		httpSwagger.URL("doc.json"),
		// httpSwagger.DeepLinking(true),
		// httpSwagger.DocExpansion("none"),
		// httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)
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

// --------------------------------------------------
// annotations		: Annotation used for Swagger-UI
//					and will be mapping to folder and files (./root/docs/**)
// docs import		: import ( _ "denitiawan/research-swagger-gomod-gin/docs" )
//					will be used for update all values on all files inside that folder
//					when you run syntax (swag init)
// url swagger-ui 	: http://localhost:5050/nexsoft/doc/api/swagger/index.html
// --------------------------------------------------
//	@version		1.1.0
//	@title			Demo Swagger-UI (GO+GORILLA MUX) for Nexsoft Project
//	@description	Implement swagger-ui on Go project with Gorilla Mux (web framework) + JWT Authorization
//	@host			localhost:5050
//	@BasePath		/

// ------showing authorize button (but validation jwt is not working)---------
//	@Security					Authorization
//	@securityDefinitions.apikey	Authorization
//	@in							header
//	@name						Authorization
//	@schemes					http
// ------showing authorize button (but validation jwt is not working)---------
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

	prefix := "/api/v2"
	mux := mux.NewRouter()
	mux.HandleFunc(prefix+"/search", search.HandleSearchRequest)
	mux.HandleFunc(prefix+"/basic_chat", rag.HandleBasicChatRequest)

	muxWithMiddlewares := http.TimeoutHandler(mux, time.Second*3, "Timeout!")
	handler := corsMiddleware(muxWithMiddlewares)

	server := &http.Server{
		Addr:         ":4041",
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	SwaggerRouting(mux)

	log.Println("Starting server on :4041")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Webserver Failed: %s", err)
	}
}
