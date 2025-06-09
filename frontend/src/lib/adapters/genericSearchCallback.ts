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
import { contextualApiUrl, getEnvConfig } from "../env_variables/env_variables";
import assert from "assert";
import { BackendFilterObject } from "../filters";

export enum GenericSearchType {
  Filling = "filing",
  Organization = "organization",
  Docket = "docket",
  Dummy = "dummy",
}

export interface GenericSearchInfo {
  search_type: GenericSearchType;
  query: string;
  filters?: BackendFilterObject;
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

/**
 * Generic function to perform network search requests and handle response checks.
 */
async function performSearchRequest<Req, Res, Item extends SearchResult>(
  url: string,
  method: 'get' | 'post',
  transform: (raw: Res) => Item[],
  pagination?: PaginationData,
  requestData?: Req,
): Promise<Item[]> {
  try {
    const response =
      method === 'get'
        ? await axios.get<Res>(url)
        : await axios.post<Res>(url, requestData);

    if (!response.data) {
      throw new Error(
        `Search data returned from backend URL ${url} was undefined (requestData=${JSON.stringify(requestData)})`
      );
    }

    const data = response.data as unknown as Res;

    // Handle empty or unexpected data shapes
    if (Array.isArray(data) && (data as any[]).length === 0) {
      console.warn(`Response from ${url} is an empty array`);
      return [];
    }
    if (typeof data === 'string') {
      console.warn(`Response from ${url} is a string, expected structured data`);
      return [];
    }

    const items = transform(data);
    if (pagination) {
      mutateIndexifySearchResults(items, pagination);
    }
    return items;
  } catch (error) {
    if (axios.isAxiosError(error)) {
      console.error(`Axios error in search (${method.toUpperCase()} ${url}):`, error.message);
      if (error.response) {
        console.error('Response status:', error.response.status);
        console.error('Response data:', error.response.data);
      }
    } else {
      console.error(`Unexpected error in search (${method.toUpperCase()} ${url}):`, error);
    }
    throw error;
  }
}

export const createGenericSearchCallback = (
  info: GenericSearchInfo,
): SearchResultsGetter => {
  // debug and default to dummy search results for stylistic changes.
  info.search_type = GenericSearchType.Dummy as GenericSearchType;
  console.log("All searches are dummys for momentary testing purposes");

  const api_url = contextualApiUrl(getEnvConfig());
  if (!api_url) {
    throw new Error("API URL cannot be undefined");
  }
  console.log("searching with api_url:", api_url);
  console.log("info:", info);
  console.log("search type:", info.search_type);

  switch (info.search_type) {
    case GenericSearchType.Dummy:
      return async (pagination: PaginationData): Promise<SearchResult[]> => {
        // Dummy fetch just to simulate network call
        const url = `${api_url}/v2/version_hash`;
        await performSearchRequest<void, any, SearchResult>(url, 'get', () => []);
        const results = await generateFakeResults(pagination);
        console.log("Got fake results");
        mutateIndexifySearchResults(results, pagination);
        return results;
      };

    case GenericSearchType.Filling:
      return async (pagination: PaginationData): Promise<SearchResult[]> => {
        const paginationQueryString = queryStringFromPagination(pagination);
        const { query: searchQuery } = info;
        const url = `${api_url}/v2/search/file${paginationQueryString}`;
        const requestData: SearchRequest = { query: searchQuery };

        return performSearchRequest<SearchRequest, any[], DocumentCardData>(
          url,
          'post',
          (raw) => {
            const filings = hydratedSearchResultsToFilings(raw);
            console.log(`Successfully got ${filings.length} search results`);
            return filings.map(adaptFilingToCard);
          },
          pagination,
          requestData,
        );
      };

    case GenericSearchType.Organization:
      return async (): Promise<SearchResult[]> => {
        throw new Error("Organization search not implemented");
      };

    case GenericSearchType.Docket:
      return async (): Promise<SearchResult[]> => {
        throw new Error("Docket search not implemented");
      };

    default:
      const _exhaustive: never = info.search_type;
      throw new Error(`Unhandled search type: ${info.search_type}`);
  }
};
