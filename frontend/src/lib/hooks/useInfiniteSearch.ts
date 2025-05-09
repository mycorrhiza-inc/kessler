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
  initialData: SearchResult[];
  initialPage: number;
  pageSize?: number;
}

export function useInfiniteSearch({
  searchCallback,
  initialData,
  initialPage,
  pageSize,
}: UseInfiniteSearchParams) {
  const [data, setData] = useState<SearchResult[]>(initialData);
  const [page, setPage] = useState<number>(initialPage);
  const actualPageSize = pageSize || DEFAULT_PAGE_SIZE;
  const hasMore = initialData.length === actualPageSize * page;
  const [hasReset, setHasReset] = useState(initialData.length == 0);

  const loadMore = async () => {
    const newResults = await searchCallback({
      limit: actualPageSize,
      page: page,
    });
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
    setData(newResults);
    setPage(INITIAL_PAGES);
    setHasReset(true);
  };

  return { data, hasMore, loadMore, hasReset, loadInitial };
}
