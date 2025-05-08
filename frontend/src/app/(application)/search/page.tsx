import React, { Suspense } from 'react';
import HomeSearchBar from '@/components/NewSearch/HomeSearch';
import SearchResultsServer from '@/components/Search/SearchResultsServer';

interface SearchPageProps {
  searchParams: { q?: string };
}

export default function Page({ searchParams }: SearchPageProps) {
  const initialQuery = (searchParams.q || '').trim();

  const handleSearch = (query: string) => {
    const q = query.trim();
    if (q) window.location.href = `/search?q=${encodeURIComponent(q)}`;
  };

  return (
    <>
      <div className="flex flex-col items-center justify-center bg-base-100 p-4">
        {/* Client-side search bar */}
        <HomeSearchBar setTriggeredQuery={handleSearch} initialState={initialQuery} />
      </div>
      <Suspense fallback={<div>Loading results...</div>}>
        <SearchResultsServer q={initialQuery} />
      </Suspense>
    </>
  );
}