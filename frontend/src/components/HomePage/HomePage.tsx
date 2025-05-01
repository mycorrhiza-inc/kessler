"use client";
import Link from "next/link";
import { useState, useEffect } from "react";
import { motion, AnimatePresence } from "framer-motion";
import ConversationTableInfiniteScroll from "../LookupPages/ConvoLookup/ConversationTable";
import { ExperimentalChatModalClickDiv } from "../Chat/ChatModal";
import OrganizationTableInfiniteScroll from "../LookupPages/OrgLookup/OrganizationTable";
import FileSearchView from "../Search/FileSearch/FileSearchView";
import HomeSearchBar from "../NewSearch/HomeSearch";
import DummyResults from "../NewSearch/DummySearchResults";
export default function HomePage() {
  const [searchQuery, setSearchQuery] = useState("");
  const [isSearching, setIsSearching] = useState(false);

  // Handle search submission
  const handleSearch = (query: string) => {
    if (query.trim()) {
      setSearchQuery(query);
      setIsSearching(true);

      // Update URL without navigation
      window.history.pushState(
        null,
        "",
        `/search?text=${encodeURIComponent(query)}`,
      );
    }
  };

  // Reset search state when URL changes (back/forward navigation)
  useEffect(() => {
    const handlePopState = () => {
      setIsSearching(false);
    };

    window.addEventListener("popstate", handlePopState);
    return () => window.removeEventListener("popstate", handlePopState);
  }, []);

  return (
    <>
      <motion.div
        initial={{ height: "70vh" }}
        animate={{ height: isSearching ? "30vh" : "70vh" }}
        transition={{ duration: 0.5, ease: "easeInOut" }}
        className="flex flex-col items-center justify-center bg-base-100 p-4"
        style={{ overflow: "hidden" }}
      >
        <HomeSearchBar onSubmit={handleSearch} />
      </motion.div>
      <div>
        {/* Animated content container */}

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
                <div className="z-[1]">
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
        {/* Search results container */}
        <div
          className={`transition-all duration-500 ease-in-out ${
            isSearching
              ? "opacity-100 translate-y-0"
              : "opacity-0 -translate-y-4"
          }`}
        >
          <DummyResults />
        </div>
      </div>
    </>
  );
}
