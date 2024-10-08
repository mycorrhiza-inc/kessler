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

func TestPostgresConnection() (string, error) {
	ctx := context.Background()
	config := pgx.ConnConfig{}

	conn, err := pgx.Connect(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return "", fmt.Errorf("Unable to connect to database")

	}
	defer conn.Close()
	queries := dbstore.New(conn)
	files, err := queries.ListFiles(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing files: %v\n", err)
		return "", fmt.Errorf("Error Found")

	}
	truncatedFiles := files[:100]
	fmt.Println("Successfully listed files:", truncatedFiles)
	return "Success", nil
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
