package crud

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mycorrhiza-inc/kessler/backend-go/gen/dbstore"
)

func upsertFileTexts(ctx context.Context, q dbstore.Queries, doc_uuid uuid.UUID, texts []FileChildTextSource, insert bool) error {
	doc_pgUUID := pgtype.UUID{Bytes: doc_uuid, Valid: true}

	if len(texts) == 0 {
		return nil
	}
	if !insert {
		// TODO: Implement this func to Nuke all the previous texts
		err := NukePriPubFileTexts(q, ctx, doc_pgUUID)
		if err != nil {
			fmt.Print("Error deleting old texts, proceeding with new editions")
			return err
		}
	}
	// TODO : Make Async at some point in future
	error_list := []error{}
	for _, text := range texts {
		textRaw := FileTextSchema{
			FileID:         doc_pgUUID,
			IsOriginalText: text.IsOriginalText,
			Language:       text.Language,
			Text:           text.Text,
		}
		err := InsertPriPubFileText(q, ctx, textRaw, false)
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

func upsertFileMetadata(ctx context.Context, q dbstore.Queries, doc_uuid uuid.UUID, metadata FileMetadataSchema, insert bool) error {
	doc_pgUUID := pgtype.UUID{Bytes: doc_uuid, Valid: true}
	metadata["id"] = doc_uuid.String()
	json_obj, err := json.Marshal(metadata)
	if err != nil {
		return err
	}
	pgPrivate := pgtype.Bool{
		Bool:  false,
		Valid: true,
	}
	if insert {
		insert_args := dbstore.InsertMetadataParams{ID: doc_pgUUID, Isprivate: pgPrivate, Mdata: json_obj}
		_, err := q.InsertMetadata(
			ctx, insert_args)
		if err != nil {
			return err
		}
	}
	args := dbstore.UpdateMetadataParams{ID: doc_pgUUID, Isprivate: pgPrivate, Mdata: json_obj}
	_, err = q.UpdateMetadata(
		ctx, args)
	if err != nil {
		return err
	}

	return nil
}

func juristictionFileUpsert(ctx context.Context, q dbstore.Queries, doc_uuid uuid.UUID, juristiction_info JuristictionInformation, insert bool) error {
	json_obj, err := json.Marshal(juristiction_info.ExtraObject)
	if err != nil {
		return err
	}
	doc_pgUUID := pgtype.UUID{Bytes: doc_uuid, Valid: true}
	country := pgtype.Text{String: juristiction_info.Country}
	state := pgtype.Text{String: juristiction_info.State}
	municipality := pgtype.Text{String: juristiction_info.Municipality}
	agency := pgtype.Text{String: juristiction_info.Agency}
	proceeding_name := pgtype.Text{String: juristiction_info.ProceedingName}

	if insert {
		insert_args := dbstore.JuristictionFileInsertParams{
			ID:             doc_pgUUID,
			Country:        country,
			State:          state,
			Municipality:   municipality,
			Agency:         agency,
			ProceedingName: proceeding_name,
			Extra:          json_obj,
		}
		_, err := q.JuristictionFileInsert(ctx, insert_args)
		if err != nil {
			return err
		}
	}
	update_args := dbstore.JuristictionFileUpdateParams{
		ID:             doc_pgUUID,
		Country:        country,
		State:          state,
		Municipality:   municipality,
		Agency:         agency,
		ProceedingName: proceeding_name,
		Extra:          json_obj,
	}
	_, err = q.JuristictionFileUpdate(ctx, update_args)
	if err != nil {
		return err
	}

	return nil
}

func upsertFileExtras(ctx context.Context, q dbstore.Queries, doc_uuid uuid.UUID, extras FileGeneratedExtras, insert bool) error {
	doc_pgUUID := pgtype.UUID{Bytes: doc_uuid, Valid: true}
	extras_json_obj, err := json.Marshal(extras)
	if err != nil {
		err = fmt.Errorf("error marshalling extras json object, to my understanding this should be absolutely impossible: %v", err)
		fmt.Println(err)
		panic(err)
	}
	pgPrivate := pgtype.Bool{
		Bool:  false,
		Valid: true,
	}
	if insert {
		args := dbstore.ExtrasFileCreateParams{ID: doc_pgUUID, Isprivate: pgPrivate, ExtraObj: extras_json_obj}
		_, err := q.ExtrasFileCreate(
			ctx, args)
		if err != nil {
			return err
		}
		return nil
	}
	args := dbstore.ExtrasFileUpdateParams{ID: doc_pgUUID, Isprivate: pgPrivate, ExtraObj: extras_json_obj}
	_, err = q.ExtrasFileUpdate(
		ctx, args)
	if err != nil {
		return err
	}
	return nil
}

func fileStatusInsert(ctx context.Context, q dbstore.Queries, doc_uuid uuid.UUID, stage DocProcStage, insert bool) error {
	stage_json, err := json.Marshal(stage)
	if err != nil {
		return err
	}
	params := dbstore.StageLogAddParams{
		FileID: pgtype.UUID{Bytes: doc_uuid, Valid: true},
		Status: dbstore.NullStageState{StageState: dbstore.StageState(stage.PGStage)},
		Log:    stage_json,
	}
	_, err = q.StageLogAdd(ctx, params)
	return err
}
