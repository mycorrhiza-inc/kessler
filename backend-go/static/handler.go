package static

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/mux"
	"github.com/mycorrhiza-inc/kessler/backend-go/crud"
	"github.com/mycorrhiza-inc/kessler/backend-go/gen/dbstore"
)

func HandleStaticGenerationRouting(router *mux.Router, dbtx_val dbstore.DBTX) {
	admin_subrouter := router.PathPrefix("/api/v2/admin").Subrouter()
	admin_subrouter.HandleFunc("/generate-static-site", renderStaticSitemapmMakeHandler(dbtx_val))
}

func renderStaticSitemapmMakeHandler(dbtx_val dbstore.DBTX) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := RenderStaticSitemap(dbtx_val, 10)
		if err != nil {
			error_string := fmt.Sprintf("Encountered error while building static site map %v", err)
			http.Error(w, error_string, http.StatusInternalServerError)
			return
		}
		w.Write([]byte("Sucessfully built static site map"))
	}
}

func RenderStaticKesslerObj(obj crud.KesslerObject, dbtx_val dbstore.DBTX, ctx context.Context) error {
	q := dbstore.New(dbtx_val)
	file_path := crud.GetBaseFilePath(obj)

	os.Remove(file_path)
	// This can throw an error, but do nothing, file just didnt exist in all likelyhood.

	html_file, err := os.Create(file_path)
	if err != nil {
		return err
	}
	defer html_file.Close()
	err = obj.WriteHTMLString(html_file, *q, ctx)
	if err != nil {
		return err
	}
	return nil
}

func RenderStaticSitemap(dbtx_val dbstore.DBTX, max_docs int) error {
	ctx := context.Background()
	chanFileList := make(chan []crud.RawFileSchema)
	go func() {
		q := *dbstore.New(dbtx_val)
		list_all_files, err := crud.GetListAllRawFiles(ctx, q)
		if err != nil {
			fmt.Printf("Error encountered while getting all files %s", err)
		}
		var filteredFiles []crud.RawFileSchema
		for _, file := range list_all_files {
			if file.Stage != "completed" {
				filteredFiles = append(filteredFiles, file)
			}
		}
		chanFileList <- filteredFiles
	}()
	filteredFiles := <-chanFileList
	fmt.Printf("Generating %v static document pages\n", len(filteredFiles))

	var wg sync.WaitGroup
	fileChan := make(chan int)

	// Worker function
	worker := func() {
		defer wg.Done()
		for index := range fileChan {
			fileSchema := filteredFiles[index]
			err := RenderStaticKesslerObj(fileSchema, dbtx_val, ctx)
			if err != nil {
				fmt.Printf("Encountered error on document %v, with error %v ", index, err)
				// Handle error or return from here if needed
			}
		}
	}

	numGoroutines := 20 // You can change this number based on your needs or make it configurable
	// Start workers
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go worker()
	}

	// Send work to workers
	max_files := min(max_docs, len(filteredFiles))
	for index := range max_files {
		fileChan <- index
	}

	// Close channel and wait for workers to finish
	close(fileChan)
	wg.Wait()
	fmt.Printf("Successfully built site map\n")
	return nil
}
