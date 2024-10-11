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

	public_subrouter.HandleFunc("/files/{uuid}", makeFileHandler(
		FileHandlerInfo{dbtx_val: dbtx_val, private: false, return_type: "object"}))
	public_subrouter.HandleFunc("/files/{uuid}/markdown", makeFileHandler(
		FileHandlerInfo{dbtx_val: dbtx_val, private: false, return_type: "markdown"}))
	private_subrouter := router.PathPrefix("/private").Subrouter()
	private_subrouter.HandleFunc("/files/{uuid}", makeFileHandler(
		FileHandlerInfo{dbtx_val: dbtx_val, private: true, return_type: "object"}))
	private_subrouter.HandleFunc("/files/{uuid}/markdown", makeFileHandler(
		FileHandlerInfo{dbtx_val: dbtx_val, private: true, return_type: "markdown"}))
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

			ctx := r.Context()

			get_file_texts := func(pgUUID pgtype.UUID) ([]FileTextSchema, error) {
				if private {
					texts, err := q.ListPrivateTextsOfFile(ctx, pgUUID)
					if err != nil {
						http.Error(w, "Error retrieving texts", http.StatusInternalServerError)
						return []FileTextSchema{}, err
					}
					schemas := make([]FileTextSchema, len(texts))
					for i, text := range texts {
						schemas[i] = PrivateTextToSchema(text)
					}
					return schemas, nil
				}
				texts, err := q.ListTextsOfFile(ctx, pgUUID)
				if err != nil {
					http.Error(w, "Error retrieving texts", http.StatusInternalServerError)
					return []FileTextSchema{}, err
				}
				schemas := make([]FileTextSchema, len(texts))
				for i, text := range texts {
					schemas[i] = PublicTextToSchema(text)
				}
				return schemas, nil
			}
			texts, err := get_file_texts(pgUUID)
			if err != nil {
				http.Error(w, "Error retrieving texts", http.StatusInternalServerError)
				return
			}
			if len(texts) == 0 {
				http.Error(w, "No texts found for document", http.StatusNotFound)
				return
			}
			// TODO: Add suport for non english text retrieval and original text retrieval
			// originalLang := r.URL.Query().Get("original_lang") == "true"
			// matchLang := r.URL.Query().Get("match_lang")
			markdownText := texts[0].Text
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte(markdownText))
		}
	}
	return return_func
}

type UpsertHandlerInfo struct {
	dbtx_val dbstore.DBTX
	private  bool
	insert   bool
}

// func makeUpsertHandler(info UpsertHandlerInfo) func(w http.ResponseWriter, r *http.Request) {
// 	dbtx_val := info.dbtx_val
// 	private := info.private
// 	return_func := func(w http.ResponseWriter, r *http.Request) {
// 		q := *dbstore.New(dbtx_val)
// 		ctx := r.Context()
// 		token := r.Header.Get("Authorization")
// 		if !strings.HasPrefix(token, "Authorized ") {
// 			http.Error(w, "Forbidden", http.StatusForbidden)
// 			return
// 		}
// 		userID := strings.TrimPrefix(token, "Authorized ")
// 		if !private && userID != "thaumaturgy" {
// 			http.Error(w, "Forbidden", http.StatusForbidden)
// 			return
// 		}
// 		return
// 	}
// 	return return_func
// }
