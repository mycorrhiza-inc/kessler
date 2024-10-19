package static

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path"
	"sync"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mycorrhiza-inc/kessler/backend-go/crud"
	"github.com/mycorrhiza-inc/kessler/backend-go/gen/dbstore"
	"github.com/yuin/goldmark"
)

func HandleStaticGenerationRouting(router *mux.Router, dbtx_val dbstore.DBTX) {
	admin_subrouter := router.PathPrefix("/api/v2/admin").Subrouter()
	admin_subrouter.HandleFunc("/generate-static-site", renderStaticSitemapmMakeHandler(dbtx_val))
}

func GetStaticDir() (string, error) {
	const static_dir = "/static/assets"
	wd, err := os.Getwd()
	return_dir := path.Join(wd, static_dir)
	if err != nil {
		return return_dir, err
	}

	if _, err := os.Stat(return_dir); os.IsNotExist(err) {
		if err := os.MkdirAll(return_dir, os.ModePerm); err != nil {
			return return_dir, err
		}
	}
	return return_dir, nil
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
		chanFileList <- list_all_files
	}()
	allFiles := <-chanFileList
	var filteredFiles []crud.FileSchema
	for _, file := range allFiles {
		if file.Stage != "completed" {
			filteredFiles = append(filteredFiles, file)
		}
	}
	fmt.Printf("Generating %v static document pages\n", len(filteredFiles))
	proc_func := func(fileSchema crud.FileSchema) error {
		fmt.Printf("Found processed file %s with stage %s doing something.\n", fileSchema.ID, fileSchema.Stage)
		q := *dbstore.New(dbtx_val)
		params := crud.GetFileParam{
			Queries: q,
			Context: ctx,
			PgUUID:  pgtype.UUID{Bytes: fileSchema.ID, Valid: true},
			Private: false,
		}
		text, err := crud.GetSpecificFileText(params, "", false)
		if err != nil {
			return fmt.Errorf("encountered error processing file with uuid %s: %v", fileSchema.ID, err)
		}
		text_bytes := []byte(text)
		var html_buffer bytes.Buffer
		err = goldmark.Convert(text_bytes, &html_buffer)
		if err != nil {
			return fmt.Errorf("Error Converting Markdown to HTML", err)
		}
		static_doc_data := StaticDocData{
			HTML:  template.HTML(html_buffer.String()),
			Title: "Test Title",
			Date:  "Test Date",
		}
		static_dir, _ := GetStaticDir()
		file_path := path.Join(static_dir, "/"+fileSchema.ID.String())

		err = os.Remove(file_path)
		if err != nil {
			// return err
		}

		html_file, err := os.Create(file_path)
		if err != nil {
			return err
		}
		defer html_file.Close()
		err = tmpl.Execute(html_file, static_doc_data)
		if err != nil {
			return err
		}
		return nil
	}

	var wg sync.WaitGroup
	fileChan := make(chan int)

	// Worker function
	worker := func() {
		defer wg.Done()
		for index := range fileChan {
			fileSchema := filteredFiles[index]
			err := proc_func(fileSchema)
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
