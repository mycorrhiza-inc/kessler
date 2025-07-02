// objects/files/handler/handler.go
package handler

import (
	"context"
	"kessler/internal/dbstore"
	"kessler/internal/objects/files"
	"net/http"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("file-handler")

// FileHandler holds dependencies for file operations
type FileHandler struct {
	db dbstore.DBTX
}

// NewFileHandler creates a new file handler with the given database connection
func NewFileHandler(db dbstore.DBTX) *FileHandler {
	return &FileHandler{
		db: db,
	}
}

// DefineFileRoutes registers all file-related routes
func DefineFileRoutes(r *mux.Router, db dbstore.DBTX) {
	// Create handler instance with database
	handler := NewFileHandler(db)

	// Insert endpoint

	r.HandleFunc(
		"/{uuid}/card",
		handler.FileCardGet,
	).Methods(http.MethodGet)

	r.HandleFunc(
		"/{uuid}/pageinfo",
		handler.FilePageInfoGet,
	).Methods(http.MethodGet)
	// Minimal file endpoint

	// Markdown file endpoint
	r.HandleFunc(
		"/{uuid}/markdown",
		handler.ReadFileHandler(
			FileHandlerConfig{
				private:     false,
				return_type: "markdown",
			},
		)).Methods(http.MethodGet)

	// Raw file endpoint
	r.HandleFunc(
		"/{uuid}/raw",
		handler.ReadFileHandler(
			FileHandlerConfig{
				private:     false,
				return_type: "raw",
			},
		)).Methods(http.MethodGet)

	// DO NOT TOUCH. this is necessary for well named downloaded files
	r.HandleFunc(
		"/{uuid}/raw/{filename}",
		handler.ReadFileHandler(
			FileHandlerConfig{
				private:     false,
				return_type: "raw",
			},
		)).Methods(http.MethodGet)

	// Metadata endpoint
	r.HandleFunc(
		"/{uuid}/metadata",
		handler.FileWithMetaGetHandler,
	).Methods(http.MethodGet)
}

// CONVERT TO UPPER CASE IF YOU EVER WANT TO USE IT OUTSIDE OF THIS CONTEXT
type FileHandlerConfig struct {
	private     bool
	return_type string // Can be either markdown, object or raw
}

type FileUpsertHandlerConfig struct {
	Private bool
	Insert  bool
}

type ReturnFilesSchema struct {
	Files []files.FileSchema `json:"files"`
}

// GetListAllRawFiles retrieves all raw files from the database
func (h *FileHandler) GetListAllRawFiles(ctx context.Context) ([]files.FileSchema, error) {
	q := dbstore.New(h.db)
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

// GetListAllFiles retrieves all files from the database
func (h *FileHandler) GetListAllFiles(ctx context.Context) ([]files.FileSchema, error) {
	db_files, err := h.GetListAllRawFiles(ctx)
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

// Note: The actual implementations of these methods are in:
// - file_read_handler.go (for read operations)
// - file_write_handler.go (for write operations)
