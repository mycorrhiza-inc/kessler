package rag

import (
	"context"
	"fmt"
	"kessler/crud"
	"kessler/gen/dbstore"

	"github.com/google/uuid"
)

type ObjectInfo struct {
	UUID        uuid.UUID              `json:"uuid"`
	Name        string                 `json:"name"`
	ObjectType  string                 `json:"object_type"`
	Description string                 `json:"description"`
	Extras      map[string]interface{} `json:"extras"`
}

func getObjectInformation(obj_uuid_string string, obj_named_lookup string, obj_type string, q dbstore.Queries, ctx context.Context) (ObjectInfo, error) {
	obj_uuid := uuid.Nil
	err := error(nil)
	if obj_uuid_string == "" {
		if obj_named_lookup == "" {
			return ObjectInfo{}, fmt.Errorf("Object UUID and name lookup were empty")
		}
		return ObjectInfo{}, fmt.Errorf("Named lookup not implemented yet")

	} else {
		obj_uuid, err = uuid.Parse(obj_uuid_string)
		if err != nil {
			return ObjectInfo{}, err
		}

	}
	return getObjectInformationUUID(obj_uuid, obj_type, q, ctx)
}

func getObjectInformationUUID(obj_uuid uuid.UUID, obj_type string, q dbstore.Queries, ctx context.Context) (ObjectInfo, error) {
	if obj_uuid == uuid.Nil {
		return ObjectInfo{}, fmt.Errorf("Object UUID was nil")
	}
	exampleInfo := ObjectInfo{UUID: obj_uuid, ObjectType: obj_type, Name: "Example", Description: "Example Description"}
	switch obj_type {
	case "file":
		return getFileInformationUUID(obj_uuid, q, ctx)
	case "org":
		return getOrgInformation(obj_uuid, q, ctx)
	case "docket":
		return getDocketInformation(obj_uuid, q, ctx)
	}
	return exampleInfo, fmt.Errorf("Object type %s did not match anything known", obj_type)
}

func getFileInformationUUID(file_uuid uuid.UUID, q dbstore.Queries, ctx context.Context) (ObjectInfo, error) {
	returnInfo := ObjectInfo{UUID: file_uuid, ObjectType: "file", Name: "Example File", Description: "Example File Description"}
	file, err := crud.SemiCompleteFileGetFromUUID(ctx, q, file_uuid)
	if err != nil {
		return ObjectInfo{}, err
	}
	returnInfo.Name = file.Name
	returnInfo.Description = file.Extra.Summary
	returnInfo.Extras["date"] = file.Mdata["date"]
	returnInfo.Extras["file_extension"] = file.Extension
	returnInfo.Extras["parent_docket_name"] = file.Conversation.Name
	returnInfo.Extras["parent_docket_goverment_id"] = file.Conversation.DocketID
	returnInfo.Extras["parent_docket_uuid"] = file.Conversation.ID

	return returnInfo, nil
}

func getDocketInformationNamed(conversation_uuid uuid.UUID, q dbstore.Queries, ctx context.Context) (ObjectInfo, error) {
	convo, err := crud.ConversationGetByName()
	exampleInfo := ObjectInfo{UUID: conversation_uuid, ObjectType: "docket", Name: "Example Docket", Description: "Example Docket Description"}
	return exampleInfo, nil
}

func getOrgInformation(org_uuid uuid.UUID, q dbstore.Queries, ctx context.Context) (ObjectInfo, error) {
	exampleInfo := ObjectInfo{UUID: org_uuid, ObjectType: "org", Name: "Example Organization", Description: "Example Organization Description"}
	return exampleInfo, nil
}
