package crud

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mycorrhiza-inc/kessler/backend-go/gen/dbstore"
)

type FileSchema struct {
	ID        uuid.UUID `json:"id"`
	Verified  bool      `json:"verified"`
	Extension string    `json:"extension"`
	Lang      string    `json:"lang"`
	Name      string    `json:"name"`
	Hash      string    `json:"hash"`
	IsPrivate bool      `json:"is_private"`
}

// A UUID is a 128 bit (16 byte) Universal Unique IDentifier as defined in RFC
// 4122.
// type UUID [16]byte
func pguuidToString(uuid_pg pgtype.UUID) string {
	return uuid.UUID(uuid_pg.Bytes).String()
}

func PublicFileToSchema(file dbstore.File) FileSchema {
	return FileSchema{
		ID:        file.ID.Bytes,
		Verified:  file.Verified.Bool,
		Extension: file.Extension.String,
		Lang:      file.Lang.String,
		Name:      file.Name.String,
		Hash:      file.Hash.String,
		IsPrivate: file.Isprivate.Bool,
	}
}

type FileTextSchema struct {
	FileID         pgtype.UUID `json:"file_id"`
	IsOriginalText bool        `json:"is_original_text"`
	Text           string      `json:"text"`
	Language       string      `json:"language"`
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
	q := params.Queries
	ctx := params.Context
	pgUUID := params.PgUUID
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

func GetFileObjectRaw(params GetFileParam) (FileSchema, error) {
	private := params.Private
	q := params.Queries
	ctx := params.Context
	pgUUID := params.PgUUID

	if !private {
		file, err := q.ReadFile(ctx, pgUUID)
		if err != nil {
			return FileSchema{}, err
		}
		return PublicFileToSchema(file), nil
	}
	return FileSchema{}, fmt.Errorf("private files not implemented")
}

type FileCreationDataRaw struct {
	Extension pgtype.Text
	Lang      pgtype.Text
	Name      pgtype.Text
	Hash      pgtype.Text
	IsPrivate pgtype.Bool
	Verified  pgtype.Bool
}

func InsertPubPrivateFileObj(q dbstore.Queries, ctx context.Context, fileCreation FileCreationDataRaw, private bool) (FileSchema, error) {
	params := dbstore.CreateFileParams{
		Extension: fileCreation.Extension,
		Lang:      fileCreation.Lang,
		Name:      fileCreation.Name,
		Isprivate: fileCreation.IsPrivate,
		Hash:      fileCreation.Hash,
		Verified:  fileCreation.Verified,
	}
	resultID, err := q.CreateFile(ctx, params)
	if err != nil {
		return FileSchema{ID: resultID.Bytes}, err
	}
	resultFile, err := q.ReadFile(ctx, resultID)
	return PublicFileToSchema(resultFile), err
}

func UpdatePubPrivateFileObj(q dbstore.Queries, ctx context.Context, fileCreation FileCreationDataRaw, private bool, pgUUID pgtype.UUID) (FileSchema, error) {
	params := dbstore.UpdateFileParams{
		Extension: fileCreation.Extension,
		Lang:      fileCreation.Lang,
		Name:      fileCreation.Name,
		Isprivate: fileCreation.IsPrivate,
		Hash:      fileCreation.Hash,
		Verified:  fileCreation.Verified,
	}
	resultID, err := q.UpdateFile(ctx, params)
	if err != nil {
		fmt.Printf("Error updating file: %v\n", err)
		return FileSchema{ID: resultID.Bytes}, err
	}
	resultFile, err := q.ReadFile(ctx, resultID)
	if err != nil {
		fmt.Printf("Error reading file after update: %v\n", err)
		return FileSchema{ID: resultID.Bytes}, err
	}
	return PublicFileToSchema(resultFile), nil
}

func NukePriPubFileTexts(q dbstore.Queries, ctx context.Context, pgUUID pgtype.UUID) error {
	return nil
}

func InsertPriPubFileText(q dbstore.Queries, ctx context.Context, text FileTextSchema, private bool) error {
	args := dbstore.CreateFileTextSourceParams{
		FileID:         text.FileID,
		IsOriginalText: text.IsOriginalText,
		Text:           pgtype.Text{String: text.Text, Valid: true},
		Language:       text.Language,
	}
	_, err := q.CreateFileTextSource(ctx, args)
	return err
}

func HashGetUUIDsFile(q dbstore.Queries, ctx context.Context, hash string) ([]uuid.UUID, error) {
	pgHash := pgtype.Text{String: hash, Valid: true}
	filePGUUIDs, err := q.HashGetFileID(ctx, pgHash)
	if err != nil {
		return nil, err
	}
	var fileIDs []uuid.UUID
	for _, file := range filePGUUIDs {
		fileUUID := uuid.UUID(file.Bytes)
		fileIDs = append(fileIDs, fileUUID)
	}
	return fileIDs, nil
}
