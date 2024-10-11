package crud

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mycorrhiza-inc/kessler/backend-go/gen/dbstore"
)

func DefineCrudRoutes(router *mux.Router, dbtx_val dbstore.DBTX) {
	public_subrouter := router.PathPrefix("/public").Subrouter()

  pub_file := 
	public_subrouter.HandleFunc("/files/{uuid}", makeFileHandler(FileHandlerInfo{dbtx_val: dbtx_val, private: false, return_type: "object"}))
	public_subrouter.HandleFunc("/files/{uuid}/markdown", makeFileHandler(FileHandlerInfo{dbtx_val: dbtx_val, private: false, markdown: false}))
	// private_subrouter := router.PathPrefix("/private").Subrouter()
	// private_subrouter.HandleFunc("/files/{uuid}", getPrivateFileHandler)
}

type FileHandlerInfo{
	dbtx_val dbstore.DBTX
	private bool
  return_type string
}

func makeFileHandler(info FileHandlerInfo) func(w http.ResponseWriter, r *http.Request) {
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
		}

		getFile := func(uuid pgtype.UUID) (rawFileSchema, error) {
			if !private {
				file, err := q.ReadFile(ctx, pgUUID)
				if err != nil {
					return rawFileSchema{}, err
				}
				return PublicFileToSchema(file), nil
			} else {
				file, err := q.ReadPrivateFile(ctx, pgUUID)
				if err != nil {
					return rawFileSchema{}, err
				}
				return PrivateFileToSchema(file), nil
			}
		}

		file, err := getFile(pgUUID)
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		// fileSchema := fileToSchema(file)
		fileSchema, err := RawToFileSchema(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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

type UpsertHandlerInfo struct {
  dbtx_val dbstore.DBTX
  private bool
	insert bool
}

func makeUpsertHandler(info UpsertHandlerInfo) func(w http.ResponseWriter, r *http.Request) {
	return_func := func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(token, "Authorized ") {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		if !private && token != "Authorized thaumaturgy" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
	}
	return return_func
}
