# Generic Card Architecture

## Overview

This document outlines the high-level architecture and design strategy for migrating existing search and lookup components to leverage the new frontend search infrastructure described in `frontend_search_current_design.md`. We will unify search fetchers with the generic card adapter layer to provide consistent, SSR-capable lookup pages.

## Goals

- Reuse existing fetchers to retrieve search data.
- Plug data into `genericCardAdapters` to map hits to card components.
- Provide a shared search core that handles paging, sorting, filtering, and SSR.
- Maintain pluggable adapters for each lookup page (FileSearch, ConversationLookup, OrganizationLookup).
- Simplify SSR integration to a minimal common wrapper.

## Components to Migrate

1. **FileSearchView** (`components/Search/FileSearch/FileSearchView.tsx`)
2. **ConversationTable** (`components/LookupPages/ConvoLookup/ConversationTable.tsx`)
3. **OrganizationTable** (`components/LookupPages/OrgLookup/OrganizationTable.tsx`)

## New Architecture Components

### 1. Search Core

A React Server Component (RSC) / Client wrapper that:
- Reads query parameters (page, sort, filter) from props or URL.
- Invokes the appropriate fetcher with search parameters.
- Passes raw hits and metadata (total, page) into the adapter layer.
- Renders the adapter’s output (cards, table rows).

### 2. Fetchers

Keep existing fetcher functions in `frontend/src/lib/searchFetchers.ts`, for example:
- `fileSearchFetcher(params): Promise<SearchResponse<FileHit>>`
- `conversationLookupFetcher(params): Promise<SearchResponse<ConversationHit>>`
- `organizationLookupFetcher(params): Promise<SearchResponse<OrganizationHit>>`

### 3. Generic Card Adapters

In `lib/adapters/genericCardAdapters.ts`, define an Adapter interface:

```ts
interface SearchAdapter<Hit> {
  /** Map a single hit into card props for rendering */
  mapHitToCard: (hit: Hit) => CardProps;

  /** Optional: Provide an empty state when no hits are returned */
  getEmptyState?: () => React.ReactNode;

  /** Optional: Transform raw query params to fetcher-compatible shape */
  transformQueryParams?: (rawParams: Record<string, any>) => SearchParams;
}
```

Implement adapter instances per page:
- `fileSearchAdapter`
- `conversationLookupAdapter`
- `organizationLookupAdapter`

### 4. Generic Search Page Template

Create a template component in `components/Search/GenericSearchPage.tsx`:

```tsx
export function GenericSearchPage<Hit>({
  adapter,        // SearchAdapter<Hit>
  fetcher,       // (params: SearchParams) => Promise<SearchResponse<Hit>>
  title,         // Page title
}) {
  // 1. Parse SSR props or URL params
  const rawParams = useSearchParams();
  const params = adapter.transformQueryParams
    ? adapter.transformQueryParams(rawParams)
    : normalizeParams(rawParams);

  // 2. Fetch data (inside RSC or useSWR in CSR part)
  const { data, error } = useServerData(() => fetcher(params));

  // 3. Handle loading, error, empty states
  if (error) return <ErrorBoundary error={error} />;
  if (!data) return <LoadingSpinner />;
  if (data.hits.length === 0 && adapter.getEmptyState) {
    return adapter.getEmptyState();
  }

  // 4. Map hits to cards
  const cards = data.hits.map(adapter.mapHitToCard);

  // 5. Render UI
  return (
    <>
      <PageHeader title={title} />
      <SortControls fields={params.sortableFields} />
      <CardGrid items={cards} />
      <Pagination total={data.total} page={params.page} size={params.size} />
    </>
  );
}
```

This template handles:
- Pagination controls
- Sorting UI
- SSR data-fetch invocation
- Error, loading, and empty states

## SSR Integration

Re-use the pattern from `app/(application)/search/page.tsx`:
- Accept `searchParams` as RSC props.
- Call `await fetcher(...)` at the top level of the server component.
- Pass data to the client component (`GenericSearchPage`) for rendering.
- Use `useServerData` or a similar hook to revalidate or client-navigate.

### Simplifications / Generalizations

- **Unified Query Param Shape**: Normalize `page`, `size`, `sort`, and `filter` across adapters.
- **Centralized Pagination Component**: One shared `Pagination` driven by `total`, `page`, `size`.
- **Shared Sorting Controls**: Accept a list of fields/directions from adapters.
- **Generic Error Boundary**: Wrap template to catch and display fetch errors.

## Adapter API Details

```ts
interface SearchAdapter<Hit> {
  /** Map a hit to a card props */
  mapHitToCard(hit: Hit): CardProps;

  /** Optional: Transform raw query params */
  transformQueryParams?(raw: Record<string, any>): SearchParams;

  /** Optional: Custom empty state UI */
  getEmptyState?(): React.ReactNode;
}
```

### Example: Conversation Lookup Adapter

```ts
export const conversationLookupAdapter: SearchAdapter<ConversationHit> = {
  mapHitToCard(hit) {
    return {
      title: hit.subject,
      subtitle: hit.participants.join(", "),
      link: `/conversations/${hit.id}`,
      details: hit.updatedAt,
    };
  },
  transformQueryParams(raw) {
    return {
      q: raw.q as string,
      page: Number(raw.page) || 1,
      size: Number(raw.size) || 20,
      sort: raw.sort || 'updatedAt:desc',
      sortableFields: ['updatedAt', 'createdAt'],
    };
  },
  getEmptyState() {
    return <div>No conversations found.</div>;
  },
};
```

## Migration Steps

1. **Create GenericSearchPage** template.
2. **Implement Adapters** for FileSearch, ConversationLookup, OrganizationLookup.
3. **Wire up Pages** to use `GenericSearchPage` with respective adapter and fetcher.
4. **Extract Shared Components**: `Pagination`, `SortControls`, `ErrorBoundary`.
5. **Validate SSR Flow** and client interactivity.

---

This architecture provides a unified, SSR-friendly search platform and isolates page-specific logic into adapters. New search pages can be added by defining a fetcher and an adapter, then wiring them into `GenericSearchPage` with minimal boilerplate.
