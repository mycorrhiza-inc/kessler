import React from "react";
import SearchResultsClient from "./SearchResultsClient";
import RawSearchResults from "./RawSearchResults";
import { generateFakeResults } from "@/lib/search/search_utils";
import { SearchResult } from "@/lib/types/new_search_types";
import {
  GenericSearchInfo,
  GenericSearchType,
  searchInvoke,
} from "@/lib/adapters/genericSearchCallback";

interface SearchResultsServerProps {
  q: string;
  filters?: any;
}

/**
 * Server Component: Fetches initial results and renders the client component.
 */
export default async function SearchResultsServer({
  q,
  filters,
}: SearchResultsServerProps) {
  // Fetch two pages worth of data server-side
  const initialPages = 2;
  const PAGE_SIZE = 40;
  const intiialPagination = { limit: PAGE_SIZE * initialPages, page: 0 };
  const searchCallbackInfo: GenericSearchInfo = {
    search_type: GenericSearchType.Filling,
    query: q,
    filters: filters,
  };

  const initialResults: SearchResult[] = await searchInvoke(
    searchCallbackInfo,
    intiialPagination,
  );
  const reloadOnChange = 0;

  return (
    <SearchResultsClient
      reloadOnChange={reloadOnChange}
      genericSearchInfo={searchCallbackInfo}
      initialData={initialResults}
      initialPage={2}
    >
      <RawSearchResults data={initialResults} />
    </SearchResultsClient>
  );
}
