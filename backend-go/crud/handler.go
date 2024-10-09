package crud

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/mycorrhiza-inc/kessler/backend-go/gen/dbstore"
)

var pgConnString = os.Getenv("DATABASE_CONNECTION_STRING")

func TestPostgresConnection() (string, error) {
	ctx := context.Background()

	// conn, err := pgx.Connect(ctx, "user=pqgotest dbname=pqgotest sslmode=verify-full")
	conn, err := pgx.Connect(ctx, pgConnString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return "", fmt.Errorf("Unable to connect to database")
	}
	defer conn.Close(ctx)
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
	params := mux.Vars(r)
	fileID := params["uuid"]

	parsedUUID, err := parseUUID(fileID) // hypothetical helper function
	if err != nil {
		http.Error(w, "Invalid File ID format", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	q := NewQueries(nil) // assuming the database connection is managed better in the real app
	file, err := q.ReadFile(ctx, parsedUUID)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	fileSchema := fileToSchema(file)
	response, _ := json.Marshal(fileSchema)

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func getMarkdownHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fileID := params["uuid"]

	parsedUUID, err := parseUUID(fileID) // hypothetical helper function
	if err != nil {
		http.Error(w, "Invalid File ID format", http.StatusBadRequest)
		return
	}

	originalLang := r.URL.Query().Get("original_lang") == "true"
	matchLang := r.URL.Query().Get("match_lang")

	ctx := r.Context()
	q := NewQueries(nil) // assuming the database connection is managed better in the real app

	var texts []FileTextSource
	if originalLang {
		texts, err = q.ListTextsOfFileOriginal(ctx, parsedUUID)
	} else if matchLang != "" {
		texts, err = q.ListTextsOfFileWithLanguage(ctx, ListTextsOfFileWithLanguageParams{
			FileID:   parsedUUID,
			Language: matchLang,
		})
	} else {
		texts, err = q.ListTextsOfFile(ctx, parsedUUID)
		matchLang = "en"
	}

	if err != nil {
		http.Error(w, "Error retrieving texts", http.StatusInternalServerError)
		return
	}

	if len(texts) == 0 {
		http.Error(w, "No texts found for document", http.StatusNotFound)
		return
	}

	markdownText := texts[0].Text
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(markdownText))
}
