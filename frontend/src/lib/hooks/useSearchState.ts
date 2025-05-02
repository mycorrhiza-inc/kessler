import { useState, useEffect, useMemo } from "react";
import {
  PaginationData,
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
  handleSearch: (query: string) => void;
  resetToInitial: () => void;
  triggerSearch: () => void;
}

export const useSearchState = ({}: {}): SearchStateExport => {
  const [searchQuery, setSearchQuery] = useState("");
  const { filters, setFilter, deleteFilter } = useFilterState([]);
  const [isSearching, setIsSearching] = useState(false);
  const trimmedQuery = useMemo(() => searchQuery.trim(), [searchQuery]);
  const getResultsCallback = useMemo(() => {
    return async (pagination: PaginationData) => {
      const results = await generateFakeResults(pagination);
      return results;
    };
  }, [trimmedQuery, filters]);

  // Handle search submission
  const initializeSearch = (query: string, filters: Filters) => {
    if (query.trim()) {
      setSearchQuery(query);
      setIsSearching(true);
      window.history.pushState(
        { search: query },
        "",
        `/search?text=${encodeURIComponent(query)}`,
      );
    }
  };

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
    handleSearch,
    getResultsCallback,
    resetToInitial: resetSearch,
    triggerSearch,
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
