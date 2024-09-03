"use  client"
import axios from "axios";
import { useState } from "react";
import { CenteredFloatingSearhBox } from "@/components/SearchBox";
import SearchResultBox from "@/components/SearchResultBox";



export default function Home() {
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