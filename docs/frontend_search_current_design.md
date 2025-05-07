# Current Frontend Search Architecture and Design Decisions

This document describes the current implementation of the generic search architecture in the frontend, covering key components, data flow, state management, URL synchronization, and infinite scrolling. It also highlights design decisions, placeholders, and areas for future improvement.

---

## 1. High-Level Overview

The front-end search functionality is structured into three main layers:

1. **Page Container (`page.tsx`)**
   - The top-level React Client Component that hosts the search bar and results.
   - Handles layout animations (e.g., expanding/collapsing the search area with Framer Motion).
   - Reads initial query from URL and triggers the first search.

2. **Search State Hook (`useSearchState.ts`)**
   - Encapsulates all search-related state: `searchQuery`, `filters`, `isSearching`, and a trigger indicator.
   - Exposes actions: `triggerSearch`, `resetToInitial`, and a memoized `getResultsCallback`.
   - Synchronizes query parameters with browser history (`pushState`/`replaceState`) and listens for `popstate`.

3. **Results Rendering & Infinite Scroll (`SearchResults.tsx` + `InfiniteScrollPlus`)**
   - Renders search results using a reusable infinite-scroll component.
   - Supports loading spinners, error handling, and incremental data fetches.
   - Uses a `SearchResultsGetter` function to fetch page-based data.

All pieces are designed to be generic so that they can power multiple search contexts:
- Home page search bar (`/`)
- Dedicated search page (`/search?query=<q>`)
- Global “Command-K” modal search
- Document-scoped searches with fixed filters
- “Recently updated documents” list with infinite scroll

---

## 2. Page Container: `frontend/src/app/(application)/search/page.tsx`
```tsx
"use client";
import HomeSearchBar from "@/components/NewSearch/HomeSearch";
import { motion } from "framer-motion";
import { useSearchState } from "@/lib/hooks/useSearchState";
import { SearchResultsComponent } from "@/components/Search/SearchResults";

export default function Page() {
  const { isSearching, triggerSearch, getResultsCallback, searchTriggerIndicator } = useSearchState();

  // On mount, read `?q=<query>` and trigger search
  useEffect(() => {
    const q = new URLSearchParams(window.location.search).get("q");
    if (q) triggerSearch({ query: q.trim() });
  }, []);

  // Adjust search bar height based on `isSearching`
  return (
    <>
      <motion.div initial={{ height: "70vh" }} animate={{ height: isSearching ? "30vh" : "70vh" }}>
        <HomeSearchBar setTriggeredQuery={query => triggerSearch({ query: query.trim() })} />
      </motion.div>

      {/* Animate results in/out */}
      <SearchResultsComponent
        isSearching={isSearching}
        searchGetter={getResultsCallback}
        reloadOnChange={searchTriggerIndicator}
      />
    </>
  );
}
```

**Key Points**
- Uses `motion.div` for smooth height transitions.
- Reads initial URL query only once on mount.
- Delegates all state logic to `useSearchState`.
- Passes the `searchTriggerIndicator` to force reload when query or filters change.

---

## 3. Search State Hook: `frontend/src/lib/hooks/useSearchState.ts`

```ts
export const useSearchState = (): SearchStateExport => {
  const [searchQuery, setSearchQuery] = useState("");
  const { filters, setFilter, deleteFilter, clearFilters, replaceFilters } = useFilterState([]);
  const [isSearching, setIsSearching] = useState(false);
  const [searchTriggerIndicator, setSearchTriggerIndicator] = useState(0);

  // Trigger a new search: update URL, state, filters, and indicator
  const triggerSearch = ({ query, filters }: TriggerSearchObject = {}) => {
    const q = query !== undefined ? query.trim() : searchQuery;
    setSearchTriggerIndicator(i => (i + 1) % 1024);
    setSearchUrl(q, filters || currentFilters);
    setIsSearching(true);
    if (query !== undefined) setSearchQuery(q);
    if (filters) replaceFilters(filters);
  };

  // Memoize result fetcher so it changes only on indicator bump
  const getResultsCallback = useMemo(
    () => generateSearchFunctions({ query: searchQuery.trim(), filters }),
    [searchTriggerIndicator]
  );

  // URL management: pushState vs replaceState depending on route
  const setSearchUrl = (q: string, filters: Filters) => { /*...*/ };

  // Listen to pathname and popstate to reset or re-trigger
  useEffect(/* reset on leave `/search` */, [pathname]);
  useEffect(/* handle back/forward */, []);

  return { searchQuery, filters, isSearching, getResultsCallback, triggerSearch, searchTriggerIndicator, /*...*/ };
};
```

**Key Points**
- **`searchTriggerIndicator`**: Integer that flips on every search to signal downstream components to reload.
- **URL Sync**: Maintains a shareable `/search?query=` URL; uses `pushState` or `replaceState` depending on context.
- **`generateSearchFunctions`**: Placeholder that currently returns either `nilSearchResultsGetter` or `generateFakeResults`.
- **Popstate Handling**: Re-applies search when navigating browser history and resets state when leaving the search route.

---

## 4. Infinite Scroll & Results: `SearchResults.tsx` + `InfiniteScrollPlus`

### 4.1 `SearchResultsComponent`
```tsx
export function SearchResultsComponent({ isSearching, searchGetter, reloadOnChange }) {
  return (
    <AnimatePresence>
      {isSearching && (
        <motion.div /* fade in/out */>
          <SearchResultsInfiniteScroll
            searchGetter={searchGetter}
            reloadOnChange={reloadOnChange}
          />
        </motion.div>
      )}
    </AnimatePresence>
  );
}
```

### 4.2 `SearchResultsInfiniteScroll`
```tsx
function SearchResultsInfiniteScroll({ searchGetter, reloadOnChange }) {
  const [data, setData] = useState<SearchResult[]>([]);
  const [page, setPage] = useState(0);
  const [hasMore, setHasMore] = useState(true);
  const pageSize = 40;

  const loadInitial = async () => {
    setData([]); setPage(0); setHasMore(true);
    // Fetch 2 pages worth of data for quicker above-the-fold
    const initialLimit = pageSize * 2;
    const firstBatch = await searchGetter({ page: 0, limit: initialLimit });
    setData(firstBatch);
    setPage(2);
    if (firstBatch.length < initialLimit) setHasMore(false);
  };

  const fetchMore = async () => {
    const batch = await searchGetter({ page, limit: pageSize });
    setData(d => d.concat(batch));
    setPage(p => p + 1);
    if (batch.length < pageSize) setHasMore(false);
  };

  return (
    <InfiniteScrollPlus
      loadInitial={loadInitial}
      getMore={fetchMore}
      reloadOnChange={reloadOnChange}
      hasMore={hasMore}
      dataLength={data.length}
    >
      <RawSearchResults searchResults={data} />
    </InfiniteScrollPlus>
  );
}
```

### 4.3 `InfiniteScrollPlus`
```tsx
// Wraps `react-infinite-scroll-component`
// - Shows loading spinner on initial load
// - Catches errors and shows retry button
// - Delegates `next` to `getMore` when scrolled to bottom
```

**Key Points**
- Loads two pages up-front for a snappier first impression.
- Uses `reloadOnChange` to clear and re-invoke `loadInitial` on query/filter updates.
- Abstracts common infinite-scroll concerns (loading state, errors) into a reusable component.

---

## 5. Placeholder & Utility Code

- **`generateFakeResults`** (in `search_utils.ts`): Simulates network latency with a delay and returns faux data generated via Faker.
- **Type Definitions** (`new_search_types.ts`): Defines `SearchResult`, `PaginationData`, and `SearchResultsGetter`.
- **Filters**: Managed via a parallel `useFilterState` hook (not detailed here).

---

## 6. Design Decisions & Trade-offs

- **Separation of Concerns**: Splitting state logic (`useSearchState`), layout (`page.tsx`), and data rendering (`SearchResults.tsx`) keeps components focused and testable.
- **URL Synchronization**: Provides shareable links and browser navigation support, at the cost of additional complexity managing `popstate`.
- **Infinite Scroll**: Improves UX for large result sets but requires careful state resets on query/filter changes.
- **Placeholders**: Using `generateFakeResults` allows front-end work to progress in parallel with back-end search API development.
- **SSR & Performance**: Not yet implemented. Future work includes:
  - Adopting Next.js Server Components to fetch initial search results server-side.
  - Streaming rendered HTML inside a `<Suspense>` boundary.
  - Hydration of interactive infinite scroll on the client.

---

## 7. Next Steps

1. **Integrate Real API**: Replace `generateFakeResults` with real `searchGetter` implementations calling `frontend/src/lib/requests/search.ts`.
2. **Server-Side Rendering**: Move initial queries into React Server Components under `/search` and organization pages, wrap in `<Suspense>`.
3. **Command-K Modal**: Reuse `useSearchState` and `SearchResultsComponent` inside a modal with URL rollback on close.
4. **Document-Scoped Filters**: Extend `triggerSearch` to accept default filters that do not update the URL.
5. **Performance Optimization**: Cache server-rendered results for common queries (e.g., recent updates) and leverage Next.js caching headers.

