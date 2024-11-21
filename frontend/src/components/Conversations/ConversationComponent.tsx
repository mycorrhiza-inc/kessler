"use client";
import React, {
  Dispatch,
  SetStateAction,
  useState,
  useMemo,
  useEffect,
  Suspense,
} from "react";
import { BasicDocumentFiltersList } from "@/components/DocumentFilters";
import {
  emptyQueryOptions,
  QueryFilterFields,
  CaseFilterFields,
  InheritedFilterValues,
  FilterField,
  QueryDataFile,
} from "@/lib/filters";
import { Filing } from "@/lib/types/FilingTypes";
import { AnimatePresence, motion } from "framer-motion";
import axios from "axios";
import { FilingTable } from "./FilingTable";
import { getSearchResults, getFilingMetadata } from "@/lib/requests/search";
import FilingTableQuery from "./FilingTableQuery";
import { ConversationHeader } from "../NavigationHeader";
import { PageContext } from "@/lib/page_context";
import InfiniteScroll from "react-infinite-scroll-component";
import LoadingSpinner from "@/components/styled-components/LoadingSpinner";
import { set } from "date-fns";

const testFiling: Filing = {
  id: "0",
  url: "https://documents.dps.ny.gov/public/Common/ViewDoc.aspx?DocRefId={7F4AA7FC-CF71-4C2B-8752-A1681D8F9F46}",
  date: "05/12/2022",
  lang: "en",
  title: "Press Release - PSC Announces CLCPA Tracking Initiative",
  author: "Public Service Commission",
  source: "Public Service Commission",
  language: "en",
  extension: "pdf",
  file_class: "Press Releases",
  item_number: "3",
  author_organisation: "Public Service Commission",
  // uuid?: "3c4ba5f3-febc-41f2-aa86-2820db2b459a",
};

const TableFilters = ({
  searchFilters,
  setSearchFilters,
  disabledFilters,
  toggleFilters,
}: {
  searchFilters: QueryFilterFields;
  setSearchFilters: Dispatch<SetStateAction<QueryFilterFields>>;
  disabledFilters: FilterField[];
  toggleFilters: () => void;
}) => {
  return (
    <>
      <p className="text-xl font-bold">Search Text</p>
      <input
        type="text"
        placeholder="Type here"
        className="input input-bordered w-full max-w-xs"
      />
      <p className="text-lg font-bold">Filter Documents by: </p>
      <BasicDocumentFiltersList
        queryOptions={searchFilters}
        setQueryOptions={setSearchFilters}
        showQueries={CaseFilterFields}
        disabledQueries={disabledFilters}
      />
    </>
  );
};

const ConversationComponent = ({
  inheritedFilters,
  pageContext,
}: {
  inheritedFilters: InheritedFilterValues;
  pageContext: PageContext;
}) => {
  const disabledFilters = useMemo(() => {
    return inheritedFilters.map((val) => {
      return val.filter;
    });
  }, [inheritedFilters]);

  const initialFilterState = useMemo(() => {
    var initialFilters = emptyQueryOptions;
    inheritedFilters.map((val) => {
      initialFilters[val.filter] = val.value;
    });
    return initialFilters;
  }, [inheritedFilters]);
  const [searchFilters, setSearchFilters] =
    useState<QueryFilterFields>(initialFilterState);
  // const [searchResults, setSearchResults] = useState<string[]>([]);
  const [filing_ids, setFilingIds] = useState<string[]>([]);
  const [filings, setFilings] = useState<Filing[]>([]);
  const [page, setPage] = useState(0);
  const [isSearching, setIsSearching] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");
  const queryData: QueryDataFile = useMemo(() => {
    return {
      filters: searchFilters,
      query: searchQuery,
      start_offset: 0,
    };
  }, [searchFilters]);

  const getUpdates = async () => {
    setIsSearching(true);
    console.log("getting recent updates");
    const data = await getSearchResults(queryData);
    console.log();

    const ids = data.map((item: any) => item.id);
    console.log("ids", ids);
    setFilingIds(ids);
    setIsSearching(false);
  };

  const getMore = async () => {
    setIsSearching(true);
    try {
      const data = await getSearchResults({
        ...queryData,
        start_offset: page + 20,
      });
      setPage(page + 20);
      console.log("data", data);
      if (data.length > 0) {
        setFilingIds([...filing_ids, ...data.map((item: any) => item.id)]);
      }
    } catch (error) {
      console.log(error);
    } finally {
      setIsSearching(false);
    }
  };

  useEffect(() => {
    setIsSearching(true);
    console.log("search filters changed", searchFilters);
    setFilingIds([]);
    setFilings([]);
    getUpdates();
  }, [searchFilters]);

  useEffect(() => {
    if (!filing_ids || isSearching) {
      return;
    }

    const fetchFilings = async () => {
      const newFilings = await Promise.all(
        filing_ids.map(async (id) => {
          const filing_data = await getFilingMetadata(id);
          console.log("new filings", filing_data);
          return filing_data;
        }),
      );

      setFilings((previous) => {
        const existingIds = new Set(previous.map((f) => f.id));
        const uniqueNewFilings = newFilings.filter(
          (f) => !existingIds.has(f.id),
        );
        console.log(" uniques: ", uniqueNewFilings);
        console.log("all data: ", [...previous, ...uniqueNewFilings]);
        return [...previous, ...uniqueNewFilings];
      });
    };

    fetchFilings();
  }, [filing_ids]);

  const [isFocused, setIsFocused] = useState(true);
  const toggleFilters = () => {
    setIsFocused(!isFocused);
  };

  return (
    <div className="drawer drawer-end">
      <input id="my-drawer" type="checkbox" className="drawer-toggle" />
      <div className="drawer-content">
        <div
          id="conversation-header"
          className="flex justify-between items-center mb-4"
        >
          <label htmlFor="my-drawer" className="btn btn-primary drawer-button">
            Filters
          </label>
          <button
            onClick={toggleFilters}
            className="btn btn-outline"
            style={{
              display: !isFocused ? "inline-block" : "none",
            }}
          >
            Filters
          </button>
        </div>
        <div className="w-full h-full">
          <InfiniteScroll
            dataLength={filings.length}
            next={getMore}
            hasMore={true}
            loader={
              <div onClick={getMore}>
                <LoadingSpinner loadingText="Loading Files" />
              </div>
            }
          >
            <FilingTable filings={filings} scroll={false} />
          </InfiniteScroll>
        </div>
      </div>
      <div className="drawer-side">
        <label
          htmlFor="my-drawer"
          aria-label="close sidebar"
          className="drawer-overlay"
        ></label>
        <ul className="menu bg-base-200 text-base-content min-h-full w-90 p-4">
          <TableFilters
            searchFilters={searchFilters}
            setSearchFilters={setSearchFilters}
            disabledFilters={disabledFilters}
            toggleFilters={toggleFilters}
          />
        </ul>
      </div>
    </div>
  );
};

export default ConversationComponent;

// ----------
// OLD ANIMATED FILTER VIEW
// ----------
// <div className="w-full h-full p-10 card relative box-border border-4 border-black flex flex-row overflow-hidden">
//   <AnimatePresence mode="sync">
//     {isFocused && (
//       <motion.div
//         className="flex-none"
//         style={{ padding: "10px" }}
//         initial={{ width: 0, opacity: 0 }}
//         animate={{ width: "500px", opacity: 1 }}
//         exit={{ width: 0, opacity: 0 }}
//       >
//         <TableFilters
//           searchFilters={searchFilters}
//           setSearchFilters={setSearchFilters}
//           disabledFilters={disabledFilters}
//           toggleFilters={toggleFilters}
//         />
//       </motion.div>
//     )}
//     <motion.div
//       className="flex-grow p-10"
//       layout
//       transition={{ duration: 0.3, ease: "easeInOut" }}
//     >
//       <div
//         id="conversation-header"
//         className="flex justify-between items-center mb-4"
//       >
//         <h1 className="text-2xl font-bold">Conversation</h1>
//         <button
//           onClick={toggleFilters}
//           className="btn btn-outline"
//           style={{
//             display: !isFocused ? "inline-block" : "none",
//           }}
//         >
//           Filters
//         </button>
//       </div>
//       <div className="w-full h-full">
//         <FilingTableQuery queryData={queryData} scroll={true} />
//       </div>
//     </motion.div>
//   </AnimatePresence>
// </div>
