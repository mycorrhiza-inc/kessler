"use client";
import HomeSearchBar from "@/components/NewSearch/HomeSearch";
import React, { useEffect, useState } from "react";

import { motion, AnimatePresence } from "framer-motion";
import { useSearchState } from "@/lib/hooks/useSearchState";
import { SearchResultsComponent } from "@/components/Search/SearchResults";
export default function Page() {
  const searchState = useSearchState();

  const isSearching = searchState.isSearching;

  const setTriggeredQuery = (query: string) => {
    if (query.trim() != searchState.searchQuery) {
      searchState.triggerSearch({ query: query.trim() });
    }
  };
  useEffect(() => {
    const urlParams = new URLSearchParams(window.location.search);
    const initialQuery = urlParams.get("q");
    if (initialQuery) {
      searchState.triggerSearch({ query: initialQuery.trim() });
    }
  }, []);

  // Reset search state when URL changes (back/forward navigation)
  return (
    <>
      <motion.div
        initial={{ height: "70vh" }}
        animate={{ height: isSearching ? "30vh" : "70vh" }}
        transition={{ duration: 0.5, ease: "easeInOut" }}
        className="flex flex-col items-center justify-center bg-base-100 p-4"
        style={{ overflow: "hidden" }}
      >
        <HomeSearchBar setTriggeredQuery={setTriggeredQuery} />
      </motion.div>

      <div
        className={`transition-all duration-500 ease-in-out ${
          isSearching ? "opacity-100 translate-y-0" : "opacity-0 -translate-y-4"
        }`}
      >
        <SearchResultsComponent
          isSearching={isSearching}
          searchGetter={searchState.getResultsCallback}
          reloadOnChange={searchState.searchTriggerIndicator}
        ></SearchResultsComponent>
      </div>
    </>
  );
}
