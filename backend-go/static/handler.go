package static

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/mycorrhiza-inc/kessler/backend-go/crud"
	"github.com/mycorrhiza-inc/kessler/backend-go/gen/dbstore"
)

//	params := crud.GetFileParam{
//		q:       q,
//		ctx:     ctx,
//		pgUUID:  pgtype.UUID{Bytes: fileSchema.ID, Valid: true},
//		private: false,
//	}

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

func RenderStaticKesslerObj(obj crud.KesslerObject, dbtx_val dbstore.DBTX) error {
	q := dbstore.New(dbtx_val)

	return nil
}

func RenderStaticSitemap(dbtx_val dbstore.DBTX, max_docs int) error {
	tmpl, err := template.ParseFiles("static/templates/post.html")
	if err != nil {
		return err
	}
	ctx := context.Background()
	chanFileList := make(chan []crud.FileSchema)
	go func() {
		q := *dbstore.New(dbtx_val)
		list_all_files, err := crud.GetListAllFiles(ctx, q)
		if err != nil {
			fmt.Printf("Error encountered while getting all files %s", err)
		}
		var filteredFiles []crud.FileSchema
		for _, file := range list_all_files {
			if file.Stage != "completed" {
				filteredFiles = append(filteredFiles, file)
			}
		}
		chanFileList <- filteredFiles
	}()
	filteredFiles := <-chanFileList
	fmt.Printf("Generating %v static document pages\n", len(filteredFiles))

	// // Could you split this for loop so that the task of processing each element with proc_func is evenly distributed across an arbitrary number of goroutines, s
	// for index, fileSchema := range allFiles {
	// 	err = proc_func(fileSchema)
	// 	if err != nil {
	// 		fmt.Printf("Encountered error on document %v, with error %v ", index, err)
	// 		return fmt.Errorf("encountered error on document %v, with error %s ", index, err)
	// 	}
	// }
	var wg sync.WaitGroup
	fileChan := make(chan int)

	// Worker function
	worker := func() {
		defer wg.Done()
		for index := range fileChan {
			fileSchema := filteredFiles[index]
			err := RenderStaticKesslerObj(fileSchema, dbtx_val)
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
