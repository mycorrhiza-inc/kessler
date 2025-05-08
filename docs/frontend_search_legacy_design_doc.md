# Legacy Frontend Search Architecture

This document describes the existing (legacy) search system in the Kessler frontend. It focuses on three request modules, their corresponding API endpoints, and how their results are consumed by upper-level search components. The goal is to provide a clear map of the bottom part of the architecture before migrating to the newer, unified search format.

---

## 1. Request Modules and Endpoints

### 1.1 Conversations (`frontend/src/lib/requests/conversations.ts`)

- **Function**: `GetConversationInformation(conversation_id: string)`
- **Endpoint**: `GET ${runtimeConfig.public_api_url}/v2/public/conversations/{conversation_id}`
- **Usage**: Fetches metadata for a single conversation (docket) by ID.
- **Returns**: `conversation.data` (untyped JSON).

```ts
export const GetConversationInformation = async (conversation_id: string) => {
  const runtimeConfig = getRuntimeEnv();
  const conversation = await axios.get(
    `${runtimeConfig.public_api_url}/v2/public/conversations/${conversation_id}`
  );
  return conversation.data;
};
```

---

### 1.2 Organizations (`frontend/src/lib/requests/organizations.ts`)

- **Function**: `getOrganizationInfo(orgID: string)`
- **Endpoint**: `GET ${internalAPIURL}/v2/public/organizations/{orgID}`
- **Usage**: Retrieves organization details and file counts by organization ID.
- **Returns**: `response.data` (currently logged then returned as `any`).

```ts
export const getOrganizationInfo = async (orgID: string) => {
  const response = await axios.get(
    `${internalAPIURL}/v2/public/organizations/${orgID}`
  );
There is currently a bunch of documentation for the newer search architecture in 

kessler/docs/frontend_search_current_design.md

The legacy search system is kinda sprawling throughout the codebase at the moment. But starts at the root of the project with these files that describe how the search requests hit various endpoints.

In the beginning could you just search the project and try to document the bottom part of the search architecture in an attempt to try and move everything into the new format.

kessler/frontend/src/lib/requests/conversations.ts
kessler/frontend/src/lib/requests/organizations.ts
kessler/frontend/src/lib/requests/search.ts

and document how these endpoints eventually end up getting consumed by the upper level components in 

kessler/frontend/src/components/Search/FileSearch/FileSearchView.tsx
kessler/frontend/src/components/LookupPages/ConvoLookup/ConversationTable.tsx
kessler/frontend/src/components/LookupPages/OrgLookup/OrganizationTable.tsx

One big part of this is the necessary conversion of the backend schemas in all these different files into a unified format that could get used by our filter card system. But all of that is a much lower priority, just try to understand and summarize how everythingt works and throw it in that design document.

Could you go ahead and throw all this documentation and throw it into docs/frontend_search_legacy_design_doc.md?
  console.log("organization data", response.data);
  return response.data;
};
```

---

### 1.3 File Search (`frontend/src/lib/requests/search.ts`)

**Primary functions**:

- **`getSearchResults(queryData: QueryDataFile, page: number, maxHits: number): Promise<Filing[]>`**
  - Builds a `BackendFilterObject` via `backendFilterGenerate(queryData.filters)`.
  - Constructs a pagination query string (`?offset=...&limit=...`).
  - `POST` to `${public_api_url}/v2/search/file{paginationQueryString}` with payload `{ query, filters }`.
  - On success, calls `hydratedSearchResultsToFilings(response.data)`.
  - Returns an array of typed `Filing` objects.

- **`hydratedSearchResultsToFilings(hydratedSearchResults: any): Filing[]`**
  - Iterates `hydratedSearchResults`, extracts `file` payload.
  - Validates each file against `CompleteFileSchemaValidator` (Zod).
  - Converts valid schemas into `Filing` via `generateFilingFromFileSchema()`.

- **`getRecentFilings(page?: number, page_size?: number)`**
  - Fetches most recent filings: `GET ${public_api_url}/v2/search/file/recent_updates{queryString}`.
  - Parses and returns `Filing[]` via the same hydration logic.

- **`completeFileSchemaGet(url: string)`**
  - Generic `GET` on an arbitrary URL returning a `CompleteFileSchema`.

- **`generateFilingFromFileSchema(file_schema: CompleteFileSchema): Filing`**
  - Maps backend schema fields (`id`, `name`, `mdata`, etc.) to UI-friendly `Filing`.

```ts
// Example field mapping
return {
  id: file_schema.id,
  title: file_schema.name,
  source: file_schema.mdata.docID,
  lang: file_schema.lang,
  date: file_schema.mdata.date,
  author: file_schema.mdata.author,
  authors_information: file_schema.authors,
  item_number: file_schema.mdata.item_number,
  file_class: file_schema.mdata.file_class,
  docket_id: file_schema.mdata.docket_id,
  url: file_schema.mdata.url,
  extension: file_schema.extension,
};
```

---

## 2. Search UI Components

### 2.1 File Search View (`frontend/src/components/Search/FileSearch/FileSearchView.tsx`)

- **Responsibilities**:
  1. Initialize filter state via `initialFiltersFromInherited(inheritedFilters)`.
  2. Manage `queryData: QueryDataFile` (`{ query, filters }`).
  3. Fetch filings:
     - **Initial load**: `getInitialUpdates()` requests first N pages.
     - **Infinite scroll**: `getMore()` calls `getSearchResults` with updated page and pageSize.
  4. Render:
     - `<SearchBox>` updates `queryData` on submit.
     - `<InfiniteScrollPlus>` wraps `<FilingTable>` to handle pagination UI.

- **Data Flow**: User types ➔ `SearchBox` sets `queryData` ➔ `InfiniteScrollPlus.loadInitial` or `getMore` calls `getSearchResults` ➔ `Filing[]` updates state ➔ `FilingTable` renders rows.

```tsx
// Fetching a page of results
const new_filings = await getSearchResults(queryData, page, pageSize);
setFilings(prev => [...prev, ...new_filings]);
```


---

### 2.2 Conversation Lookup Table (`frontend/src/components/LookupPages/ConvoLookup/ConversationTable.tsx`)

- **Components**:
  - `ConversationTableInfiniteScroll` (default export)
  - `ConversationTable` (renders an HTML table)

- **Search Logic** (`conversationSearchGet`):
  - Normalizes search payload via `InstanitateConversationSearchSchema`.
  - `POST` to `${public_api_url}/v2/search/conversation?offset={offset}&limit={limit}`.
  - Cleans empty responses into `[]`.
  - Returns typed `ConversationTableSchema[]`.

- **Infinite Scroll**:
  - `getPageResults(page, limit)` handles offset, URL, and state updates.
  - `InfiniteScrollPlus` triggers `getInitialData` and `getMore`.

- **UI**:
  - `ConversationTable` maps each `convo` to a table row, extracting fields from JSON-encoded `metadata` for `date_filed` and `matter_type`.
  - Rows are clickable to navigate to `/dockets/{docket_gov_id}`.


---

### 2.3 Organization Lookup Table (`frontend/src/components/LookupPages/OrgLookup/OrganizationTable.tsx`)

- **Components**:
  - `OrganizationTableInfiniteScroll` (default export)
  - `OrganizationTable` (renders rows of org name and file count)

- **Search Logic** (`organizationsListGet`):
  - Uses `InstanitateOrganizationSearchSchema` to fill `{ query?: string }`.
  - `POST` to `${public_api_url}/v2/search/organization?offset={offset}&limit={limit}`.
  - Returns `OrganizationTableSchema[]` (`{ name, id, aliases, files_authored_count }`).

- **Infinite Scroll**:
  - Identical pattern to conversation lookup: `getInitialData`, `getMore`, `InfiniteScrollPlus`.

- **UI**:
  - `OrganizationTable` links each row to `/orgs/{org.id}` and displays `files_authored_count`.

---

## 3. Data Flow Summary

1. **Input**: User enters search terms or filter values.
2. **State Update**: Client component (`SearchBox` or lookup inputs) updates local search state.
3. **API Request**: Infinite scroll or manual trigger calls the respective request function:
   - Files: `getSearchResults`
   - Conversations: `conversationSearchGet`
   - Organizations: `organizationsListGet`
4. **Response Parsing**:
   - Files: Zod-based schema validation + mapping to `Filing` type.
   - Conversations & Orgs: Minimal cleaning; no schema validation.
5. **Render**: Table or list component consumes typed arrays and displays rows.

---

## 4. Notes on Schema Conversion and Filters

- **BackendFilterObject**: Generated by `backendFilterGenerate()` from UI filter fields; drives server-side filtering logic.
- **Zod Validation**: Only the file-search path uses `CompleteFileSchemaValidator` for type safety.
- **Unification Goal**: Migrate all search paths to a shared filter-card system and standardized result type before deprecating these modules.

---

## 5. Next Steps for Migration

1. **Centralize Filter Logic**: Replace per-module instantiation with a shared filter-card API.
2. **Schema Unification**: Introduce a common search result interface for all contexts (files, convos, orgs).
3. **Client Hooks**: Extract infinite scroll + fetch into a generic `useSearch` hook.
4. **URL Sync & SSR**: Adopt the newer search architecture (Current Design) for lookup pages where applicable.

---

*Document generated to assist in migrating legacy code into the new `frontend_search_current_design` format.*
