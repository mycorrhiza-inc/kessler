import {
  PaginationData,
  SearchResult,
  SearchResultsGetter,
} from "../types/new_search_types";
import { generateFakeResults } from "../search/search_utils";
import axios from "axios";
import { DocumentCardData } from "../types/generic_card_types";
import { encodeUrlParams, TypedUrlParams } from "../types/url_params";
import { DEFAULT_PAGE_SIZE } from "../constants";
import { getContextualAPIUrl } from "../env_variables";

export enum GenericSearchType {
  Filling = "filing",
  Organization = "organization",
  Docket = "docket",
  Dummy = "dummy",
}

export interface GenericSearchInfo {
  query: string;
  filters?: Record<string, string>;
}

export const searchInvokeFromUrlParams = async (urlParams: TypedUrlParams, inheritedFilters: Record<string, string>) => {
  const searchInfo: GenericSearchInfo = {
    query: urlParams.queryData.query || "",
    filters: urlParams.queryData.filters
  }
  const pagination: PaginationData = {
    page: urlParams.paginationData.page || 0,
    limit: urlParams.paginationData.limit || DEFAULT_PAGE_SIZE
  }
  const callback = createGenericSearchCallback(searchInfo);
  return await callback(pagination);
}

export const mutateIndexifySearchResults = (
  results: SearchResult[],
  pagination: PaginationData
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
  info: GenericSearchInfo
): SearchResultsGetter => {
  // debug and default to dummy search results for stylistic changes.
  // info.search_type = GenericSearchType.Filling as GenericSearchType;

  const api_url = getContextualAPIUrl();

  console.log("searching with api_url:", api_url);
  console.log("info:", info);

  return async (pagination: PaginationData): Promise<SearchResult> => {
    const urlParams: TypedUrlParams = {
      paginationData: pagination,
      queryData: {
        query: info.query,
        filters: info.filters,
      }
    }
    console.log("returning a filling async callback:")

    console.log("query data:", info);
    console.log("API URL:", api_url);
    const encodedQueryParams = encodeUrlParams(urlParams)
    const url = `${api_url}/search/${encodedQueryParams}`
    console.log("ENDPOINT: ", url)


    return performSearchRequest<SearchRequest, { data: any[] }, DocumentCardData>(
      url,
      'get',
      (raw_results): DocumentCardData[] => {
        return raw_results.data.map((raw): DocumentCardData => {
          // return DocumentCardDataValidator.parse(raw)
          console.log("data::: ", raw);
          return raw as DocumentCardData
        });
      },
    );
  };
};
