package quickwit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"kessler/objects/timestamp"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
)

// Maybe this should go in its own module base class at some point to avoid recursive dependancies
var QuickwitURL = os.Getenv("QUICKWIT_ENDPOINT")

type QuickwitSearchRequest struct {
	Query         string `json:"query,omitempty"`
	SnippetFields string `json:"snippet_fields,omitempty"`
	MaxHits       int    `json:"max_hits"`
	StartOffset   int    `json:"start_offset"`
	SortBy        string `json:"sort_by,omitempty"`
}

func ConstructDateQuery(DateFrom timestamp.KesslerTime, DateTo timestamp.KesslerTime) string {
	// construct date query
	fromDate := "*"
	toDate := "*"
	log.Printf("building date from: %s\n", DateFrom)
	log.Printf("building date to: %s\n", DateTo)

	if !(DateFrom.IsZero()) {
		fromDate = DateFrom.String()
	}
	if !(DateTo.IsZero()) {
		fromDate = DateTo.String()
	}
	dateQuery := fmt.Sprintf("date_filed:[%s TO %s]", fromDate, toDate)
	return dateQuery
}

func ConstructDateTextQuery(DateFrom timestamp.KesslerTime, DateTo timestamp.KesslerTime, query string) string {
	if DateFrom.IsZero() && DateTo.IsZero() {
		dateQueryString := fmt.Sprintf("(text:(%s) OR name:(%s))", query, query)
		return dateQueryString
	}
	var dateQueryString string
	dateQuery := ConstructDateQuery(DateFrom, DateTo)
	if len(query) >= 0 {
		dateQueryString = fmt.Sprintf("((text:(%s) OR name:(%s)) AND %s)", query, query, dateQuery)
		return dateQueryString
		// queryString = fmt.Sprintf("((text:(%s) OR name:(%s)) AND verified:true AND %s)", query, query, dateQuery)
	}
	return dateQuery
}

func ConstructGenericFilterQuery(values reflect.Value, types reflect.Type, useQuotes bool) string {
	var filterQuery string
	filters := []string{}

	// fmt.Printf("values: %v\n", values)
	// fmt.Printf("types: %v\n", types)

	// ===== iterate over metadata for filter =====
	for i := 0; i < types.NumField(); i++ {
		// get the field and value
		field := types.Field(i)
		value := values.Field(i)
		tag := field.Tag.Get("json")
		if strings.Contains(tag, ",omitempty") {
			tag = strings.Split(tag, ",")[0]
		}

		// fmt.Printf("tag: %v\nfield: %v\nvalue: %v\n", tag, field, value)

		if tag == "fileuuid" {
			tag = "source_id"
		}
		s := fmt.Sprintf("metadata.%s:(%s)", tag, value)
		if useQuotes && !(value.IsZero()) {
			s = fmt.Sprintf("metadata.%s:(\"%s\")", tag, value)
		}

		// exlude empty values
		if strings.Contains(s, "00000000-0000-0000-0000-000000000000") {
			continue
		}
		// log.Printf("new filter: %s\n", s)
		filters = append(filters, s)
	}
	// concat all filters with AND clauses
	for _, f := range filters {
		filterQuery += fmt.Sprintf(" AND (%s)", f)
	}
	fmt.Printf("filter query: %s\n", filterQuery)
	return filterQuery
}

func PerformGenericQuickwitRequest(request QuickwitSearchRequest, search_index string) ([]byte, error) {
	jsonData, err := json.Marshal(request)
	// ===== submit request to quickwit =====
	// log.Printf("jsondata: \n%s", jsonData)
	if err != nil {
		log.Printf("Error Marshalling quickwit request: %s", err)
		return []byte{}, err
	}

	request_url := fmt.Sprintf("%s/api/v1/%s/search", QuickwitURL, search_index)
	resp, err := http.Post(
		request_url,
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	curlCmd := fmt.Sprintf("curl -X POST -H 'Content-Type: application/json' -d '%s' %s", string(jsonData), request_url)
	if err != nil {
		log.Printf("Error sending request to quickwit: %s\n", err)
		log.Printf("Replay with: %s\n", curlCmd)
		return []byte{}, err
	}

	defer resp.Body.Close()

	// ===== handle response =====
	if resp.StatusCode != http.StatusOK {
		log.Printf("Error: received status code %v, with body: %s", resp.StatusCode, resp.Body)
		log.Printf("Error: received status code %v", resp.StatusCode)
		log.Printf("Replay with: %s\n", curlCmd)
		return []byte{}, fmt.Errorf("received status code %v", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to read response body: %v", err)
	}
	return bodyBytes, nil
}

func SearchHitsQuickwitGeneric[V GenericQuickwitSearchSchema](return_hits *[]V, request QuickwitSearchRequest, search_index string) error {
	type QuickwitHit struct {
		Hits []V `json:"hits"`
	}
	results, err := PerformGenericQuickwitRequest(request, search_index)
	if err != nil {
		return err
	}
	var testReturnHits QuickwitHit
	err = json.Unmarshal(results, &testReturnHits)
	if err != nil {
		return err
	}
	*return_hits = testReturnHits.Hits

	return nil
}
