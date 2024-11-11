"use client";
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
import LoadingSpinner from "../styled-components/LoadingSpinner";
import { getSearchResults, getFilingMetadata } from "@/lib/requests/search";

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
      <button
        onClick={toggleFilters}
        className="btn "
        style={{
          alignItems: "center",
          justifyContent: "center",
        }}
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="32"
          height="32"
          viewBox="0 0 512 512"
        >
          <polygon points="400 145.49 366.51 112 256 222.51 145.49 112 112 145.49 222.51 256 112 366.51 145.49 400 256 289.49 366.51 400 400 366.51 289.49 256 400 145.49" />
        </svg>
      </button>
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
}: {
  inheritedFilters: InheritedFilterValues;
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
  const [isSearching, setIsSearching] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");

  const getUpdates = async () => {
    setIsSearching(true);
    console.log("getting recent updates");
    const data = await getSearchResults(queryData);
    console.log();

    const ids = data.map((item: any) => item.sourceID);
    console.log("ids", ids);
    setFilingIds(ids);
    setIsSearching(false);
  };
  useEffect(() => {
    if (!filing_ids) {
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

  // const handleSearch = async () => {
  //   setSearchResults([]);
  //   console.log(`searching for ${searchQuery}`);
  //   try {
  //     const response = await axios.post("https://api.kessler.xyz/v2/search", {
  //       query: searchQuery,
  //       filters: {
  //         name: searchFilters.match_name,
  //         author: searchFilters.match_author,
  //         docket_id: searchFilters.match_docket_id,
  //         doctype: searchFilters.match_doctype,
  //         source: searchFilters.match_source,
  //       },
  //     });
  //     if (response.data.length === 0) {
  //       return;
  //     }
  //     if (typeof response.data === "string") {
  //       setSearchResults([]);
  //       return;
  //     }
  //     console.log("getting data");
  //     console.log(response.data);
  //     const ids = response.data.map((filing: any) => filing.id);
  //     setFilingIds(ids);
  //   } catch (error) {
  //     console.log(error);
  //   } finally {
  //     setIsSearching(false);
  //   }
  // };

  const [isFocused, setIsFocused] = useState(true);
  const toggleFilters = () => {
    setIsFocused(!isFocused);
  };
  const queryData: QueryDataFile = useMemo(() => {
    return {
      filters: searchFilters,
      query: "",
    };
  }, [searchFilters]);

  return (
    <div className="w-full h-full p-10 card relative box-border border-2 border-black flex flex-row">
      <AnimatePresence mode="sync">
        {isFocused && (
          <motion.div
            className="flex-none"
            style={{ padding: "10px" }}
            initial={{ width: 0, opacity: 0 }}
            animate={{ width: "500px", opacity: 1 }}
            exit={{ width: 0, opacity: 0 }}
          >
            <TableFilters
              searchFilters={searchFilters}
              setSearchFilters={setSearchFilters}
              disabledFilters={disabledFilters}
              toggleFilters={toggleFilters}
            />
          </motion.div>
        )}
        <motion.div
          className="flex-grow p-10"
          layout
          transition={{ duration: 0.3, ease: "easeInOut" }}
        >
          <div
            id="conversation-header"
            className="flex justify-between items-center mb-4"
          >
            <h1 className="text-2xl font-bold">Conversation</h1>
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
          <div className="w-full overflow-x-scroll">
            <Suspense
              fallback={
                <LoadingSpinner loadingText="Loading Search Results..." />
              }
            >
              <FilingTable filings={filings} />
            </Suspense>
          </div>
        </motion.div>
      </AnimatePresence>
    </div>
  );
};

export default ConversationComponent;
