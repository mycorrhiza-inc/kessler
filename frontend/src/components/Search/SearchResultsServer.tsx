import React from 'react';
import SearchResultsClient from './SearchResultsClient';
import RawSearchResults from './RawSearchResults';
import { generateFakeResults } from '@/lib/search/search_utils';
import { SearchResult } from '@/lib/types/new_search_types';

const PAGE_SIZE = 40;

interface SearchResultsServerProps {
  q: string;
  filters?: any;
}

/**
 * Server Component: Fetches initial results and renders the client component.
 */
export default async function SearchResultsServer({ q, filters }: SearchResultsServerProps) {
  // Fetch two pages worth of data server-side
  const initialLimit = PAGE_SIZE * 2;
  const initialResults: SearchResult[] = await generateFakeResults({ page: 0, limit: initialLimit });

  return (
    <SearchResultsClient
      q={q}
      filters={filters}
      initialData={initialResults}
      initialPage={2}
    >
      <RawSearchResults data={initialResults} />
    </SearchResultsClient>
  );
}
