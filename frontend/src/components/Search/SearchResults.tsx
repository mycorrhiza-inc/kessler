'use client'

import React, { useState, useEffect, useRef } from 'react'
import { QueryDataFile } from '@/lib/types/new_search_types'
import { getSearchResults } from '@/lib/requests/search'

interface SearchResultComponentProps {
  initialResults: any[]
  initialQuery: string
  initialPage: number
  pageSize: number
}

export default function SearchResults({
  initialResults,
  initialQuery,
  initialPage,
  pageSize,
}: SearchResultComponentProps) {
  const [results, setResults] = useState<any[]>(initialResults)
  const [page, setPage] = useState<number>(initialPage)
  const [isLoading, setIsLoading] = useState<boolean>(false)
  const [hasMore, setHasMore] = useState<boolean>(initialResults.length === pageSize)
  const [error, setError] = useState<string | null>(null)
  const sentinelRef = useRef<HTMLDivElement | null>(null)

  const loadNextPage = async () => {
    if (isLoading || !hasMore) return
    setIsLoading(true)
    try {
      const queryData: QueryDataFile = { query: initialQuery, filters: [] }
      const nextPageIndex = page
      const newResults = await getSearchResults(queryData, nextPageIndex, pageSize)
      setResults((prev) => [...prev, ...newResults])
      setPage(nextPageIndex + 1)
      if (newResults.length < pageSize) {
        setHasMore(false)
      }
    } catch (err: any) {
      setError(err.message || 'Error loading more results')
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    const sentinel = sentinelRef.current
    if (!sentinel) return
    const observer = new IntersectionObserver((entries) => {
      if (entries[0].isIntersecting) {
        loadNextPage()
      }
    })
    observer.observe(sentinel)
    return () => observer.disconnect()
  }, [sentinelRef.current, isLoading, hasMore])

  return (
    <div className="search-results">
      <div className="grid grid-cols-1 gap-4 p-8">
        {results.map((item, idx) => (
          <div key={item.id || idx} className="search-hit">
            {/* Customize rendering as needed */}
            <h2 className="font-semibold">{item.title}</h2>
          </div>
        ))}
      </div>

      {error && <div className="text-red-600 text-center">{error}</div>}
      {isLoading && <div className="text-center py-4">Loading more…</div>}

      <div ref={sentinelRef} style={{ height: '1px' }} />

      {!hasMore && <div className="text-center py-4">No more results</div>}
    </div>
  )
}