package search

import (
	"context"
	"fmt"
	"kessler/gen/dbstore"
	"kessler/objects/files"

	"github.com/charmbracelet/log"

	"github.com/google/uuid"
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
			log.Info(fmt.Sprintf("Error parsing uuid: %v\n", err))
		}
		uuid_list[i] = uuid
	}
	log.Info(fmt.Sprintf("Hydrating %v files\n", len(uuid_list)))
	files_complete, err := doShittyHydration(uuid_list, ctx, q)
	if err != nil {
		log.Info(fmt.Sprintf("Error hydrating search results: %v\n", err))
		return []SearchDataHydrated{}, err
	}
	log.Info(fmt.Sprintf("Got back %v complete files for hydration, out of %v requested files\n", len(files_complete), len(results)))

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
			log.Info(fmt.Sprintf("Error hydrating search result: %v\n", uuid))
			res.Name = "ERROR HYDRATING SEARCH RESULT W/ PG DATA: " + res.Name
			results_hydrated[i] = setupSearchDataHydrated(res, files.CompleteFileSchema{})
		}
	}
	log.Info(fmt.Sprintf("Hydrated %v of %v files\n", files_actually_hydrated, len(results)))
	return results_hydrated, nil
}

func SearchThenHydrate(r SearchRequest, ctx context.Context) ([]files.FileMetadataSchema, error) {
	data, err := SearchQuickwit(r)
	if err != nil {
		log.Info(fmt.Sprintf("search-then-hydrate: error getting quickwit data", err))
	}
	idList := make([]string, len(data))
	for i, d := range data {
		idList[i] = d.SourceID
	}

	// hydratedData := make([]files.FileMetadataSchema, len(idList))

	return []files.FileMetadataSchema{}, nil
}
