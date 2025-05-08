"use client";
import React, { useEffect } from 'react';
import InfiniteScrollPlus from '../InfiniteScroll/InfiniteScroll';
import RawSearchResults from './RawSearchResults';
import { useInfiniteSearch } from '@/lib/hooks/useInfiniteSearch';
import { SearchResult } from '@/lib/types/new_search_types';

interface SearchResultsClientProps {
  q: string;
  filters?: any;
  initialData: SearchResult[];
  initialPage: number;
  children: React.ReactNode; // SSR seed
}

const PAGE_SIZE = 40;

export default function SearchResultsClient({ q, filters, initialData, initialPage, children }: SearchResultsClientProps) {
  const { data, hasMore, loadMore, reset } = useInfiniteSearch({ q, filters, initialData, initialPage });

  // Reset when query or filters change
  useEffect(() => {
    reset({ data: initialData, page: initialPage });
  }, [q, filters, initialData, initialPage, reset]);

  return (
    <InfiniteScrollPlus
      loadInitial={() => { /* no-op: SSR seed covers initial */ }}
      getMore={loadMore}
      reloadOnChange={0}
      dataLength={data.length}
      hasMore={hasMore}
    >
      {data.length === initialData.length ? (
        children
      ) : (
        <RawSearchResults data={data} />
      )}
    </InfiniteScrollPlus>
  );
}