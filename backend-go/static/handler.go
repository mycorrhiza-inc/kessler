package static

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mycorrhiza-inc/kessler/backend-go/crud"
	"github.com/mycorrhiza-inc/kessler/backend-go/gen/dbstore"
	"github.com/yuin/goldmark"
)

//	params := crud.GetFileParam{
//		q:       q,
//		ctx:     ctx,
//		pgUUID:  pgtype.UUID{Bytes: fileSchema.ID, Valid: true},
//		private: false,
//	}
type StaticDocData struct {
	HTML  string
	Title string
	Date  string
}

func HandleStaticGenerationRouter(router *mux.Router, dbtx_val dbstore.DBTX) {
	admin_subrouter := router.PathPrefix("/admin").Subrouter()
	admin_subrouter.HandleFunc("/generate-static-site", renderStaticSitemapmMakeHandler(dbtx_val))
}

func GetStaticDir() (string, error) {
	const static_dir = "/static"
	wd, err := os.Getwd()
	return_dir := path.Join(wd, static_dir)
	return return_dir, err
}

func renderStaticSitemapmMakeHandler(dbtx_val dbstore.DBTX) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := RenderStaticSitemap(dbtx_val)
		if err != nil {
			error_string := fmt.Sprintf("Encountered error while building static site map %v", err)
			http.Error(w, error_string, http.StatusInternalServerError)
			return
		}
		w.Write([]byte("Sucessfully built static site map"))
	}
}

func RenderStaticSitemap(dbtx_val dbstore.DBTX) error {
	tmpl := template.Must(template.ParseFiles("templates/post.html"))
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
	proc_func := func(fileSchema crud.FileSchema) error {
		q := *dbstore.New(dbtx_val)
		params := crud.GetFileParam{
			Queries: q,
			Context: ctx,
			PgUUID:  pgtype.UUID{Bytes: fileSchema.ID, Valid: true},
			Private: false,
		}
		text, err := crud.GetSpecificFileText(params, "", false)
		if err != nil {
			fmt.Printf("encountered error processing file with uuid %s", fileSchema.ID)
		}
		text_bytes := []byte(text)
		var html_buffer bytes.Buffer
		err = goldmark.Convert(text_bytes, &html_buffer)
		if err != nil {
			fmt.Printf("Error Converting Markdown to HTML", err)
		}
		static_doc_data := StaticDocData{
			HTML:  html_buffer.String(),
			Title: "Test Title",
			Date:  "Test Date",
		}
		static_dir, _ := GetStaticDir()
		file_path := path.Join(static_dir, "/"+fileSchema.ID.String())

		html_file, err := os.Open(file_path)
		if err != nil {
			fmt.Print(err)
			log.Fatal(err)
		}
		defer html_file.Close()

		if err := tmpl.Execute(html_file, static_doc_data); err != nil {
			return err
		}
		return nil
	}

	for index, fileSchema := range allFiles {
		err := proc_func(fileSchema)
		if err != nil {
			fmt.Printf("Encountered error on document %v, with error %v ", index, err)
			return fmt.Errorf("Encountered error on document %v, with error %v ", index, err)
		}
	}
	return nil
}
