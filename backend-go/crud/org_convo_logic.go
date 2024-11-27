package crud

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mycorrhiza-inc/kessler/backend-go/gen/dbstore"
)

// Make this generic
func getFirstElement(array []dbstore.Organization) (dbstore.Organization, error) {
	if len(array) == 0 {
		return dbstore.Organization{}, fmt.Errorf("List had no elements")
	}
	test := array[0] // Is there a way to do this safely?
	return test, nil
}

type OrganizationRequest struct {
	OrganizationName string `json:"organization_name"`
	IsPerson         bool   `json:"is_person"`
}

func OrganizationVerifyHandlerFactory(dbtx_val dbstore.DBTX) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			errorstring := fmt.Sprintf("Error reading request body: %v\n", err)
			fmt.Println(errorstring)
			http.Error(w, errorstring, http.StatusBadRequest)
			return
		}
		var req OrganizationRequest
		err = json.Unmarshal(bodyBytes, &req)
		if err != nil {
			errorstring := fmt.Sprintf("Error decoding JSON: %v\n", err)
			fmt.Println(errorstring)
			http.Error(w, errorstring, http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		q := *dbstore.New(dbtx_val)
		author_info := AuthorInformation{AuthorName: req.OrganizationName, IsPerson: req.IsPerson}
		author_info, err = verifyAuthorOrganizationUUID(ctx, q, &author_info)
		if err != nil {
			errorstring := fmt.Sprintf("Error verifying author %v: %v\n", req.OrganizationName, err)
			fmt.Println(errorstring)
			http.Error(w, errorstring, http.StatusInternalServerError)
			return
		}
		// No error handling since we always want it to retun a 200 at this point.
		response, _ := json.Marshal(author_info)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}

type ConversationRequest struct {
	DocketID string `json:"docket_id"`
}

func ConversationVerifyHandlerFactory(dbtx_val dbstore.DBTX) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			errorstring := fmt.Sprintf("Error reading request body: %v\n", err)
			fmt.Println(errorstring)
			http.Error(w, errorstring, http.StatusBadRequest)
			return
		}
		var req ConversationRequest
		err = json.Unmarshal(bodyBytes, &req)
		if err != nil {
			errorstring := fmt.Sprintf("Error decoding JSON: %v\n", err)
			fmt.Println(errorstring)
			http.Error(w, errorstring, http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		q := *dbstore.New(dbtx_val)
		conversation_info := ConversationInformation{DocketID: req.DocketID}
		conversation_info, err = verifyConversationUUID(ctx, q, &conversation_info)
		if err != nil {
			errorstring := fmt.Sprintf("Error verifying conversation %v: %v\n", req.DocketID, err)
			fmt.Println(errorstring)
			http.Error(w, errorstring, http.StatusInternalServerError)
			return
		}
		// No error handling since we always want it to retun a 200 at this point.
		response, _ := json.Marshal(conversation_info)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}

func verifyAuthorOrganizationUUID(ctx context.Context, q dbstore.Queries, author_info *AuthorInformation) (AuthorInformation, error) {
	if author_info.AuthorID != uuid.Nil {
		return *author_info, nil
	}
	// TODO: Change the sql so that this also matches IsPerson, but for now it shouldnt matter.
	results, err := q.OrganizationFetchByNameMatch(ctx, author_info.AuthorName)
	if err != nil {
		fmt.Println(err)
		return *author_info, err
	}
	org_return_info, err := getFirstElement(results) // Gives an error if list is empty
	if err == nil {
		org_uuid := org_return_info.ID
		author_info.AuthorID = org_uuid
		author_info.IsPerson = org_return_info.IsPerson.Bool
		return *author_info, nil
	}
	org_create_params := dbstore.CreateOrganizationParams{
		OrganizationAlias: author_info.AuthorName,
		Description:       "", // I should make this fixable at some point, but for now it will kinda work (tm)
		IsPerson:          pgtype.Bool{Bool: author_info.IsPerson, Valid: true},
	}
	pg_org_id, err := q.CreateOrganization(ctx, org_create_params)
	if err != nil {
		fmt.Println(err)
		return *author_info, err
	}
	org_uuid := pg_org_id
	author_info.AuthorID = org_uuid
	return *author_info, nil
}

func fileAuthorsUpsert(ctx context.Context, q dbstore.Queries, doc_uuid uuid.UUID, authors_info []AuthorInformation, insert bool) error {
	if !insert {
		err := q.AuthorshipDocumentDeleteAll(ctx, doc_uuid)
		if err != nil {
			return err
		}
	}
	fileAuthorInsert := func(author_info AuthorInformation) error {
		new_author_info, err := verifyAuthorOrganizationUUID(ctx, q, &author_info)
		if err != nil {
			return err
		}
		if new_author_info.AuthorID == uuid.Nil {
			return fmt.Errorf("ASSERT FAILURE: verifyAuthorOrganizationUUID should never return a null uuid.")
		}

		author_params := dbstore.AuthorshipDocumentOrganizationInsertParams{
			DocumentID:      doc_uuid,
			OrganizationID:  new_author_info.AuthorID,
			IsPrimaryAuthor: pgtype.Bool{Bool: new_author_info.IsPrimaryAuthor, Valid: true},
		}
		_, err = q.AuthorshipDocumentOrganizationInsert(ctx, author_params)
		if err != nil {
			return err
		}
		return nil
	}
	// Maybe m,ake async at some point, low priority since it isnt latency sensitive.
	for _, author_info := range authors_info {
		err := fileAuthorInsert(author_info)
		if err != nil {
			fmt.Printf("Encountered error while inserting author for document %s, ignoring and continuing: %v", doc_uuid, err)
		}
	}

	return nil
}

func verifyConversationUUID(ctx context.Context, q dbstore.Queries, conv_info *ConversationInformation) (ConversationInformation, error) {
	if conv_info.ID != uuid.Nil {
		return *conv_info, nil
	}

	// Try to find existing conversation for this docket
	// TODO: Change query to also match state if state exists
	results, err := q.DocketConversationFetchByDocketIdMatch(ctx, conv_info.DocketID)
	if err != nil {
		return *conv_info, err
	}

	// If conversation exists, return it
	if len(results) > 0 {
		conv := results[0]
		conv_info.ID = conv.ID
		conv_info.State = conv.State
		return *conv_info, nil
	}

	// Create new conversation if none exists
	create_params := dbstore.DocketConversationCreateParams{
		DocketID: conv_info.DocketID,
		State:    conv_info.State,
	}

	conv_id, err := q.DocketConversationCreate(ctx, create_params)
	if err != nil {
		return *conv_info, err
	}

	conv_info.ID = conv_id
	return *conv_info, nil
}

func fileConversationUpsert(ctx context.Context, q dbstore.Queries, file_id uuid.UUID, conv_info ConversationInformation, insert bool) error {
	// Sometimes this is getting called with an insert when the metadata already exists in the table, this causes a PGERROR, since it violates uniqueness. However, setting it up so it tries to update will fall back to insert if the file doesnt exist. Its probably a good idea to remove this and debug what is causing the new file thing at some point.
	insert = false
	shouldnt_process := conv_info.ID == uuid.Nil && conv_info.DocketID == ""
	if shouldnt_process {
		return nil
	}
	new_conv_info, err := verifyConversationUUID(ctx, q, &conv_info)
	if err != nil {
		return err
	}
	if new_conv_info.ID == uuid.Nil {
		return fmt.Errorf("ASSERT FAILURE: verifyConversationUUID should never return a null uuid")
	}
	if !insert {
		args := dbstore.DocketDocumentUpdateParams{
			DocketID: new_conv_info.ID,
			FileID:   file_id,
		}
		_, err = q.DocketDocumentUpdate(ctx, args)
		// If encounter a not found error, break error handling control flow and inserting object
		if err == nil {
			return nil
		}
		if err.Error() != "no rows in result set" {
			// If the error is nil, this still returns the error
			return err
		}
	}

	insert_params := dbstore.DocketDocumentInsertParams{
		DocketID: new_conv_info.ID,
		FileID:   file_id,
	}
	_, err = q.DocketDocumentInsert(ctx, insert_params)
	return err
}
