import { CardData } from "./generic_card_types";

export type SearchResult = any;
export interface PaginationData {
  page: number;
  limit: number;
}

interface PaginationDataRaw {
  offset: number;
  limit: number;
}

const queryStringFromPaginationRaw = (pagination: PaginationDataRaw) => {
  return `?limit=${pagination.limit}&offset=${pagination.offset}`;
};

export const queryStringFromPagination = (pagination: PaginationData) => {
  return queryStringFromPaginationRaw({
    offset: pagination.page * pagination.limit,
    limit: pagination.limit,
  });
};
export type SearchResultsGetter = (
  data: PaginationData,
) => Promise<SearchResult[]>;

export const nilSearchResultsGetter: SearchResultsGetter = async (
  data: PaginationData,
) => [];
