"use client";
import React, {
  Dispatch,
  SetStateAction,
  useState,
  useMemo,
  useEffect,
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
import { FilingTable } from "./FilingTable";
import { getSearchResults, getFilingMetadata } from "@/lib/requests/search";
import { PageContext } from "@/lib/page_context";
import InfiniteScroll from "react-infinite-scroll-component";
import LoadingSpinnerTimeout from "../styled-components/LoadingSpinnerTimeout";

import ConversationDescription from "./ConversationDescription";

const TableFilters = ({
  searchQuery,
  setSearchQuery,
  searchFilters,
  setSearchFilters,
  disabledFilters,
  toggleFilters,
}: {
  searchQuery: string;
  setSearchQuery: Dispatch<SetStateAction<string>>;
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
        value={searchQuery}
        onChange={(e) => setSearchQuery(e.target.value)}
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
      const newFilings_raw = await Promise.all(
        filing_ids.map(async (id): Promise<Filing | null> => {
          const filing_data = await getFilingMetadata(id);
          console.log("new filings", filing_data);
          return filing_data;
        }),
      );
      const notnull_newFilings = newFilings_raw.filter(
        (f: Filing | null) => f !== null && f !== undefined,
      );
      const newFilings: Filing[] = notnull_newFilings as Filing[];

      setFilings((previous: Filing[]): Filing[] => {
        return [...previous, ...newFilings];
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
        <div id="conversation-header" className="mb-4 flex justify-end">
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
                <LoadingSpinnerTimeout
                  loadingText="Loading Files"
                  timeoutSeconds={3}
                />
              </div>
            }
          >
            <FilingTable filings={filings} />
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
            searchQuery={searchQuery}
            setSearchQuery={setSearchQuery}
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
