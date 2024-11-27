package crud

import (
	"context"
	"strings"

	"github.com/mycorrhiza-inc/kessler/backend-go/gen/dbstore"
)

func deduplicateAllOrganizationsInitial(ctx context.Context, q dbstore.Queries, author_info *AuthorInformation) error {
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
