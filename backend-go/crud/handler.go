package crud

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mycorrhiza-inc/kessler/backend-go/gen/dbstore"
)

func DefineCrudRoutes(router *mux.Router, dbtx_val dbstore.DBTX) {
	public_subrouter := router.PathPrefix("/api/v2/public").Subrouter()

	public_subrouter.HandleFunc("/files/insert", makeUpsertHandler(
		UpsertHandlerInfo{dbtx_val: dbtx_val, private: false, insert: true})).Methods(http.MethodPost)

	public_subrouter.HandleFunc("/files/{uuid}", makeUpsertHandler(
		UpsertHandlerInfo{dbtx_val: dbtx_val, private: false, insert: false})).Methods(http.MethodPost)

	public_subrouter.HandleFunc("/files/{uuid}", makeFileHandler(
		FileHandlerInfo{dbtx_val: dbtx_val, private: false, return_type: "object"})).Methods(http.MethodGet)

	public_subrouter.HandleFunc("/files/{uuid}/markdown", makeFileHandler(
		FileHandlerInfo{dbtx_val: dbtx_val, private: false, return_type: "markdown"})).Methods(http.MethodGet)

	public_subrouter.HandleFunc("/files/{uuid}/raw", makeFileHandler(
		FileHandlerInfo{dbtx_val: dbtx_val, private: false, return_type: "raw"})).Methods(http.MethodGet)

	private_subrouter := router.PathPrefix("/api/v2/private").Subrouter()

	private_subrouter.HandleFunc("/files/insert", makeUpsertHandler(
		UpsertHandlerInfo{dbtx_val: dbtx_val, private: true, insert: true})).Methods(http.MethodPost)

	private_subrouter.HandleFunc("/files/{uuid}", makeUpsertHandler(
		UpsertHandlerInfo{dbtx_val: dbtx_val, private: true, insert: false})).Methods(http.MethodPost)

	private_subrouter.HandleFunc("/files/{uuid}", makeFileHandler(
		FileHandlerInfo{dbtx_val: dbtx_val, private: true, return_type: "object"})).Methods(http.MethodGet)

	private_subrouter.HandleFunc("/files/{uuid}/markdown", makeFileHandler(
		FileHandlerInfo{dbtx_val: dbtx_val, private: true, return_type: "markdown"})).Methods(http.MethodGet)

	private_subrouter.HandleFunc("/files/{uuid}/raw", makeFileHandler(
		FileHandlerInfo{dbtx_val: dbtx_val, private: true, return_type: "raw"})).Methods(http.MethodGet)
}

// CONVERT TO UPPER CASE IF YOU EVER WANT TO USE IT OUTSIDE OF THIS CONTEXT
type FileHandlerInfo struct {
	dbtx_val    dbstore.DBTX
	private     bool
	return_type string // Can be either markdown, object or raw
}

func checkPrivateFileAuthorization(q dbstore.Queries, ctx context.Context, objectID uuid.UUID, token string) (bool, error) {
	if !strings.HasPrefix(token, "Authenticated") {
		return false, nil
	}
	viewerID := strings.TrimPrefix(token, "Authenticated ")
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

			isAuthorized, err := checkPrivateFileAuthorization(q, ctx, parsedUUID, token)
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
			// Please go into the file, try to infer the mime type, and then return the file in binary to the user, using something like
			// w.Header().Set("Content-Type", <Whatever MIME TYPE>)
			// w.Write(file binary)
			file, err := GetFileObjectRaw(file_params)
			if err != nil {
				http.Error(w, "File not found", http.StatusNotFound)
				return
			}
			filehash := file.Hash
			kefiles := NewKeFileManager()
			file_path, err := kefiles.downloadFileFromS3(filehash)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error encountered when getting file from s3:%v", err), http.StatusInternalServerError)
				return
			}
			content, err := ioutil.ReadFile(file_path)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error reading file: %v", err), http.StatusInternalServerError)
				return
			}

			mimeType := http.DetectContentType(content)
			// if mimeType == "application/octet-stream" {
			// 	mimeType = "application/pdf" // Default to PDF if mime type can't be determined
			// }

			w.Header().Set("Content-Type", mimeType)
			w.Write(content)
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
			originalLang := r.URL.Query().Get("original_lang") == "true"
			matchLang := r.URL.Query().Get("match_lang")
			// TODO: Add suport for non english text retrieval and original text retrieval
			markdownText, err := GetSpecificFileText(file_params, matchLang, originalLang)
			if err != nil {
				http.Error(w, "Error retrieving texts or no texts found that mach query params", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte(markdownText))
		default:
			fmt.Printf("Encountered unreachable code with file return type %v", return_type)
			http.Error(w, "Congradulations for encountering unreachable code about support types!", http.StatusInternalServerError)
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
	ID           uuid.UUID      `json:"id"`
	Url          string         `json:"url"`
	Doctype      string         `json:"doctype"`
	Lang         string         `json:"lang"`
	Name         string         `json:"name"`
	Source       string         `json:"source"`
	Hash         string         `json:"hash"`
	Mdata        map[string]any `json:"mdata"`
	Stage        string         `json:"stage"`
	Summary      string         `json:"summary"`
	ShortSummary string         `json:"short_summary"`
	Private      bool           `json:"private"`
	DocTexts     []DocTextInfo  `json:"doc_texts"`
}

func ConvertToCreationData(updateInfo UpdateDocumentInfo) FileCreationDataRaw {
	mdata_string, _ := json.Marshal(updateInfo.Mdata)
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
	return creationData
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
		if !strings.HasPrefix(token, "Authenticated ") {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		userID := strings.TrimPrefix(token, "Authenticated ")
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
		rawFileData := ConvertToCreationData(newDocInfo)
		var fileSchema RawFileSchema
		if insert {
			fileSchema, err = InsertPubPrivateFileObj(q, ctx, rawFileData, private)
		} else {
			pgUUID := pgtype.UUID{doc_uuid, true}
			fileSchema, err = UpdatePubPrivateFileObj(q, ctx, rawFileData, private, pgUUID)
		}
		if err != nil {
			http.Error(w, "Error inserting/updating document", http.StatusInternalServerError)
		}
		texts := newDocInfo.DocTexts
		doc_uuid = fileSchema.ID.Bytes // Ensure UUID is same as one returned from database
		doc_pgUUID := pgtype.UUID{Bytes: doc_uuid, Valid: true}
		if len(texts) != 0 {
			if !insert {
				// TODO: Implement this func to Nuke all the previous texts
				err := NukePriPubFileTexts(q, ctx, doc_pgUUID)
				if err != nil {
					fmt.Print("Error deleting old texts, proceeding with new editions")
				}
			}
			// TODO : Make Async at some point in future
			for _, text := range texts {
				textRaw := FileTextSchema{
					FileID:         doc_pgUUID,
					IsOriginalText: text.IsOriginalText,
					Language:       text.Language,
					Text:           text.Text,
				}
				err = InsertPriPubFileText(q, ctx, textRaw, private)
				if err != nil {
					fmt.Print("Error adding a text value, not doing anything and procceeding since error handling is hard.")
				}
			}
		}
		response, _ := json.Marshal(fileSchema)

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}

func checkPrivateUploadAbility(token string) bool {
	if !strings.HasPrefix(token, "Authenticated ") {
		return false
	}
	viewerID := strings.TrimPrefix(token, "Authenticated ")
	return viewerID != "anonomous" && viewerID != "thaumaturgy"
}

type FileProcessRequest struct {
	ID                 uuid.UUID         `json:"id"`
	Hash               string            `json:"hash"`
	Private            bool              `json:"private"`
	FileUploadName     string            `json:"file_upload_name"`
	UserID             uuid.UUID         `json:"user_id"`
	Mdata              map[string]string `json:"mdata"`
	ExistingFileSchema FileSchema        `json:"existing_file_schema"`
}

func sendFileProcessRequest(req FileProcessRequest) error {
	return nil
}

func generateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func makePrivateUploadHandler(dbtx_val dbstore.DBTX) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		file, _, err := r.FormFile("file")
		fileName := r.FormValue("file_name")
		if err != nil {
			return
		}
		defer file.Close()
		randomFileName := generateRandomString(10) // Function to generate a random string
		f, err := os.OpenFile(randomFileName, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return
		}
		defer f.Close()
		io.Copy(f, file)
		keFileMan := NewKeFileManager()
		hash, err := keFileMan.uploadFileToS3(randomFileName)
		if err != nil {
			fmt.Printf("Error uploading to s3, %v", err)
		}
		fmt.Fprintf(w, "File %s uploaded successfully with hash %s", fileName, hash)
	}
}

type ReturnFilesSchema struct {
	Files []FileSchema `json:"files"`
}

func GetListAllRawFiles(ctx context.Context, q dbstore.Queries) ([]RawFileSchema, error) {
	files, err := q.ListFiles(ctx)
	if err != nil {
		return []RawFileSchema{}, err
	}
	var fileSchemas []RawFileSchema
	for _, fileRaw := range files {
		rawSchema := PublicFileToSchema(fileRaw)
		fileSchemas = append(fileSchemas, rawSchema)
	}
	return fileSchemas, nil
}

func GetListAllFiles(ctx context.Context, q dbstore.Queries) ([]FileSchema, error) {
	files, err := GetListAllRawFiles(ctx, q)
	if err != nil {
		return []FileSchema{}, err
	}
	var fileSchemas []FileSchema
	for _, rawSchema := range files {
		fileSchema, _ := RawToFileSchema(rawSchema)
		fileSchemas = append(fileSchemas, fileSchema)
	}
	return fileSchemas, nil
}

func getListOfAllPublicFilesHandler(dbtx_val dbstore.DBTX) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		q := *dbstore.New(dbtx_val)
		ctx := r.Context()
		fileSchemas, err := GetListAllFiles(ctx, q)
		if err != nil {
			http.Error(w, "Encountered db error reading files", http.StatusInternalServerError)
			return
		}
		return_schema := ReturnFilesSchema{Files: fileSchemas}
		response, _ := json.Marshal(return_schema)

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}
