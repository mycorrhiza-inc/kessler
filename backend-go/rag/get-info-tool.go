package rag

import (
	"kessler/gen/dbstore"

	"github.com/google/uuid"
)

func getObjectInformation(obj_uuid uuid.UUID, obj_type string, q dbstore.Queries) (map[string]interface{}, error) {
	exampleInfo := map[string]interface{}{}
	switch obj_type {
	case "file":
		return exampleInfo, nil
	case "org":
		return exampleInfo, nil
	case "docket":
		return exampleInfo, nil
	}
	return exampleInfo, nil
}
