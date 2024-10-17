package static

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

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
func RenderStaticSitemap(dbtx_val dbstore.DBTX) {
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
	proc_func := func(fileSchema crud.FileSchema) {
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
		var html_buffer bytes.Buffer
		err := goldmark.Convert(text, &html_buffer)
		if err != nil {
			fmt.Printf("Error Converting Markdown to HTML", err)
		}
	}
	for _, fileSchema := range allFiles {
		proc_func(fileSchema)
	}
}
