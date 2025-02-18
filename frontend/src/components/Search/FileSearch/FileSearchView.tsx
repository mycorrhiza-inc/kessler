"use client";
import React, {
  Dispatch,
  SetStateAction,
  useState,
  useMemo,
  useEffect,
} from "react";
import { BasicDocumentFiltersList } from "@/components/Filters/DocumentFilters";
import {
  QueryFileFilterFields,
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
import LoadingSpinnerTimeout from "@/components/styled-components/LoadingSpinnerTimeout";

import { ChatModalClickDiv } from "@/components/Chat/ChatModal";
import { useKesslerStore } from "@/lib/store";
import SearchBox from "@/components/Search/SearchBox";
import { FileSearchBoxProps, PageContextMode } from "@/lib/types/SearchTypes";
import InfiniteScrollPlus from "@/components/InfiniteScroll/InfiniteScroll";

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
  searchFilters: QueryFileFilterFields;
  setSearchFilters: Dispatch<SetStateAction<QueryFileFilterFields>>;
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

const FileSearchView = ({
  inheritedFilters,
  DocketColumn,
}: {
  inheritedFilters: InheritedFilterValues;
  DocketColumn?: boolean;
}) => {
  // filter data
  const initialFilterState = useMemo(() => {
    return initialFiltersFromInherited(inheritedFilters);
  }, [inheritedFilters]);

  const disabledFilters = useMemo(() => {
    return disableListFromInherited(inheritedFilters);
  }, [inheritedFilters]);

  const [page, setPage] = useState(0);
  const [isSearching, setIsSearching] = useState(false);

  // query data
  const [queryData, setQueryData] = useState<QueryDataFile>({
    filters: initialFilterState,
    query: "",
  });
  
  const SetComplexQuery = (querydata : QueryDataFile) => {

  }

  // query results
  const [filings, setFilings] = useState<Filing[]>([]);
  const pageSize = 40;

  const getInitialUpdates = async () => {
    const load_initial_pages = 2;
    setPage(0);
    setIsSearching(true);
    setFilings([]);
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

  const [isFocused, setIsFocused] = useState(true);
  const toggleFilters = () => {
    setIsFocused(!isFocused);
  };
  const globalStore = useKesslerStore();
  const experimentalFeatures = globalStore.experimentalFeaturesEnabled;
  const searchBoxProp: FileSearchBoxProps = {
    pageContext: PageContextMode.Files,
    setSearchData: setQueryData,
    inheritedFileFilters: inheritedFilters,
  };

  return (
    <>
      <div id="conversation-header" className=" flex justify-between">
        {experimentalFeatures ? (
          <ChatModalClickDiv
            className="btn btn-accent"
            inheritedFilters={inheritedFiltersFromValues(queryData.filters)}
          >
            Chat with Document List
          </ChatModalClickDiv>
        ) : (
          <div></div>
        )}
      </div>
      <div className="w-full h-full z-50">
        <SearchBox input={searchBoxProp} />
        <InfiniteScrollPlus
          dataLength={filings.length}
          getMore={getMore}
          hasMore={true}
          loadInitial={getInitialUpdates}
          reloadOnChangeObj={queryData}
        >
          <FilingTable filings={filings} DocketColumn={DocketColumn} />
        </InfiniteScrollPlus>
      </div>
    </>
  );
};

export default FileSearchView;
