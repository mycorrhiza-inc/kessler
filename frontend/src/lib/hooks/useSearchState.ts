import { useState, useEffect, useMemo } from "react";
import {
  SearchResultsGetter,
  nilSearchResultsGetter,
} from "../types/new_search_types";
import { generateFakeResults } from "../search/search_utils";

interface SearchStateExport {
  searchQuery: string;
  setSearchQuery: (query: string) => void;
  filters;
  isSearching: boolean;
  handleSearch: (query: string) => void;
  resetToInitial: () => void;
  triggerSearch: () => void;
}

export const useSearchState = ({
  executeSearchOnObjectUpdate,
}: {
  executeSearchOnObjectUpdate: any;
}): SearchStateExport => {
  const [searchQuery, setSearchQuery] = useState("");
  const [filters, setFilters] = useState<any[]>([]);
  const [isSearching, setIsSearching] = useState(false);
  const trimmedQuery = useMemo(() => searchQuery.trim(), [searchQuery]);

  // Handle search submission
  const handleSearch = (query: string) => {
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

  return {
    searchQuery,
    setSearchQuery,
    isSearching,
    handleSearch,
    resetSearch,
  };
};

const generateSearchFunctions = ({
  query,
}: {
  query: string;
}): SearchResultsGetter => {
  if (query == "") {
    return nilSearchResultsGetter;
  }
  return generateFakeResults;
};
