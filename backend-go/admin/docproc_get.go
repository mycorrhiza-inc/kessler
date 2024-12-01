package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/mycorrhiza-inc/kessler/backend-go/crud"
	"github.com/mycorrhiza-inc/kessler/backend-go/gen/dbstore"
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

func UnverifedCompleteFileSchemaList(ctx context.Context, q dbstore.Queries, max_responses uint) ([]crud.CompleteFileSchema, error) {
	fmt.Printf("Getting %d unverified files\n", max_responses)
	files, err := q.FilesListUnverified(ctx, int32(max_responses)*2)
	// If postgres return randomization doesnt work, then you can still get it to kinda work by returning double the results, randomizing and throwing away half.
	if err != nil {
		return []crud.CompleteFileSchema{}, err
	}
	fmt.Printf("Got unverified files, randomizing uuids\n")
	unverified_raw_uuids := make([]uuid.UUID, len(files))
	for i, file := range files {
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
	fmt.Printf("Trimmed unverified uuids to %d, getting complete files for all remaining\n", len(unverified_raw_uuids))

	complete_files := []crud.CompleteFileSchema{}
	fileChan := make(chan crud.CompleteFileSchema)
	errChan := make(chan error)
	var wg sync.WaitGroup

	for _, uuid := range unverified_raw_uuids {
		wg.Add(1)
		go func() {
			defer wg.Done()
			complete_file, err := crud.CompleteFileSchemaGetFromUUID(ctx, q, uuid)
			fmt.Printf("Got complete file %v\n", uuid)
			if err != nil {
				fmt.Printf("Error getting file %v: %v\n", uuid, err)
				// errChan <- err
				return
			}
			fileChan <- complete_file
		}()
	}

	// Close channels when all goroutines complete
	go func() {
		wg.Wait()
		close(fileChan)
		close(errChan)
	}()

	// Collect results
	for file := range fileChan {
		complete_files = append(complete_files, file)
	}
	return complete_files, nil
}
