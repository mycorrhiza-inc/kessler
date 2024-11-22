package crud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mycorrhiza-inc/kessler/backend-go/gen/dbstore"
)

func GetFileWithMeta(config FileHandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := *dbstore.New(config.dbtx_val)
		params := mux.Vars(r)
		fileID := params["uuid"]
		parsedUUID, err := uuid.Parse(fileID)
		if err != nil {
			errorstring := fmt.Sprintf("Error parsing file %v: %v\n", fileID, err)
			fmt.Println(errorstring)
			http.Error(w, errorstring, http.StatusBadRequest)
			return
		}
		pgUUID := pgtype.UUID{Bytes: parsedUUID, Valid: true}
		ctx := r.Context()
		file_raw, err := q.GetFileWithMetadata(ctx, pgUUID)
		if err != nil {

			errorstring := fmt.Sprintf("Error Retriving file %v: %v\n", fileID, err)
			fmt.Println(errorstring)
			http.Error(w, errorstring, http.StatusNotFound)
			return
		}
		var mdata_obj map[string]interface{}
		err = json.Unmarshal(file_raw.Mdata, &mdata_obj)
		if err != nil {
			errorstring := fmt.Sprintf("Error Unmarshalling file %v: %v\n", fileID, err)
			fmt.Println(errorstring)
			http.Error(w, errorstring, http.StatusInternalServerError)
			return
		}
		file := CompleteFileSchema{
			ID:        file_raw.ID.Bytes,
			Verified:  file_raw.Verified.Bool,
			Extension: file_raw.Extension.String,
			Lang:      file_raw.Lang.String,
			Name:      file_raw.Name.String,
			Hash:      file_raw.Hash.String,
			IsPrivate: file_raw.Isprivate.Bool,
			Mdata:     mdata_obj,
		}

		response, _ := json.Marshal(file)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}

func FileSemiCompleteGetFactory(dbtx_val dbstore.DBTX) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := *dbstore.New(dbtx_val)
		params := mux.Vars(r)
		fileID := params["uuid"]
		parsedUUID, err := uuid.Parse(fileID)
		if err != nil {
			errorstring := fmt.Sprintf("Error parsing file %v: %v\n", fileID, err)
			fmt.Println(errorstring)
			http.Error(w, errorstring, http.StatusBadRequest)
			return
		}
		pgUUID := pgtype.UUID{Bytes: parsedUUID, Valid: true}
		ctx := r.Context()
		files_raw, err := q.SemiCompleteFileGet(ctx, pgUUID)
		if err != nil {
			errorstring := fmt.Sprintf("Error Retriving file %v: %v\n", fileID, err)
			fmt.Println(errorstring)
			http.Error(w, errorstring, http.StatusNotFound)
			return
		}
		if len(files_raw) == 0 {
			errorstring := fmt.Sprintf("Error No Files Found for a list of length zero.\n")
			fmt.Println(errorstring)
			http.Error(w, errorstring, http.StatusNotFound)
			return
		}
		file_raw := files_raw[0]
		var mdata_obj map[string]interface{}
		err = json.Unmarshal(file_raw.Mdata, &mdata_obj)
		if err != nil {
			errorstring := fmt.Sprintf("Error Unmarshalling file metadata %v: %v\n", fileID, err)
			fmt.Println(errorstring)
			http.Error(w, errorstring, http.StatusInternalServerError)
			return
		}
		var extra_obj FileGeneratedExtras
		err = json.Unmarshal(file_raw.ExtraObj, &extra_obj)
		if err != nil {
			errorstring := fmt.Sprintf("Error Unmarshalling file extras %v: %v\n", fileID, err)
			fmt.Println(errorstring)
			http.Error(w, errorstring, http.StatusInternalServerError)
			return
		}
		// Missing info here, it doesnt have the name.
		conv_info := ConversationInformation{ID: file_raw.DocketUuid.Bytes}

		file := CompleteFileSchema{
			ID:           file_raw.ID.Bytes,
			Verified:     file_raw.Verified.Bool,
			Extension:    file_raw.Extension.String,
			Lang:         file_raw.Lang.String,
			Name:         file_raw.Name.String,
			Hash:         file_raw.Hash.String,
			Mdata:        mdata_obj,
			Extra:        extra_obj,
			Conversation: conv_info,
		}

		response, _ := json.Marshal(file)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}

func ReadFileHandlerFactory(config FileHandlerConfig) http.HandlerFunc {
	private := config.private
	dbtx_val := config.dbtx_val
	return_type := config.return_type

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
			file, err := GetFileObjectRaw(file_params)
			if err != nil {
				error_string := fmt.Sprintf("Error retrieving file object %v", err)
				fmt.Println(error_string)
				http.Error(w, error_string, http.StatusNotFound)
				return
			}
			filehash := file.Hash
			kefiles := NewKeFileManager()
			file_path, err := kefiles.downloadFileFromS3(filehash)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error encountered when getting file with hash %v from s3:%v", filehash, err), http.StatusInternalServerError)
				return
			}
			content, err := os.ReadFile(file_path)
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
		case "markdown":
			originalLang := r.URL.Query().Get("original_lang") == "true"
			matchLang := r.URL.Query().Get("match_lang")
			// TODO: Add suport for non english text retrieval and original text retrieval
			markdownText, err := GetSpecificFileText(file_params, matchLang, originalLang)
			if err != nil {
				http.Error(w, "Error retrieving texts or no texts found that mach query params", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte(markdownText))
		case "object-minimal":
			file, err := GetFileObjectRaw(file_params)
			if err != nil {
				http.Error(w, "File not found", http.StatusNotFound)
				return
			}
			inflated_schema := CompleteFileSchemaInflateFromPartialSchema(file)

			response, _ := json.Marshal(inflated_schema)

			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
		default:
			fmt.Printf("Encountered unreachable code with file return type %v", return_type)
			http.Error(w, "Congradulations for encountering unreachable code about support types!", http.StatusInternalServerError)
		}
	}
}
