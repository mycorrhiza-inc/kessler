package crud

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	"github.com/mycorrhiza-inc/kessler/backend-go/gen/dbstore"
)

var pgConnString = os.Getenv("DATABASE_CONNECTION_STRING")

func testPostgresConnection() {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, pgConnString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return

	}
	defer conn.Close(ctx)
	queries := dbstore.New(conn)
	files, err := queries.ListFiles(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing files: %v\n", err)
		return
	}
	truncatedFiles := files[:100]
	fmt.Println("Successfully listed files:", truncatedFiles)
}

func defineCrudRoutes(router *mux.Router) {
	s := router.PathPrefix("/crud").Subrouter()
	s.HandleFunc("/files/{uuid}", getFileHandler)
}

func getFileHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		fmt.Fprint(w, "Hi there!")
		return
	case http.MethodPost:
		fmt.Fprint(w, "POST request")
		return
	case http.MethodPut:
		fmt.Fprintf(w, "PUT request")
	case http.MethodDelete:
		fmt.Fprintf(w, "DELETE request")
	default:
		http.Error(w, "Unsupported request method", http.StatusMethodNotAllowed)
	}
}
