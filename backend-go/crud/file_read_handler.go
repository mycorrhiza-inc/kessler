package crud

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"kessler/common/hashes"
	"kessler/common/objects/authors"
	"kessler/common/objects/conversations"
	"kessler/common/objects/files"
	"kessler/common/s3utils"
	"kessler/db"
	"kessler/gen/dbstore"
	"net/http"
	"os"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func FileWithMetaGetHandler(w http.ResponseWriter, r *http.Request) {
	q := db.GetTx()
	params := mux.Vars(r)
	fileID := params["uuid"]
	parsedUUID, err := uuid.Parse(fileID)
	if err != nil {
		errorstring := fmt.Sprintf("Error parsing file %v: %v\n", fileID, err)
		log.Info(errorstring)
		http.Error(w, errorstring, http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	file_raw, err := q.GetFileWithMetadata(ctx, parsedUUID)
	if err != nil {
		errorstring := fmt.Sprintf("Error Retriving file %v: %v\n", fileID, err)
		log.Info(errorstring)
		http.Error(w, errorstring, http.StatusNotFound)
		return
	}
	var mdata_obj map[string]interface{}
	err = json.Unmarshal(file_raw.Mdata, &mdata_obj)
	if err != nil {
		errorstring := fmt.Sprintf("Error Unmarshalling file %v: %v\n", fileID, err)
		log.Info(errorstring)
		http.Error(w, errorstring, http.StatusInternalServerError)
		return
	}
	file := files.CompleteFileSchema{
		ID:        file_raw.ID,
		Verified:  file_raw.Verified.Bool,
		Extension: file_raw.Extension,
		Lang:      file_raw.Lang,
		Name:      file_raw.Name,
		Hash:      file_raw.Hash,
		IsPrivate: file_raw.Isprivate.Bool,
		Mdata:     mdata_obj,
	}

	response, _ := json.Marshal(file)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func FileSemiCompleteGet(w http.ResponseWriter, r *http.Request) {
	q := *db.GetTx()

	params := mux.Vars(r)
	fileID := params["uuid"]
	parsedUUID, err := uuid.Parse(fileID)
	if err != nil {
		errorstring := fmt.Sprintf("Error parsing file %v: %v\n", fileID, err)
		log.Info(errorstring)
		http.Error(w, errorstring, http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	file, err := SemiCompleteFileGetFromUUID(ctx, q, parsedUUID)
	if err != nil {
		log.Info(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response, _ := json.Marshal(file)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// TODO: refactor config into a middleware pattern
func SemiCompleteFileGetFromUUID(ctx context.Context, q dbstore.Queries, uuid uuid.UUID) (files.CompleteFileSchema, error) {
	files_raw, err := q.SemiCompleteFileGet(ctx, uuid)
	if err != nil {
		errorstring := fmt.Sprintf("Error Retriving file %v: %v\n", uuid, err)
		return files.CompleteFileSchema{}, errors.New(errorstring)
	}
	if len(files_raw) == 0 {
		errorstring := fmt.Sprintf("Error No Files Found for a list of length zero.\n")
		return files.CompleteFileSchema{}, errors.New(errorstring)
	}
	file_raw := files_raw[0]

	var mdata_obj map[string]interface{}
	nilSchema := files.CompleteFileSchema{}
	err = json.Unmarshal(file_raw.Mdata, &mdata_obj)
	if err != nil {
		errorstring := fmt.Sprintf("Error Unmarshalling file metadata %v: %v\n", uuid, err)
		return nilSchema, errors.New(errorstring)
	}
	var extra_obj files.FileGeneratedExtras
	err = json.Unmarshal(file_raw.ExtraObj, &extra_obj)
	if err != nil {
		errorstring := fmt.Sprintf("Error Unmarshalling file extras %v: %v\n", uuid, err)
		return nilSchema, errors.New(errorstring)
	}
	// Missing info here, it doesnt have the name.
	conv_info := conversations.ConversationInformation{ID: file_raw.OrganizationID.Bytes}
	author_info := make([]authors.AuthorInformation, len(files_raw))
	for i, author_file_raw := range files_raw {
		author_info[i] = authors.AuthorInformation{
			AuthorName:      author_file_raw.OrganizationName.String,
			IsPerson:        author_file_raw.IsPerson.Bool,
			IsPrimaryAuthor: author_file_raw.IsPrimaryAuthor.Bool,
			AuthorID:        author_file_raw.OrganizationID.Bytes,
		}
	}

	file := files.CompleteFileSchema{
		ID:           file_raw.ID,
		Verified:     file_raw.Verified.Bool,
		Extension:    file_raw.Extension,
		Lang:         file_raw.Lang,
		Name:         file_raw.Name,
		Hash:         file_raw.Hash,
		Mdata:        mdata_obj,
		Extra:        extra_obj,
		Conversation: conv_info,
		Authors:      author_info,
	}
	return file, nil
}

func FileTextsGetAll(ctx context.Context, q dbstore.Queries, uuid uuid.UUID) ([]files.FileChildTextSource, error) {
	texts, err := q.ListTextsOfFile(ctx, uuid)
	if err != nil {
		return make([]files.FileChildTextSource, 0), err
	}
	return_texts := make([]files.FileChildTextSource, len(texts))
	for i, text := range texts {
		return_texts[i] = files.FileChildTextSource{
			IsOriginalText: text.IsOriginalText,
			Text:           text.Text,
			Language:       text.Language,
		}
	}
	return return_texts, nil
}

func FileStageGet(ctx context.Context, q dbstore.Queries, uuid uuid.UUID) (files.DocProcStage, error) {
	stage_str, err := q.StageLogFileGetLatest(ctx, uuid)
	if err != nil {
		return files.DocProcStage{}, err
	}
	stage := files.DocProcStage{}
	err = json.Unmarshal(stage_str.Log, &stage)
	if err != nil {
		return stage, err
	}
	return stage, nil
}

func CompleteFileSchemaGetFromUUID(ctx context.Context, q dbstore.Queries, uuid uuid.UUID) (files.CompleteFileSchema, error) {
	file, err := SemiCompleteFileGetFromUUID(ctx, q, uuid)
	nilSchema := files.CompleteFileSchema{}
	if err != nil {
		return nilSchema, err
	}
	texts, err := FileTextsGetAll(ctx, q, uuid)
	if err != nil {
		return nilSchema, err
	}
	file.DocTexts = texts
	stage, err := FileStageGet(ctx, q, uuid)
	if err != nil {
		return nilSchema, err
	}
	file.Stage = stage
	return file, nil
}

// TODO: refactor config into a middleware pattern
func ReadFileHandler(config FileHandlerConfig) http.HandlerFunc {
	private := config.private
	return_type := config.return_type

	return func(w http.ResponseWriter, r *http.Request) {
		q := *db.GetTx()

		// token := r.Header.Get("Authorization")
		params := mux.Vars(r)
		fileID := params["uuid"]
		parsedUUID, err := uuid.Parse(fileID)
		if err != nil {
			http.Error(w, "Invalid File ID format", http.StatusBadRequest)
			return
		}
		ctx := r.Context()
		// if private {
		//
		// 	isAuthorized, err := checkPrivateFileAuthorization(q, ctx, parsedUUID, token)
		// 	if !isAuthorized {
		// 		http.Error(w, "Forbidden", http.StatusForbidden)
		// 	}
		// 	if err != nil {
		// 		log.Info(fmt.Sprintf("Ran into the follwing error with authentication $v", err))
		// 	}
		// }
		// Since all three of these methods share the same authentication and database connection prerecs
		// switching functionality using an if else, or a cases switch lets code get reused
		// TODO: This is horrible, I need to refactor
		file_params := GetFileParam{
			Queries: q,
			Context: ctx,
			PgUUID:  parsedUUID,
			Private: private,
		}
		switch return_type {
		case "raw":
			file, err := GetFileObjectRaw(file_params)
			if err != nil {
				error_string := fmt.Sprintf("Error retrieving file object %v", err)
				log.Info(error_string)
				http.Error(w, error_string, http.StatusNotFound)
				return
			}
			filehash, err := hashes.HashFromString(file.Hash)
			if err != nil {
				error_string := fmt.Sprintf("ASSERTION ERROR: File hash could not be decoded: %v", err)
				log.Error(error_string)
				http.Error(w, error_string, http.StatusInternalServerError)
				return
			}
			kefiles := s3utils.NewKeFileManager()
			file_path, err := kefiles.DownloadFileFromS3(filehash)
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
			inflated_schema := file.CompleteFileSchemaInflateFromPartialSchema()

			response, _ := json.Marshal(inflated_schema)

			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
		default:
			log.Info(fmt.Sprintf("Encountered unreachable code with file return type %v", return_type))
			http.Error(w, "Congradulations for encountering unreachable code about support types!", http.StatusInternalServerError)
		}
	}
}
