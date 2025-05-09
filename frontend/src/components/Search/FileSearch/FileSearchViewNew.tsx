"use client";
import React, { useState, useMemo, useEffect } from "react";
import { getSearchResults } from "@/lib/requests/search";
import InfiniteScrollPlus from "@/components/InfiniteScroll/InfiniteScroll";
import SearchBox from "@/components/Search/SearchBox";
import { FileSearchBoxProps, PageContextMode } from "@/lib/types/SearchTypes";
import { adaptFilingToCard } from "@/lib/adapters/genericCardAdapters";
import Card, { CardSize } from "@/components/NewSearch/GenericResultCard";
import {
  QueryDataFile,
  InheritedFilterValues,
  initialFiltersFromInherited,
} from "@/lib/filters";
import { Filing } from "@/lib/types/FilingTypes";

interface FileSearchViewNewProps {
  /** Initial results for SSR */
  initialData?: Filing[];
  /** Initial page index for SSR */
  initialPage?: number;
  inheritedFilters: InheritedFilterValues;
  DocketColumn?: boolean;
}

const FileSearchViewNew: React.FC<FileSearchViewNewProps> = ({
  initialData = [],
  initialPage = 2,
  inheritedFilters,
}) => {
  // Initialize filters from inherited values
  const initialFilterState = useMemo(
    () => initialFiltersFromInherited(inheritedFilters),
    [inheritedFilters],
  );

  // Query state
  const [queryData, setQueryData] = useState<QueryDataFile>({
    filters: initialFilterState,
    query: "",
  });

  // Search results state, seeded from SSR
  const pageSize = 40;
  const [searchData, setSearchData] = useState<Filing[]>(initialData);
  const [page, setPage] = useState<number>(initialPage);
  const [hasMore, setHasMore] = useState<boolean>(
    initialData.length === pageSize * initialPage,
  );

  const [oldQueryData, setOldQueryData] = useState<QueryDataFile>(queryData);
  if (oldQueryData != queryData) {
    setSearchData(initialData);
    setPage(initialPage);
    setHasMore(initialData.length === pageSize * initialPage);
    setOldQueryData(queryData);
  }

  // Fetch a page of results
  const getPageResults = async (
    pageNum: number,
    limit: number,
  ): Promise<Filing[]> => {
    const newFilings = await getSearchResults(queryData, pageNum, limit);
    setFilings((prev) => [...prev, ...newFilings]);
    if (newFilings.length < limit) {
      setHasMore(false);
    }
    return newFilings;
  };

  // Load initial on client is no-op; SSR seeded data covers initial load
  const loadInitial = async (): Promise<void> => {};

  // Load more on scroll
  const loadMore = async (): Promise<void> => {
    await getPageResults(page, pageSize);
    setPage((prev) => prev + 1);
  };

  // Props for the SearchBox
  const searchBoxProp: FileSearchBoxProps = {
    pageContext: PageContextMode.Files,
    setSearchData: setQueryData,
    inheritedFileFilters: inheritedFilters,
  };

  return (
    <div className="w-full h-full">
      <SearchBox input={searchBoxProp} />
      <InfiniteScrollPlus
        dataLength={filings.length}
        hasMore={hasMore}
        loadInitial={loadInitial}
        getMore={loadMore}
        reloadOnChangeObj={0}
      >
        <div className="grid grid-cols-1 gap-4 p-4">
          {filings.map((filing, index) => {
            return <Card key={index} data={cardData} size={CardSize.Medium} />;
          })}
        </div>
      </InfiniteScrollPlus>
    </div>
  );
};

export default FileSearchViewNew;
