package crud

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mycorrhiza-inc/kessler/backend-go/gen/dbstore"
)

func DefineCrudRoutes(router *mux.Router, dbtx_val dbstore.DBTX) {
	public_subrouter := router.PathPrefix("/public").Subrouter()

	pub_file := public_subrouter.HandleFunc("/files/{uuid}", makeFileHandler(FileHandlerInfo{dbtx_val: dbtx_val, private: false, return_type: "object"}))
	public_subrouter.HandleFunc("/files/{uuid}/markdown", makeFileHandler(FileHandlerInfo{dbtx_val: dbtx_val, private: false, markdown: false}))
	// private_subrouter := router.PathPrefix("/private").Subrouter()
	// private_subrouter.HandleFunc("/files/{uuid}", getPrivateFileHandler)
}

type FileHandlerInfo struct {
	dbtx_val    dbstore.DBTX
	private     bool
	return_type string // Can be either markdown or object
}

func checkPrivateFileAuthorization(q dbstore.Queries, ctx context.Context, objectID uuid.UUID, viewerID string) (bool, error) {
	if viewerID == "thaumaturgy" {
		return true, nil
	}
	viewerUUID, err := uuid.Parse(viewerID)
	if err != nil {
		return false, err
	}
	viewerPgUUID := pgtype.UUID{Bytes: viewerUUID, Valid: true}
	objectPgUUID := pgtype.UUID{Bytes: objectID, Valid: true}
	auth_params := dbstore.CheckOperatorAccessToObjectParams{viewerPgUUID, objectPgUUID}
	auth_result, err := q.CheckOperatorAccessToObject(ctx, auth_params)
	if err != nil {
		return false, err
	}
	return auth_result, nil
}

func makeFileHandler(info FileHandlerInfo) func(w http.ResponseWriter, r *http.Request) {
	private := info.private
	dbtx_val := info.dbtx_val
	return_type := info.return_type
	return_func := func(w http.ResponseWriter, r *http.Request) {
		q := *dbstore.New(dbtx_val)
		token := r.Header.Get("Authorization")
		params := mux.Vars(r)
		fileID := params["uuid"]
		parsedUUID, err := uuid.Parse(fileID)
		if err != nil {
			http.Error(w, "Invalid File ID format", http.StatusBadRequest)
			return
		}
		pgUUID := pgtype.UUID{Bytes: parsedUUID, Valid: true}
		ctx := r.Context()
		if private {

			userID := strings.TrimPrefix(token, "Authorized ")
			isAuthorized, err := checkPrivateFileAuthorization(q, ctx, parsedUUID, userID)
			if !isAuthorized {
				http.Error(w, "Forbidden", http.StatusForbidden)
			}
			if err != nil {
				fmt.Printf("Ran into the follwing error with authentication $v", err)
			}
		}
		// TODO: This is horrible, I need to refactor
		if return_type == "object" {
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
		} else if return_type == "markdown" {
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
	}
	return return_func
}

type UpsertHandlerInfo struct {
	dbtx_val dbstore.DBTX
	private  bool
	insert   bool
}

func makeUpsertHandler(info UpsertHandlerInfo) func(w http.ResponseWriter, r *http.Request) {
	dbtx_val := info.dbtx_val
	private := info.private
	return_func := func(w http.ResponseWriter, r *http.Request) {
		q := *dbstore.New(dbtx_val)
		ctx := r.Context()
		token := r.Header.Get("Authorization")
		if !strings.HasPrefix(token, "Authorized ") {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		userID := strings.TrimPrefix(token, "Authorized ")
		if !private && userID != "thaumaturgy" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		private_auth := checkPrivateFileAuthorization(q, ctx)
	}
	return return_func
}
