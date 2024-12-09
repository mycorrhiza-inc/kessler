export enum FilterField {
  MatchName = "match_name",
  MatchSource = "match_source",
  MatchDoctype = "match_doctype",
  MatchDocketId = "match_docket_id",
  MatchDocumentClass = "match_file_class",
  MatchAuthor = "match_author",
  MatchBeforeDate = "match_before_date",
  MatchAfterDate = "match_after_date",
  MatchAuthorUUID = "match_author_uuid",
}

export type InheritedFilterValues = Array<{
  filter: FilterField;
  value: string;
}>;

export type QueryFilterFields = {
  [key in FilterField]: string;
};

export type QueryDataFile = {
  filters: QueryFilterFields;
  query: string;
  start_offset: number;
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
//   MatchDoctype = "match_doctype",
//   MatchDocketId = "match_docket_id",
//   MatchDocumentClass = "match_file_class",
//   MatchAuthor = "match_author",
//   MatchBeforeDate = "match_before_date",
//   MatchAfterDate = "match_after_date",
// }
export const emptyQueryOptions: QueryFilterFields = {
  [FilterField.MatchName]: "",
  [FilterField.MatchSource]: "",
  [FilterField.MatchDoctype]: "",
  [FilterField.MatchDocketId]: "",
  [FilterField.MatchDocumentClass]: "",
  [FilterField.MatchAuthor]: "",
  [FilterField.MatchBeforeDate]: "",
  [FilterField.MatchAfterDate]: "",
  [FilterField.MatchAuthorUUID]: "",
};
