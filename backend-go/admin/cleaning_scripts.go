package admin

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mycorrhiza-inc/kessler/backend-go/gen/dbstore"
)

func deduplicateOrganizationsOnNames(ctx context.Context, q dbstore.Queries) error {
	all_orgs, err := q.OrganizationList(ctx)
	if err != nil {
		return err
	}

	// Create a map to track organization names
	orgMap := make(map[string]dbstore.Organization)

	for _, org := range all_orgs {
		orgname := org.Name
		orgNameTrimmed := strings.TrimSpace(orgname)
		if orgNameTrimmed != orgname {
			org.Name = orgNameTrimmed
			args := dbstore.OrganizationUpdateParams{ID: org.ID, Name: orgNameTrimmed, IsPerson: org.IsPerson, Description: org.Description}
			q.OrganizationUpdate(ctx, args)
		}
		if old_org, exists := orgMap[orgname]; exists {
			// Found a duplicate organization name
			args := dbstore.OrganizationDeduplicateCascadeParams{
				ID:             org.ID,     // Org to be deleted
				OrganizationID: old_org.ID, // Org to be kept
			}

			err := q.OrganizationDeduplicateCascade(ctx, args)
			if err != nil {
				return err
			}
		} else {
			orgMap[orgname] = org
		}
	}
	return nil
}

func organizationsNameAsAlias(ctx context.Context, q dbstore.Queries) error {
	all_orgs, err := q.OrganizationList(ctx)
	if err != nil {
		return err
	}
	for _, org := range all_orgs {
		orgname := org.Name
		orgNameTrimmed := strings.TrimSpace(orgname)
		arg := dbstore.OrganizationAliasIdNameGetParams{OrganizationID: org.ID, OrganizationAlias: pgtype.Text{String: orgNameTrimmed, Valid: true}}
		org_matched_aliases, err := q.OrganizationAliasIdNameGet(ctx, arg)
		if err != nil {
			return err
		}
		if len(org_matched_aliases) == 0 {
			org_alias := dbstore.AliasOrganizationCreateParams{
				OrganizationID: org.ID,
				OrganizationAlias: pgtype.Text{
					String: orgNameTrimmed,
					Valid:  true,
				},
			}
			_, err := q.AliasOrganizationCreate(ctx, org_alias)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func completeCleanDatabaseFactory(dbtx_val dbstore.DBTX) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		q := *dbstore.New(dbtx_val)
		err := deduplicateOrganizationsOnNames(ctx, q)
		if err != nil {
			errorstring := fmt.Sprintf("Error deduping orgs: %v\n", err)
			fmt.Println(errorstring)
			http.Error(w, errorstring, http.StatusInternalServerError)
			return
		}
		err = organizationsNameAsAlias(ctx, q)
		if err != nil {
			errorstring := fmt.Sprintf("Error ensuring organization aliases: %v\n", err)
			fmt.Println(errorstring)
			http.Error(w, errorstring, http.StatusInternalServerError)
			return
		}
	}
}
