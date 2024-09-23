"use client";
import axios from "axios";
import { useRef, useState } from "react";

import { CenteredFloatingSearhBox } from "@/components/SearchBox";
import SearchResultBox from "@/components/SearchResultBox";
import ChatBoxInternals from "./ChatBoxInternals";

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
  const divRef = useRef<HTMLDivElement>(null);

  return (
    <div
      className="searchContainer"
      ref={divRef}
      style={{
        position: "relative",
        width: "99vw",
        height: "90vh",
        padding: "20px",
        overflow: "hidden",
      }}
    >
      <CenteredFloatingSearhBox
        handleSearch={handleSearch}
        searchQuery={searchQuery}
        setSearchQuery={setSearchQuery}
        setChatVisible={setChatVisible}
        inSearchSession={resultView}
      />
      <div
        className="results-container"
        style={{
          display: "flex",
          height: "calc(100% - 20px)",
        }}
      >
        <div
          className="search-results"
          style={{
            flex: 1,
            overflowY: "auto",
          }}
        >
          <SearchResultBox
            searchResults={searchResults}
            isSearching={isSearching}
          />
        </div>
        {chatVisible && (
          <div
            className="chat-box"
            style={{
              flex: "0 0 35%",
              overflowY: "visible",
            }}
          >
            <ChatBoxInternals
              setCitations={setSearchResults}
            ></ChatBoxInternals>
          </div>
        )}
      </div>
    </div>
  );
}
