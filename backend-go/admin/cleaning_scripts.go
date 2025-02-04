package admin

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"kessler/gen/dbstore"
	"kessler/util"
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
			fmt.Printf("Deduplicated org name: %s\n", orgname)
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

	errChan := make(chan error, len(all_orgs))
	var wg sync.WaitGroup

	for _, org := range all_orgs {
		wg.Add(1)
		go func(org dbstore.Organization) {
			defer wg.Done()

			orgname := org.Name
			orgNameTrimmed := strings.Trim(strings.TrimSpace(orgname), ".,!?:;\"'()[]{}")
			orgNameTrimmed = strings.TrimSpace(orgNameTrimmed)

			arg := dbstore.OrganizationAliasIdNameGetParams{
				OrganizationID:    org.ID,
				OrganizationAlias: orgNameTrimmed,
			}

			org_matched_aliases, err := q.OrganizationAliasIdNameGet(ctx, arg)
			if err != nil {
				errChan <- err
				return
			}

			if len(org_matched_aliases) == 0 {
				org_alias := dbstore.AliasOrganizationCreateParams{
					OrganizationID:    org.ID,
					OrganizationAlias: orgNameTrimmed,
				}
				_, err := q.AliasOrganizationCreate(ctx, org_alias)
				if err != nil {
					errChan <- err
					return
				}
				fmt.Printf("Ensured Aliases for org : %s\n", orgname)
			}
		}(org)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(errChan)

	// Check for any errors
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func completeCleanDatabaseHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Starting complete clean of database\n")
	ctx := context.Background()
	q := *util.DBQueriesFromRequest(r)

	err := deduplicateOrganizationsOnNames(ctx, q)
	if err != nil {
		errorstring := fmt.Sprintf("Error deduping orgs: %v\n", err)
		fmt.Println(errorstring)
		http.Error(w, errorstring, http.StatusInternalServerError)
		return
	}
	fmt.Printf("Sucessfully deduplicated orgs\n")
	err = organizationsNameAsAlias(ctx, q)
	if err != nil {
		errorstring := fmt.Sprintf("Error ensuring organization aliases: %v\n", err)
		fmt.Println(errorstring)
		http.Error(w, errorstring, http.StatusInternalServerError)
		return
	}
	fmt.Printf("Sucessfully ensured org aliases\n")
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Successfully cleaned database"))
}
