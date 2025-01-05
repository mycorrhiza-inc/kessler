func getObjectInformation(obj_uuid uuid.UUID, obj_type string, q dbstore.Queries) (map[string]interface{}, error) {
  switch obj_type {
  case "file":
    return getFileInformation(obj_uuid, q)
  case "org":
	return getOrgInformation(obj_uuid, q)
	case "docket":
	return getDocketInformation(obj_uuid, q)
}
