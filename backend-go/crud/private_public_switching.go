package crud

import (
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mycorrhiza-inc/kessler/backend-go/gen/dbstore"
)

type rawFileSchema struct {
	ID           pgtype.UUID
	Url          string
	Doctype      string
	Lang         string
	Name         string
	Source       string
	Hash         string
	Mdata        string
	Stage        string
	Summary      string
	ShortSummary string
}
type FileSchema struct {
	ID           pgtype.UUID
	Url          string
	Doctype      string
	Lang         string
	Name         string
	Source       string
	Hash         string
	Mdata        map[string]string
	Stage        string
	Summary      string
	ShortSummary string
}

func RawToFileSchema(file rawFileSchema) (FileSchema, error) {
	var new_mdata map[string]string
	err := json.Unmarshal([]byte(file.Mdata), new_mdata)
	if err != nil {
		return FileSchema{}, fmt.Errorf("error unmarshaling metadata: %v", err) // err
	}
	return FileSchema{
		ID:           file.ID,
		Url:          file.Url,
		Doctype:      file.Doctype,
		Lang:         file.Lang,
		Name:         file.Name,
		Source:       file.Source,
		Hash:         file.Hash,
		Mdata:        new_mdata,
		Stage:        file.Stage,
		Summary:      file.Summary,
		ShortSummary: file.ShortSummary,
	}, nil
}

func PrivateFileToSchema(file dbstore.UserfilesPrivateFile) rawFileSchema {
	return rawFileSchema{
		ID:           file.ID,
		Url:          file.Url.String,
		Doctype:      file.Doctype.String,
		Lang:         file.Lang.String,
		Name:         file.Name.String,
		Source:       file.Source.String,
		Hash:         file.Hash.String,
		Mdata:        file.Mdata.String,
		Stage:        file.Stage.String,
		Summary:      file.Summary.String,
		ShortSummary: file.ShortSummary.String,
	}
}

func PublicFileToSchema(file dbstore.File) rawFileSchema {
	return rawFileSchema{
		ID:           file.ID,
		Url:          file.Url.String,
		Doctype:      file.Doctype.String,
		Lang:         file.Lang.String,
		Name:         file.Name.String,
		Source:       file.Source.String,
		Hash:         file.Hash.String,
		Mdata:        file.Mdata.String,
		Stage:        file.Stage.String,
		Summary:      file.Summary.String,
		ShortSummary: file.ShortSummary.String,
	}
}

type FileTextSchema struct {
	ID             pgtype.UUID
	FileID         pgtype.UUID
	IsOriginalText bool
	Text           string
}

func PrivateTextToSchema(file dbstore.UserfilesPrivateFileTextSource) FileTextSchema {
	return FileTextSchema{
		ID:             file.ID,
		FileID:         file.FileID,
		IsOriginalText: file.IsOriginalText,
		Text:           file.Text.String,
	}
}

func PublicTextToSchema(file dbstore.FileTextSource) FileTextSchema {
	return FileTextSchema{
		ID:             file.ID,
		FileID:         file.FileID,
		IsOriginalText: file.IsOriginalText,
		Text:           file.Text.String,
	}
}
