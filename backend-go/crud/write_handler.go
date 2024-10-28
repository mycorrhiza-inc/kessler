package crud

import (
	"encoding/json"
	"fmt"
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

func makeFileUpsertHandler(info UpsertHandlerInfo) func(w http.ResponseWriter, r *http.Request) {
	dbtx_val := info.dbtx_val
	private := info.private
	insert := info.insert
	return func(w http.ResponseWriter, r *http.Request) {
		var doc_uuid uuid.UUID
		var err error
		if !insert {
			params := mux.Vars(r)
			fileID := params["uuid"]

			doc_uuid, err = uuid.Parse(fileID)
			if err != nil {
				http.Error(w, "Error parsing uuid", http.StatusBadRequest)
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
		var newDocInfo CompleteFileSchema
		if err := json.NewDecoder(r.Body).Decode(&newDocInfo); err != nil {

			http.Error(w, "Invalid request payload", http.StatusBadRequest)
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
			fmt.Printf("Error inserting/updating document: %v", err)
			http.Error(w, fmt.Sprintf("Error inserting/updating document: %v", err), http.StatusInternalServerError)
		}
		texts := newDocInfo.DocTexts
		doc_uuid = fileSchema.ID // Ensure UUID is same as one returned from database
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
				err = InsertPriPubFileText(q, ctx, textRaw, private)
				if err != nil {
					fmt.Print("Error adding a text value, not doing anything and procceeding since error handling is hard.")
				}
			}
		}
		response, _ := json.Marshal(fileSchema)

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}
