import { useState, useEffect } from "react";

export function useSearchState() {
  const [searchQuery, setSearchQuery] = useState("");
  const [isSearching, setIsSearching] = useState(false);
  const [searchResults, setSearchResults] = useState<any[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

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

  return {
    searchQuery,
    isSearching,
    searchResults,
    isLoading,
    error,
    handleSearch,
    setSearchResults,
    setIsLoading,
    setError,
    resetSearch: () => {
      setSearchQuery("");
      setIsSearching(false);
      setSearchResults([]);
      setError(null);
      window.history.pushState(null, "", window.location.pathname);
    },
  };
}
