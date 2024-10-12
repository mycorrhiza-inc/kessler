package crud

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
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
	ID           string
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

// A UUID is a 128 bit (16 byte) Universal Unique IDentifier as defined in RFC
// 4122.
// type UUID [16]byte
func pguuidToString(uuid_pg pgtype.UUID) string {
	return uuid.UUID(uuid_pg.Bytes).String()
}

func RawToFileSchema(file rawFileSchema) (FileSchema, error) {
	var new_mdata map[string]string
	err := json.Unmarshal([]byte(file.Mdata), &new_mdata)
	if err != nil {
		return FileSchema{}, fmt.Errorf("error unmarshaling metadata: %v", err) // err
	}
	return FileSchema{
		ID:           pguuidToString(file.ID),
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

type GetFileParam struct {
	q       dbstore.Queries
	ctx     context.Context
	pgUUID  pgtype.UUID
	private bool
}

func GetTextSchemas(params GetFileParam) ([]FileTextSchema, error) {
	private := params.private
	q := params.q
	ctx := params.ctx
	pgUUID := params.pgUUID
	if private {
		texts, err := q.ListPrivateTextsOfFile(ctx, pgUUID)
		schemas := make([]FileTextSchema, len(texts))
		if err != nil {
			return []FileTextSchema{}, err
		}
		for i, text := range texts {
			schemas[i] = PrivateTextToSchema(text)
		}
		return schemas, nil
	}
	texts, err := q.ListTextsOfFile(ctx, pgUUID)
	schemas := make([]FileTextSchema, len(texts))
	if err != nil {
		return []FileTextSchema{}, err
	}
	for i, text := range texts {
		schemas[i] = PublicTextToSchema(text)
	}
	return schemas, nil
}

func GetFileObjectRaw(params GetFileParam) (rawFileSchema, error) {
	private := params.private
	q := params.q
	ctx := params.ctx
	pgUUID := params.pgUUID

	if !private {
		file, err := q.ReadFile(ctx, pgUUID)
		if err != nil {
			return rawFileSchema{}, err
		}
		return PublicFileToSchema(file), nil
	}
	file, err := q.ReadPrivateFile(ctx, pgUUID)
	if err != nil {
		return rawFileSchema{}, err
	}
	return PrivateFileToSchema(file), nil
}

type FileCreationDataRaw struct {
	Url          pgtype.Text
	Doctype      pgtype.Text
	Lang         pgtype.Text
	Name         pgtype.Text
	Source       pgtype.Text
	Hash         pgtype.Text
	Mdata        pgtype.Text
	Stage        pgtype.Text
	Summary      pgtype.Text
	ShortSummary pgtype.Text
}

func InsertPubPrivateFileObj(q dbstore.Queries, ctx context.Context, fileCreation FileCreationDataRaw, private bool) (rawFileSchema, error) {
	if private {
		params := dbstore.CreatePrivateFileParams{
			Url:          fileCreation.Url,
			Doctype:      fileCreation.Doctype,
			Lang:         fileCreation.Lang,
			Name:         fileCreation.Name,
			Source:       fileCreation.Source,
			Hash:         fileCreation.Hash,
			Mdata:        fileCreation.Mdata,
			Stage:        fileCreation.Stage,
			Summary:      fileCreation.Summary,
			ShortSummary: fileCreation.ShortSummary,
		}
		result, err := q.CreatePrivateFile(ctx, params)
		resultSchema := PrivateFileToSchema(result)
		return resultSchema, err
	}
	params := dbstore.CreateFileParams{
		Url:          fileCreation.Url,
		Doctype:      fileCreation.Doctype,
		Lang:         fileCreation.Lang,
		Name:         fileCreation.Name,
		Source:       fileCreation.Source,
		Hash:         fileCreation.Hash,
		Mdata:        fileCreation.Mdata,
		Stage:        fileCreation.Stage,
		Summary:      fileCreation.Summary,
		ShortSummary: fileCreation.ShortSummary,
	}
	result, err := q.CreateFile(ctx, params)
	resultSchema := PublicFileToSchema(result)
	return resultSchema, err
}
