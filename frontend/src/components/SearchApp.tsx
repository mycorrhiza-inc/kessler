"use client";
import axios from "axios";
import SearchResult from "@/components/SearchResult";
import { useState } from "react";
import SearchBox, { CenteredFloatingSearhBox } from "../components/SearchBox";
import SearchResultBox from "@/components/SearchResultBox";
import { Grid, Box, Stack } from "@mui/joy";

export default function SearchApp() {
  const iOS =
    typeof navigator !== "undefined" &&
    /iPad|iPhone|iPod/.test(navigator.userAgent);

  const [isSearching, setIsSearching] = useState(false);
  const [searchResults, setSearchResults] = useState<any[]>([]);
  const [searchQuery, setSearchQuery] = useState("");
  const [resultView, setResultView] = useState(false);
  // ...

  const handleSearch = async () => {
    setSearchResults([]);
    setIsSearching(true);
    try {
      const response = await axios.post("http://localhost:4041/search", {
        query: searchQuery,
      });
      if (response.data.length === 0) {
        return;
      }
      setSearchResults(response.data);
    } catch (error) {
      console.log(error);
    } finally {
      setIsSearching(false);
    }
    setResultView(true);
  };

  /*
   */

  return (
    <main className="flex min-w-screen h-100vh justify-center">
      <CenteredFloatingSearhBox
        handleSearch={handleSearch}
        searchQuery={searchQuery}
        setSearchQuery={setSearchQuery}
        inSearchSession={resultView}
      />
      <SearchResultBox
        searchResults={searchResults}
        isSearching={isSearching}
      />
    </main>
  );
}
