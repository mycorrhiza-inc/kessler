import React, { Suspense } from "react";
import SearchResultsClient from "./SearchResultsClient";
import RawSearchResults from "./RawSearchResults";
import { SearchResult } from "@/lib/types/new_search_types";
import {
  GenericSearchInfo,
  GenericSearchType,
  searchInvoke,
} from "@/lib/adapters/genericSearchCallback";
import ErrorMessage from "../messages/ErrorMessage";
import LoadingSpinner from "../styled-components/LoadingSpinner";

interface SearchResultsServerProps {
  searchInfo: GenericSearchInfo;
}

/**
 * Server Component: Fetches initial results and renders the client component.
 */

async function SearchResultsServerStandalone({
  searchInfo,
}: SearchResultsServerProps) {
  // Fetch two pages worth of data server-side
  const initialPages = 2;
  const PAGE_SIZE = 40;
  const intiialPagination = { limit: PAGE_SIZE * initialPages, page: 0 };

  try {
    const initialResults: SearchResult[] = await searchInvoke(
      searchInfo,
      intiialPagination,
    );

    return (
      <Suspense
        fallback={<LoadingSpinner loadingText="Getting Results From Server" />}
      >
        <SearchResultsClient
          reloadOnChange={0}
          searchInfo={searchInfo}
          initialData={initialResults}
        >
          <RawSearchResults data={initialResults} />
        </SearchResultsClient>
      </Suspense>
    );
  } catch (error) {
    return (
      <ErrorMessage
        message={`Error getting search results from the server: ${error}`}
        error={`Error getting search results from the server: ${error}`}
      />
    );
  }
}
export default SearchResultsServerStandalone;
