// objects/files/handler/file_read_handler.go
package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"kessler/internal/dbstore"
	"kessler/internal/objects/authors"
	"kessler/internal/objects/conversations"
	"kessler/internal/objects/files"
	"kessler/pkg/hashes"
	"kessler/pkg/logger"
	"kessler/pkg/s3utils"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// FileWithMetaGetHandler handles GET requests for files with metadata
func (h *FileHandler) FileWithMetaGetHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "files:FileWithMetaGetHandler")
	defer span.End()
	log := logger.FromContext(ctx)

	params := mux.Vars(r)
	fileID := params["uuid"]
	parsedUUID, err := uuid.Parse(fileID)
	if err != nil {
		errorstring := fmt.Sprintf("Error parsing file %v: %v", fileID, err)
		log.Info(errorstring)
		http.Error(w, errorstring, http.StatusBadRequest)
		return
	}

	q := dbstore.New(h.db)
	file_raw, err := q.GetFileWithMetadata(ctx, parsedUUID)
	if err != nil {
		errorstring := fmt.Sprintf("Error retrieving file %v: %v", fileID, err)
		log.Info(errorstring)
		http.Error(w, errorstring, http.StatusNotFound)
		return
	}

	var mdata_obj map[string]interface{}
	err = json.Unmarshal(file_raw.Mdata, &mdata_obj)
	if err != nil {
		errorstring := fmt.Sprintf("Error unmarshalling file %v: %v", fileID, err)
		log.Info(errorstring)
		http.Error(w, errorstring, http.StatusInternalServerError)
		return
	}

	file := files.CompleteFileSchema{
		ID:        file_raw.ID,
		Verified:  file_raw.Verified.Bool,
		Lang:      file_raw.Lang,
		Name:      file_raw.Name,
		IsPrivate: file_raw.Isprivate.Bool,
		Mdata:     mdata_obj,
	}

	response, _ := json.Marshal(file)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// FileMarkdownByHashHandler retrieves markdown text for a file by its hash.
// It selects the first matching file UUID and returns its markdown text.
func (h *FileHandler) FileMarkdownByHashHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "files:FileMarkdownByHashHandler")
	defer span.End()
	log := logger.FromContext(ctx)

	params := mux.Vars(r)
	hash := params["hash"]

	q := dbstore.New(h.db)

	// Map hash to file UUIDs
	uuids, err := files.HashGetUUIDsFile(*q, ctx, hash)
	if err != nil || len(uuids) == 0 {
		log.Info("no file found for hash", zap.String("hash", hash))
		http.Error(w, fmt.Sprintf("No file found for hash %v", hash), http.StatusNotFound)
		return
	}

	// Use the first matching file
	fileParams := files.GetFileParam{Queries: *q, Context: ctx, PgUUID: uuids[0], Private: false}

	// Query parameters for language filtering
	original := r.URL.Query().Get("original_lang") == "true"
	matchLang := r.URL.Query().Get("match_lang")

	markdown, err := files.GetSpecificFileText(fileParams, matchLang, original)
	if err != nil {
		log.Info("error retrieving text", zap.Error(err))
		http.Error(w, "Error retrieving text or no matching text found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(markdown))
}

// FileSemiCompleteGet handles GET requests for semi-complete file information
func (h *FileHandler) FileSemiCompleteGet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "files:FileSemiCompleteGet")
	defer span.End()
	log := logger.FromContext(ctx)

	params := mux.Vars(r)
	fileID := params["uuid"]
	parsedUUID, err := uuid.Parse(fileID)
	if err != nil {
		errorstring := fmt.Sprintf("Error parsing file %v: %v", fileID, err)
		log.Info(errorstring)
		http.Error(w, errorstring, http.StatusBadRequest)
		return
	}

	q := dbstore.New(h.db)
	file, err := h.SemiCompleteFileGetFromUUID(ctx, q, parsedUUID)
	if err != nil {
		log.Info("encountered error getting file from uuid", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, _ := json.Marshal(file)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// SemiCompleteFileGetFromUUID retrieves semi-complete file information
func (h *FileHandler) SemiCompleteFileGetFromUUID(ctx context.Context, q *dbstore.Queries, uuid uuid.UUID) (files.CompleteFileSchema, error) {
	files_raw, err := q.SemiCompleteFileGet(ctx, uuid)
	if err != nil {
		return files.CompleteFileSchema{}, fmt.Errorf("error retrieving file %v: %w", uuid, err)
	}
	if len(files_raw) == 0 {
		return files.CompleteFileSchema{}, errors.New("no files found")
	}
	file_raw := files_raw[0]

	var mdata_obj map[string]interface{}
	nilSchema := files.CompleteFileSchema{}
	err = json.Unmarshal(file_raw.Mdata, &mdata_obj)
	if err != nil {
		return nilSchema, fmt.Errorf("error unmarshalling file metadata %v: %w", uuid, err)
	}

	var extra_obj files.FileGeneratedExtras
	err = json.Unmarshal(file_raw.ExtraObj, &extra_obj)
	if err != nil {
		return nilSchema, fmt.Errorf("error unmarshalling file extras %v: %w", uuid, err)
	}

	// Missing info here, it doesn't have the name.
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
		Lang:         file_raw.Lang,
		Name:         file_raw.Name,
		Mdata:        mdata_obj,
		Extra:        extra_obj,
		Conversation: conv_info,
		Authors:      author_info,
	}
	return file, nil
}

// FileStageGet retrieves the processing stage for a file
func (h *FileHandler) FileStageGet(ctx context.Context, q *dbstore.Queries, uuid uuid.UUID) (files.DocProcStage, error) {
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

// ReadFileHandler creates a handler for reading files in various formats
func (h *FileHandler) ReadFileHandler(config FileHandlerConfig) http.HandlerFunc {
	private := config.private
	return_type := config.return_type

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logger.FromContext(ctx)

		q := dbstore.New(h.db)

		// token := r.Header.Get("Authorization")
		params := mux.Vars(r)
		fileID := params["uuid"]
		parsedUUID, err := uuid.Parse(fileID)
		if err != nil {
			http.Error(w, "Invalid File ID format", http.StatusBadRequest)
			return
		}

		// if private {
		// 	isAuthorized, err := checkPrivateFileAuthorization(q, ctx, parsedUUID, token)
		// 	if !isAuthorized {
		// 		http.Error(w, "Forbidden", http.StatusForbidden)
		// 	}
		// 	if err != nil {
		// 		log.Info(fmt.Sprintf("Ran into the following error with authentication %v", err))
		// 	}
		// }

		file_params := files.GetFileParam{
			Queries: *q,
			Context: ctx,
			PgUUID:  parsedUUID,
			Private: private,
		}

		switch return_type {
		case "raw":
			file, err := files.GetFileObjectRaw(file_params)
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
			markdownText, err := files.GetSpecificFileText(file_params, matchLang, originalLang)
			if err != nil {
				http.Error(w, "Error retrieving texts or no texts found that match query params", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte(markdownText))

		case "object-minimal":
			file, err := files.GetFileObjectRaw(file_params)
			if err != nil {
				http.Error(w, "File not found", http.StatusNotFound)
				return
			}
			inflated_schema := file.CompleteFileSchemaInflateFromPartialSchema()

			response, _ := json.Marshal(inflated_schema)

			w.Header().Set("Content-Type", "application/json")
			w.Write(response)

		default:
			log.Info("encountered unreachable code", zap.String("return_type", return_type))
			http.Error(w, "Congratulations for encountering unreachable code about support types!", http.StatusInternalServerError)
		}
	}
}

// Helper functions that don't depend on the handler

func AttachmentSchemaCompleteFill(ctx context.Context, q dbstore.Queries, attachment *files.CompleteAttachmentSchema) (files.CompleteAttachmentSchema, error) {
	texts, err := q.AttachmentTextList(ctx, attachment.ID)
	if err != nil {
		return *attachment, err
	}
	return_texts := make([]files.AttachmentChildTextSource, len(texts))
	for i, text := range texts {
		return_texts[i] = files.AttachmentChildTextSource{
			IsOriginalText: text.IsOriginalText,
			Text:           text.Text,
			Language:       text.Language,
		}
	}
	attachment.Texts = return_texts

	return *attachment, nil
}

func AttachmentFromDBStore(attach dbstore.Attachment) (files.CompleteAttachmentSchema, error) {
	valid_hash, err := hashes.HashFromString(attach.Hash)
	if err != nil {
		return files.CompleteAttachmentSchema{}, err
	}

	return files.CompleteAttachmentSchema{
		ID:        attach.ID,
		Name:      attach.Name,
		Extension: attach.Extension,
		Hash:      valid_hash,
	}, nil
}

func AttachmentsCompleteGet(ctx context.Context, q dbstore.Queries, fileUUID uuid.UUID) ([]files.CompleteAttachmentSchema, error) {
	raw_attachments, err := q.AttachmentListByFileId(ctx, fileUUID)
	if err != nil {
		return []files.CompleteAttachmentSchema{}, err
	}
	complete_attachments := make([]files.CompleteAttachmentSchema, len(raw_attachments))
	for i, raw_attachment := range raw_attachments {
		kinda_raw_attachment, err := AttachmentFromDBStore(raw_attachment)
		if err != nil {
			return []files.CompleteAttachmentSchema{}, err
		}
		attachment, err := AttachmentSchemaCompleteFill(ctx, q, &kinda_raw_attachment)
		if err != nil {
			return []files.CompleteAttachmentSchema{}, err
		}
		complete_attachments[i] = attachment
	}
	return complete_attachments, nil
}

func (h *FileHandler) CompleteFileSchemaGetFromUUID(ctx context.Context, q *dbstore.Queries, uuid uuid.UUID) (files.CompleteFileSchema, error) {
	file, err := h.SemiCompleteFileGetFromUUID(ctx, q, uuid)
	nilSchema := files.CompleteFileSchema{}
	if err != nil {
		return nilSchema, err
	}
	stage, err := h.FileStageGet(ctx, q, uuid)
	if err != nil {
		return nilSchema, err
	}
	file.Stage = stage
	return file, nil
}
