package search

import (
	"context"
	"kessler/internal/dbstore"
	"kessler/internal/objects/files"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func doShittyHydration(uuidList []uuid.UUID, ctx context.Context, q dbstore.Queries) ([]files.CompleteFileSchema, error) {
	return []files.CompleteFileSchema{}, nil
}

type SearchDataHydrated struct {
	Name     string                   `json:"name"`
	Snippet  string                   `json:"snippet,omitempty"`
	DocID    string                   `json:"docID"`
	SourceID string                   `json:"sourceID"`
	File     files.CompleteFileSchema `json:"file"`
}

func setupSearchDataHydrated(s SearchData, f files.CompleteFileSchema) SearchDataHydrated {
	return SearchDataHydrated{
		Name:     s.Name,
		Snippet:  s.Snippet,
		DocID:    s.DocID,
		SourceID: s.SourceID,
		File:     f,
	}
}

func HydrateSearchResults(results []SearchData, ctx context.Context, q dbstore.Queries) ([]SearchDataHydrated, error) {
	uuid_list := make([]uuid.UUID, len(results))
	for i, r := range results {
		uuid, err := uuid.Parse(r.SourceID)
		if err != nil {
			log.Info("Error parsing uuid", zap.Error(err))
		}
		uuid_list[i] = uuid
	}
	log.Info("Hydrating files", zap.Int("count", len(uuid_list)))
	files_complete, err := doShittyHydration(uuid_list, ctx, q)
	if err != nil {
		log.Info("Error hydrating search results", zap.Error(err))
		return []SearchDataHydrated{}, err
	}
	log.Info("Got complete files for hydration",
		zap.Int("complete_files", len(files_complete)),
		zap.Int("requested_files", len(results)))

	// Create map of file ID to complete file
	fileMap := make(map[uuid.UUID]files.CompleteFileSchema)
	for _, f := range files_complete {
		fileMap[f.ID] = f
	}

	results_hydrated := make([]SearchDataHydrated, len(results))
	files_actually_hydrated := 0

	for i, res := range results {
		uuid, _ := uuid.Parse(res.SourceID)
		if file, ok := fileMap[uuid]; ok {
			files_actually_hydrated++
			results_hydrated[i] = setupSearchDataHydrated(res, file)
		} else {
			log.Info("Error hydrating search result", zap.String("uuid", uuid.String()))
			res.Name = "ERROR HYDRATING SEARCH RESULT W/ PG DATA: " + res.Name
			results_hydrated[i] = setupSearchDataHydrated(res, files.CompleteFileSchema{})
		}
	}
	log.Info("Hydration complete",
		zap.Int("hydrated_files", files_actually_hydrated),
		zap.Int("total_files", len(results)))
	return results_hydrated, nil
}

func SearchThenHydrate(r SearchRequest, ctx context.Context) ([]files.FileMetadataSchema, error) {
	data, err := SearchQuickwit(r)
	if err != nil {
		log.Info("search-then-hydrate: error getting quickwit data", zap.Error(err))
	}
	idList := make([]string, len(data))
	for i, d := range data {
		idList[i] = d.SourceID
	}

	// hydratedData := make([]files.FileMetadataSchema, len(idList))

	return []files.FileMetadataSchema{}, nil
}
