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
	public_subrouter.HandleFunc("/files/{uuid}/raw", makeFileHandler(
		FileHandlerInfo{dbtx_val: dbtx_val, private: false, return_type: "raw"}))
	private_subrouter := router.PathPrefix("/private").Subrouter()
	private_subrouter.HandleFunc("/files/{uuid}", makeFileHandler(
		FileHandlerInfo{dbtx_val: dbtx_val, private: true, return_type: "object"}))
	private_subrouter.HandleFunc("/files/{uuid}/markdown", makeFileHandler(
		FileHandlerInfo{dbtx_val: dbtx_val, private: true, return_type: "markdown"}))
	private_subrouter.HandleFunc("/files/{uuid}/raw", makeFileHandler(
		FileHandlerInfo{dbtx_val: dbtx_val, private: true, return_type: "raw"}))
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
	return func(w http.ResponseWriter, r *http.Request) {
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
		// Since all three of these methods share the same authentication and database connection prerecs
		// switching functionality using an if else, or a cases switch lets code get reused
		// TODO: This is horrible, I need to refactor
		file_params := GetFileParam{
			q, ctx, pgUUID, private,
		}
		switch return_type {
		case "raw":
			http.Error(w, "Retriving raw files from s3 not implemented", http.StatusNotImplemented)
		case "object":
			file, err := GetFileObjectRaw(file_params)
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
		case "markdown":
			texts, err := GetTextSchemas(file_params)
			if err != nil || len(texts) == 0 {
				http.Error(w, "Error retrieving texts or no texts found.", http.StatusInternalServerError)
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
}

type UpsertHandlerInfo struct {
	dbtx_val dbstore.DBTX
	private  bool
	insert   bool
}
type DocTextInfo struct {
	Language       string `json:"language"`
	Text           string `json:"text"`
	IsOriginalText bool   `json:"is_original_text"`
}

type UpdateDocumentInfo struct {
	Url          string            `json:"url"`
	Doctype      string            `json:"doctype"`
	Lang         string            `json:"lang"`
	Name         string            `json:"name"`
	Source       string            `json:"source"`
	Hash         string            `json:"hash"`
	Mdata        map[string]string `json:"mdata"`
	Stage        string            `json:"stage"`
	Summary      string            `json:"summary"`
	ShortSummary string            `json:"short_summary"`
	Private      bool              `json:"private"`
	DocTexts     []DocTextInfo     `json:"doc_texts"`
}

func ConvertToCreationData(updateInfo UpdateDocumentInfo) (FileCreationDataRaw, error) {
	mdata_string, err := json.Marshal(updateInfo.Mdata)
	if err != nil {
		return FileCreationDataRaw{}, nil
	}
	creationData := FileCreationDataRaw{
		Url:          pgtype.Text{String: updateInfo.Url, Valid: true},
		Doctype:      pgtype.Text{String: updateInfo.Doctype, Valid: true},
		Lang:         pgtype.Text{String: updateInfo.Lang, Valid: true},
		Name:         pgtype.Text{String: updateInfo.Name, Valid: true},
		Source:       pgtype.Text{String: updateInfo.Source, Valid: true},
		Hash:         pgtype.Text{String: updateInfo.Hash, Valid: true},
		Stage:        pgtype.Text{String: updateInfo.Stage, Valid: true},
		Summary:      pgtype.Text{String: updateInfo.Summary, Valid: true},
		ShortSummary: pgtype.Text{String: updateInfo.ShortSummary, Valid: true},
		Mdata:        pgtype.Text{String: string(mdata_string), Valid: true},
	}
	return creationData, nil
}

func makeUpsertHandler(info UpsertHandlerInfo) func(w http.ResponseWriter, r *http.Request) {
	dbtx_val := info.dbtx_val
	private := info.private
	insert := info.insert
	return func(w http.ResponseWriter, r *http.Request) {
		var doc_uuid uuid.UUID
		var err error
		if !insert {
			params := mux.Vars(r)
			fileID := params["uuid"]

			doc_uuid, err = uuid.Parse(fileID)
			if err != nil {
				http.Error(w, "Error parsing uuid", http.StatusBadRequest)
				return
			}
		}
		q := *dbstore.New(dbtx_val)
		ctx := r.Context()
		token := r.Header.Get("Authorization")
		if !strings.HasPrefix(token, "Authorized ") {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		userID := strings.TrimPrefix(token, "Authorized ")
		forbiddenPublic := !private && userID != "thaumaturgy"
		if forbiddenPublic || userID == "anonomous" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		if !insert {
			authorized, err := checkPrivateFileAuthorization(q, ctx, doc_uuid, userID)
			if !authorized || err == nil {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
		}
		// TODO: IF user is not a paying user, disable insert functionality
		var newDocInfo UpdateDocumentInfo
		if err := json.NewDecoder(r.Body).Decode(&newDocInfo); err != nil {

			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		w.Write([]byte("Sucessfully inserted"))
	}
}
