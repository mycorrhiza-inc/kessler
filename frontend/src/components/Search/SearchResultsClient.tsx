"use client";
import React, { useCallback, useEffect } from "react";
import InfiniteScrollPlus from "../InfiniteScroll/InfiniteScroll";
import RawSearchResults from "./RawSearchResults";
import { useInfiniteSearch } from "@/lib/hooks/useInfiniteSearch";
import {
  SearchResult,
  SearchResultsGetter,
} from "@/lib/types/new_search_types";
import {
  GenericSearchInfo,
  createGenericSearchCallback,
} from "@/lib/adapters/genericSearchCallback";

interface SearchResultsClientProps {
  genericSearchInfo: GenericSearchInfo;
  reloadOnChange: number;
  initialData: SearchResult[];
  initialPage: number;
  children: React.ReactNode; // SSR seed
}

export default function SearchResultsClient({
  genericSearchInfo,
  reloadOnChange,
  initialData,
  initialPage,
  children,
}: SearchResultsClientProps) {
  const searchCallback: SearchResultsGetter = useCallback(
    createGenericSearchCallback(genericSearchInfo),
    [reloadOnChange],
  );
  const { data, hasMore, loadMore, hasReset, loadInitial } = useInfiniteSearch({
    searchCallback,
    initialData,
    initialPage,
  });

  const displayInitalChildren = !hasReset && data.length === initialData.length;

  return (
    <InfiniteScrollPlus
      loadInitial={loadInitial}
      getMore={loadMore}
      reloadOnChange={reloadOnChange}
      dataLength={data.length}
      hasMore={hasMore}
    >
      {displayInitalChildren ? children : <RawSearchResults data={data} />}
    </InfiniteScrollPlus>
  );
}
