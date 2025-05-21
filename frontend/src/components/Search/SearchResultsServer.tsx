import React from "react";
import SearchResultsClient from "./SearchResultsClient";
import RawSearchResults from "./RawSearchResults";
import { SearchResult } from "@/lib/types/new_search_types";
import {
  GenericSearchInfo,
  GenericSearchType,
  searchInvoke,
} from "@/lib/adapters/genericSearchCallback";

interface SearchResultsServerProps {
  searchInfo: GenericSearchInfo;
}

/**
 * Server Component: Fetches initial results and renders the client component.
 */
export default async function SearchResultsServer({
  searchInfo,
}: SearchResultsServerProps) {
  // Fetch two pages worth of data server-side
  const initialPages = 2;
  const PAGE_SIZE = 40;
  const intiialPagination = { limit: PAGE_SIZE * initialPages, page: 0 };

  const initialResults: SearchResult[] = await searchInvoke(
    searchInfo,
    intiialPagination,
  );
  const reloadOnChange = 0;

  return (
    <SearchResultsClient
      reloadOnChange={reloadOnChange}
      searchInfo={searchInfo}
      initialData={initialResults}
    >
      <RawSearchResults data={initialResults} />
    </SearchResultsClient>
  );
}
