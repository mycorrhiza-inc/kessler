package crud

import (
	"context"
	"fmt"

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

func verifyAuthorOrganizationUUID(ctx context.Context, q dbstore.Queries, author_info *AuthorInformation) (AuthorInformation, error) {
	empty_uuid, _ := uuid.Parse("00000000-0000-0000-0000-000000000000")
	if author_info.AuthorID != empty_uuid {
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
		org_uuid := uuid.UUID(org_return_info.ID.Bytes)
		author_info.AuthorID = org_uuid
		author_info.IsPerson = org_return_info.IsPerson.Bool
		return *author_info, nil
	}
	org_create_params := dbstore.CreateOrganizationParams{
		Name:        author_info.AuthorName,
		Description: pgtype.Text{String: "", Valid: true}, // I should make this fixable at some point, but for now it will kinda work (tm)
		IsPerson:    pgtype.Bool{Bool: author_info.IsPerson, Valid: true},
	}
	pg_org_id, err := q.CreateOrganization(ctx, org_create_params)
	if err != nil {
		fmt.Println(err)
		return *author_info, err
	}
	org_uuid := uuid.UUID(pg_org_id.Bytes)
	author_info.AuthorID = org_uuid
	return *author_info, nil
}

func fileAuthorsUpsert(ctx context.Context, q dbstore.Queries, doc_uuid uuid.UUID, authors_info []AuthorInformation, insert bool) error {
	if !insert {
		err := q.AuthorshipDocumentDeleteAll(ctx, pgtype.UUID{Bytes: doc_uuid, Valid: true})
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
			DocumentID:      pgtype.UUID{Bytes: doc_uuid, Valid: true},
			OrganizationID:  pgtype.UUID{Bytes: new_author_info.AuthorID, Valid: true},
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
		conv_info.ID = uuid.UUID(conv.ID.Bytes)
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

	conv_info.ID = uuid.UUID(conv_id.Bytes)
	return *conv_info, nil
}

func fileConversationUpsert(ctx context.Context, q dbstore.Queries, file_id uuid.UUID, conv_info ConversationInformation, insert bool) error {
	// Sometimes this is getting called with an insert when the metadata already exists in the table, this causes a PGERROR, since it violates uniqueness. However, setting it up so it tries to update will fall back to insert if the file doesnt exist. Its probably a good idea to remove this and debug what is causing the new file thing at some point.
	insert = false
	shouldnt_process := conv_info.ID == uuid.Nil && conv_info.DocketID == ""
	if shouldnt_process {
		return nil
	}
	if !insert {
		err := q.DocketConversationDelete(ctx, pgtype.UUID{Bytes: file_id, Valid: true})
		if err != nil {
			return err
		}
	}

	new_conv_info, err := verifyConversationUUID(ctx, q, &conv_info)
	if err != nil {
		return err
	}

	if new_conv_info.ID == uuid.Nil {
		return fmt.Errorf("ASSERT FAILURE: verifyConversationUUID should never return a null uuid")
	}

	insert_params := dbstore.DocketDocumentInsertParams{
		DocketID: pgtype.UUID{Bytes: new_conv_info.ID, Valid: true},
		FileID:   pgtype.UUID{Bytes: file_id, Valid: true},
	}
	_, err = q.DocketDocumentInsert(ctx, insert_params)
	return err
}
