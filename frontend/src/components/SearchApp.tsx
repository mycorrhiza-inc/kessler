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
    console.log(`searchhing for ${searchQuery}`);
    try {
      const response = await axios.post("http://localhost/api/v2/search", {
        query: searchQuery,
      });
      if (response.data.length === 0) {
        return;
      }
      if (typeof response.data === "string") {
        setSearchResults([]);
        return;
      }
      console.log("getting data");
      console.log(response.data);
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
