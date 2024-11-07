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
import { AnimatePresence, motion } from "framer-motion";
import FilingTableQuery from "./FilingTable";

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
      <div className=" p-10">
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
