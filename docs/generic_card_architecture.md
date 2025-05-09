# Generic Result Card Adapter Architecture

This document outlines the proposed adapter layer to convert API endpoint responses into the standardized `CardData` shape consumed by the `GenericResultCard` component.

---

## Goals

- **Decouple** raw API response formats from the UI component.
- **Normalize** data from diverse endpoints (conversations, organizations, filings).
- **Reuse** a single `CardData` interface for all result types.
- **Simplify** UI rendering logic by centralizing transformation.

---

## 1. Export `CardData` Types from Component

Update `src/components/NewSearch/GenericResultCard.tsx` to export the data types:

```ts
export type { CardData, CardSize, AuthorCardData, DocketCardData, DocumentCardData };
```

This allows adapters to import the exact interfaces.

---

## 2. Create Adapter Module

**File:** `src/lib/adapters/genericCardAdapters.ts`

```ts
import {
  CardData,
  AuthorCardData,
  DocketCardData,
  DocumentCardData,
} from "@/components/NewSearch/GenericResultCard";
import type { Filing } from "@/lib/types/FilingTypes";
import type { OrganizationInfo } from "@/lib/requests/organizations";
import type { ConversationType } from "@/lib/types/ConversationTypes";

/**
 * Adapter: Filing → DocumentCardData
 */
export function adaptFilingToCard(filing: Filing): DocumentCardData {
  return {
    type: "document",
    name: filing.title,
    description: filing.file_class,
    timestamp: formatDate(filing.date),
    authors: [filing.author, ...(filing.authors_information || [])],
    extraInfo: `Docket: ${filing.docket_id}`,
  };
}

/**
 * Adapter: Organization → DocketCardData
 */
export function adaptOrganizationToCard(
  org: OrganizationInfo
): DocketCardData {
  return {
    type: "docket",
    name: org.name,
    description: org.description || "",
    timestamp: formatDate(org.updated_at || org.created_at),
    authors: org.admins?.map((u) => u.username),
    extraInfo: org.location,
  };
}

/**
 * Adapter: Conversation → AuthorCardData
 */
export function adaptConversationToCard(
  convo: ConversationType
): AuthorCardData {
  return {
    type: "author",
    name: convo.title || convo.id,
    description: convo.subject || "Conversation thread",
    timestamp: formatDate(convo.last_active_at),
    authors: convo.participants,
    extraInfo: convo.peek || undefined,
  };
}

/**
 * Helper: ISO → "MMM DD, YYYY" display
 */
function formatDate(isoString: string): string {
  return new Date(isoString).toLocaleDateString();
}

/**
 * Unified adapter for mixed search results
 */
export function adaptResult(
  item: any,
  type: "filing" | "organization" | "conversation"
): CardData {
  switch (type) {
    case "filing":
      return adaptFilingToCard(item as Filing);
    case "organization":
      return adaptOrganizationToCard(item as OrganizationInfo);
    case "conversation":
      return adaptConversationToCard(item as ConversationType);
  }
}
```

---

## 3. Use Adapters in UI

1. **Search Hook / Page** invokes the relevant request:
   - `getSearchResults` → returns `Filing[]` → map via `adaptFilingToCard`.
   - `getOrganizationInfo` → returns `OrganizationInfo` → map via `adaptOrganizationToCard`.
   - `GetConversationInformation` → returns `ConversationType` → map via `adaptConversationToCard`.

2. **Rendering**:

```tsx
import Card, { CardSize } from "@/components/NewSearch/GenericResultCard";
import { adaptResult } from "@/lib/adapters/genericCardAdapters";

const results = [ /* heterogeneous API data */ ];

return (
  <>
    {results.map(({ item, type }) => {
      const data = adaptResult(item, type);
      return <Card key={data.name} data={data} size={CardSize.Medium} />;
    })}
  </>
);
```

---

## 4. Directory Structure

```
src/
├── components/
│   └── NewSearch/
│       └── GenericResultCard.tsx
├── lib/
│   ├── adapters/
│   │   └── genericCardAdapters.ts
│   ├── requests/
│   │   ├── conversations.ts
│   │   ├── organizations.ts
│   │   └── search.ts
│   └── types/
│       └── ConversationTypes.ts
└── pages/
    └── search.tsx
```

---

## Benefits

- **Single source** for UI shape transformations.
- **Type-safe** mapping using TS interfaces.
- **Easily extendable** for new result types.
- **Clear separation** between data-fetching and presentation.

---

_Last updated: 2025-05-08_