package crud

import (
	"context"
	"encoding/json"
	"fmt"
	"kessler/common/objects/files"
	"kessler/internal/dbstore"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func upsertFileAttachmentTexts(ctx context.Context, q dbstore.Queries, attachment_uuid uuid.UUID, texts []files.AttachmentChildTextSource, insert bool) error {
	error_list := []error{}
	for _, text := range texts {
		textRaw := dbstore.AttachmentTextCreateParams{
			AttachmentID:   attachment_uuid,
			Language:       text.Language,
			IsOriginalText: text.IsOriginalText,
			Text:           text.Text,
		}
		_, err := q.AttachmentTextCreate(ctx, textRaw)
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

func upsertFileAttachments(ctx context.Context, q dbstore.Queries, doc_uuid uuid.UUID, attachments []files.CompleteAttachmentSchema, insert bool) error {
	if !insert {
		// Delete previous file attachments
	}
	for _, attachment := range attachments {
		attachment_insert_args := dbstore.AttachmentCreateParams{
			FileID:    doc_uuid,
			Name:      attachment.Name,
			Extension: attachment.Extension,
			Hash:      attachment.Hash,
			Lang:      attachment.Lang,
			Mdata:     []byte{},
		}
		pg_attachment, err := q.AttachmentCreate(ctx, attachment_insert_args)
		if err != nil {
			return err
		}
		err = upsertFileAttachmentTexts(ctx, q, pg_attachment.ID, attachment.Texts, insert)
		if err != nil {
			return err
		}
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
