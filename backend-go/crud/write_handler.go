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

type UpsertHandlerConfig struct {
	dbtx_val dbstore.DBTX
	private  bool
	insert   bool
}

func makeFileUpsertHandler(config UpsertHandlerConfig) func(w http.ResponseWriter, r *http.Request) {
	dbtx_val := config.dbtx_val
	private := config.private
	insert := config.insert
	deduplicate_with_respect_to_hash := true
	empty_uuid, _ := uuid.Parse("00000000-0000-0000-0000-000000000000")
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

		rawFileData := ConvertToCreationData(newDocInfo)
		var fileSchema FileSchema
		if insert {
			fileSchema, err = InsertPubPrivateFileObj(q, ctx, rawFileData, private)
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
		if newDocInfo.ID == empty_uuid {
			err := fmt.Errorf("ASSERT FAILURE: docUUID should never have a null uuid.")
			errorstring := fmt.Sprint(err)
			fmt.Println(errorstring)
			http.Error(w, errorstring, http.StatusInternalServerError)
			return
		}

		// TODO: Implement error handling for any of these.
		if err := fileStatusInsert(ctx, q, doc_uuid, newDocInfo.Stage, insert); err != nil {
			errorstring := fmt.Sprintf("Error in fileStatusInsert: %v", err)
			fmt.Println(errorstring)
			// http.Error(w, errorstring, http.StatusInternalServerError)
			// return
		}

		if err := upsertFileTexts(ctx, q, doc_uuid, newDocInfo.DocTexts, insert); err != nil {
			errorstring := fmt.Sprintf("Error in upsertFileTexts: %v", err)
			fmt.Println(errorstring)
			// http.Error(w, errorstring, http.StatusInternalServerError)
			// return
		}

		if err := upsertFileMetadata(ctx, q, doc_uuid, newDocInfo.Mdata, insert); err != nil {
			errorstring := fmt.Sprintf("Error in upsertFileMetadata: %v", err)
			fmt.Println(errorstring)
			// http.Error(w, errorstring, http.StatusInternalServerError)
			// return
		}

		if err := upsertFileExtras(ctx, q, doc_uuid, newDocInfo.Extra, insert); err != nil {
			errorstring := fmt.Sprintf("Error in upsertFileExtras: %v", err)
			fmt.Println(errorstring)
			// http.Error(w, errorstring, http.StatusInternalServerError)
			// return
		}

		if err := fileAuthorsUpsert(ctx, q, doc_uuid, newDocInfo.Authors, insert); err != nil {
			errorstring := fmt.Sprintf("Error in fileAuthorsUpsert: %v", err)
			fmt.Println(errorstring)
			// http.Error(w, errorstring, http.StatusInternalServerError)
			// return
		}

		if err := juristictionFileUpsert(ctx, q, doc_uuid, newDocInfo.Juristiction, insert); err != nil {
			errorstring := fmt.Sprintf("Error in juristictionFileUpsert: %v", err)
			fmt.Println(errorstring)
			// http.Error(w, errorstring, http.StatusInternalServerError)
			// return
		}

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
