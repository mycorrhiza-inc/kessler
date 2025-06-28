import axios from "axios";
import {
  PaginationData,
  SearchResult,
} from "../types/new_search_types";
import { DocumentCardData, DocumentCardDataValidator } from "../types/generic_card_types";
import { TypedUrlParams } from "../types/url_params";
import { DEFAULT_PAGE_SIZE } from "../constants";
import { getContextualAPIUrl } from "../env_variables";

export enum GenericSearchType {
  Filing = "filing",
  Organization = "organization",
  Docket = "docket",
  Dummy = "dummy",
}

export interface GenericSearchInfo {
  query: string;
  filters?: Record<string, string>;
}

// Request structure matching Go backend SearchRequest
interface SearchRequest {
  query: string;
  filters?: Record<string, string>;
  page: number;
  per_page: number;
  namespace?: string;
}

// Response structure matching Go backend SearchResponse  
interface SearchResponse {
  data: any[];
  total?: number;
  page?: number;
  per_page?: number;
}

// Main search function that works with the Go POST /search/ endpoint
export const searchWithUrlParams = async (
  urlParams: TypedUrlParams,
  inheritedFilters: Record<string, string> = {}
): Promise<DocumentCardData[]> => {

  // Handle the converstion from beginning at 1 indexing, to beginning at 0 indexing here. Backend should never worry about this
  const actual_page = (urlParams.paginationData.page || 1) - 1;
  const requestBody: SearchRequest = {
    query: urlParams.queryData.query || "",
    filters: { ...inheritedFilters, ...urlParams.queryData.filters },
    page: actual_page,
    per_page: urlParams.paginationData.limit || DEFAULT_PAGE_SIZE,
    namespace: ""
  };

  try {
    const response = await axios.post<SearchResponse>(
      `${getContextualAPIUrl()}/search/`,
      requestBody,
      {
        headers: {
          'Content-Type': 'application/json'
        }
      }
    );

    return response.data.data.map((item): DocumentCardData => {
      // Uncomment when validator is ready: return DocumentCardDataValidator.parse(item);
      return item as DocumentCardData;
    });
  } catch (error) {
    console.error("Search request failed:", error);
    throw error;
  }
};

// Alternative function for namespace-specific searches
export const searchNamespace = async (
  query: string,
  namespace: "conversations" | "organizations" | "",
  filters: Record<string, string> = {},
  pagination: PaginationData = { page: 0, limit: DEFAULT_PAGE_SIZE }
): Promise<DocumentCardData[]> => {
  const endpoint = namespace ? `/search/${namespace}` : '/search/all';

  const requestBody: SearchRequest = {
    query,
    filters,
    page: pagination.page || 0,
    per_page: pagination.limit || DEFAULT_PAGE_SIZE
  };

  try {
    const response = await axios.post<SearchResponse>(
      `${getContextualAPIUrl()}${endpoint}`,
      requestBody,
      {
        headers: {
          'Content-Type': 'application/json'
        }
      }
    );

    return response.data.data.map((item): DocumentCardData => {
      return item as DocumentCardData;
    });
  } catch (error) {
    console.error(`Search failed for namespace ${namespace}:`, error);
    throw error;
  }
};

// Simple GET search function (uses query parameters)
export const searchWithGet = async (
  query: string,
  filters: Record<string, string> = {},
  pagination: PaginationData = { page: 0, limit: DEFAULT_PAGE_SIZE },
  namespace?: string
): Promise<DocumentCardData[]> => {
  const params = new URLSearchParams({
    q: query,
    page: pagination.page?.toString() || '0',
    per_page: pagination.limit?.toString() || DEFAULT_PAGE_SIZE.toString(),
    ...filters
  });

  if (namespace) {
    params.set('namespace', namespace);
  }

  try {
    const response = await axios.get<SearchResponse>(
      `${getContextualAPIUrl()}/search/?${params.toString()}`
    );

    return response.data.data.map((item): DocumentCardData => {
      return item as DocumentCardData;
    });
  } catch (error) {
    console.error("GET search request failed:", error);
    throw error;
  }
};

// Utility functions
export const mutateSearchResultsWithIndex = (
  results: SearchResult[],
  pagination: PaginationData
): void => {
  const offset = pagination.limit * pagination.page;
  results.forEach((result, index) => {
    result.index = index + offset;
  });
};

export const validateSearchResultsOrder = (results: SearchResult[]): boolean => {
  if (results.length <= 1) return true;

  return results.every((result, index) =>
    index === 0 || result.index === results[index - 1].index + 1
  );
};

// Get available filters
export const getAvailableFilters = async (): Promise<any> => {
  try {
    const response = await axios.get(`${getContextualAPIUrl()}/search/filters`);
    return response.data;
  } catch (error) {
    console.error("Failed to get available filters:", error);
    throw error;
  }
};

// Health check
export const checkSearchHealth = async (): Promise<any> => {
  try {
    const response = await axios.get(`${getContextualAPIUrl()}/search/health`);
    return response.data;
  } catch (error) {
    console.error("Search health check failed:", error);
    throw error;
  }
};
