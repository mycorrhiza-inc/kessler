import React, { Suspense } from "react";
import HomeSearchBar, {
  HomeSearchBarClientBaseUrl,
} from "@/components/NewSearch/HomeSearch";
import SearchResultsServer from "@/components/Search/SearchResultsServer";
import LoadingSpinner from "@/components/styled-components/LoadingSpinner";

interface SearchPageProps {
  searchParams: { q?: string };
}

export default function Page({ searchParams }: SearchPageProps) {
  const initialQuery = (searchParams.q || "").trim();

  return (
    <>
      <div className="flex flex-col items-center justify-center bg-base-100 p-4">
        <HomeSearchBarClientBaseUrl
          baseUrl="/search"
          initialState={initialQuery}
        />
      </div>
      <Suspense
        fallback={
          <LoadingSpinner loadingText="Fetching results from server." />
        }
      >
        <SearchResultsServer q={initialQuery} />
      </Suspense>
    </>
  );
}

