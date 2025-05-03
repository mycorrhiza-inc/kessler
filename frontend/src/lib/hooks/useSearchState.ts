import { useState, useEffect, useMemo } from "react";
import {
  SearchResultsGetter,
  nilSearchResultsGetter,
} from "../types/new_search_types";
import { generateFakeResults } from "../search/search_utils";
import { Filters, useFilterState } from "../types/new_filter_types";

import { usePathname } from "next/navigation";

interface SearchStateExport {
  searchQuery: string;
  setSearchQuery: (query: string) => void;
  filters: Filters;
  setFilter: any;
  deleteFilter: any;
  isSearching: boolean;
  getResultsCallback: SearchResultsGetter;
  resetToInitial: () => void;
  triggerSearch: (trigger?: TriggerSearchObject) => void;
  searchTriggerIndicator: number;
}

export interface TriggerSearchObject {
  query?: string;
  filters?: Filters;
}

export const useSearchState = (): SearchStateExport => {
  const [searchQuery, setSearchQuery] = useState("");
  const { filters, setFilter, deleteFilter, clearFilters, replaceFilters } =
    useFilterState([]);
  const [isSearching, setIsSearching] = useState(false);
  const trimmedQuery = useMemo(() => searchQuery.trim(), [searchQuery]);
  const [searchTriggerIndicator, setSearchTriggerIndicator] = useState(0);

  const triggerSearch = (trigger?: TriggerSearchObject) => {
    const exec_query =
      trigger?.query === "" ? "" : (trigger?.query ?? searchQuery);
    const exec_filters = trigger?.filters || filters;
    setSearchTriggerIndicator((prev) => (prev + 1) % 1024);
    setSearchUrl(exec_query, exec_filters);
    setIsSearching(true);
    if (trigger?.query || trigger?.query === "") {
      setSearchQuery(trigger.query);
    }
    if (trigger?.filters) {
      replaceFilters(trigger.filters);
    }
  };
  const getResultsCallback = useMemo(
    () => generateSearchFunctions({ query: trimmedQuery, filters: filters }),
    [searchTriggerIndicator],
  );

  // Handle search submission
  const setSearchUrl = (trimmedQuery: string, filters: Filters) => {
    const method =
      window.location.pathname === "/search" ? "replaceState" : "pushState";
    window.history[method](
      { search: trimmedQuery },
      "",
      `/search?query=${encodeURIComponent(trimmedQuery)}`,
    );
  };

  const pathname = usePathname();
  useEffect(() => {
    // Reset search when leaving search page
    if (pathname !== "/search") {
      resetSearchNoNav();
    }
  }, [pathname]);

  // Update the popstate handler to handle home navigation
  useEffect(() => {
    const handlePopState = (e: PopStateEvent) => {
      const isSearchPage = window.location.pathname === "/search";

      if (isSearchPage && e.state?.search) {
        triggerSearch();
      } else {
        resetSearchNoNav();
      }
    };

    window.addEventListener("popstate", handlePopState);
    return () => window.removeEventListener("popstate", handlePopState);
  }, []);

  const resetSearchNoNav = () => {
    setSearchQuery("");
    clearFilters();
    setIsSearching(false);
  };
  const originalPathname = useMemo(() => {
    return window.location.pathname;
  }, []);
  const resetSearch = () => {
    resetSearchNoNav();
    console.log("Resetting Search to None.");
    window.history.pushState(null, "", originalPathname);
  };

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
    searchTriggerIndicator,
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
