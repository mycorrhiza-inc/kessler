package crud

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"kessler/common/objects/authors"
	"kessler/database"
	"kessler/gen/dbstore"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// Make this generic
// func getFirstElement(array []dbstore.Organization) (dbstore.Organization, error) {
// 	if len(array) == 0 {
// 		return dbstore.Organization{}, fmt.Errorf("List had no elements")
// 	}
// 	test := array[0] // Is there a way to do this safely?
// 	return test, nil
// }

type OrganizationRequest struct {
	OrganizationName string `json:"organization_name"`
	IsPerson         bool   `json:"is_person"`
}

func OrganizationVerifyHandler(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		errorstring := fmt.Sprintf("Error reading request body: %v\n", err)
		log.Info(errorstring)
		http.Error(w, errorstring, http.StatusBadRequest)
		return
	}
	var req OrganizationRequest
	err = json.Unmarshal(bodyBytes, &req)
	if err != nil {
		errorstring := fmt.Sprintf("Error decoding JSON: %v\n", err)
		log.Info(errorstring)
		http.Error(w, errorstring, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	q := *database.GetTx()

	author_info := authors.AuthorInformation{AuthorName: req.OrganizationName, IsPerson: req.IsPerson}
	author_info, err = verifyAuthorOrganizationUUID(ctx, q, &author_info)
	if err != nil {
		errorstring := fmt.Sprintf("Error verifying author %v: %v\n", req.OrganizationName, err)
		log.Info(errorstring)
		http.Error(w, errorstring, http.StatusInternalServerError)
		return
	}
	// No error handling since we always want it to retun a 200 at this point.
	response, _ := json.Marshal(author_info)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func verifyAuthorOrganizationUUID(ctx context.Context, q dbstore.Queries, author_info *authors.AuthorInformation) (authors.AuthorInformation, error) {
	if author_info.AuthorID != uuid.Nil {
		return *author_info, nil
	}
	// TODO: Change the sql so that this also matches IsPerson, but for now it shouldnt matter.
	org_return_info, err := q.OrganizationFetchByAliasMatchSingle(ctx, author_info.AuthorName)
	if err == nil {
		org_uuid := org_return_info.ID.Bytes
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
		log.Info(err)
		return *author_info, err
	}
	org_uuid := pg_org_id
	author_info.AuthorID = org_uuid
	return *author_info, nil
}

func fileAuthorsUpsert(ctx context.Context, q dbstore.Queries, doc_uuid uuid.UUID, authors_info []authors.AuthorInformation, insert bool) error {
	if !insert {
		err := q.AuthorshipDocumentDeleteAll(ctx, doc_uuid)
		if err != nil {
			return err
		}
	}
	fileAuthorInsert := func(author_info authors.AuthorInformation) error {
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
			log.Info(fmt.Sprintf("Encountered error while inserting author for document %s, ignoring and continuing: %v", doc_uuid, err))
		}
	}

	return nil
}
