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
  QueryFilterFields,
  CaseFilterFields,
  InheritedFilterValues,
  FilterField,
  QueryDataFile,
  disableListFromInherited,
  initialFiltersFromInherited,
  inheritedFiltersFromValues,
} from "@/lib/filters";
import { Filing } from "@/lib/types/FilingTypes";
import { FilingTable } from "@/components/Tables/FilingTable";
import { getSearchResults } from "@/lib/requests/search";
import InfiniteScroll from "react-infinite-scroll-component";
import LoadingSpinnerTimeout from "../styled-components/LoadingSpinnerTimeout";

import { ChatModalClickDiv } from "../Chat/ChatModal";
import { useKesslerStore } from "@/lib/store";

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
        className="input input-bordered w-full "
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
  const initialFilterState = useMemo(() => {
    return initialFiltersFromInherited(inheritedFilters);
  }, [inheritedFilters]);
  const disabledFilters = useMemo(() => {
    return disableListFromInherited(inheritedFilters);
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
    };
  }, [searchFilters]);
  const pageSize = 40;

  const getUpdates = async () => {
    const load_initial_pages = 2;
    setIsSearching(true);
    console.log("getting recent updates");
    const filings: Filing[] = await getSearchResults(
      queryData,
      0,
      pageSize * load_initial_pages,
    );

    setFilings(filings);
    setPage(load_initial_pages);
    setIsSearching(false);
  };

  const getMore = async () => {
    setIsSearching(true);
    try {
      const new_filings: Filing[] = await getSearchResults(
        queryData,
        page,
        pageSize,
      );
      setPage((prev) => prev + 1);
      if (new_filings.length > 0) {
        setFilings((prev: Filing[]): Filing[] => [...prev, ...new_filings]);
      }
    } catch (error) {
      console.log(error);
    } finally {
      setIsSearching(false);
    }
  };
  useEffect(() => {
    getUpdates();
  }, []);

  const [isFocused, setIsFocused] = useState(true);
  const toggleFilters = () => {
    setIsFocused(!isFocused);
  };
  const globalStore = useKesslerStore();
  const experimentalFeatures = globalStore.experimentalFeaturesEnabled;

  return (
    <div className="drawer drawer-end">
      <input id="my-drawer" type="checkbox" className="drawer-toggle" />
      <div className="drawer-content">
        <div id="conversation-header" className="mb-4 flex justify-between">
          {experimentalFeatures ? (
            <ChatModalClickDiv
              className="btn btn-accent"
              inheritedFilters={inheritedFiltersFromValues(searchFilters)}
            >
              Chat with Document List
            </ChatModalClickDiv>
          ) : (
            <div></div>
          )}
          <div>
            <label
              htmlFor="my-drawer"
              className="btn btn-primary drawer-button"
            >
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
                  replacement={
                    filings.length == 0 ? <p>No Documents Found</p> : <></>
                  }
                />
              </div>
            }
          >
            <FilingTable filings={filings} DocketColumn />
          </InfiniteScroll>
        </div>
      </div>
      <div className="drawer-side">
        <label
          htmlFor="my-drawer"
          aria-label="close sidebar"
          className="drawer-overlay"
        ></label>
        <ul className="menu bg-base-200 text-base-content min-h-full w-90 p-4 w-1/5">
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
