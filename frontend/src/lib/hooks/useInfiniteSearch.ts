import { useState, useCallback } from "react";
import { generateFakeResults } from "@/lib/search/search_utils";
import {
  PaginationData,
  SearchResult,
  SearchResultsGetter,
} from "@/lib/types/new_search_types";
import {
  GenericSearchInfo,
  createGenericSearchCallback,
} from "../adapters/genericSearchCallback";
import { DEFAULT_PAGE_SIZE } from "../constants";

export interface UseInfiniteSearchParams {
  searchCallback: SearchResultsGetter;
  initialData?: SearchResult[];
  pageSize?: number;
}

export function useInfiniteSearch({
  searchCallback,
  initialData,
  pageSize,
}: UseInfiniteSearchParams) {
  const [data, setData] = useState<SearchResult[]>(initialData || []);
  const actualPageSize = pageSize || DEFAULT_PAGE_SIZE;
  const initialPageEstimate = initialData
    ? initialData.length / actualPageSize
    : 0;
  const [page, setPage] = useState<number>(initialPageEstimate);
  const hasMore = data.length === actualPageSize * page;
  const isDataDefined = initialData ? initialData.length != 0 : false;
  const [hasReset, setHasReset] = useState(!isDataDefined);

  const loadMore = async () => {
    const newResults = await searchCallback({
      limit: actualPageSize,
      page: page,
    });
    console.log(`Got ${newResults.length} more search results`);
    setData((prev) => [...prev, ...newResults]);
    setPage((prev) => prev + 1);
  };

  const loadInitial = async () => {
    if (hasReset) {
      return;
    }
    const INITIAL_PAGES = 2;
    const newResults = await searchCallback({
      limit: actualPageSize * INITIAL_PAGES,
      page: 0,
    });

    console.log(`Got ${newResults.length} initial search results`);
    setData(newResults);
    setPage(INITIAL_PAGES);
    setHasReset(true);
  };

  return { data, hasMore, loadMore, hasReset, loadInitial };
}
