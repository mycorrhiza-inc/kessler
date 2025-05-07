import React from 'react'
import { Suspense } from 'react'
import HomeSearchBar from '@/components/NewSearch/HomeSearch'
import SearchResultsWrapper from '@/components/Search/SearchResultsWrapper'

interface SearchPageProps {
  searchParams: {
    q?: string
  }
}

export default function Page({ searchParams }: SearchPageProps) {
  const initialQuery = searchParams.q ?? ''

  return (
    <main className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">Search</h1>

      {/* Search input / filters */}
      <HomeSearchBar defaultValue={initialQuery} />

      {/* Server-side results streaming + client hydration */}
      <Suspense fallback={<div className="py-8 text-center">Loading results…</div>}>
        {/* @ts-expect-error Async Server Component */}
        <SearchResultsWrapper initialQuery={initialQuery} />
      </Suspense>
    </main>
  )
}