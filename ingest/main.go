package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"thaumaturgy/routes"
	"thaumaturgy/tasks"

	"github.com/gorilla/mux"
	"github.com/hibiken/asynq"
	_ "github.com/swaggo/http-swagger/example/gorilla/docs" // docs is generated by Swag CLI, you have to import it.
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

//	@title			Swagger Example API
//	@version		1.0
//	@description	This is a sample server Petstore server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		petstore.swagger.io
// @BasePath	/v2

const (
	redisAddr   = "127.0.0.1:6379"
	concurrency = 30 // Max concurrent tasks
)

// In main.go add this middleware
func clientMiddleware(client *asynq.Client) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := tasks.WithClient(r.Context(), client)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func main() {
	// Initialize router and root path

	r := mux.NewRouter()
	root := "/ingest_v1"
	// Set up Swagger endpoint
	r.PathPrefix(root + "/swagger").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost/ingest_v1/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	// Create asynq client
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	defer client.Close()

	// Create API subrouter with client middleware
	api := r.PathPrefix(root).Subrouter()
	api.Use(clientMiddleware(client))
	routes.DefineGlobalRouter(api) // Pass the subrouter to routes package
	// Create asynq client

	// Create and start worker
	worker := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: concurrency,
			Queues: map[string]int{
				"default": 1,
			},
		},
	)

	// Create mux and register handlers
	asyncq_client := asynq.NewServeMux()
	// mux.Use(tasksMiddleware(client))
	// mux.HandleFunc(tasks.TypeAddFileScraper, tasks.HandleAddFileScraperTask)
	// mux.HandleFunc(tasks.TypeProcessExistingFile, tasks.HandleProcessFileTask)

	// Run worker in separate goroutine
	go func() {
		if err := worker.Run(asyncq_client); err != nil {
			log.Fatalf("Failed to start worker: %v", err)
		}
	}()

	// Start HTTP server in a goroutine
	go func() {
		if err := http.ListenAndServe(":4042", r); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Set up shutdown handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down...")
	worker.Shutdown()
}
