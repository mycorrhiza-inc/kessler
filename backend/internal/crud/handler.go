package crud

import (
	"context"
	"kessler/common/objects/files"
	"kessler/internal/dbstore"
	"net/http"

	"github.com/gorilla/mux"
)

func DefineCrudRoutes(r *mux.Router) {
	filesRoute := r.PathPrefix("/files").Subrouter()
	filesRoute.HandleFunc(
		"/insert",
		makeFileUpsertHandler(
			FileUpsertHandlerConfig{
				private: false,
				insert:  true,
			},
		)).Methods(http.MethodPost)

	filesRoute.HandleFunc(
		"/{uuid}/update",
		makeFileUpsertHandler(
			FileUpsertHandlerConfig{
				private: false,
				insert:  false,
			},
		)).Methods(http.MethodPost)
	filesRoute.HandleFunc(
		"/{uuid}",
		FileSemiCompleteGet,
	).Methods(http.MethodGet)

	filesRoute.HandleFunc(
		"/{uuid}/minimal",
		ReadFileHandler(
			FileHandlerConfig{
				private:     false,
				return_type: "object-minimal",
			},
		)).Methods(http.MethodGet)

	filesRoute.HandleFunc(
		"/{uuid}/markdown",
		ReadFileHandler(
			FileHandlerConfig{
				private:     false,
				return_type: "markdown",
			},
		)).Methods(http.MethodGet)

	filesRoute.HandleFunc(
		"/{uuid}/raw",
		ReadFileHandler(
			FileHandlerConfig{
				private:     false,
				return_type: "raw",
			},
		)).Methods(http.MethodGet)

	// DO NOT TOUCH. this is necessary for well named downloaded files
	filesRoute.HandleFunc(
		"/{uuid}/raw/{filename}",
		ReadFileHandler(
			FileHandlerConfig{
				private:     false,
				return_type: "raw",
			},
		)).Methods(http.MethodGet)

	filesRoute.HandleFunc(
		"/{uuid}/metadata",
		FileWithMetaGetHandler,
	).Methods(http.MethodGet)

	// ------- organizations subroute -------
	organizationsRouter := r.PathPrefix("/organizations").Subrouter()

	organizationsRouter.HandleFunc(
		"/organizations/list",
		OrgSemiCompletePaginated,
	).Methods(http.MethodGet)

	organizationsRouter.HandleFunc(
		"/organizations/{uuid}",
		OrgGetWithFilesHandler,
	).Methods(http.MethodGet)

	organizationsRouter.HandleFunc(
		"/organizations/verify",
		OrganizationVerifyHandler,
	).Methods(http.MethodPost)

	// ------- conversations subroute -------
	conversationsRouter := r.PathPrefix("/conversations").Subrouter()

	conversationsRouter.HandleFunc(
		"/conversations/list",
		ConversationSemiCompletePaginatedList,
	).Methods(http.MethodGet)

	conversationsRouter.HandleFunc(
		"/conversations/named-lookup/{name}",
		ConversationGetByUnknownHandler,
	).Methods(http.MethodGet)

	conversationsRouter.HandleFunc(
		"/conversations/verify",
		ConversationVerifyHandler,
	).Methods(http.MethodPost)

	conversationsRouter.HandleFunc(
		"/conversations/list/semi-complete",
		ConversationSemiCompleteListAll,
	).Methods(http.MethodGet)

}

// CONVERT TO UPPER CASE IF YOU EVER WANT TO USE IT OUTSIDE OF THIS CONTEXT
type FileHandlerConfig struct {
	private     bool
	return_type string // Can be either markdown, object or raw
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
		rawSchema := PublicFileToSchema(fileRaw)
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
