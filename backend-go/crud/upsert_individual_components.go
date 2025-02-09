package crud

import (
	"context"
	"encoding/json"
	"fmt"

	"kessler/gen/dbstore"
	"kessler/objects/files"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func upsertFileTexts(ctx context.Context, q dbstore.Queries, doc_uuid uuid.UUID, texts []files.FileChildTextSource, insert bool) error {
	// Sometimes this is getting called with an insert when the metadata already exists in the table, this causes a PGERROR, since it violates uniqueness. However, setting it up so it tries to update will fall back to insert if the file doesnt exist. Its probably a good idea to remove this and debug what is causing the new file thing at some point.
	// I think I might have solved the error, it was only happenining in the other ones so its added here for an abundance of safety and only degrades perf slightly
	insert = false
	if len(texts) == 0 {
		return nil
	}
	if !insert {
		// TODO: Implement this func to Nuke all the previous texts
		// err := NukePriPubFileTexts(q, ctx, doc_pgUUID)
		// if err != nil {
		// 	fmt.Print("Error deleting old texts, proceeding with new editions")
		// 	return err
		// }
	}
	// TODO : Make Async at some point in future
	error_list := []error{}
	for _, text := range texts {
		textRaw := files.FileTextSchema{
			FileID:         doc_uuid,
			IsOriginalText: text.IsOriginalText,
			Language:       text.Language,
			Text:           text.Text,
		}
		err := files.InsertPriPubFileText(q, ctx, textRaw, false)
		if err != nil {
			fmt.Print("Error adding a text value, not doing anything and procceeding since error handling is hard.")
			error_list = append(error_list, err)
		}
	}
	if len(error_list) > 0 {
		return error_list[0]
	}
	return nil
}

func upsertFileMetadata(ctx context.Context, q dbstore.Queries, doc_uuid uuid.UUID, metadata files.FileMetadataSchema, insert bool) error {
	// Sometimes this is getting called with an insert when the metadata already exists in the table, this causes a PGERROR, since it violates uniqueness. However, setting it up so it tries to update will fall back to insert if the file doesnt exist. Its probably a good idea to remove this and debug what is causing the new file thing at some point.
	// UPDATE: I am pretty sure I solved it this should be safe to take out soon - nic
	insert = false
	metadata["id"] = doc_uuid.String()
	// log.Info(fmt.Sprintf("Is it the json marshall?\n"))

	json_obj, err := json.Marshal(metadata)
	if err != nil {
		return err
	}
	pgPrivate := pgtype.Bool{
		Bool:  false,
		Valid: true,
	}

	// log.Info(fmt.Sprintf("Wasnt it, is it any of the db operations?\n"))
	if !insert {
		args := dbstore.UpdateMetadataParams{ID: doc_uuid, Isprivate: pgPrivate, Mdata: json_obj}
		// log.Info(fmt.Sprintf("Defined args successfully\n"))
		_, err = q.UpdateMetadata(ctx, args)
		// If encounter a not found error, break error handling control flow and inserting the file metadata.
		if err == nil {
			return nil
		}
		if err.Error() != "no rows in result set" {
			return err
		}
	}
	// log.Info(fmt.Sprintf("Failed and trying to insert metadta instead"))
	insert_args := dbstore.InsertMetadataParams{ID: doc_uuid, Isprivate: pgPrivate, Mdata: json_obj}
	_, err = q.InsertMetadata(ctx, insert_args)
	// log.Info(fmt.Sprintf("What the actual fuck is going on?"))
	return err
}

func upsertFileExtras(ctx context.Context, q dbstore.Queries, doc_uuid uuid.UUID, extras files.FileGeneratedExtras, insert bool) error {
	// Sometimes this is getting called with an insert when the metadata already exists in the table, this causes a PGERROR, since it violates uniqueness. However, setting it up so it tries to update will fall back to insert if the file doesnt exist. Its probably a good idea to remove this and debug what is causing the new file thing at some point.
	// UPDATE: I am pretty sure I solved it this should be safe to take out soon - nic
	insert = false
	extras_json_obj, err := json.Marshal(extras)
	if err != nil {
		err = fmt.Errorf("error marshalling extras json object, to my understanding this should be absolutely impossible: %v", err)
		log.Info(err)
		panic(err)
	}
	pgPrivate := pgtype.Bool{
		Bool:  false,
		Valid: true,
	}
	if !insert {
		args := dbstore.ExtrasFileUpdateParams{ID: doc_uuid, Isprivate: pgPrivate, ExtraObj: extras_json_obj}
		_, err = q.ExtrasFileUpdate(ctx, args)
		// If encounter a not found error, break error handling control flow and inserting object
		if err == nil {
			return nil
		}
		if err.Error() != "no rows in result set" {
			// If the error is nil, this still returns the error
			return err
		}
	}
	args := dbstore.ExtrasFileCreateParams{ID: doc_uuid, Isprivate: pgPrivate, ExtraObj: extras_json_obj}
	_, err = q.ExtrasFileCreate(
		ctx, args)
	return err
}

func fileStatusInsert(ctx context.Context, q dbstore.Queries, doc_uuid uuid.UUID, stage files.DocProcStage) error {
	stage_json, err := json.Marshal(stage)
	if err != nil {
		return err
	}
	params := dbstore.StageLogAddParams{
		FileID: doc_uuid,
		Status: dbstore.NullStageState{StageState: dbstore.StageState(stage.PGStage)},
		Log:    stage_json,
	}
	_, err = q.StageLogAdd(ctx, params)
	return err
}

func fileStatusGetLatestStage(ctx context.Context, q dbstore.Queries, doc_uuid uuid.UUID) (files.DocProcStage, error) {
	result_pgschema, err := q.StageLogFileGetLatest(ctx, doc_uuid)
	return_obj := files.DocProcStage{}
	if err != nil {
		return files.DocProcStage{}, err
	}
	result_json := result_pgschema.Log
	err = json.Unmarshal(result_json, &return_obj)
	if err != nil {
		return files.DocProcStage{}, err
	}
	return return_obj, nil
}
