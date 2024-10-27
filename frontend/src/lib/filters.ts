export enum FilterField {
  MatchName = "match_name",
  MatchSource = "match_source",
  MatchDoctype = "match_doctype",
  MatchDocketId = "match_docket_id",
  MatchDocumentClass = "match_document_class",
  MatchAuthor = "match_author",
  MatchBeforeDate = "match_before_date",
  MatchAfterDate = "match_after_date",
}
export type QueryFilterFields = {
  [key in FilterField]: string;
};
export const allFilterFields: FilterField[] = Object.values(FilterField);
export const emptyQueryOptions: QueryFilterFields = {
  [FilterField.MatchName]: "",
  [FilterField.MatchSource]: "",
  [FilterField.MatchDoctype]: "",
  [FilterField.MatchDocketId]: "",
  [FilterField.MatchDocumentClass]: "",
  [FilterField.MatchAuthor]: "",
  [FilterField.MatchBeforeDate]: "",
  [FilterField.MatchAfterDate]: "",
};
