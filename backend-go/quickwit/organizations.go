package quickwit

import (
	"context"
	"encoding/json"
	"fmt"
	"kessler/gen/dbstore"
	"kessler/objects/networking"
	"kessler/objects/organizations"
	"log"
	"net/http"
)

type OrgSearchRequestData struct {
	Search OrganizationSearchSchema `json:"search"`
	Limit  int                      `json:"limit"`
	Offset int                      `json:"offset"`
}

type OrganizationSearchSchema struct {
	Query            string `json:"query"`
	MinFilesAuthored int64  `json:"min_files_authored"`
}

func HandleOrgSearch(w http.ResponseWriter, r *http.Request) {
	log.Println("Received organization search request")

	var orgSearch OrganizationSearchSchema
	err := json.NewDecoder(r.Body).Decode(&orgSearch)
	if err != nil {
		errorstring := fmt.Sprintf("Error decoding JSON: %v", err)
		log.Println(errorstring)
		http.Error(w, errorstring, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	pagination := networking.PaginationFromUrlParams(r)
	searchData := OrgSearchRequestData{
		Search: orgSearch,
		Limit:  int(pagination.Limit),
		Offset: int(pagination.Offset),
	}

	results, err := SearchOrganizations(searchData, r.Context())
	if err != nil {
		errorstring := fmt.Sprintf("Error searching organizations: %v", err)
		log.Println(errorstring)
		http.Error(w, errorstring, http.StatusInternalServerError)
		return
	}

	respString, err := json.Marshal(results)
	if err != nil {
		log.Println("Error marshaling response data")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respString)
}

func SearchOrganizations(searchData OrgSearchRequestData, ctx context.Context) ([]organizations.OrganizationQuickwitSchema, error) {
	baseQuery := fmt.Sprintf("(name:(%s))", searchData.Search.Query)
	if searchData.Search.MinFilesAuthored > 0 {
		baseQuery += fmt.Sprintf(" AND files_authored_count:>=%d", searchData.Search.MinFilesAuthored)
	}

	searchRequest := QuickwitSearchRequest{
		Query:       baseQuery,
		MaxHits:     searchData.Limit,
		StartOffset: searchData.Offset,
	}

	var searchResults []organizations.OrganizationQuickwitSchema
	err := SearchHitsQuickwitGeneric(&searchResults, searchRequest, NYOrganizationIndex)
	return searchResults, err
}

func IndexOrganizations(orgs []organizations.OrganizationQuickwitSchema, indexName string) error {
	if indexName == "" {
		indexName = NYOrganizationIndex
	}
	IngestIntoIndex(indexName, orgs, true)
	return nil
}

func ReindexAllOrganizations(ctx context.Context, q dbstore.Queries, indexName string) error {
	orgs, err := q.OrganizationCompleteQuickwitListGet(ctx)
	if err != nil {
		return err
	}

	quickwitOrgs := make([]organizations.OrganizationQuickwitSchema, len(orgs))
	for i, org := range orgs {
		// TODO: Cache PG query to get aliases and number of documents for each org.
		// complete_org, err := crud.OrgWithFilesGetByID(ctx, &q, org.ID)
		if err != nil {
			fmt.Printf("Error getting org with files: %v\n", err)
		}
		quickwitOrgs[i] = organizations.OrganizationQuickwitSchema{
			ID:                 org.ID,
			Name:               org.Name,
			Aliases:            org.OrganizationAliases,
			FilesAuthoredCount: int(org.TotalDocumentsAuthored),
		}
	}

	return IndexOrganizations(quickwitOrgs, indexName)
}

func CreateQuickwitOrganizationsIndex(indexName string) error {
	if indexName == "" {
		indexName = NYOrganizationIndex
	}
	requestData := QuickwitIndex{
		Version: "0.7",
		IndexID: indexName,
		DocMapping: DocMapping{
			Mode: "dynamic",
			DynamicMapping: DynamicMapping{
				Indexed:    true,
				Stored:     true,
				Tokenizer:  "default",
				Record:     "basic",
				ExpandDots: true,
				Fast:       true,
			},
			FieldMappings: []FieldMapping{
				{
					Name: "name",
					Type: "text",
					Fast: true,
				},
				{
					Name: "aliases",
					Type: "array<text>",
					Fast: true,
				},
				{
					Name: "id",
					Type: "text",
					Fast: true,
				},
				{
					Name:          "timestamp",
					Type:          "datetime",
					Fast:          true,
					InputFormats:  []string{"unix_timestamp"},
					FastPrecision: "seconds",
				},
			},
			TimestampField: "timestamp",
		},
		SearchSettings: SearchSettings{
			DefaultSearchFields: []string{"name"},
		},
		IndexingSettings: IndexingSettings{
			MergePolicy: MergePolicy{
				Type:           "limit_merge",
				MaxMergeOps:    3,
				MergeFactor:    10,
				MaxMergeFactor: 12,
			},
			Resources: Resources{
				MaxMergeWriteThroughput: "80mb",
			},
		},
		Retention: Retention{
			Period:   "10 years",
			Schedule: "yearly",
		},
	}
	err := CreateIndex(requestData)
	return err
}
