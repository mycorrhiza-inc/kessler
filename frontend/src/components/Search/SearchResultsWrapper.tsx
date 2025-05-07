import React from "react";
import SearchResults from "./SearchResults";
import { getSearchResults } from "@/lib/requests/search";
import { QueryDataFile } from "@/lib/types/new_search_types";

type Props = {
  initialQuery: string;
};

// Server Component: fetches initial page and streams
export default async function SearchResultsWrapper({ initialQuery }: Props) {
  const initialPage = 1;
  const pageSize = 40;

  // QueryDataFile: query + filters; here no filters
  const queryData: QueryDataFile = { query: initialQuery, filters: [] };
  const initialResults = await getSearchResults(
    queryData,
    initialPage - 1,
    pageSize,
  );

  return (
    <SearchResults
      initialResults={initialResults}
      initialQuery={initialQuery}
      initialPage={initialPage}
      pageSize={pageSize}
    />
  );
}

