package crud

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mycorrhiza-inc/kessler/backend-go/gen/dbstore"
)

type UpsertHandlerInfo struct {
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
	json_obj := []byte(metadata.JsonObj)
	pgPrivate := pgtype.Bool{false, true}
	if insert {
		args := dbstore.InsertMetadataParams{ID: doc_pgUUID, Isprivate: pgPrivate, Mdata: json_obj}
		metadata_id, err := q.InsertMetadata(
			ctx, args)
		if err != nil {
			return err
		}
		args2 := dbstore.AddMetadataToFileParams{ID: doc_pgUUID, MetadataID: metadata_id}
		_, err = q.AddMetadataToFile(ctx, args2)
		if err != nil {
			return err
		}
		return nil
	}
	args := dbstore.UpdateMetadataParams{ID: doc_pgUUID, Isprivate: pgPrivate, Mdata: json_obj}
	_, err := q.UpdateMetadata(
		ctx, args)
	if err != nil {
		return err
	}

	return nil
}

func upsertFileExtras(ctx context.Context, q dbstore.Queries, doc_uuid uuid.UUID, extras FileGeneratedExtras, insert bool) error {
	doc_pgUUID := pgtype.UUID{Bytes: doc_uuid, Valid: true}
	summary_text := pgtype.Text{String: extras.Summary}
	short_summaries_text := pgtype.Text{String: extras.ShortSummary}
	purpose_text := pgtype.Text{String: extras.Purpose}
	pgPrivate := pgtype.Bool{false, true}
	if insert {
		args := dbstore.ExtrasFileCreateParams{ID: doc_pgUUID, Isprivate: pgPrivate, Summary: summary_text, ShortSummary: short_summaries_text, Purpose: purpose_text}
		_, err := q.ExtrasFileCreate(
			ctx, args)
		if err != nil {
			return err
		}
		return nil
	}
	args := dbstore.ExtrasFileUpdateParams{ID: doc_pgUUID, Isprivate: pgPrivate, Summary: summary_text, ShortSummary: short_summaries_text, Purpose: purpose_text}
	_, err := q.ExtrasFileUpdate(
		ctx, args)
	if err != nil {
		return err
	}
	return nil
}

func makeFileUpsertHandler(info UpsertHandlerInfo) func(w http.ResponseWriter, r *http.Request) {
	dbtx_val := info.dbtx_val
	private := info.private

	insert := info.insert
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
		token := r.Header.Get("Authorization")
		isAuthorizedFunc := func() bool {
			// Enable insert auth at some point
			return true
			if !strings.HasPrefix(token, "Authenticated ") {
				return true
			}
			userID := strings.TrimPrefix(token, "Authenticated ")
			forbiddenPublic := !private && userID != "thaumaturgy"
			if forbiddenPublic || userID == "anonomous" {
				return false
			}
			if !insert {
				authorized, err := checkPrivateFileAuthorization(q, ctx, doc_uuid, userID)
				if !authorized || err == nil {
					return false
				}
			}
			return true
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
			pgUUID := pgtype.UUID{doc_uuid, true}
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

		upsertFileTexts(ctx, q, doc_uuid, newDocInfo.DocTexts, insert)

		upsertFileMetadata(ctx, q, doc_uuid, newDocInfo.Mdata, insert)
		upsertFileExtras(ctx, q, doc_uuid, newDocInfo.Extra, insert)

		response, _ := json.Marshal(fileSchema)

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}
