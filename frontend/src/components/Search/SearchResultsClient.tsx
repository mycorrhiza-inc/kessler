"use client";
import React, { useCallback, useEffect, useState } from "react";
import InfiniteScrollPlus from "../InfiniteScroll/InfiniteScrollPlus ";
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
  searchInfo: GenericSearchInfo;
  reloadOnChange: number;
  initialData?: SearchResult[];
  children?: React.ReactNode; // SSR seed
}

export default function SearchResultsClient({
  searchInfo: genericSearchInfo,
  reloadOnChange,
  initialData,
  children,
}: SearchResultsClientProps) {
  const searchCallback: SearchResultsGetter = useCallback(
    createGenericSearchCallback(genericSearchInfo),
    [reloadOnChange],
  );
  const { data, hasMore, loadMore, hasReset, loadInitial } = useInfiniteSearch({
    searchCallback,
    initialData,
  });

  const displayInitalChildren =
    children && !hasReset && initialData && data.length === initialData.length;

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
