package crud

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mycorrhiza-inc/kessler/backend-go/gen/dbstore"
)

type RawFileSchema struct {
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
	ID           uuid.UUID
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

func RawToFileSchema(file RawFileSchema) (FileSchema, error) {
	// fmt.Println(file.ID)
	var new_mdata map[string]string
	err := json.Unmarshal([]byte(file.Mdata), &new_mdata)
	if err != nil {
		// fmt.Printf("Error unmarhalling Metadata: %v\n", err)
		return FileSchema{
			ID:           file.ID.Bytes,
			Url:          file.Url,
			Doctype:      file.Doctype,
			Lang:         file.Lang,
			Name:         file.Name,
			Source:       file.Source,
			Hash:         file.Hash,
			Mdata:        map[string]string{},
			Stage:        file.Stage,
			Summary:      file.Summary,
			ShortSummary: file.ShortSummary,
		}, fmt.Errorf("error unmarshaling metadata: %v", err) // err
	}
	return FileSchema{
		ID:           file.ID.Bytes,
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

func PrivateFileToSchema(file dbstore.UserfilesPrivateFile) RawFileSchema {
	return RawFileSchema{
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

func PublicFileToSchema(file dbstore.File) RawFileSchema {
	return RawFileSchema{
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
	FileID         pgtype.UUID
	IsOriginalText bool
	Text           string
	Language       string
}

func PrivateTextToSchema(file dbstore.UserfilesPrivateFileTextSource) FileTextSchema {
	return FileTextSchema{
		FileID:         file.FileID,
		IsOriginalText: file.IsOriginalText,
		Text:           file.Text.String,
		Language:       file.Language,
	}
}

func PublicTextToSchema(file dbstore.FileTextSource) FileTextSchema {
	return FileTextSchema{
		FileID:         file.FileID,
		IsOriginalText: file.IsOriginalText,
		Text:           file.Text.String,
		Language:       file.Language,
	}
}

type GetFileParam struct {
	Queries dbstore.Queries
	Context context.Context
	PgUUID  pgtype.UUID
	Private bool
}

func GetTextSchemas(params GetFileParam) ([]FileTextSchema, error) {
	// params_1 := GetFileParam{
	// 	q:       q,
	// 	ctx:     ctx,
	// 	pgUUID:  pgtype.UUID{Bytes: fileSchema.UUID, Valid: true},
	// 	private: false,
	// }
	private := params.Private
	q := params.Queries
	ctx := params.Context
	pgUUID := params.PgUUID
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

func GetSpecificFileText(params GetFileParam, lang string, original bool) (string, error) {
	prioritize_en := !original && lang == ""

	texts, err := GetTextSchemas(params) // Returns a slice of FileTextSchema
	if err != nil || len(texts) == 0 {
		return "", fmt.Errorf("Error retrieving texts or no texts found, error: %v", err)
	}
	// TODO: Add suport for non english text retrieval and original text retrieval
	var filteredTexts []FileTextSchema

	for _, text := range texts {
		if prioritize_en && text.Language == "en" {
			return text.Text, nil
		}
		originalIfUserCares := !original || text.IsOriginalText
		matchLangIfUserCares := lang == "" || text.Language == lang
		if originalIfUserCares && matchLangIfUserCares {
			filteredTexts = append(filteredTexts, text)
		}
	}
	if prioritize_en {
		return texts[0].Text, nil
	}
	if len(filteredTexts) > 0 {
		return filteredTexts[0].Text, nil
	}
	return "", fmt.Errorf("No texts found that mach filter criterion")
}

func GetFileObjectRaw(params GetFileParam) (RawFileSchema, error) {
	private := params.Private
	q := params.Queries
	ctx := params.Context
	pgUUID := params.PgUUID

	if !private {
		file, err := q.ReadFile(ctx, pgUUID)
		if err != nil {
			return RawFileSchema{}, err
		}
		return PublicFileToSchema(file), nil
	}
	file, err := q.ReadPrivateFile(ctx, pgUUID)
	if err != nil {
		return RawFileSchema{}, err
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

func InsertPubPrivateFileObj(q dbstore.Queries, ctx context.Context, fileCreation FileCreationDataRaw, private bool) (RawFileSchema, error) {
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

func UpdatePubPrivateFileObj(q dbstore.Queries, ctx context.Context, fileCreation FileCreationDataRaw, private bool, pgUUID pgtype.UUID) (RawFileSchema, error) {
	if private {
		params := dbstore.UpdatePrivateFileParams{
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
			ID:           pgUUID,
		}
		result, err := q.UpdatePrivateFile(ctx, params)
		resultSchema := PrivateFileToSchema(result)
		return resultSchema, err
	}
	params := dbstore.UpdateFileParams{
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
		ID:           pgUUID,
	}
	result, err := q.UpdateFile(ctx, params)
	resultSchema := PublicFileToSchema(result)
	return resultSchema, err
}

func NukePriPubFileTexts(q dbstore.Queries, ctx context.Context, pgUUID pgtype.UUID) error {
	return nil
}

func InsertPriPubFileText(q dbstore.Queries, ctx context.Context, text FileTextSchema, private bool) error {
	if private {
		args := dbstore.CreatePrivateFileTextSourceParams{
			FileID:         text.FileID,
			IsOriginalText: text.IsOriginalText,
			Text:           pgtype.Text{String: text.Text, Valid: true},
			Language:       text.Language,
		}
		_, err := q.CreatePrivateFileTextSource(ctx, args)
		return err
	}
	args := dbstore.CreateFileTextSourceParams{
		FileID:         text.FileID,
		IsOriginalText: text.IsOriginalText,
		Text:           pgtype.Text{String: text.Text, Valid: true},
		Language:       text.Language,
	}
	_, err := q.CreateFileTextSource(ctx, args)
	return err
}
