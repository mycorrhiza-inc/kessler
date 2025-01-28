package search

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

func ConstructGenericFilterQuery(values reflect.Value, types reflect.Type) string {
	var filterQuery string
	filters := []string{}

	fmt.Printf("values: %v\n", values)
	fmt.Printf("types: %v\n", types)

	// ===== iterate over metadata for filter =====
	for i := 0; i < types.NumField(); i++ {
		// get the field and value
		field := types.Field(i)
		value := values.Field(i)
		tag := field.Tag.Get("json")
		if strings.Contains(tag, ",omitempty") {
			tag = strings.Split(tag, ",")[0]
		}

		fmt.Printf("tag: %v\nfield: %v\nvalue: %v\n", tag, field, value)

		if tag == "fileuuid" {
			tag = "source_id"
		}
		s := fmt.Sprintf("metadata.%s:(%s)", tag, value)

		// exlude empty values
		if strings.Contains(s, "00000000-0000-0000-0000-000000000000") {
			continue
		}
		log.Printf("new filter: %s\n", s)
		filters = append(filters, s)
	}
	// concat all filters with AND clauses
	for _, f := range filters {
		filterQuery += fmt.Sprintf(" AND (%s)", f)
	}
	fmt.Printf("filter query: %s\n", filterQuery)
	return filterQuery
}
