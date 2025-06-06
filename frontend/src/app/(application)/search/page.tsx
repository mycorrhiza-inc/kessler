import React, { Suspense } from "react";
import HomeSearchBar, {
  HomeSearchBarClientBaseUrl,
} from "@/components/NewSearch/HomeSearch";
import SearchResultsServerStandalone from "@/components/Search/SearchResultsServer";
import LoadingSpinner from "@/components/styled-components/LoadingSpinner";
import {
  GenericSearchInfo,
  GenericSearchType,
} from "@/lib/adapters/genericSearchCallback";
import AIOServerSearch from "@/components/NewSearch/AIOServerSearch";

interface SearchPageProps {
  searchParams: { q?: string };
}

export default function Page({ searchParams }: SearchPageProps) {
  const initialQuery = (searchParams.q || "").trim();

  return (
    <AIOServerSearch searchType={GenericSearchType.Filling} initialQuery={initialQuery} initialFilters={[]} />
  );
}
