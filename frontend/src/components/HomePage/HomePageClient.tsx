"use client";
import { motion, AnimatePresence } from "framer-motion";
import StandardSearchbar from "../NewSearch/HomeSearch";

import { useSearchState } from "@/lib/hooks/useSearchState";
import { SearchResultsHomepageComponent } from "@/components/Search/SearchResults";
import {
  GenericSearchInfo,
  GenericSearchType,
} from "@/lib/adapters/genericSearchCallback";
import SearchResultsClient from "../Search/SearchResultsClient";

export default function HomePageClient({
  serverRecentTables,
}: {
  serverRecentTables: React.ReactNode;
}) {
  const searchState = useSearchState();
  const searchInfo: GenericSearchInfo = {
    search_type: GenericSearchType.Filling,
    query: searchState.searchQuery,
    filters: searchState.filters,
  };

  const isSearching = searchState.isSearching;

  const setTriggeredQuery = (query: string) => {
    if (query.trim() != searchState.searchQuery) {
      searchState.triggerSearch({ query: query.trim() });
    }
  };

  return (
    <>
      <motion.div
        initial={{ height: "70vh" }}
        animate={{ height: isSearching ? "30vh" : "70vh" }}
        transition={{ duration: 0.5 }}
        className="flex flex-col items-center justify-center  bg-base-100 p-4"
        style={{ overflow: "visible" }}
      >
        <StandardSearchbar setTriggeredQuery={setTriggeredQuery} />
      </motion.div>

      <AnimatePresence mode="wait">
        {!isSearching && (
          <motion.div
            key="homepage-content"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            transition={{ duration: 0.5, ease: "easeInOut" }}
          >
            {serverRecentTables}
          </motion.div>
        )}
      </AnimatePresence>

      <SearchResultsHomepageComponent
        isSearching={isSearching}
        searchInfo={searchInfo}
        reloadOnChange={searchState.searchTriggerIndicator}
      ></SearchResultsHomepageComponent>
    </>
  );
}
