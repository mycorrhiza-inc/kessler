package crud

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mycorrhiza-inc/kessler/backend-go/gen/dbstore"
)

type FileUpsertHandlerConfig struct {
	dbtx_val dbstore.DBTX
	private  bool
	insert   bool
}

func makeFileUpsertHandler(config FileUpsertHandlerConfig) func(w http.ResponseWriter, r *http.Request) {
	dbtx_val := config.dbtx_val
	private := config.private
	insert_parent := config.insert
	deduplicate_with_respect_to_hash := true
	empty_uuid, _ := uuid.Parse("00000000-0000-0000-0000-000000000000")
	return func(w http.ResponseWriter, r *http.Request) {
		// Maybe mutating the underlying parent value is a bit of a problem when it comes to unreachable control code pathways
		insert := insert_parent && true
		var doc_uuid uuid.UUID
		var err error
		if !insert && r.URL.Path == "/v2/public/files/insert" {
			errorstring := "UNREACHABLE CODE: Using insert endpoint with update configuration"
			fmt.Println(errorstring)
			http.Error(w, errorstring, http.StatusInternalServerError)
			return
		}
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
		}
		isAuthorized := isAuthorizedFunc()

		// Usage:
		if !isAuthorized {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		// Proceed with the write operation
		defer r.Body.Close()
		var newDocInfo CompleteFileSchema
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			errorstring := fmt.Sprintf("Error reading request body: %v\n", err)
			fmt.Println(errorstring)
			http.Error(w, errorstring, http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(bodyBytes, &newDocInfo)
		if err != nil {
			errorstring := fmt.Sprintf("Error reading request body json: %v\n", err)
			fmt.Println(errorstring)
			http.Error(w, errorstring, http.StatusBadRequest)
			return
		}
		// Deduplicate with respect to hash
		hash := newDocInfo.Hash
		if hash == "" { // Maybe replace with a comprehensive check to see if the hash is a valid 256 bit base 64 hash
			err := fmt.Errorf("hash is empty, cannot insert document without hash")
			fmt.Println(err)
			http.Error(w, fmt.Sprintf("Error inserting/updating document: %v\n", err), http.StatusBadRequest)
			return
		}
		if insert && deduplicate_with_respect_to_hash {
			ids, err := HashGetUUIDsFile(q, ctx, hash)
			if err != nil {
				errorstring := fmt.Sprintf("Error getting document ids from hash for deduplication: %v\n", err)
				fmt.Println(errorstring)
				http.Error(w, errorstring, http.StatusInternalServerError)
				return
			}
			if len(ids) > 0 {
				id := ids[0]
				if len(ids) > 1 {
					fmt.Printf("More than one document with this hash, this shouldnt happen, continuing deduplication with id: %s\n", id)
				}
				insert = false
				doc_uuid = id
				newDocInfo.ID = id
			}
		}
		rawFileCreationData := ConvertToCreationData(newDocInfo)
		// If we complete all other parts of the file upload process we can set this to true
		// but assuming some parts fail we want the process to fail safe.
		newDocInfo.Verified = false
		rawFileCreationData.Verified = pgtype.Bool{Bool: false, Valid: true}
		var fileSchema FileSchema
		// TODO : For print debugging only, might be a good idea to put these in a deubug logger with lowest priority??
		fmt.Printf("Inserting document with uuid: %s\n", doc_uuid)
		if insert {
			fileSchema, err = InsertPubPrivateFileObj(q, ctx, rawFileCreationData, private)
		} else {
			if doc_uuid == empty_uuid {
				err := fmt.Errorf("ASSERT FAILURE: docUUID should never have a null uuid, when updating document.")
				errorstring := fmt.Sprint(err)
				fmt.Println(errorstring)
				http.Error(w, errorstring, http.StatusInternalServerError)
				return
			}
			pgUUID := pgtype.UUID{
				Bytes: doc_uuid,
				Valid: true,
			}
			fileSchema, err = UpdatePubPrivateFileObj(q, ctx, rawFileCreationData, private, pgUUID)
		}
		if err != nil {
			errorstring := fmt.Sprintf("Error inserting/updating document: %v\n", err)
			fmt.Print(errorstring)
			http.Error(w, errorstring, http.StatusInternalServerError)
			return
		}
		doc_uuid = fileSchema.ID // Ensure UUID is same as one returned from database
		newDocInfo.ID = doc_uuid
		if newDocInfo.ID == empty_uuid {
			err := fmt.Errorf("ASSERT FAILURE: docUUID should never have a null uuid.")
			errorstring := fmt.Sprint(err)
			fmt.Println(errorstring)
			http.Error(w, errorstring, http.StatusInternalServerError)
			return
		}
		has_db_errored := false
		db_error_string := ""

		// TODO : For print debugging only, might be a good idea to put these in a deubug logger with lowest priority??
		fmt.Printf("Attempting to insert all file extras, texts, metadata for uuid: %s\n", doc_uuid)
		if err := upsertFileTexts(ctx, q, doc_uuid, newDocInfo.DocTexts, insert); err != nil {
			errorstring := fmt.Sprintf("Error in upsertFileTexts: %v", err)
			fmt.Println(errorstring)
			has_db_errored = true
			db_error_string = db_error_string + errorstring + "\n"
		}

		// fmt.Printf("Starting upsertFileMetadata for uuid: %s\n", doc_uuid)
		if err := upsertFileMetadata(ctx, q, doc_uuid, newDocInfo.Mdata, insert); err != nil {
			fmt.Printf("Is it getting past the if block?")
			errorstring := fmt.Sprintf("Error in upsertFileMetadata: %v", err)
			fmt.Println(errorstring)
			has_db_errored = true
			db_error_string = db_error_string + errorstring + "\n"
		}

		if err := upsertFileExtras(ctx, q, doc_uuid, newDocInfo.Extra, insert); err != nil {
			errorstring := fmt.Sprintf("Error in upsertFileExtras: %v", err)
			fmt.Println(errorstring)
			has_db_errored = true
			db_error_string = db_error_string + errorstring + "\n"
		}

		// fmt.Printf("Starting fileAuthorsUpsert for uuid: %s\n", doc_uuid)
		if err := fileAuthorsUpsert(ctx, q, doc_uuid, newDocInfo.Authors, insert); err != nil {
			errorstring := fmt.Sprintf("Error in fileAuthorsUpsert: %v", err)
			fmt.Println(errorstring)
			has_db_errored = true
			db_error_string = db_error_string + errorstring + "\n"
		}
		if err := fileConversationUpsert(ctx, q, doc_uuid, newDocInfo.Conversation, insert); err != nil {
			errorstring := fmt.Sprintf("Error in fileConversationUpsert: %v", err)
			fmt.Println(errorstring)
			has_db_errored = true
			db_error_string = db_error_string + errorstring + "\n"
		}

		// fmt.Printf("Starting juristictionFileUpsert for uuid: %s\n", doc_uuid)
		// This doesnt do anything for the time being, so its better to exclude imho
		// if err := juristictionFileUpsert(ctx, q, doc_uuid, newDocInfo.Juristiction, insert); err != nil {
		// 	errorstring := fmt.Sprintf("Error in juristictionFileUpsert: %v", err)
		// 	fmt.Println(errorstring)
		// 	has_db_errored = true
		// 	db_error_string = db_error_string + errorstring + "\n"
		// }
		// Incorperate DB errors into filestatus
		newDocInfo.Stage.IsErrored = newDocInfo.Stage.IsErrored || has_db_errored
		newDocInfo.Stage.DatabaseErrorMsg = db_error_string

		if err := fileStatusInsert(ctx, q, doc_uuid, newDocInfo.Stage); err != nil {
			errorstring := fmt.Sprintf("Error in fileStatusInsert: %v", err)
			fmt.Println(errorstring)
			has_db_errored = true
			// db_error_string = db_error_string + errorstring + "\n"
		}
		encountered_error := newDocInfo.Stage.IsErrored || has_db_errored
		completed_successfully := !encountered_error && newDocInfo.Stage.IsCompleted
		if completed_successfully {
			newDocInfo.Verified = true
			params := dbstore.FileVerifiedUpdateParams{
				Verified: pgtype.Bool{Bool: true, Valid: true},
				ID:       pgtype.UUID{Bytes: doc_uuid, Valid: true},
			}
			_, err := q.FileVerifiedUpdate(ctx, params)
			if err != nil {
				errorstring := fmt.Sprintf("Error in FileVerifiedUpdate, this shouldnt effect anything, but might mean something weird is going on, since this code is only called if every other DB operation succeeded: %v", err)
				fmt.Println(errorstring)
			}
		}

		// TODO : For print debugging only, might be a good idea to put these in a deubug logger with lowest priority??
		fmt.Printf("Completed all DB Operations for uuid, returning schema: %s\n", doc_uuid)

		response, err := json.Marshal(newDocInfo)
		if err != nil {
			errorstring := fmt.Sprintf("Error marshalling response: %v", err)
			fmt.Println(errorstring)
			// http.Error(w, errorstring, http.StatusInternalServerError)
			// return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}