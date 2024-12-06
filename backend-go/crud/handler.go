package crud

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"kessler/gen/dbstore"
	"kessler/objects/files"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func DefineCrudRoutes(public_subrouter *mux.Router, dbtx_val dbstore.DBTX) {
	public_subrouter.HandleFunc(
		"/files/insert",
		makeFileUpsertHandler(
			FileUpsertHandlerConfig{
				dbtx_val: dbtx_val,
				private:  false,
				insert:   true,
			},
		)).Methods(http.MethodPost)

	public_subrouter.HandleFunc(
		"/files/{uuid}/update",
		makeFileUpsertHandler(
			FileUpsertHandlerConfig{
				dbtx_val: dbtx_val,
				private:  false,
				insert:   false,
			},
		)).Methods(http.MethodPost)
	public_subrouter.HandleFunc(
		"/files/{uuid}",
		FileSemiCompleteGetFactory(dbtx_val),
	).Methods(http.MethodGet)

	public_subrouter.HandleFunc(
		"/files/{uuid}/minimal",
		ReadFileHandlerFactory(
			FileHandlerConfig{
				dbtx_val:    dbtx_val,
				private:     false,
				return_type: "object-minimal",
			},
		)).Methods(http.MethodGet)

	public_subrouter.HandleFunc(
		"/files/{uuid}/markdown",
		ReadFileHandlerFactory(
			FileHandlerConfig{
				dbtx_val:    dbtx_val,
				private:     false,
				return_type: "markdown",
			},
		)).Methods(http.MethodGet)
	// These shouldnt have to be duplicated, but such is life.
	public_subrouter.HandleFunc(
		"/files/{uuid}/raw/{filename}",
		ReadFileHandlerFactory(
			FileHandlerConfig{
				dbtx_val:    dbtx_val,
				private:     false,
				return_type: "raw",
			},
		)).Methods(http.MethodGet)

	public_subrouter.HandleFunc(
		"/files/{uuid}/raw",
		ReadFileHandlerFactory(
			FileHandlerConfig{
				dbtx_val:    dbtx_val,
				private:     false,
				return_type: "raw",
			},
		)).Methods(http.MethodGet)

	public_subrouter.HandleFunc(
		"/files/{uuid}/metadata",
		GetFileWithMeta(
			FileHandlerConfig{
				dbtx_val: dbtx_val,
				private:  false,
			},
		)).Methods(http.MethodGet)
	// TODO : Split out the organizations into their own crud handler module
	public_subrouter.HandleFunc(
		"/organizations/list",
		OrgListAllFactory(dbtx_val),
	).Methods(http.MethodGet)

	public_subrouter.HandleFunc(
		"/conversations/list",
		ConversationListAllFactory(dbtx_val),
	).Methods(http.MethodGet)

	public_subrouter.HandleFunc(
		"/conversations/named-lookup/{name}",
		ConversationGetByNameFactory(dbtx_val),
	).Methods(http.MethodGet)

	public_subrouter.HandleFunc(
		"/organizations/{uuid}",
		GetOrgWithFilesFactory(dbtx_val),
	).Methods(http.MethodGet)

	public_subrouter.HandleFunc(
		"/organizations/verify",
		OrganizationVerifyHandlerFactory(dbtx_val),
	).Methods(http.MethodPost)
	public_subrouter.HandleFunc(
		"/conversations/verify",
		ConversationVerifyHandlerFactory(dbtx_val),
	).Methods(http.MethodPost)

	public_subrouter.HandleFunc(
		"/conversations/list/semi-complete",
		ConversationSemiCompleteListAllFactory(dbtx_val),
	).Methods(http.MethodGet)

	//
	// private_subrouter := router.PathPrefix("/v2/private").Subrouter()
	//
	// private_subrouter.HandleFunc("/files/insert", makeFileUpsertHandler(
	// 	UpsertHandlerInfo{dbtx_val: dbtx_val, private: true, insert: true})).Methods(http.MethodPost)
	//
	// private_subrouter.HandleFunc("/files/{uuid}", makeFileUpsertHandler(
	// 	UpsertHandlerInfo{dbtx_val: dbtx_val, private: true, insert: false})).Methods(http.MethodPost)
	//
	// private_subrouter.HandleFunc("/files/{uuid}", makeReadFileHandler(
	// 	FileHandlerInfo{dbtx_val: dbtx_val, private: true, return_type: "object"})).Methods(http.MethodGet)
	//
	// private_subrouter.HandleFunc("/files/{uuid}/markdown", makeReadFileHandler(
	// 	FileHandlerInfo{dbtx_val: dbtx_val, private: true, return_type: "markdown"})).Methods(http.MethodGet)
	//
	// private_subrouter.HandleFunc("/files/{uuid}/raw", makeReadFileHandler(
	// 	FileHandlerInfo{dbtx_val: dbtx_val, private: true, return_type: "raw"})).Methods(http.MethodGet)
}

// CONVERT TO UPPER CASE IF YOU EVER WANT TO USE IT OUTSIDE OF THIS CONTEXT
type FileHandlerConfig struct {
	dbtx_val    dbstore.DBTX
	private     bool
	return_type string // Can be either markdown, object or raw
}

func checkPrivateFileAuthorization(q dbstore.Queries, ctx context.Context, objectID uuid.UUID, token string) (bool, error) {
	if !strings.HasPrefix(token, "Authenticated") {
		return false, nil
	}
	// viewerID := strings.TrimPrefix(token, "Authenticated ")
	// if viewerID == "thaumaturgy" {
	// 	return true, nil
	// }
	// viewerUUID, err := uuid.Parse(viewerID)
	// if err != nil {
	// 	return false, err
	// }
	// viewerPgUUID := pgtype.UUID{Bytes: viewerUUID, Valid: true}
	// objectPgUUID := pgtype.UUID{Bytes: objectID, Valid: true}
	// auth_params := dbstore.CheckOperatorAccessToObjectParams{viewerPgUUID, objectPgUUID}
	// auth_result, err := q.CheckOperatorAccessToObject(ctx, auth_params)
	// if err != nil {
	// 	return false, err
	// }
	auth_result := false
	return auth_result, nil
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
	ExistingFileSchema files.FileSchema  `json:"existing_file_schema"`
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

func privateUploadFactory(dbtx_val dbstore.DBTX) func(w http.ResponseWriter, r *http.Request) {
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
	Files []files.FileSchema `json:"files"`
}

func GetListAllRawFiles(ctx context.Context, q dbstore.Queries) ([]files.FileSchema, error) {
	db_files, err := q.FilesList(ctx)
	if err != nil {
		return []files.FileSchema{}, err
	}
	var fileSchemas []files.FileSchema
	for _, fileRaw := range db_files {
		rawSchema := files.PublicFileToSchema(fileRaw)
		fileSchemas = append(fileSchemas, rawSchema)
	}
	return fileSchemas, nil
}

func GetListAllFiles(ctx context.Context, q dbstore.Queries) ([]files.FileSchema, error) {
	db_files, err := GetListAllRawFiles(ctx, q)
	if err != nil {
		return []files.FileSchema{}, err
	}
	var fileSchemas []files.FileSchema
	for _, rawSchema := range db_files {
		fileSchema := rawSchema
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
