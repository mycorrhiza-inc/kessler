"use client";
import axios from "axios";
import { useRef, useState, useEffect } from "react";

import { AnimatePresence, motion } from "framer-motion";
import { CenteredFloatingSearhBox } from "@/components/SearchBox";
import SearchResultBox from "@/components/SearchResultBox";
import ChatBoxInternals from "./ChatBoxInternals";

import { QueryFilterFields, emptyQueryOptions } from "@/lib/filters";
import { User } from "@supabase/supabase-js";

import { SearchRequest } from "@/utils/interfaces";
import Navbar from "./Navbar";

export default function SearchApp() {
  const [isSearching, setIsSearching] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");
  const [chatVisible, setChatVisible] = useState(false);
  const [searchResults, setSearchResults] = useState<any[]>([]);
  const [searchDisplay, setSearchDisplay] = useState<any[]>([]);
  const [resultView, setResultView] = useState(false);
  const [searchFilters, setSearchFilters] =
    useState<QueryFilterFields>(emptyQueryOptions);
  const [showCard, setShowCard] = useState("introduction");
  useEffect(() => {
    if (!chatVisible) {
      setSearchDisplay(searchResults);
    }
    // else {
    //   setSearchDisplay([]);
    // }
  }, [chatVisible]);
  useEffect(() => {
    setSearchDisplay(searchResults);
  }, [searchResults]);
  useEffect(() => {
    if (searchDisplay.length > 0) {
      setShowCard("");
    }
    // setShowCard("");
  }, [searchDisplay]);

  // Should this be refactored out of components and into the lib as a async function that returns search results?
  const handleSearch = async () => {
    setSearchResults([]);
    setIsSearching(true);
    console.log(`searchhing for ${searchQuery}`);
    try {
      const response = await axios.post("https://api.kessler.xyz/v2/search", {
        query: searchQuery,
        filters: {
          name: searchFilters.match_name,
          author: searchFilters.match_author,
          docket_id: searchFilters.match_docket_id,
          doctype: searchFilters.match_doctype,
          source: searchFilters.match_source,
        },
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

  const divRef = useRef<HTMLDivElement>(null);

  return (
    <>
      <Navbar user={user} />
      <div
        className="searchContainer"
        ref={divRef}
        style={{
          position: "relative",
          width: "99vw",
          height: "90vh",
          padding: "20px",
          overflow: "scroll",
        }}
      >
        <CenteredFloatingSearhBox
          handleSearch={handleSearch}
          searchQuery={searchQuery}
          setSearchQuery={setSearchQuery}
          chatVisible={chatVisible}
          setChatVisible={setChatVisible}
          inSearchSession={resultView}
          queryOptions={searchFilters}
          setQueryOptions={setSearchFilters}
        />
        {/* Refactor this code to use a motion.div so that SearchResultBox fills the */}
        {/* full screen when chat isnt visible, currently it only takes up 65% */}
        {/* regardless of chat-box being visible. */}
        <motion.div
          className="search-results"
          initial={{ width: "100%" }}
          animate={
            chatVisible ? { width: "calc(100% - 35%)" } : { width: "100%" }
          }
          transition={{ type: "tween", stiffness: 200 }}
          style={{
            position: "relative",
            top: 0,
            left: 0,
            height: "calc(100% - 20px)",
          }}
        >
          <SearchResultBox
            showCard={showCard}
            searchResults={searchDisplay}
            isSearching={isSearching}
          />
        </motion.div>
        {/* Remove animate presense to make chat persistent when closing app */}
        <motion.div
          key="chat-box"
          className="chat-box"
          initial={{ x: "110%" }}
          animate={chatVisible ? { x: 0 } : { x: "110%" }}
          transition={{ type: "tween", stiffness: 200 }}
          style={{
            position: "fixed",
            right: 0,
            bottom: 0,
            height: "auto",
            width: "35%",
            overflowY: "visible",
          }}
        >
          <ChatBoxInternals
            setCitations={setSearchDisplay}
            ragFilters={searchFilters}
          ></ChatBoxInternals>
        </motion.div>
      </div>
    </>
  );
}
