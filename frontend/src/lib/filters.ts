export enum FilterField {
  MatchName = "match_name",
  MatchSource = "match_source",
  MatchExtension = "match_extension",
  MatchDocketId = "match_docket_id",
  MatchDocumentClass = "match_file_class",
  MatchAuthor = "match_author",
  MatchBeforeDate = "match_before_date",
  MatchAfterDate = "match_after_date",
  MatchDateRange = "match_date_range",
  MatchAuthorUUID = "match_author_uuids",
  MatchConversationUUID = "match_conversation_uuid",
  MatchFileUUID = "match_file_uuid",
}

export type InheritedFilterValues = Array<{
  filter: FilterField;
  value: any;
}>;

export type QueryFileFilterFields = {
  [key in FilterField]: string;
};

export type QueryDataFile = {
  query: string;
  filters: QueryFileFilterFields;
};

export const allFilterFields: FilterField[] = Object.values(FilterField);
export const CaseFilterFields: FilterField[] = [
  FilterField.MatchName,
  FilterField.MatchDocketId,
  FilterField.MatchAuthor,
  FilterField.MatchDocumentClass,
  FilterField.MatchBeforeDate,
  FilterField.MatchAfterDate,
];

// This seems redundant with the list of case filter fields only being referenced in the codebace, going ahead and commenting out for now.
// export enum CaseFilterField {
//   MatchName = "match_name",
//   MatchExtension = "match_extension",
//   MatchDocketId = "match_docket_id",
//   MatchDocumentClass = "match_file_class",
//   MatchAuthor = "match_author",
//   MatchBeforeDate = "match_before_date",
//   MatchAfterDate = "match_after_date",
// }
export const emptyQueryOptions: QueryFileFilterFields = {
  [FilterField.MatchName]: "",
  [FilterField.MatchSource]: "",
  [FilterField.MatchExtension]: "",
  [FilterField.MatchDocketId]: "",
  [FilterField.MatchDocumentClass]: "",
  [FilterField.MatchAuthor]: "",
  [FilterField.MatchBeforeDate]: "",
  [FilterField.MatchAfterDate]: "",
  [FilterField.MatchDateRange]: "",
  [FilterField.MatchAuthorUUID]: "",
  [FilterField.MatchConversationUUID]: "",
  [FilterField.MatchFileUUID]: "",
};

export const disableListFromInherited = (
  inheritedFilters: InheritedFilterValues
): FilterField[] => {
  return inheritedFilters.map((val) => {
    return val.filter;
  });
};

export const initialFiltersFromInherited = (
  inheritedFilters: InheritedFilterValues
): QueryFileFilterFields => {
  var initialFilters = emptyQueryOptions;
  inheritedFilters.map((val) => {
    initialFilters[val.filter] = val.value;
  });
  return initialFilters;
};

export const inheritedFiltersFromValues = (
  filters: QueryFileFilterFields
): InheritedFilterValues => {
  if (filters == null) {
    return [];
  }
  var inheritedFilters: InheritedFilterValues = [];
  for (const [key, value] of Object.entries(filters)) {
    if (value != "") {
      inheritedFilters.push({ filter: key as FilterField, value: value });
    }
  }
  return inheritedFilters;
};

export interface BackendFilterObject {
  metadata_filters: {
    name: string;
    author: string;
    docket_id: string;
    file_class: string;
    extension: string;
    source: string;
    date_from: string;
    date_to: string;
    date_range?: string;
  };
  uuid_filters: {
    author_uuids: string;
    conversation_uuid: string;
    file_uuid: string;
  };
}

export const backendFilterGenerate = (
  filters: QueryFileFilterFields
): BackendFilterObject => {
  let from = filters.match_after_date;
  let to = filters.match_before_date;
  if (filters.match_date_range !== undefined) {
    
  }

  const metadataFilters: BackendFilterObject["metadata_filters"] = {
    name: filters.match_name,
    author: filters.match_author,
    docket_id: filters.match_docket_id,
    file_class: filters.match_file_class,
    extension: filters.match_extension,
    source: filters.match_source,
    date_from: from,
    date_to: to,
  };

  const uuidFilters: BackendFilterObject["uuid_filters"] = {
    author_uuids: filters.match_author_uuids,
    conversation_uuid: filters.match_conversation_uuid,
    file_uuid: filters.match_file_uuid,
  };

  if (filters.match_author_uuids !== "") {
    // If filtering by author uuid, remove author name
    metadataFilters.author = "";
  }
  
  if (filters.match_conversation_uuid !== "") {
    metadataFilters.docket_id = "";
  }

  if (filters.match_file_uuid !== "") {
    // Since the filters only match files and not text, then set all the other filter values to "" and return only the file uuid
    metadataFilters.name = "";
    metadataFilters.author = "";
    metadataFilters.docket_id = "";
    metadataFilters.file_class = "";
    metadataFilters.extension = "";
    metadataFilters.source = "";
    metadataFilters.date_from = "";
    metadataFilters.date_to = "";
    metadataFilters.source = "";
    uuidFilters.author_uuids = "";
    uuidFilters.conversation_uuid = "";
  }
  return { metadata_filters: metadataFilters, uuid_filters: uuidFilters };
};
