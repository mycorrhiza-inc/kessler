package rawattachments

import (
	"fmt"
	"kessler/pkg/hashes"
	"kessler/pkg/logger"
	"kessler/pkg/s3utils"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func DefineRawAttachmentRoutes(r *mux.Router) {
	// Create handler instance with database

	// Insert endpoint

	r.HandleFunc(
		"/{hash}/raw",
		RawAttachmentRawBytesGet,
	).Methods(http.MethodGet)
}

// ReadFileHandler creates a handler for reading files in various formats
func RawAttachmentRawBytesGet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	// token := r.Header.Get("Authorization")
	params := mux.Vars(r)
	rawHash := params["hash"]
	hash, err := hashes.HashFromString(rawHash)
	if err != nil {
		log.Error("Could not validate hash", zap.Error(err), zap.String("raw_hash", rawHash))
		http.Error(w, "Invalid Hash format", http.StatusBadRequest)
		return
	}
	kefiles := s3utils.NewKeFileManager()
	file_path, err := kefiles.DownloadFileFromS3(hash)
	if err != nil {
		log.Error("Could not get file from s3", zap.Error(err), zap.String("hash", hash.String()))
		http.Error(w, fmt.Sprintf("Error encountered when getting file with hash %v from s3:%v", hash, err), http.StatusInternalServerError)
		return
	}
	content, err := os.ReadFile(file_path)
	if err != nil {
		log.Error("Encountered error reading downloaded file from s3", zap.Error(err), zap.String("hash", hash.String()))
		http.Error(w, fmt.Sprintf("Error reading file: %v", err), http.StatusInternalServerError)
		return
	}

	mimeType := http.DetectContentType(content)
	// if mimeType == "application/octet-stream" {
	// 	mimeType = "application/pdf" // Default to PDF if mime type can't be determined
	// }

	w.Header().Set("Content-Type", mimeType)
	w.Write(content)
}
