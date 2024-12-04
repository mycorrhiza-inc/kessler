package crud

import (
	"context"
	"html/template"
	"io"
	"os"
	"path"

	"kessler/gen/dbstore"
)

type StaticDocData struct {
	HTML  template.HTML
	Title string
	Date  string
}

type KesslerObject interface {
	GetShortPath() string
	WriteHTMLString(wr io.Writer, q dbstore.Queries, ctx context.Context) error
}

var doc_template = template.Must(template.ParseFiles("crud/templates/doc_template.html"))

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

// TODO : Refactor with new file schemas
// func (rawFile FileSchema) GetShortPath() string {
// 	uuid_bytes := uuid.UUID(rawFile.ID.Bytes)
// 	short_id_bytes := uuid_bytes[:6]
// 	short_id_string := base64.URLEncoding.EncodeToString(short_id_bytes)
// 	source := rawFile.Source
// 	name := rawFile.Name
// 	return "docs/" + source + "/" + name + "-" + short_id_string
// }
//
// func (fileSchema FileSchema) WriteHTMLString(wr io.Writer, q dbstore.Queries, ctx context.Context) error {
// 	fmt.Printf("Found processed file %s with stage %s doing something.\n", fileSchema.ID, fileSchema.Stage)
// 	params := GetFileParam{
// 		Queries: q,
// 		Context: ctx,
// 		PgUUID:  fileSchema.ID,
// 		Private: false,
// 	}
// 	text, err := GetSpecificFileText(params, "", false)
// 	if err != nil {
// 		return fmt.Errorf("encountered error processing file with uuid %s: %v", fileSchema.ID, err)
// 	}
// 	text_bytes := []byte(text)
// 	var html_buffer bytes.Buffer
// 	err = goldmark.Convert(text_bytes, &html_buffer)
// 	html_string := html_buffer.String()
// 	if err != nil {
// 		io.WriteString(wr, html_string)
// 		return fmt.Errorf("error Converting Markdown to HTML: %v", err)
// 	}
// 	static_doc_data := StaticDocData{
// 		HTML:  template.HTML(html_string),
// 		Title: "Test Title",
// 		Date:  "Test Date",
// 	}
// 	err = doc_template.Execute(wr, static_doc_data)
// 	return err
// }
//
// func GetUrl(obj KesslerObject) string {
// 	shortPath := obj.GetShortPath()
// 	urlPath := "/static/" + shortPath
// 	return urlPath
// }
//
// func GetBaseFilePath(obj KesslerObject) string {
// 	static_dir, err := GetStaticDir() // Innefficent
// 	if err != nil {
// 		fmt.Printf("Error getting base directory %v\n", err)
// 	}
// 	shortPath := obj.GetShortPath()
// 	return path.Join(static_dir, shortPath)
// }
