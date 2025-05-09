"use client";
import Link from "next/link";
import { useState, useEffect } from "react";
import { motion, AnimatePresence } from "framer-motion";
import ConversationTableInfiniteScroll from "../LookupPages/ConvoLookup/ConversationTable";
import { ExperimentalChatModalClickDiv } from "../Chat/ChatModal";
import OrganizationTableInfiniteScroll from "../LookupPages/OrgLookup/OrganizationTable";
import FileSearchView from "../Search/FileSearch/FileSearchView";
import HomeSearchBar from "../NewSearch/HomeSearch";

// So I have this problem, where lots of pages are going to want to have this search bar be present on the page, as well as in other stuff like command-k popups
//
// The current system would be to implement this somewhat sketchy url redirection code that sets the url for search, and then sets it back once you do any other navigation or finish your search. This code has the potential to introduce an absolute ton of bugs, so standardizing that code so that it had a very standardized behavior across all impelmentations would be nice.
//
// Another potential thing that could work is having whatever manages the search state essentially just return an element of the search results, and then the page would be in chqarge of making sure everything was animated correctly.
//
// Any thoughts on how you would architect this?

// src/components/HomePage/HomePage.tsx
import { useSearchState } from "@/lib/hooks/useSearchState";
import { SearchResultsComponent } from "@/components/Search/SearchResults";
import {
  GenericSearchInfo,
  GenericSearchType,
} from "@/lib/adapters/genericSearchCallback";

export default function HomePage() {
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
        className="flex flex-col items-center justify-center bg-base-100 p-4"
        style={{ overflow: "hidden" }}
      >
        <HomeSearchBar setTriggeredQuery={setTriggeredQuery} />
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
            <div className="grid grid-cols-2 w-full z-1">
              <div>
                <Link
                  className="text-3xl font-bold hover:underline"
                  href="/dockets"
                >
                  Dockets
                </Link>
                <div className="max-h-[600px] overflow-x-hidden border-r pr-4">
                  <ConversationTableInfiniteScroll
                    truncate
                    lookup_data={{ query: "" }}
                  />
                </div>
              </div>
              <div className="z-1">
                <Link
                  className="text-3xl font-bold hover:underline mb-5 p-10"
                  href="/orgs"
                >
                  Organizations
                </Link>
                <div className="max-h-[600px] overflow-x-hidden pl-4">
                  <OrganizationTableInfiniteScroll />
                </div>
              </div>
            </div>
            <ExperimentalChatModalClickDiv
              className="btn btn-accent w-full"
              inheritedFilters={[]}
            >
              Unsure of what to do? Try chatting with the entire New York PUC
            </ExperimentalChatModalClickDiv>

            <h1 className=" text-2xl font-bold">Newest Docs</h1>
            <FileSearchView inheritedFilters={[]} />
          </motion.div>
        )}
      </AnimatePresence>

      <SearchResultsComponent
        isSearching={isSearching}
        searchInfo={searchInfo}
        reloadOnChange={searchState.searchTriggerIndicator}
      ></SearchResultsComponent>
    </>
  );
}
