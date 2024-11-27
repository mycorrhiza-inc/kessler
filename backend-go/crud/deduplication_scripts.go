package crud

import (
	"context"

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
		if old_org, exists := orgMap[orgname]; exists {
			// Found a duplicate organization name
			args := dbstore.OrganizationDeduplicateCascadeParams{
				ID:             old_org.ID, // Org to be deleted
				OrganizationID: org.ID,     // Org to be kept
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
