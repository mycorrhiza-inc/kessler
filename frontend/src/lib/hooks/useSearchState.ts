import { useState, useEffect, useMemo } from "react";
import {
  SearchResultsGetter,
  nilSearchResultsGetter,
} from "../types/new_search_types";
import { generateFakeResults } from "../search/search_utils";
import { Filters, useFilterState } from "../types/new_filter_types";

interface SearchStateExport {
  searchQuery: string;
  setSearchQuery: (query: string) => void;
  filters: Filters;
  setFilter: any;
  deleteFilter: any;
  isSearching: boolean;
  getResultsCallback: SearchResultsGetter;
  resetToInitial: () => void;
  triggerSearch: () => void;
  reloadOnChange: number;
}

export const useSearchState = (): SearchStateExport => {
  const [searchQuery, setSearchQuery] = useState("");
  const { filters, setFilter, deleteFilter } = useFilterState([]);
  const [isSearching, setIsSearching] = useState(false);
  const trimmedQuery = useMemo(() => searchQuery.trim(), [searchQuery]);
  const getResultsCallback = useMemo(
    () => generateSearchFunctions({ query: trimmedQuery, filters: filters }),
    [trimmedQuery, filters],
  );
  const [reloadOnChange, setReloadOnChange] = useState(0);

  // Handle search submission
  const setSearchUrl = (trimmedQuery: string, filters: Filters) => {
    window.history.pushState(
      { search: trimmedQuery },
      "",
      `/search?text=${encodeURIComponent(trimmedQuery)}`,
    );
  };
  // UseEffect that Executes a Search and Updates all the associated state
  useEffect(() => {
    setSearchUrl(trimmedQuery, filters);
    setReloadOnChange((prev) => (prev + 1) % 1024);
  }, [trimmedQuery, filters]); // Dont use a dependancy for filters since they arent encoded yet.

  // Handle browser navigation
  useEffect(() => {
    const handlePopState = (e: PopStateEvent) => {
      if (e.state?.search) {
        setSearchQuery(e.state.search);
        setIsSearching(true);
      } else {
        setSearchQuery("");
        setIsSearching(false);
      }
    };

    window.addEventListener("popstate", handlePopState);
    return () => window.removeEventListener("popstate", handlePopState);
  }, []);

  const resetSearch = () => {
    setSearchQuery("");
    setIsSearching(false);
    window.history.pushState(null, "", window.location.pathname);
  };
  const triggerSearch = () => {};

  return {
    searchQuery,
    setSearchQuery,
    filters,
    setFilter,
    deleteFilter,
    isSearching,
    // handleSearch,
    getResultsCallback,
    resetToInitial: resetSearch,
    triggerSearch,
    reloadOnChange,
  };
};

const generateSearchFunctions = ({
  query,
  filters,
}: {
  query: string;
  filters: Filters;
}): SearchResultsGetter => {
  if (query == "") {
    return nilSearchResultsGetter;
  }
  return generateFakeResults;
};
