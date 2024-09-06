"use client";
import axios from "axios";
import { useState } from "react";

import  { CenteredFloatingSearhBox } from "@/components/SearchBox";
import SearchResultBox from "@/components/SearchResultBox";
import Chatbox from "@/components/Chatbox";

export default function SearchApp() {
  const [isSearching, setIsSearching] = useState(false);
  const [searchResults, setSearchResults] = useState<any[]>([]);
  const [searchQuery, setSearchQuery] = useState("");
  const [chatVisible, setChatVisible] = useState(false);
  const [resultView, setResultView] = useState(false);

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

  const chatMobileStyle = {};
  const chatDesktopStyle = {};

  return (
    <main className="flex min-w-screen h-100vh justify-center">
      <CenteredFloatingSearhBox
        handleSearch={handleSearch}
        searchQuery={searchQuery}
        setSearchQuery={setSearchQuery}
        setChatVisible={setChatVisible}
        inSearchSession={resultView}
      />
      <SearchResultBox
        searchResults={searchResults}
        isSearching={isSearching}
      />
      <Chatbox chatVisible={chatVisible} setChatVisible={setChatVisible} />
    </main>
  );
}
