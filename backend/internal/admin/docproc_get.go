package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"kessler/internal/dbstore"
	"kessler/internal/objects/files"
	FileHandler "kessler/internal/objects/files/handler"
	"kessler/pkg/util"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func HandleUnverifedCompleteFileSchemaList(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ctx := r.Context()
	// ctx := context.Background()
	max_responses_str := params["max_responses"]
	max_responses, err := strconv.Atoi(max_responses_str)
	if err != nil || max_responses < 0 {
		errorstring := fmt.Sprintf("Error parsing max responses: %v", err)
		log.Info(errorstring)
		http.Error(w, errorstring, http.StatusBadRequest)
		return
	}
	files, err := UnverifedCompleteFileSchemaRandomList(ctx, uint(max_responses))
	if err != nil {
		errorstring := fmt.Sprintf("Error getting unverified files: %v", err)
		log.Info(errorstring)
		http.Error(w, errorstring, http.StatusInternalServerError)
		return
	}
	response, _ := json.Marshal(files)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func UnverifedCompleteFileSchemaRandomList(ctx context.Context, max_responses uint) ([]files.CompleteFileSchema, error) {
	q := util.DBQueriesFromContext(ctx)
	log.Info(fmt.Sprintf("Getting %d unverified files\n", max_responses))
	db_files, err := q.FilesListUnverified(ctx, int32(max_responses)*2)
	// If postgres return randomization doesnt work, then you can still get it to kinda work by returning double the results, randomizing and throwing away half.
	if err != nil {
		return []files.CompleteFileSchema{}, err
	}
	log.Info(fmt.Sprintf("Got unverified files, randomizing uuids\n"))
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
	complete_files, err := CompleteFileSchemasFromUUIDs(ctx, unverified_raw_uuids)
	if err != nil {
		return []files.CompleteFileSchema{}, err
	}
	return complete_files, nil
}

func CompleteFileSchemasFromUUIDs(ctx context.Context, uuids []uuid.UUID) ([]files.CompleteFileSchema, error) {
	dbtx_val := util.DBTXFromContext(ctx)

	complete_start := time.Now()
	complete_files := []files.CompleteFileSchema{}
	fileChan := make(chan files.CompleteFileSchema)

	// errChan := make(chan error)
	var wg sync.WaitGroup
	for _, file_uuid := range uuids {
		wg.Add(1)
		go func(file_uuid uuid.UUID, dbtx_val dbstore.DBTX) {
			defer wg.Done()
			q := *dbstore.New(dbtx_val)
			// start := time.Now()
			complete_file, err := FileHandler.CompleteFileSchemaGetFromUUID(ctx, q, file_uuid)
			// elapsed := time.Since(start)
			// TODO: Debug why these loading times are so fucking slow.
			// if elapsed > 10*time.Second {
			// 	log.Info(fmt.Sprintf("Slow query for file %v, took %v seconds.\n", file_uuid, elapsed/time.Second))
			// }
			if err != nil {
				log.Info(fmt.Sprintf("Error getting file %v: %v\n", file_uuid, err))
				// errChan <- err
				return
			}
			// log.Info(fmt.Sprintf("Got complete file %v\n", file_uuid))
			fileChan <- complete_file
		}(file_uuid, dbtx_val)
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
	elapsed := time.Since(complete_start)
	log.Info(fmt.Sprintf("Got %d complete files in %v\n", len(complete_files), elapsed))
	return complete_files, nil
}
