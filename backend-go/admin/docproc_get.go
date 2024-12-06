package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"kessler/crud"
	"kessler/gen/dbstore"
	"kessler/objects/files"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func UnverifedCompleteFileSchemaListFactory(dbtx_val dbstore.DBTX) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := *dbstore.New(dbtx_val)
		params := mux.Vars(r)
		ctx := r.Context()
		// ctx := context.Background()
		max_responses_str := params["max_responses"]
		max_responses, err := strconv.Atoi(max_responses_str)
		if err != nil || max_responses < 0 {
			errorstring := fmt.Sprintf("Error parsing max responses: %v", err)
			fmt.Println(errorstring)
			http.Error(w, errorstring, http.StatusBadRequest)
			return
		}
		files, err := UnverifedCompleteFileSchemaList(ctx, q, uint(max_responses))
		if err != nil {
			errorstring := fmt.Sprintf("Error getting unverified files: %v", err)
			fmt.Println(errorstring)
			http.Error(w, errorstring, http.StatusInternalServerError)
			return
		}
		response, _ := json.Marshal(files)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}

func UnverifedCompleteFileSchemaList(ctx context.Context, q dbstore.Queries, max_responses uint) ([]files.CompleteFileSchema, error) {
	fmt.Printf("Getting %d unverified files\n", max_responses)
	db_files, err := q.FilesListUnverified(ctx, int32(max_responses)*2)
	// If postgres return randomization doesnt work, then you can still get it to kinda work by returning double the results, randomizing and throwing away half.
	if err != nil {
		return []files.CompleteFileSchema{}, err
	}
	fmt.Printf("Got unverified files, randomizing uuids\n")
	unverified_raw_uuids := make([]uuid.UUID, len(db_files))
	for i, file := range db_files {
		unverified_raw_uuids[i] = file.ID
	}
	// Shuffle the uuids around to get a random selection while processing
	for i := range unverified_raw_uuids {
		j := rand.Intn(i + 1) // Want to get a range of [0,i] so that there is a possibility of the null swap.
		// Inductive proof this distributes the elements randomly at step k:
		// The element at index k is evenly distributed, since it has a 1/k chance of being the element at k, and a k-1/k chance of selecting from an even distribution of k-1 elements, thus meaning it has an even distribution of 1/k probability of selecting k elements.
		// Same thing for other elements, it has a k-1/k chance of sampling from its EXISTING even distribution of k-1 elements, and a 1/k chance of swapping with the k'th element. Thus it has a even 1/k chance of being any of k elements.
		unverified_raw_uuids[i], unverified_raw_uuids[j] = unverified_raw_uuids[j], unverified_raw_uuids[i]
	}
	if len(unverified_raw_uuids) > int(max_responses) {
		unverified_raw_uuids = unverified_raw_uuids[:max_responses]
	}
	complete_files, err := CompleteFileSchemasFromUUIDs(ctx, q, unverified_raw_uuids)
	if err != nil {
		return []files.CompleteFileSchema{}, err
	}
	return complete_files, nil
}

func CompleteFileSchemasFromUUIDs(ctx context.Context, q dbstore.Queries, uuids []uuid.UUID) ([]files.CompleteFileSchema, error) {
	complete_files := []files.CompleteFileSchema{}
	fileChan := make(chan files.CompleteFileSchema)
	// errChan := make(chan error)
	var wg sync.WaitGroup
	for _, file_uuid := range uuids {
		wg.Add(1)
		go func(file_uuid uuid.UUID) {
			defer wg.Done()
			start := time.Now()
			complete_file, err := crud.CompleteFileSchemaGetFromUUID(ctx, q, file_uuid)
			elapsed := time.Since(start)
			if elapsed > time.Second {
				fmt.Printf("Slow query for file %v, took %v seconds.\n", file_uuid, elapsed/time.Second)
			}
			if err != nil {
				fmt.Printf("Error getting file %v: %v\n", file_uuid, err)
				// errChan <- err
				return
			}
			fmt.Printf("Got complete file %v\n", file_uuid)
			fileChan <- complete_file
		}(file_uuid)
	}

	// Close channels when all goroutines complete
	go func() {
		wg.Wait()
		close(fileChan)
		// close(errChan)
	}()

	// Collect results
	for file := range fileChan {
		complete_files = append(complete_files, file)
	}
	return complete_files, nil
}
