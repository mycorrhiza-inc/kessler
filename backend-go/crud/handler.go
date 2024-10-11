package crud

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mycorrhiza-inc/kessler/backend-go/gen/dbstore"
)

func DefineCrudRoutes(router *mux.Router, dbtx_val dbstore.DBTX) {
	public_subrouter := router.PathPrefix("/public").Subrouter()
	public_subrouter.HandleFunc("/files/{uuid}", makeFileHandler(dbtx_val))
	public_subrouter.HandleFunc("/files/{uuid}/markdown", makeMarkdownHandler(dbtx_val))
	// private_subrouter := router.PathPrefix("/private").Subrouter()
	// private_subrouter.HandleFunc("/files/{uuid}", getPrivateFileHandler)
}

func makeFileHandler(dbtx_val dbstore.DBTX, private bool) func(w http.ResponseWriter, r *http.Request) {
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

		if !private {
			file, err := q.ReadFile(ctx, pgUUID)
		} else {
			file, err := q.ReadPrivateFile(ctx, pgUUID)
		}

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

func makeUpsertHandler(dbtx_val dbstore.DBTX, private bool) func(w http.ResponseWriter, r *http.Request) {
	return_func := func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(token, "Authorized ") {
			return UserValidation{validated: false}
		}
		userID , err :=
		if !private && userID !="thaumaturgy"{

		}
	}
	return return_func
}
