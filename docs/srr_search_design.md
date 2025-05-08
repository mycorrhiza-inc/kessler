# Server-Side Rendered Search Design (SSR + Infinite Scroll)

Date: 2025-05-08
Author: Goose LLM

---

This document proposes an architecture for adding Server-Side Rendering (SSR) to our existing frontend search, leveraging React Server Components (RSC) and Suspense. It preserves our infinite scroll UX and minimal client-side complexity.

## 1. Goals

- SSR initial search results (two pages) for improved performance, SEO, and TTFP.
- Seamless hydration without double-fetching initial data.
- Client-side infinite scroll continuation for subsequent pages.
- Maintain URL sync, filter support, and existing hook abstractions.
- Minimal refactoring: keep `RawSearchResults` and `InfiniteScrollPlus` unchanged.

## 2. High-Level Architecture

```
/app/search/page.tsx       (Server Component)
└── <Suspense> → SearchResultsServer (Server)
     └── <SearchResultsClient> (Client)
           ├── RawSearchResults (SSR seed)
           └── InfiniteScrollPlus (client paging)
```

### 2.1 `/app/search/page.tsx` (Server)
```tsx
// app/search/page.tsx (auto server)
import SearchResultsServer from '@/components/Search/SearchResultsServer';

export default async function SearchPage({ searchParams }) {
  const q = (searchParams.q || '').trim();
  const filters = parseFilters(searchParams);
  // Fetch initial 2 pages server-side
  const initialLimit = PAGE_SIZE * 2;
  const initialResults = await fetchSearchAPI({ q, filters, page: 0, limit: initialLimit });

  return (
    <Suspense fallback={<Loading />}>  
      <SearchResultsServer q={q} filters={filters} initialResults={initialResults} />
    </Suspense>
  );
}
```

### 2.2 `SearchResultsServer` (Server)
```tsx
// components/Search/SearchResultsServer.tsx
import SearchResultsClient from './SearchResultsClient';
import RawSearchResults from './RawSearchResults';

export default function SearchResultsServer({ q, filters, initialResults }) {
  // Create SSR seed subcomponent
  const seed = <RawSearchResults data={initialResults} />;
  return (
    <SearchResultsClient
      q={q}
      filters={filters}
      initialData={initialResults}
      initialPage={2}
    >
      {seed}
    </SearchResultsClient>
  );
}
```

### 2.3 `SearchResultsClient` (`'use client'`)
```tsx
// components/Search/SearchResultsClient.tsx
'use client';
import { useInfiniteSearch } from '@/lib/hooks/useInfiniteSearch';
import InfiniteScrollPlus from './InfiniteScrollPlus';

export default function SearchResultsClient({
  q,
  filters,
  initialData,
  initialPage,
  children // RawSearchResults SSR seed
}) {
  const { data, hasMore, loadMore, isLoading, reset } = useInfiniteSearch({
    q,
    filters,
    initialData,
    initialPage
  });

  // Reset when q or filters change
  useEffect(() => reset({ data: initialData, page: initialPage }), [q, filters]);

  return (
    <InfiniteScrollPlus
      dataLength={data.length}
      hasMore={hasMore}
      loadInitial={() => {}}   
      getMore={loadMore}
      reloadOnChange={false}
      loader={<Loading />}
      endMessage={<EndMessage />}
    >
      {data.length === initialData.length ? children : <RawSearchResults data={data} />}
    </InfiniteScrollPlus>
  );
}
```

## 3. Hook: `useInfiniteSearch`

A specialized hook seeding SSR data and managing client paging.

```ts
// lib/hooks/useInfiniteSearch.ts
import { useState, useCallback } from 'react';
import fetchSearchAPI from '@/lib/requests/search';

export function useInfiniteSearch({ q, filters, initialData, initialPage }) {
  const [data, setData] = useState(initialData);
  const [page, setPage] = useState(initialPage);
  const [hasMore, setHasMore] = useState(initialData.length === PAGE_SIZE * 2);

  const loadMore = useCallback(async () => {
    const batch = await fetchSearchAPI({ q, filters, page, limit: PAGE_SIZE });
    setData(prev => [...prev, ...batch]);
    setPage(p => p + 1);
    if (batch.length < PAGE_SIZE) setHasMore(false);
  }, [q, filters, page]);

  const reset = useCallback(({ data: newData, page: newPage }) => {
    setData(newData);
    setPage(newPage);
    setHasMore(newData.length === PAGE_SIZE * 2);
  }, []);

  return { data, hasMore, loadMore, reset };
}
```

## 4. URL Sync & Filters

- The server component reads `searchParams` for `q` and filter values; no URL mutation.
- Client only updates URL on further filter or query changes via `useSearchState` when used in client-only contexts (e.g., modal).
- For full-page search, we rely on server context; client hook does not pushState for initial mount.

## 5. Streaming & Suspense

- Wrap `SearchResultsServer` in `<Suspense>` to stream UI progressively.
- Use an error boundary around server fetch for graceful fallback.

## 6. Backward Compatibility & Impact

- **`page.tsx`**: Convert to server component; remove initial `useEffect` trigger logic.
- **`SearchResultsComponent`**: Deprecated; replaced by `SearchResultsServer` + client split.
- **`RawSearchResults`** & **`InfiniteScrollPlus`**: No changes.
- **`useSearchState`**: Remains for non-SSR contexts (modal, page reloads outside `/search`).
- **Routing**: Next.js App Router is required.

## 7. Next Steps

1. Implement and test SSR fetch and streaming.
2. Build and validate `SearchResultsServer`/`SearchResultsClient` split.
3. Migrate existing `/search` page to new SSR version.
4. QA infinite scroll behavior: no flash, no double-fetch.
5. Measure performance gains (TTFP, SEO).

---

*End of Proposal.*
