package crud

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mycorrhiza-inc/kessler/backend-go/gen/dbstore"
)

type UserInfo struct {
	validated     bool
	userID        string
	orgID         string
	isThaumaturgy bool
	paymentTier   string
}

func TestPostgresConnection() (string, error) {
	pgConnString := os.Getenv("DATABASE_CONNECTION_STRING")
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

func DefineCrudRoutes(router *mux.Router, dbtx_val dbstore.DBTX) {
	public_subrouter := router.PathPrefix("/public").Subrouter()
	public_subrouter.HandleFunc("/files/{uuid}", makeFileHandler(dbtx_val))
	public_subrouter.HandleFunc("/files/{uuid}/markdown", makeMarkdownHandler(dbtx_val))
	// private_subrouter := router.PathPrefix("/private").Subrouter()
	// private_subrouter.HandleFunc("/files/{uuid}", getPrivateFileHandler)
}

func makeFileHandler(dbtx_val dbstore.DBTX) func(w http.ResponseWriter, r *http.Request) {
	return_func := func(w http.ResponseWriter, r *http.Request) {
		q := *dbstore.New(dbtx_val)
		params := mux.Vars(r)
		fileID := params["uuid"]
		parsedUUID, err := uuid.Parse(fileID)
		if err != nil {
			http.Error(w, "Invalid File ID format", http.StatusBadRequest)
			return
		}
		pgUUID := pgtype.UUID{Bytes: parsedUUID, Valid: true}
		ctx := r.Context()

		file, err := q.ReadFile(ctx, pgUUID)
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		// fileSchema := fileToSchema(file)
		fileSchema := file
		response, _ := json.Marshal(fileSchema)

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
	return return_func
}

func makeMarkdownHandler(dbtx_val dbstore.DBTX) func(w http.ResponseWriter, r *http.Request) {
	return_func := func(w http.ResponseWriter, r *http.Request) {
		q := *dbstore.New(dbtx_val)
		params := mux.Vars(r)
		fileID := params["uuid"]

		parsedUUID, err := uuid.Parse(fileID) // hypothetical helper function
		if err != nil {
			http.Error(w, "Invalid File ID format", http.StatusBadRequest)
			return
		}
		pgUUID := pgtype.UUID{Bytes: parsedUUID, Valid: true}

		originalLang := r.URL.Query().Get("original_lang") == "true"
		matchLang := r.URL.Query().Get("match_lang")

		ctx := r.Context()

		var texts []dbstore.FileTextSource
		if originalLang {
			texts, err = q.ListTextsOfFileOriginal(ctx, pgUUID)
		} else if matchLang != "" {
			texts, err = q.ListTextsOfFileWithLanguage(ctx, dbstore.ListTextsOfFileWithLanguageParams{
				FileID:   pgUUID,
				Language: matchLang,
			})
		} else {
			texts, err = q.ListTextsOfFile(ctx, pgUUID)
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
		w.Write([]byte(markdownText.String))
	}
	return return_func
}
