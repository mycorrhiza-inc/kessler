"use client";
import axios from "axios";
import SearchResult from "@/components/SearchResult";
import { useState } from "react";
import SearchBox from "../components/SearchBox";
import SearchResultBox from "@/components/SearchResultBox";

export default function Home() {
  const [isSearching, setIsSearching] = useState(false);
  const [searchResults, setSearchResults] = useState<any[]>([]);
  const [searchQuery, setSearchQuery] = useState("");

  // ...

  const handleSearch = async () => {
    setSearchResults([]);
    setIsSearching(true);
      try {
        const response = await axios.post("http://localhost:4041/search", {
          query: searchQuery,
        });
        setSearchResults(response.data);
      } catch (error) {
        console.log(error);
      } finally {
        setIsSearching(false);
      }
    };

  return (
    <main className="flex flex-col min-w-screen min-h-screen items-center justify-center">
      <div
        className="viewport"
        style={{ width: "80vw", height: "80vh", margin: 0, padding: "50px" }}
      >
        <SearchBox
          searchQuery={searchQuery}
          setSearchQuery={setSearchQuery}
          handleSearch={handleSearch}
        />
        <SearchResultBox
          searchResults={searchResults}
          isSearching={isSearching}
        />
      </div>
    </main>
  );
}
