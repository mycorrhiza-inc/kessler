# LLM Brainstorming: SSR + Infinite Scroll Search Architecture

This document captures ideas and architectural considerations for integrating server-side rendering (SSR) using React Server Components (RSC) into our existing search frontend, while retaining infinite scroll functionality.

---

## 1. Goals

- **SSR Initial Results**: Fetch and render the first chunk (2 pages) of search results server-side to improve performance and SEO.
- **Seamless Hydration**: Hydrate the server-rendered markup on the client without re-fetching initial data.
- **Infinite Scroll on Client**: Continue fetching subsequent pages client-side with our existing infinite scroll logic.
- **Minimal Complexity Impact**: Avoid adding significant architectural overhead; leverage RSC-to-client patterns effectively.

## 2. High-Level Component Flow

1. **`/search/page.tsx` (Server Component)**
   - Reads `searchParams` (`q`, filters).
   - Calls backend search API via `fetch(...)` to get initial results (limit = 2 * pageSize).
   - Wraps a `SearchResultsServer` component in `<Suspense>` for streaming.

2. **`SearchResultsServer` (Server Component)**
   - Accepts `initialResults` and query params.
   - Constructs a server subcomponent, e.g., `<RawSearchResults data={initialResults} />`.
   - Passes this subcomponent as a prop or child to a client component.

3. **`SearchResultsClient` (`'use client'`)**
   - Receives server subcomponent via prop or `children`.
   - Renders children to display initial results (hydrated HTML).
   - Initializes client search state:
     - `useSearchState(initialQuery, initialPage=2, initialData)`
     - Sets `page` to `2`, `data` to `initialResults`, and `hasMore` flag.
   - Hooks into infinite scroll (`InfiniteScrollPlus`) to fetch further pages via `searchGetter`.

4. **`InfiniteScrollPlus`**
   - Unchanged; handles loading more pages when scrolling.

## 3. Key Technical Considerations

- **Passing RSC Subcomponents**
  - Using RSC subcomponent as `children` to `SearchResultsClient` ensures SSR-generated markup is embedded in the HTML stream.
  - Example:
    ```tsx
    const serverChunk = (<RawSearchResults data={initialResults} />);
    return <SearchResultsClient query={q}>{serverChunk}</SearchResultsClient>;
    ```

- **Seeding Client State**
  - Extend `useSearchState` or introduce a `useInfiniteSearch` hook that can accept `initialData`, `initialPage`, and `hasMore`.
  - Avoid double-fetch by disabling the initial load when `initialData` is present.

- **URL Sync and Hydration**
  - On SSR, use `searchParams` from Next.js app router; do not modify URL server-side.
  - Client hook picks up URL and syncs only on navigation or filter changes.

- **Error Handling & Fallback**
  - Wrap server-fetch in `try/catch`; provide a fallback `<Loading/>` or error UI via `<Suspense>` boundary and error boundary.

- **Streaming**
  - Leverage React Streaming for faster TTFP by splitting layout (header, search bar) and results into separate Suspense regions.

## 4. Impact on Existing Code

- **`page.tsx` Migration**
  - Convert from client to server component.
  - Remove `useEffect` for initial query trigger; rely on server fetch.

- **`SearchResultsComponent` vs. `SearchResultsServer`/`Client`**
  - Decompose into two components; mark appropriate ones client or server.
  - Keep `InfiniteScrollPlus` and `RawSearchResults` intact.

- **Hooks**
  - Possibly create a new hook, `useInfiniteSearch`, to handle seeded SSR data.
  - Retain `useSearchState` for modals/other contexts.

## 5. Open Questions

- Should filters also be fetched SSR? (Yes for initial view.)
- How to handle caching (Next.js caching options)?
- Do we need to prefetch next page on server? (E.g., 3rd page.)

---

*End of brainstorming.*
