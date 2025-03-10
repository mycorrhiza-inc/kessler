package rag

import (
	"context"
	"encoding/json"
	"fmt"
	"kessler/crud"
	"kessler/internal/dbstore"
	"kessler/internal/llm_utils"
	"kessler/util"

	"github.com/google/uuid"
	openai "github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

type ObjectInfo struct {
	UUID        uuid.UUID              `json:"uuid"`
	Name        string                 `json:"name"`
	ObjectType  string                 `json:"object_type"`
	Description string                 `json:"description"`
	Extras      map[string]interface{} `json:"extras"`
}

var more_info_func_schema = openai.FunctionDefinition{
	Name:        "get_more_info",
	Description: "If you need more context or general information about a certain object, you can query the database using a uuid or name. It should return info like a summary, metadata, and ids of other related objects.",
	Parameters: jsonschema.Definition{
		Type: jsonschema.Object,
		Properties: map[string]jsonschema.Definition{
			"query": {
				Type:        jsonschema.String,
				Description: "Query the database using a uuid (or optionally a name/id). A UUID is preferred since multiple objects can share the same name, thus leading to errors.",
			},
			"object_type": {
				Type:        jsonschema.String,
				Description: "The type of the object you want to query, currently supports: file, organization, docket",
			},
		},
		Required: []string{"query", "object_type"},
	},
}

func more_info_func_call(ctx context.Context) llm_utils.FunctionCall {
	q := *util.DBQueriesFromContext(ctx)
	return llm_utils.FunctionCall{
		Schema: rag_query_func_schema,
		Func: func(query_json string) (llm_utils.ToolCallResults, error) {
			var queryData map[string]string
			err := json.Unmarshal([]byte(query_json), &queryData)
			if err != nil {
				return llm_utils.ToolCallResults{}, fmt.Errorf("error unmarshaling query_json: %v", err)
			}
			search_query, ok := queryData["query"]
			if !ok {
				return llm_utils.ToolCallResults{}, fmt.Errorf("query field is missing in query_json")
			}
			object_type, ok := queryData["type"]
			if !ok {
				return llm_utils.ToolCallResults{}, fmt.Errorf("query field is missing in query_json")
			}
			obj_info, err := getObjectInformation(search_query, object_type, q, ctx)
			if err != nil {
				return llm_utils.ToolCallResults{}, err
			}
			// TODO: I dont know if returning yaml instead of json is good, typically yaml has better performance,
			// but the entire interface is already in json, and randomly swapping to yaml isnt probably going to improve things.
			obj_info_json, err := json.Marshal(obj_info)
			if err != nil {
				return llm_utils.ToolCallResults{}, err
			}
			return llm_utils.ToolCallResults{
				Response: string(obj_info_json),
			}, nil
		},
	}
}

func getObjectInformation(obj_query_string string, obj_type string, q dbstore.Queries, ctx context.Context) (ObjectInfo, error) {
	switch obj_type {
	case "file":
		return getFileInformationUnknown(obj_query_string, q, ctx)
	case "organization":
		return getOrgInformationUnknown(obj_query_string, q, ctx)
	case "docket":
		return getDocketInformationUnkown(obj_query_string, q, ctx)
	}
	return ObjectInfo{}, fmt.Errorf("Unknown object type: %v", obj_type)
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
	returnInfo.Extras["parent_docket_goverment_id"] = file.Conversation.DocketGovID
	returnInfo.Extras["parent_docket_uuid"] = file.Conversation.ID

	return returnInfo, nil
}

func getDocketInformationUnkown(docket_query_string string, q dbstore.Queries, ctx context.Context) (ObjectInfo, error) {
	return_obj, err := crud.ConversationGetByUnknown(ctx, &q, docket_query_string)
	if err != nil {
		return ObjectInfo{}, err
	}
	return_info := ObjectInfo{
		UUID:        return_obj.ID,
		Name:        return_obj.DocketGovID,
		ObjectType:  "docket",
		Description: return_obj.Name,
	}

	return return_info, nil
}

func getOrgInformationUnknown(org_query_string string, q dbstore.Queries, ctx context.Context) (ObjectInfo, error) {
	return_obj, err := crud.OrgWithFilesGetByUnknown(ctx, &q, org_query_string)
	if err != nil {
		return ObjectInfo{}, err
	}
	return_info := ObjectInfo{
		UUID:        return_obj.ID,
		Name:        return_obj.Name,
		ObjectType:  "organization",
		Description: "",
	}
	return return_info, nil
}
