package search

import (
	"context"
	"fmt"
	"kessler/gen/dbstore"
	"kessler/objects/files"
	"log"

	"github.com/google/uuid"
)

// func doShittyHydrationBROKEN(uuidList []uuid.UUID, ctx context.Context, q dbstore.Queries) ([]files.CompleteFileSchema, error) {
// 	fmt.Printf("Got Hydration Request for %v files", uuidList)
// 	intermediate_results, err := q.SemiCompleteFileListGet(ctx, uuidList)
// 	if err != nil {
// 		return nil, err
// 	}
// 	fmt.Printf("Got %v semi complete files\n", len(intermediate_results))
// 	results := make([]files.CompleteFileSchema, len(intermediate_results))
// 	for i, schema := range intermediate_results {
// 		file, err := unrollSingleSchema(schema)
// 		if err != nil {
// 			log.Printf("Error hydrating file: %v", err)
// 			continue
// 		}
// 		results[i] = file
// 	}
// 	if len(results) == 0 {
// 		return nil, fmt.Errorf("no valid results found after hydration")
// 	}
// 	return results, nil
// }
//
// func unrollSingleSchema(schema dbstore.SemiCompleteFileListGetRow) (files.CompleteFileSchema, error) {
// 	var mdata_obj map[string]interface{}
// 	nilSchema := files.CompleteFileSchema{}
// 	err := json.Unmarshal(schema.Mdata, &mdata_obj)
// 	if err != nil {
// 		errorstring := fmt.Sprintf("Error Unmarshalling file metadata : %v\n", schema.ID, err)
// 		return nilSchema, errors.New(errorstring)
// 	}
// 	var extra_obj files.FileGeneratedExtras
// 	err = json.Unmarshal(schema.ExtraObj, &extra_obj)
// 	if err != nil {
// 		errorstring := fmt.Sprintf("Error Unmarshalling file extras %v: %v\n", schema.ID, err)
// 		return nilSchema, errors.New(errorstring)
// 	}
// 	// Missing info here, it doesnt have the name.
// 	conv_info := conversations.ConversationInformation{ID: schema.ConversationUuid.Bytes}
// 	org_id_list := schema.OrganizationIds
// 	org_name_list := schema.OrganizationNames
// 	org_personhood_list := schema.IsPersonList
//
// 	author_info := make([]authors.AuthorInformation, len(org_id_list))
// 	for i := range org_id_list {
// 		author_info[i] = authors.AuthorInformation{
// 			AuthorName:      org_name_list[i],
// 			IsPerson:        org_personhood_list[i],
// 			IsPrimaryAuthor: false,
// 			AuthorID:        org_id_list[i],
// 		}
// 	}
//
// 	file := files.CompleteFileSchema{
// 		ID:           schema.ID,
// 		Verified:     schema.Verified.Bool,
// 		Extension:    schema.Extension,
// 		Lang:         schema.Lang,
// 		Name:         schema.Name,
// 		Hash:         schema.Hash,
// 		Mdata:        mdata_obj,
// 		Extra:        extra_obj,
// 		Conversation: conv_info,
// 		Authors:      author_info,
// 	}
// 	return file, nil
// }

// func doShittyHydration(uuidList []uuid.UUID, ctx context.Context, q dbstore.Queries) ([]files.CompleteFileSchema, error) {
// 	fmt.Printf("Got Hydration Request for %v files", uuidList)
//
// 	// Create buffered channels for results
// 	resultChan := make(chan files.CompleteFileSchema, len(uuidList))
//
// 	// Use a semaphore to limit concurrent operations
// 	semaphore := make(chan struct{}, 20) // Adjust this number based on your needs
//
// 	// Launch goroutines for each file
// 	var wg sync.WaitGroup
// 	for _, file_id := range uuidList {
// 		wg.Add(1)
// 		go func(fileID uuid.UUID) {
// 			defer wg.Done()
//
// 			// Acquire semaphore
// 			semaphore <- struct{}{}
// 			defer func() { <-semaphore }()
//
// 			// Create a new query instance for this goroutine
// 			localQ := *routing.DBQueriesFromContext(ctx)
//
// 			semi_complete_file, err := crud.SemiCompleteFileGetFromUUID(ctx, localQ, fileID)
// 			if err != nil {
// 				fmt.Printf("Encountered error getting emails from file with uuid: %v: %v\n", fileID, err)
// 				return
// 			}
//
// 			select {
// 			case resultChan <- semi_complete_file:
// 			case <-ctx.Done():
// 				return
// 			}
// 		}(file_id)
// 	}
//
// 	// Close result channel when all goroutines complete
// 	go func() {
// 		wg.Wait()
// 		close(resultChan)
// 	}()
//
// 	// Collect results
// 	var results []files.CompleteFileSchema
// 	for result := range resultChan {
// 		results = append(results, result)
// 	}
//
// 	return results, nil
// }

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
			fmt.Printf("Error parsing uuid: %v\n", err)
		}
		uuid_list[i] = uuid
	}
	fmt.Printf("Hydrating %v files\n", len(uuid_list))
	files_complete, err := doShittyHydration(uuid_list, ctx, q)
	if err != nil {
		log.Printf("Error hydrating search results: %v\n", err)
		return []SearchDataHydrated{}, err
	}
	fmt.Printf("Got back %v complete files for hydration, out of %v requested files\n", len(files_complete), len(results))

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
			fmt.Printf("Error hydrating search result: %v\n", uuid)
			res.Name = "ERROR HYDRATING SEARCH RESULT W/ PG DATA: " + res.Name
			results_hydrated[i] = setupSearchDataHydrated(res, files.CompleteFileSchema{})
		}
	}
	fmt.Printf("Hydrated %v of %v files\n", files_actually_hydrated, len(results))
	return results_hydrated, nil
}

func SearchThenHydrate(r SearchRequest, ctx context.Context) ([]files.FileMetadataSchema, error) {
	data, err := SearchQuickwit(r)
	if err != nil {
		log.Printf("search-then-hydrate: error getting quickwit data", err)
	}
	idList := make([]string, len(data))
	for i, d := range data {
		idList[i] = d.SourceID
	}

	// hydratedData := make([]files.FileMetadataSchema, len(idList))

	return []files.FileMetadataSchema{}, nil
}
