import { useState, useCallback } from "react";
import { generateFakeResults } from "@/lib/search/search_utils";
import { PaginationData, SearchResult } from "@/lib/types/new_search_types";

const PAGE_SIZE = 40;

export interface UseInfiniteSearchParams {
  q: string;
  filters?: any;
  initialData: SearchResult[];
  initialPage: number;
}

export function useInfiniteSearch({
  q,
  filters,
  initialData,
  initialPage,
}: UseInfiniteSearchParams) {
  const [data, setData] = useState<SearchResult[]>(initialData);
  const [page, setPage] = useState<number>(initialPage);
  const [hasMore, setHasMore] = useState<boolean>(
    initialData.length === PAGE_SIZE * 2,
  );

  const reset = useCallback(
    (newParams: { data: SearchResult[]; page: number }) => {
      setData(newParams.data);
      setPage(newParams.page);
      setHasMore(newParams.data.length === PAGE_SIZE * 2);
    },
    [],
  );

  return { data, hasMore, loadMore, reset };
}
