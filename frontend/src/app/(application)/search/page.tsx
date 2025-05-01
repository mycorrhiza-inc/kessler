"use client";
import HomeSearchBar from "@/components/NewSearch/HomeSearch";
import React, { useEffect, useState } from "react";

import { motion, AnimatePresence } from "framer-motion";
import DummyResults from "@/components/NewSearch/DummySearchResults";
export default function Page() {
  const [query, setSearchQuery] = useState("");
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

      <div
        className={`transition-all duration-500 ease-in-out ${
          isSearching ? "opacity-100 translate-y-0" : "opacity-0 -translate-y-4"
        }`}
      >
        <DummyResults />
      </div>
    </>
  );
}
