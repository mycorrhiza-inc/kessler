package crud

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"html/template"
	"os"
	"path"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mycorrhiza-inc/kessler/backend-go/gen/dbstore"
	"github.com/yuin/goldmark"
)

type KesslerObject interface {
	GetShortPath() string
	GetHTMLString(q dbstore.Queries, ctx context.Context) string
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

func (rawFile rawFileSchema) GetShortPath() string {
	uuid_bytes := uuid.UUID(rawFile.ID.Bytes)
	short_id_bytes := uuid_bytes[:6]
	short_id_string := base64.URLEncoding.EncodeToString(short_id_bytes)
	source := rawFile.Source
	name := rawFile.Name
	return "docs/" + source + "/" + name + "-" + short_id_string
}

func (fileSchema FileSchema) GetHTMLString(q dbstore.Queries, ctx context.Context) (string, error) {
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

func getUrl(obj KesslerObject) string {
	shortPath := obj.GetShortPath()
	urlPath := "/static/" + shortPath
	return urlPath
}

func getBaseFilePath(obj KesslerObject) string {
	static_dir, err := GetStaticDir() // Innefficent
	if err != nil {
		fmt.Printf("Error getting base directory %v\n", err)
	}
	shortPath := obj.GetShortPath()
	return path.Join(static_dir, shortPath)
}
