package files

import (
	"context"
	"fmt"

	"kessler/gen/dbstore"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func PublicTextToSchema(file dbstore.FileTextSource) FileTextSchema {
	return FileTextSchema{
		FileID:         file.FileID,
		IsOriginalText: file.IsOriginalText,
		Text:           file.Text,
		Language:       file.Language,
	}
}

type GetFileParam struct {
	Queries dbstore.Queries
	Context context.Context
	PgUUID  uuid.UUID
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
		return "", fmt.Errorf("error retrieving texts or no texts found, error: %v", err)
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
	return "", fmt.Errorf("no texts found that mach filter criterion")
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
	Extension string
	Lang      string
	Name      string
	Hash      string
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
		return FileSchema{ID: resultID}, err
	}
	resultFile, err := q.ReadFile(ctx, resultID)
	return PublicFileToSchema(resultFile), err
}

func UpdatePubPrivateFileObj(q dbstore.Queries, ctx context.Context, fileCreation FileCreationDataRaw, private bool, pgUUID uuid.UUID) (FileSchema, error) {
	params := dbstore.UpdateFileParams{
		Extension: fileCreation.Extension,
		Lang:      fileCreation.Lang,
		Name:      fileCreation.Name,
		Isprivate: fileCreation.IsPrivate,
		Hash:      fileCreation.Hash,
		Verified:  fileCreation.Verified,
		ID:        pgUUID,
	}
	err := q.UpdateFile(ctx, params)
	if err != nil {
		fmt.Printf("Error updating file: %v\n", err)
		return FileSchema{}, err
	}

	resultFile, err := q.ReadFile(ctx, pgUUID)
	if err != nil {
		fmt.Printf("Error reading file after update: %v\n", err)
		return FileSchema{}, err
	}
	update_succeded_check := func(updatedFile dbstore.File, fileCreated FileCreationDataRaw) error {
		var mismatches []string
		if updatedFile.Isprivate != fileCreated.IsPrivate {
			mismatches = append(mismatches, fmt.Sprintf("private (got: %v, want: %v)", updatedFile.Isprivate, fileCreated.IsPrivate))
		}
		if updatedFile.Lang != fileCreated.Lang {
			mismatches = append(mismatches, fmt.Sprintf("lang (got: %v, want: %v)", updatedFile.Lang, fileCreated.Lang))
		}
		if updatedFile.Name != fileCreated.Name {
			mismatches = append(mismatches, fmt.Sprintf("name (got: %v, want: %v)", updatedFile.Name, fileCreated.Name))
		}
		if updatedFile.Hash != fileCreated.Hash {
			mismatches = append(mismatches, fmt.Sprintf("hash (got: %v, want: %v)", updatedFile.Hash, fileCreated.Hash))
		}
		if updatedFile.Verified != fileCreated.Verified {
			mismatches = append(mismatches, fmt.Sprintf("verified (got: %v, want: %v)", updatedFile.Verified, fileCreated.Verified))
		}
		if len(mismatches) > 0 {
			return fmt.Errorf("encountered mismatched fields while updating: %v", mismatches)
		}
		return nil
	}
	if err := update_succeded_check(resultFile, fileCreation); err != nil {
		return FileSchema{}, err
	}

	return PublicFileToSchema(resultFile), nil
}

func InsertPriPubFileText(q dbstore.Queries, ctx context.Context, text FileTextSchema, private bool) error {
	args := dbstore.CreateFileTextSourceParams{
		FileID:         text.FileID,
		IsOriginalText: text.IsOriginalText,
		Text:           text.Text,
		Language:       text.Language,
	}
	_, err := q.CreateFileTextSource(ctx, args)
	return err
}

func HashGetUUIDsFile(q dbstore.Queries, ctx context.Context, hash string) ([]uuid.UUID, error) {
	filePGUUIDs, err := q.HashGetFileID(ctx, hash)
	if err != nil {
		return nil, err
	}
	var fileIDs []uuid.UUID
	for _, file := range filePGUUIDs {
		fileUUID := uuid.UUID(file)
		fileIDs = append(fileIDs, fileUUID)
	}
	return fileIDs, nil
}
