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
		if _, exists := orgMap[orgname]; exists {
			// Found a duplicate organization name
			err := q.ConversationDeduplicateCascade(ctx)
			if err != nil {
				return err
			}
		} else {
			orgMap[orgname] = org
		}
	}
	return nil
}
