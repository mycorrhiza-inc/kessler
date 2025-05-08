"use client";
import React, { useState, useMemo } from "react";
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
  inheritedFilters: InheritedFilterValues;
  DocketColumn?: boolean;
}

const FileSearchViewNew: React.FC<FileSearchViewNewProps> = ({
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

  // Search results state
  const [filings, setFilings] = useState<Filing[]>([]);
  const [page, setPage] = useState<number>(0);
  const [hasMore, setHasMore] = useState<boolean>(true);
  const pageSize = 40;

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

  // Load initial pages
  const loadInitial = async (): Promise<void> => {
    setHasMore(true);
    setFilings([]);
    // Preload two pages of data
    await getPageResults(0, pageSize * 2);
    setPage(2);
  };

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
        reloadOnChangeObj={queryData}
      >
        <div className="grid grid-cols-1 gap-4 p-4">
          {filings.map((filing, index) => {
            const cardData = adaptFilingToCard(filing);
            return <Card key={index} data={cardData} size={CardSize.Medium} />;
          })}
        </div>
      </InfiniteScrollPlus>
    </div>
  );
};

export default FileSearchViewNew;
