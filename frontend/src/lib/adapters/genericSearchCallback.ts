import { Filters } from "../types/new_filter_types";
import {
  PaginationData,
  SearchResult,
  SearchResultsGetter,
  queryStringFromPagination,
} from "../types/new_search_types";
import { generateFakeResults } from "../search/search_utils";
import axios from "axios";
import { Filing } from "../types/FilingTypes";
import { hydratedSearchResultsToFilings } from "../requests/search";
import { getRuntimeEnv } from "../env_variables_hydration_script";
import { adaptFilingToCard } from "./genericCardAdapters";
import { DocumentCardData } from "../types/generic_card_types";

export enum GenericSearchType {
  Filling = "filing",
  Organization = "organization",
  Docket = "docket",
  Dummy = "dummy",
}

export interface GenericSearchInfo {
  search_type: GenericSearchType;
  query: string;
  filters: Filters;
}

export const searchInvoke = async (
  info: GenericSearchInfo,
  pagination: PaginationData,
) => {
  const callback = createGenericSearchCallback(info);
  return await callback(pagination);
};

export const createGenericSearchCallback = (
  info: GenericSearchInfo,
): SearchResultsGetter => {
  const runtimeConfig = getRuntimeEnv();
  const api_url = runtimeConfig.public_api_url;
  switch (info.search_type) {
    case GenericSearchType.Dummy:
      return generateFakeResults;
    case GenericSearchType.Filling:
      return async (pagination: PaginationData): Promise<SearchResult> => {
        const paginationQueryString = queryStringFromPagination(pagination);
        const searchQuery = info.query;
        console.log("query data", searchQuery);
        const searchFilters = info.filters;
        console.log("SEARCH FILTERS DISABLED UNTIL MIRRI UPDATES THE DB");
        try {
          const filingResults: Filing[] = await axios
            .post(`${api_url}/v2/search/file${paginationQueryString}`, {
              query: searchQuery,
            })
            // check error conditions
            .then((response): Filing[] => {
              if (response.status >= 400) {
                throw new Error(
                  `Request failed with status code ${response.status}`,
                );
              }
              if (
                response.data?.length === 0 ||
                typeof response.data === "string"
              ) {
                return [];
              }

              const filings = hydratedSearchResultsToFilings(response.data);
              return filings;
            });
          console.log("getting data");
          // console.log(searchResults);
          const searchResults: DocumentCardData[] =
            filingResults.map(adaptFilingToCard);

          return searchResults;
        } catch (error) {
          console.log(error);
          throw error;
        }
      };

    case GenericSearchType.Organization:
      throw "Not Implemented";
    case GenericSearchType.Docket:
      throw "Not Implemented";
  }
  throw "Error, no type specified for generic search callback";
};
