package crud

import (
	"context"
	"fmt"
	"kessler/internal/objects/files"
	"kessler/internal/objects/timestamp"
	"kessler/internal/dbstore"
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func PublicTextToSchema(file dbstore.FileTextSource) files.FileTextSchema {
	return files.FileTextSchema{
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

func GetTextSchemas(params GetFileParam) ([]files.FileTextSchema, error) {
	q := params.Queries
	ctx := params.Context
	pgUUID := params.PgUUID
	texts, err := q.ListTextsOfFile(ctx, pgUUID)
	schemas := make([]files.FileTextSchema, len(texts))
	if err != nil {
		return []files.FileTextSchema{}, err
	}
	for i, text := range texts {
		schemas[i] = PublicTextToSchema(text)
	}
	return schemas, nil
}

func GetSpecificFileText(params GetFileParam, lang string, original bool) (string, error) {
	prioritize_en := !original && lang == ""

	texts, err := GetTextSchemas(params) // Returns a slice of files.FileTextSchema
	if err != nil || len(texts) == 0 {
		return "", fmt.Errorf("error retrieving texts or no texts found, error: %v", err)
	}
	// TODO: Add suport for non english text retrieval and original text retrieval
	var filteredTexts []files.FileTextSchema

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

func GetFileObjectRaw(params GetFileParam) (files.FileSchema, error) {
	private := params.Private
	q := params.Queries
	ctx := params.Context
	pgUUID := params.PgUUID

	if !private {
		file, err := q.ReadFile(ctx, pgUUID)
		if err != nil {
			return files.FileSchema{}, err
		}
		return PublicFileToSchema(file), nil
	}
	return files.FileSchema{}, fmt.Errorf("private files not implemented")
}

func InsertPubPrivateFileObj(q dbstore.Queries, ctx context.Context, fileCreation files.FileCreationDataRaw, private bool) (files.FileSchema, error) {
	params := dbstore.CreateFileParams{
		Extension:     fileCreation.Extension,
		Lang:          fileCreation.Lang,
		Name:          fileCreation.Name,
		Isprivate:     fileCreation.IsPrivate,
		Hash:          fileCreation.Hash,
		Verified:      fileCreation.Verified,
		DatePublished: pgtype.Timestamptz{Time: time.Time(fileCreation.DatePublished), Valid: true},
	}
	resultID, err := q.CreateFile(ctx, params)
	if err != nil {
		return files.FileSchema{ID: resultID}, err
	}
	resultFile, err := q.ReadFile(ctx, resultID)
	return PublicFileToSchema(resultFile), err
}

func UpdatePubPrivateFileObj(q dbstore.Queries, ctx context.Context, fileCreation files.FileCreationDataRaw, private bool, pgUUID uuid.UUID) (files.FileSchema, error) {
	params := dbstore.UpdateFileParams{
		Extension:     fileCreation.Extension,
		Lang:          fileCreation.Lang,
		Name:          fileCreation.Name,
		Isprivate:     fileCreation.IsPrivate,
		Hash:          fileCreation.Hash,
		Verified:      fileCreation.Verified,
		ID:            pgUUID,
		DatePublished: pgtype.Timestamptz{Time: time.Time(fileCreation.DatePublished), Valid: true},
	}
	err := q.UpdateFile(ctx, params)
	if err != nil {
		log.Info(fmt.Sprintf("Error updating file: %v\n", err))
		return files.FileSchema{}, err
	}

	resultFile, err := q.ReadFile(ctx, pgUUID)
	if err != nil {
		log.Info(fmt.Sprintf("Error reading file after update: %v\n", err))
		return files.FileSchema{}, err
	}
	update_succeded_check := func(updatedFile dbstore.File, fileCreated files.FileCreationDataRaw) error {
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
		return files.FileSchema{}, err
	}

	return PublicFileToSchema(resultFile), nil
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

func PublicFileToSchema(file dbstore.File) files.FileSchema {
	return files.FileSchema{
		ID:            file.ID,
		Verified:      file.Verified.Bool,
		Extension:     file.Extension,
		Lang:          file.Lang,
		Name:          file.Name,
		Hash:          file.Hash,
		IsPrivate:     file.Isprivate.Bool,
		DatePublished: timestamp.KesslerTime(file.DatePublished.Time),
	}
}
