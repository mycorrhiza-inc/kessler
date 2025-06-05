import { Filters } from "../types/new_filter_types";
import {
  PaginationData,
  SearchResult,
  SearchResultsGetter,
  queryStringFromPagination,
} from "../types/new_search_types";
import { generateFakeResults } from "../search/search_utils";
import axios from "axios";
import { hydratedSearchResultsToFilings } from "../requests/search";
import { adaptFilingToCard } from "./genericCardAdapters";
import { DocumentCardData } from "../types/generic_card_types";
import { getEnvConfig } from "../env_variables/env_variables";
import assert from "assert";

export enum GenericSearchType {
  Filling = "filing",
  Organization = "organization",
  Docket = "docket",
  Dummy = "dummy",
}

export interface GenericSearchInfo {
  search_type: GenericSearchType;
  query: string;
  filters?: Filters;
}

export const searchInvoke = async (
  info: GenericSearchInfo,
  pagination: PaginationData,
) => {
  const callback = createGenericSearchCallback(info);
  return await callback(pagination);
};

export const mutateIndexifySearchResults = (
  results: SearchResult[],
  pagination: PaginationData,
) => {
  const offset = pagination.limit * pagination.page;
  for (let index = 0; index < results.length; index++) {
    results[index].index = index + offset;
  }
};

export const isSearchOffsetsValid = (results: SearchResult): boolean => {
  try {
    for (let index = 0; index < results.length - 1; index++) {
      if (results[index].index + 1 != results[index + 1].index) {
        return false;
      }
    }
    return true;
  } catch {
    return false;
  }
};

interface SearchRequest {
  query: string;
}

interface SearchResponse {
  data: any[]; // Replace with actual filing data type
}

interface VersionHashResponse {
  // Define version hash response structure
  hash?: string;
  version?: string;
}

export const createGenericSearchCallback = (
  info: GenericSearchInfo,
): SearchResultsGetter => {
  const api_url = getEnvConfig().public_api_url;
  
  console.log("searching with api_url:", api_url);
  console.log("info:", info);
  console.log("search type:", info.search_type);
  
  if (!api_url) {
    throw new Error("API URL cannot be undefined");
  }

  switch (info.search_type) {
    case GenericSearchType.Dummy:
      return async (pagination: PaginationData): Promise<SearchResult> => {
        try {
          const url = `${api_url}/v2/version_hash`;
          const response = await axios.get<VersionHashResponse>(url);
          
          if (!response.data) {
            throw new Error(
              `Search data returned from backend URL ${url} was undefined`
            );
          }
          
          const results = await generateFakeResults(pagination);
          mutateIndexifySearchResults(results, pagination);
          return results;
          
        } catch (error) {
          if (axios.isAxiosError(error)) {
            console.error("Axios error in dummy search:", error.message);
            if (error.response) {
              console.error("Response status:", error.response.status);
              console.error("Response data:", error.response.data);
            }
          } else {
            console.error("Unexpected error in dummy search:", error);
          }
          throw error;
        }
      };

    case GenericSearchType.Filling:
      return async (pagination: PaginationData): Promise<SearchResult> => {
        const paginationQueryString = queryStringFromPagination(pagination);
        const { query: searchQuery, filters: searchFilters } = info;
        
        console.log("query data:", searchQuery);
        console.log("SEARCH FILTERS DISABLED UNTIL MIRRI UPDATES THE DB");
        console.log("API URL:", api_url);
        
        try {
          const url = `${api_url}/v2/search/file${paginationQueryString}`;
          const requestData: SearchRequest = { query: searchQuery };
          
          const response = await axios.post<SearchResponse>(url, requestData);
          
          if (!response.data) {
            throw new Error(
              `Search data returned from backend URL ${url} with data ${JSON.stringify(requestData)} was undefined`
            );
          }
          
          if (Array.isArray(response.data) && response.data.length === 0) {
            console.warn("Response length is zero - no results found");
            return [];
          }
          
          if (typeof response.data === "string") {
            console.warn("Received string response instead of expected data structure");
            return [];
          }
          
          const filings = hydratedSearchResultsToFilings(response.data);
          console.log(`Successfully got ${filings.length} search results`);
          console.log("Getting data");
          
          const searchResults: DocumentCardData[] = filings.map(adaptFilingToCard);
          
          mutateIndexifySearchResults(searchResults, pagination);
          return searchResults;
          
        } catch (error) {
          if (axios.isAxiosError(error)) {
            console.error("Axios error in filing search:", error.message);
            if (error.response) {
              console.error("Response status:", error.response.status);
              console.error("Response data:", error.response.data);
            }
          } else {
            console.error("Unexpected error in filing search:", error);
          }
          throw error;
        }
      };

    case GenericSearchType.Organization:
      return async (): Promise<SearchResult> => {
        throw new Error("Organization search not implemented");
      };

    case GenericSearchType.Docket:
      return async (): Promise<SearchResult> => {
        throw new Error("Docket search not implemented");
      };

    default:
      // Exhaustive check to ensure all enum values are handled
      const _exhaustive: never = info.search_type;
      throw new Error(`Unhandled search type: ${info.search_type}`);
  }
};