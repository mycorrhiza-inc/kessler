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

func getObjectInformation(obj_query_string string, obj_type string, q dbstore.Queries, ctx context.Context) (ObjectInfo, error) {
	switch obj_type {
	case "file":
		return getFileInformationUnknown(obj_query_string, q, ctx)
	case "org":
		return getOrgInformationUnknown(obj_query_string, q, ctx)
	case "docket":
		return getDocketInformationUnkown(obj_query_string, q, ctx)
	}
}

func getFileInformationUnknown(file_query_string string, q dbstore.Queries, ctx context.Context) (ObjectInfo, error) {
	file_uuid, err := uuid.Parse(file_query_string)
	if err == nil {
		return getFileInformationUUID(file_uuid, q, ctx)
	}
	return ObjectInfo{}, fmt.Errorf("UUID Failed to parse and Named File Lookup not implemented.")
}

func getFileInformationUUID(file_uuid uuid.UUID, q dbstore.Queries, ctx context.Context) (ObjectInfo, error) {
	file, err := crud.SemiCompleteFileGetFromUUID(ctx, q, file_uuid)
	if err != nil {
		return ObjectInfo{}, err
	}
	returnInfo := ObjectInfo{
		UUID:        file_uuid,
		ObjectType:  "file",
		Name:        file.Name,
		Description: file.Extra.Summary,
	}
	returnInfo.Extras["date"] = file.Mdata["date"]
	returnInfo.Extras["file_extension"] = file.Extension
	returnInfo.Extras["parent_docket_name"] = file.Conversation.Name
	returnInfo.Extras["parent_docket_goverment_id"] = file.Conversation.DocketID
	returnInfo.Extras["parent_docket_uuid"] = file.Conversation.ID

	return returnInfo, nil
}

func getDocketInformationUnkown(docket_query_string string, q dbstore.Queries, ctx context.Context) (ObjectInfo, error) {
	return_obj, err := crud.ConversationGetByUnknown(ctx, &q, docket_query_string)
	if err != nil {
		return ObjectInfo{}, err
	}
	return_info := ObjectInfo{
		Name:        return_obj.DocketID,
		ObjectType:  "docket",
		Description: return_obj.Name,
	}

	return return_info, nil
}

func getOrgInformationUnknown(org_query_string uuid.UUID, q dbstore.Queries, ctx context.Context) (ObjectInfo, error) {
	return_org, err := crud.OrganizationGetByID(ctx, &q, org_query_string)
	exampleInfo := ObjectInfo{UUID: org_uuid, ObjectType: "org", Name: "Example Organization", Description: "Example Organization Description"}
	return exampleInfo, nil
}
