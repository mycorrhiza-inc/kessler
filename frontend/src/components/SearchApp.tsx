"use client";
import axios from "axios";
import { useRef, useState, useEffect } from "react";

import { AnimatePresence, motion } from "framer-motion";
import { CenteredFloatingSearhBox } from "@/components/SearchBox";
import SearchResultBox from "@/components/SearchResultBox";
import ChatBoxInternals from "./ChatBoxInternals";

import {
  extraProperties,
  extraPropertiesInformation,
  emptyExtraProperties,
} from "@/utils/interfaces";
import Header from "./Header";
import { User } from "@supabase/supabase-js";
export default function SearchApp({ user }: { user: User }) {
  const [isSearching, setIsSearching] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");
  const [chatVisible, setChatVisible] = useState(false);
  const [searchResults, setSearchResults] = useState<any[]>([]);
  const [searchDisplay, setSearchDisplay] = useState<any[]>([]);
  const [resultView, setResultView] = useState(false);
  const [searchFilters, setSearchFilters] =
    useState<extraProperties>(emptyExtraProperties);
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
    if (searchResults.length > 0) {
      setShowCard("");
    }
  }, [searchResults]);
  useEffect(() => {
    setSearchDisplay(searchResults);
  }, [searchDisplay]);

  // Should this be refactored out of components and into the lib as a async function that returns search results?
  const handleSearch = async () => {
    setSearchResults([]);
    setIsSearching(true);
    console.log(`searchhing for ${searchQuery}`);
    try {
      const response = await axios.post("/api/v2/search", {
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

  const divRef = useRef<HTMLDivElement>(null);

  return (
    <>
      <Header user={user} />
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

// <AnimatePresence>
//   {chatVisible && (
//     <motion.div
//       key="chat-box"
//       className="chat-box"
//       initial={{ x: "110%" }}
//       animate={{ x: 0 }}
//       exit={{ x: "110%" }}
//       transition={{ type: "tween", stiffness: 200 }}
//       style={{
//         position: "fixed",
//         right: 0,
//         bottom: 0,
//         height: "auto",
//         width: "35%",
//         overflowY: "visible",
//       }}
//     >
//       <ChatBoxInternals
//         setCitations={setSearchDisplay}
//       ></ChatBoxInternals>
//     </motion.div>
//   )}
// </AnimatePresence>
