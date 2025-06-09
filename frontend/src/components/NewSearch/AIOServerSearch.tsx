import { Suspense } from "react";
import { StandardSearchBarClientBaseUrl } from "./HomeSearch";
import LoadingSpinner from "../styled-components/LoadingSpinner";
import SearchResultsServerStandalone from "../Search/SearchResultsServerStandalone";
import { GenericSearchInfo, GenericSearchType, searchInvoke } from "@/lib/adapters/genericSearchCallback";
import AIOClientSearchComponent from "./AIOClientSearch";
import { SearchResult } from "@/lib/types/new_search_types";
import { BackendFilterObject } from "@/lib/filters";
import RawSearchResults from "../Search/RawSearchResults";
import ErrorMessage from "../messages/ErrorMessage";


interface AIOServerProps {
  initialQuery: string;
  initialFilters: BackendFilterObject;
  searchType: GenericSearchType;
}
const AIOServerSearch = async ({ initialQuery, initialFilters, searchType }: AIOServerProps) => {

  // Fetch two pages worth of data server-side
  const initialPages = 2;
  const PAGE_SIZE = 40;
  const intiialPagination = { limit: PAGE_SIZE * initialPages, page: 0 };

  const searchInfo: GenericSearchInfo = {
    search_type: searchType,
    query: initialQuery,
    filters: initialFilters,
  }

  try {
    const initialResults: SearchResult[] = await searchInvoke(
      searchInfo,
      intiialPagination,
    );

    return (
      <Suspense
        fallback={<LoadingSpinner loadingText="Getting Results From Server" />}
      >
        <AIOClientSearchComponent
          initialQuery={initialQuery}
          initialFilters={initialFilters}
          searchType={searchType}
        >
          <RawSearchResults data={initialResults} />
        </AIOClientSearchComponent>
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
export default AIOServerSearch;
