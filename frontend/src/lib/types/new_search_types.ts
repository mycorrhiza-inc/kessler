export type SearchResult = any;
export interface PaginationData {
  page: number;
  limit: number;
}

export type SearchResultsGetter = (
  data: PaginationData,
) => Promise<SearchResult[]>;

export const nilSearchResultsGetter: SearchResultsGetter = async (
  data: PaginationData,
) => [];
