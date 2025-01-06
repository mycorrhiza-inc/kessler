package rag

import (
	"kessler/gen/dbstore"

	"github.com/google/uuid"
)

type ObjectInfo struct {
	UUID        uuid.UUID `json:"uuid"`
	Name        string    `json:"name"`
	ObjectType  string    `json:"object_type"`
	Description string    `json:"description"`
}

func getObjectInformation(obj_uuid uuid.UUID, obj_type string, q dbstore.Queries) (ObjectInfo, error) {
	exampleInfo := ObjectInfo{UUID: obj_uuid, ObjectType: obj_type, Name: "Example", Description: "Example Description"}
	switch obj_type {
	case "file":
		return getFileInformation(obj_uuid, q)
	case "org":
		return getOrgInformation(obj_uuid, q)
	case "docket":
		return getDocketInformation(obj_uuid, q)
	}
	return exampleInfo, nil
}

func getFileInformation(file_uuid uuid.UUID, q dbstore.Queries) (ObjectInfo, error) {
	exampleInfo := ObjectInfo{UUID: file_uuid, ObjectType: "file", Name: "Example File", Description: "Example File Description"}
	return exampleInfo, nil
}

func getOrgInformation(org_uuid uuid.UUID, q dbstore.Queries) (ObjectInfo, error) {
	exampleInfo := ObjectInfo{UUID: org_uuid, ObjectType: "org", Name: "Example Organization", Description: "Example Organization Description"}
	return exampleInfo, nil
}

func getDocketInformation(docket_uuid uuid.UUID, q dbstore.Queries) (ObjectInfo, error) {
	exampleInfo := ObjectInfo{UUID: docket_uuid, ObjectType: "docket", Name: "Example Docket", Description: "Example Docket Description"}
	return exampleInfo, nil
}
