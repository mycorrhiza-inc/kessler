package crud

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mycorrhiza-inc/kessler/backend-go/gen/dbstore"
)

type UpsertHandlerConfig struct {
	dbtx_val dbstore.DBTX
	private  bool
	insert   bool
}

func upsertFileTexts(ctx context.Context, q dbstore.Queries, doc_uuid uuid.UUID, texts []FileChildTextSource, insert bool) {
	doc_pgUUID := pgtype.UUID{Bytes: doc_uuid, Valid: true}
	if len(texts) != 0 {
		if !insert {
			// TODO: Implement this func to Nuke all the previous texts
			err := NukePriPubFileTexts(q, ctx, doc_pgUUID)
			if err != nil {
				fmt.Print("Error deleting old texts, proceeding with new editions")
			}
		}
		// TODO : Make Async at some point in future
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
			}
		}
	}
}

func upsertFileMetadata(ctx context.Context, q dbstore.Queries, doc_uuid uuid.UUID, metadata FileMetadataSchema, insert bool) error {
	doc_pgUUID := pgtype.UUID{Bytes: doc_uuid, Valid: true}
	metadata.MdataObject["id"] = doc_uuid.String()
	json_obj, err := json.Marshal(metadata.MdataObject)
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
	params := dbstore.AddStageLogParams{
		FileID: pgtype.UUID{Bytes: doc_uuid, Valid: true},
		Status: dbstore.NullStageState{StageState: dbstore.StageState(stage.PGStage)},
		Log:    stage_json,
	}
	_, err = q.AddStageLog(ctx, params)
	return err
}

func makeFileUpsertHandler(config UpsertHandlerConfig) func(w http.ResponseWriter, r *http.Request) {
	dbtx_val := config.dbtx_val
	private := config.private
	insert := config.insert
	return func(w http.ResponseWriter, r *http.Request) {
		var doc_uuid uuid.UUID
		var err error
		if !insert {
			params := mux.Vars(r)
			fileIDString := params["uuid"]

			doc_uuid, err = uuid.Parse(fileIDString)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error parsing uuid: %s\n", fileIDString), http.StatusBadRequest)
				return
			}
		}
		q := *dbstore.New(dbtx_val)
		ctx := r.Context()
		// TODO!!!!!: Enable insert auth at some point
		// token := r.Header.Get("Authorization")
		isAuthorizedFunc := func() bool {
			return true
			// if !strings.HasPrefix(token, "Authenticated ") {
			// 	return true
			// }
			// userID := strings.TrimPrefix(token, "Authenticated ")
			// forbiddenPublic := !private && userID != "thaumaturgy"
			// if forbiddenPublic || userID == "anonomous" {
			// 	return false
			// }
			// if !insert {
			// 	authorized, err := checkPrivateFileAuthorization(q, ctx, doc_uuid, userID)
			// 	if !authorized || err == nil {
			// 		return false
			// 	}
			// }
			// return true
		}
		isAuthorized := isAuthorizedFunc()

		// Usage:
		if !isAuthorized {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		// Proceed with the write operation
		// TODO: IF user is not a paying user, disable insert functionality
		defer r.Body.Close()
		var newDocInfo CompleteFileSchema
		// err = json.NewDecoder(r.Body).Decode(&newDocInfo)
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			errorstring := fmt.Sprintf("Error reading request body: %v\n", err)
			fmt.Println(errorstring)

			http.Error(w, errorstring, http.StatusBadRequest)
			return
		}
		blah := fmt.Sprintln(string(bodyBytes))
		fmt.Printf("%v\n", blah)
		// var testUnmarshal TestCompleteFileSchema
		// err = json.Unmarshal([]byte(blah), &testUnmarshal)
		err = json.Unmarshal([]byte(blah), &newDocInfo)
		if err != nil {
			errorstring := fmt.Sprintf("Error reading request body json: %v\n", err)
			fmt.Println(errorstring)

			http.Error(w, errorstring, http.StatusBadRequest)
			return
		}
		rawFileData := ConvertToCreationData(newDocInfo)
		var fileSchema FileSchema
		if insert {
			fileSchema, err = InsertPubPrivateFileObj(q, ctx, rawFileData, private)
		} else {
			pgUUID := pgtype.UUID{
				Bytes: doc_uuid,
				Valid: true,
			}
			fileSchema, err = UpdatePubPrivateFileObj(q, ctx, rawFileData, private, pgUUID)
		}
		if err != nil {
			errorstring := fmt.Sprintf("Error inserting/updating document: %v\n", err)
			fmt.Print(errorstring)
			http.Error(w, errorstring, http.StatusInternalServerError)
			return
		}
		doc_uuid = fileSchema.ID // Ensure UUID is same as one returned from database
		newDocInfo.ID = doc_uuid
		empty_uuid, _ := uuid.Parse("00000000-0000-0000-0000-000000000000")
		if newDocInfo.ID == empty_uuid {
			err := fmt.Errorf("ASSERT FAILURE: docUUID should never have a null uuid.")
			errorstring := fmt.Sprint(err)
			fmt.Println(errorstring)
			http.Error(w, errorstring, http.StatusInternalServerError)
			return
		}

		// TODO: Implement error handling for any of these.
		fileStatusInsert(ctx, q, doc_uuid, newDocInfo.Stage, insert)

		upsertFileTexts(ctx, q, doc_uuid, newDocInfo.DocTexts, insert)

		upsertFileMetadata(ctx, q, doc_uuid, newDocInfo.Mdata, insert)
		upsertFileExtras(ctx, q, doc_uuid, newDocInfo.Extra, insert)
		// TODO: Implement update functionality
		fileAuthorsUpsert(ctx, q, doc_uuid, newDocInfo.Authors, insert)
		juristictionFileUpsert(ctx, q, doc_uuid, newDocInfo.Juristiction, insert)

		response, _ := json.Marshal(newDocInfo)

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}
