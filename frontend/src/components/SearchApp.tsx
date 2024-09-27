"use client";
import axios from "axios";
import { useRef, useState, useEffect } from "react";

import { motion } from "framer-motion";
import { CenteredFloatingSearhBox } from "@/components/SearchBox";
import SearchResultBox from "@/components/SearchResultBox";
import ChatBoxInternals from "./ChatBoxInternals";

export default function SearchApp() {
  const [isSearching, setIsSearching] = useState(false);
  const [searchResults, setSearchResults] = useState<any[]>([]);
  const [searchQuery, setSearchQuery] = useState("");
  const [chatVisible, setChatVisible] = useState(false);
  const [resultView, setResultView] = useState(false);
  const [renderChat, setRenderChat] = useState(false);

  useEffect(() => {
    if (chatVisible) {
      setRenderChat(true);
    } else {
      const timer = setTimeout(() => setRenderChat(false), 2000);
      return () => clearTimeout(timer);
    }
  }, [chatVisible]);

  // Should this be refactored out of components and into the lib as a async function that returns search results?
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
      {/* Refactor this code to use a motion.div so that SearchResultBox fills the */}
      {/* full screen when chat isnt visible, currently it only takes up 65% */}
      {/* regardless of chat-box being visible. */}
      <div
        style={{
          display: "flex",
          height: "calc(100% - 20px)",
        }}
      >
        <motion.div
          className="search-results"
          initial={{ width: "100%" }}
          animate={chatVisible ? { width: "10%" } : { width: "150%" }}
          transition={{ type: "tween", stiffness: 200 }}
        >
          <SearchResultBox
            searchResults={searchResults}
            isSearching={isSearching}
          />
        </motion.div>
        <motion.div
          className="chat-box"
          initial={{ x: "110%" }}
          animate={chatVisible ? { x: 0 } : { x: "110%" }}
          transition={{ type: "tween", stiffness: 200 }}
          style={{
            flex: "0 0 35%",
            overflowY: "visible",
          }}
        >
          {renderChat && (
            <ChatBoxInternals
              setCitations={setSearchResults}
            ></ChatBoxInternals>
          )}
        </motion.div>
      </div>
    </div>
  );
}
