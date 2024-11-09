"use client";
"use client";
import React, {
  Dispatch,
  SetStateAction,
  useState,
  useMemo,
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
import {
  Filing
} from "@/lib/types/FilingTypes";
import { AnimatePresence, motion } from "framer-motion";
import axios from "axios";
import FilingTableQuery from "./FilingTable";


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
  const [searchResults, setSearchResults] = useState<string[]>([]);
  const [filing_ids, setFilingIds] = useState<string[]>([]);
  const [isSearching, setIsSearching] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");


  const handleSearch = async () => {
    setSearchResults([]);
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
      const ids = response.data.map((filing: any) => filing.id);
      setFilingIds(ids);
    } catch (error) {
      console.log(error);
    } finally {
      setIsSearching(false);
    }
  };


  const [isFocused, setIsFocused] = useState(false);
  const toggleFilters = () => {
    setIsFocused(!isFocused);
  };
  const queryData: QueryDataFile = {
    filters: searchFilters,
    query: "",
  };

  return (
    <div className="w-full h-full p-10 card grid grid-flow-col auto-cols-2 box-border border-2 border-black ">
      <AnimatePresence>
        {isFocused && (
          <motion.div
            style={{
              padding: "10px",
              transition: "width 0.3s ease-in-out",
            }}
            initial={{ x: "-50%" }}
            animate={{ x: "0" }}
            exit={{ x: "-50%", opacity: 0 }}
          >
            <TableFilters
              searchFilters={searchFilters}
              setSearchFilters={setSearchFilters}
              disabledFilters={disabledFilters}
              toggleFilters={toggleFilters}
            />
          </motion.div>
        )}
      </AnimatePresence>
      <div className="p-10">
        <div id="conversation-header p-10 justify-between"></div>
        <h1 className=" text-2xl font-bold">Conversation</h1>
        <button
          onClick={toggleFilters}
          className="btn btn-outline"
          style={{
            display: !isFocused ? "inline-block" : "none",
          }}
        >
          Filters
        </button>
        <div className="w-full overflow-x-scroll">
          <Suspense fallback={<div>Loading...</div>}>
            <FilingTableQuery queryData={queryData} />
          </Suspense>
        </div>
      </div>
    </div>
  );
};

export default ConversationComponent;
