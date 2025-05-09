# Current Frontend Search Architecture and Design Decisions (Updated for SSR Support)

This document describes the updated implementation of the generic search architecture in the frontend, now enhanced with Next.js Server Components for SSR initial fetch, covering key components, data flow, state management, URL synchronization, and infinite scrolling. It also highlights design decisions, placeholders, and areas for future improvement.

---

## 1. High-Level Overview

The front-end search functionality now comprises two distinct render phases:

1. **Server-Side (Next.js Server Component)**
   - **Page Render (`/search/page.tsx`)**: A React Server Component that reads `searchParams.q`, renders the search bar as a client component, and suspends to render initial results server-side inside `<Suspense>`.
   - **Results Fetch (`SearchResultsServer.tsx`)**: A server component that fetches the first two pages of results (using placeholder `generateFakeResults`) and seeds the client component.

2. **Client-Side (React Client Components)**
   - **Search Bar (`HomeSearchBarClientBaseUrl` / `HomeSearchBar`)**: Renders a search input, handles user input, and updates the URL (e.g., `/search?q=<query>`).
   - **Hydrated Results (`SearchResultsClient.tsx`)**: Hydrates with the SSR seed, then manages incremental infinite scrolling via a custom `useInfiniteSearch` hook.

All pieces remain generic and reusable for multiple search contexts:
- Home page search bar (`/`)
- Dedicated search page (`/search?q=<q>`)
- Global “Command-K” modal search
- Document-scoped searches with fixed filters
- “Recently updated documents” list with infinite scroll

---

## 2. Page Render: `frontend/src/app/(application)/search/page.tsx`
```tsx
// React Server Component (SSR)
import React, { Suspense } from "react";
import { HomeSearchBarClientBaseUrl } from "@/components/NewSearch/HomeSearch";
import SearchResultsServer from "@/components/Search/SearchResultsServer";
import LoadingSpinner from "@/components/styled-components/LoadingSpinner";

interface SearchPageProps {
  searchParams: { q?: string };
}

export default function Page({ searchParams }: SearchPageProps) {
  const initialQuery = (searchParams.q || "").trim();

  return (
    <>
      {/* Client-only search bar with URL sync */}
      <div className="p-4 bg-base-100 flex justify-center">
        <HomeSearchBarClientBaseUrl
          baseUrl="/search"
          initialState={initialQuery}
        />
      </div>

      {/* SSR initial fetch inside Suspense */}
      <Suspense
        fallback={
          <LoadingSpinner loadingText="Fetching results from server." />
        }
      >
        <SearchResultsServer q={initialQuery} />
      </Suspense>
    </>
  );
}
```

**Key Points**
- **Server Component**: Renders on the server, supports streaming and caching.
- Reads `?q=` from `searchParams` and trims it.
- Uses `<Suspense>` with a loading spinner while server fetch happens.
- Delegates search bar rendering and SSR results to child components.

---

## 3. SSR Results Component: `SearchResultsServer.tsx`
```tsx
import React from "react";
import SearchResultsClient from "./SearchResultsClient";
import RawSearchResults from "./RawSearchResults";
import { generateFakeResults } from "@/lib/search/search_utils";
import type { SearchResult } from "@/lib/types/new_search_types";

const PAGE_SIZE = 40;

export default async function SearchResultsServer({ q, filters }) {
  // Fetch two pages worth of data server-side
  const initialLimit = PAGE_SIZE * 2;
  const initialResults: SearchResult[] = await generateFakeResults({
    page: 0,
    limit: initialLimit,
  });

  // Seed the client component with SSR data and initial page index
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
```

**Key Points**
- Fetches initial results on the server (placeholder logic).
- Passes SSR seed to a client component via props (`initialData`, `initialPage`).
- Renders a `RawSearchResults` snapshot inside the client wrapper.

---

## 4. Hydrated Results & Infinite Scroll: `SearchResultsClient.tsx` + `useInfiniteSearch`

### 4.1 `SearchResultsClient.tsx`
```tsx
"use client";
import React, { useEffect } from "react";
import InfiniteScrollPlus from "../InfiniteScroll/InfiniteScroll";
import RawSearchResults from "./RawSearchResults";
import { useInfiniteSearch } from "@/lib/hooks/useInfiniteSearch";
import type { SearchResult } from "@/lib/types/new_search_types";

interface Props { q: string; filters?: any; initialData: SearchResult[]; initialPage: number; children: React.ReactNode; }

export default function SearchResultsClient({ q, filters, initialData, initialPage, children }: Props) {
  const { data, hasMore, loadMore, reset } = useInfiniteSearch({ q, filters, initialData, initialPage });

  // Reset infinite scroll when query/filters change
  useEffect(() => {
    reset({ data: initialData, page: initialPage });
  }, [q, filters, initialData, initialPage, reset]);

  return (
    <InfiniteScrollPlus
      loadInitial={() => {} /* SSR seed covers initial load */}
      getMore={loadMore}
      reloadOnChange={0}
      dataLength={data.length}
      hasMore={hasMore}
    >
      {data.length === initialData.length ? (
        children /* Show SSR-rendered snapshot */
      ) : (
        <RawSearchResults data={data} />
      )}
    </InfiniteScrollPlus>
  );
}
```

### 4.2 `useInfiniteSearch` Hook
- Encapsulates infinite scroll state: `data`, `page`, `hasMore`.
- Provides `loadMore()` to fetch the next page via a real API or placeholder.
- `reset()` allows jumping back to SSR seed on query/filter changes.
- Replaces the previous `useSearchState` + `generateSearchFunctions` pattern for page-based fetch.

**Key Points**
- SSR seed avoids a blank initial render.
- Hook-driven infinite scrolling with auto-reset logic.
- Simplifies client-side fetching by unifying load/more/reset operations.

---

## 5. Search Bar: `HomeSearchBarClientBaseUrl` & `HomeSearchBar`

- **`HomeSearchBarClientBaseUrl`**
  - Wraps the client `HomeSearchBar`, syncs the query to a given `baseUrl` via `window.location.href`.
  - Used in `/search/page.tsx` to update `?q=` on submit.

- **`HomeSearchBar`**
  - Renders a logo, input box, and optional filters (e.g., `StateSelector`).
  - Manages local input state and delegates submission via `setTriggeredQuery`.

*Note:* The legacy `useSearchState` hook remains for contexts like Command-K modal and homepage, but the `/search` page now uses direct URL navigation and SSR.

---

## 6. Legacy Client-Side Search State (Optional)

The previous implementation used a `useSearchState` hook to:
- Manage `searchQuery`, `filters`, `isSearching`, and a `searchTriggerIndicator`.
- Synchronize the URL via `pushState`/`replaceState` and handle browser history.
- Memoize a `getResultsCallback` for client-only fetch.

This pattern is still in use for:
- Global modal search (Command-K).
- Home page instant search without SSR.

---

## 7. Placeholder & Utility Code

- **`generateFakeResults`** (in `search_utils.ts`): Simulates data; replace with real API calls under `/lib/requests/search.ts`.
- **Type Definitions** (`new_search_types.ts`): Defines `SearchResult`, `PaginationData`, and fetcher signatures.
- **Infinite Scroll Wrapper** (`InfiniteScrollPlus`): Common loading, error, retry handling around `react-infinite-scroll-component`.

---

## 8. Design Decisions & Trade-offs

- **SSR for Initial Load**: Improves SEO and perceived performance; adds complexity around hydration and seed management.
- **Hook-based Infinite Scroll**: Encapsulates state and logic, easing reuse across contexts.
- **URL Navigation**: Simplifies SSR page but omits client-side `pushState` UX; future work could adopt shallow routing (`router.replace`) for SPA feel.
- **Placeholders**: Enable parallel development; real API integration is next.

---

## 9. Next Steps

1. **Integrate Real API**: Implement `search` requests in `/lib/requests/search.ts`, and update `useInfiniteSearch` to call them.
2. **Enhance URL Sync**: Use Next.js `useRouter` for shallow routing to avoid full-page reloads.
3. **Streaming & Suspense**: Further split SSR into chunked streams for faster TTFP.
4. **Command-K Modal**: Migrate to use `useInfiniteSearch` and reinstate URL rollback on close.
5. **Scoped Filters**: Pass filter defaults through SSR and client props without URL changes.
6. **Performance**: Leverage Next.js cache headers and preview mode for common queries.
