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
import { getContextualAPIURL } from "../env_variables/env_variables";
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

export const createGenericSearchCallback = (
  info: GenericSearchInfo,
): SearchResultsGetter => {
  const api_url = getContextualAPIURL();
  // set all search invocations to be dummy searches for now.
  info.search_type = GenericSearchType.Dummy as GenericSearchType;
  switch (info.search_type) {
    case GenericSearchType.Dummy:
      return async (pagination: PaginationData): Promise<SearchResult> => {
        const req_url = `${api_url}/v2/version_hash`;
        const response: any = await axios.get(req_url);
        // check error conditions
        if (response.status >= 400) {
          throw new Error(`Request failed with status code ${response.status}`);
        }
        if (!response.data) {
          throw new Error(
            `Search Data returned from backend url ${req_url} was undefined.`,
          );
        }
        let results = await generateFakeResults(pagination);
        mutateIndexifySearchResults(results, pagination);
        return results;
      };
    case GenericSearchType.Filling:
      return async (pagination: PaginationData): Promise<SearchResult> => {
        const paginationQueryString = queryStringFromPagination(pagination);
        const searchQuery = info.query;
        console.log("query data", searchQuery);
        const searchFilters = info.filters;
        console.log("SEARCH FILTERS DISABLED UNTIL MIRRI UPDATES THE DB");
        console.log("API URL:  ", api_url);
        assert(api_url, "Api url cannot be undefined");
        try {
          const req_url = `${api_url}/v2/search/file${paginationQueryString}`;
          const req_data = {
            query: searchQuery,
          };
          const response: any = await axios.post(req_url, req_data);
          // check error conditions
          if (response.status >= 400) {
            throw new Error(
              `Request failed with status code ${response.status}`,
            );
          }
          if (!response.data) {
            throw new Error(
              `Search Data returned from backend url ${req_url} and data ${JSON.stringify(req_data)} was undefined.`,
            );
          }
          if (
            response.data?.length === 0 ||
            typeof response.data === "string"
          ) {
            console.log("RESPONSE LENGTH IS ZERO, THIS SEEMS WEIRD");
            return [];
          }
          // console.log(
          //   `got ${response.data.length} raw results from server`,
          // );

          const filings = hydratedSearchResultsToFilings(response.data);
          console.log(`successfully got ${filings.length} search results`);
          console.log("getting data");
          // console.log(searchResults);
          const searchResults: DocumentCardData[] =
            filings.map(adaptFilingToCard);

          mutateIndexifySearchResults(searchResults, pagination);
          return searchResults;
        } catch (error) {
          console.log(error);
          throw error;
        }
      };

    case GenericSearchType.Organization:
      throw new Error("Not Implemented");
    case GenericSearchType.Docket:
      throw new Error("Not Implemented");
  }
  throw "Error, no type specified for generic search callback";
};
