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
	test := array[0]
	return test, nil
}

func verifyAuthorOrganizationUUID(ctx context.Context, q dbstore.Queries, author_info *AuthorInformation) (AuthorInformation, error) {
	empty_uuid, _ := uuid.Parse("00000000-0000-0000-0000-000000000000")
	if author_info.AuthorID != empty_uuid {
		return *author_info, nil
	}
	results, err := q.OrganizationFetchByNameMatch(ctx, author_info.AuthorName)
	if err != nil {
		fmt.Println(err)
		return *author_info, err
	}
	org_return_info, err := getFirstElement(results)
	if err == nil {
		org_uuid := uuid.UUID(org_return_info.ID.Bytes)
		author_info.AuthorID = org_uuid
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
