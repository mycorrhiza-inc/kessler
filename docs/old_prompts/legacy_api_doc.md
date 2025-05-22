# Legacy API Documentation

This document outlines the existing (legacy) search architecture and API endpoints in the `/frontend/src/lib/requests` directory and other legacy API-dependent modules. It describes how their responses can be converted into a unified format consumable by `GenericResultCard.tsx` or their respective UI components.

---

## 1. Conversations Endpoint

**Source file**: `frontend/src/lib/requests/conversations.ts`

### Endpoint
```
GET {PUBLIC_API_URL}/v2/public/conversations/{conversation_id}
```

### Request
- Path parameter:
  - `conversation_id` (string): Unique identifier for a conversation

### Response Shape
```json
{
  "id": "string",
  "title": "string",
  "participants": ["user1", "user2", ...],
  "last_message": {
    "text": "string",
    "timestamp": "ISO8601 string"
  },
  // ... other conversation metadata
}
```

### Consumption & Conversion
The helper `GetConversationInformation(conversation_id)` returns the parsed JSON. To display in `GenericResultCard`:
- **CardType**: `Author` or `Docket` (depending on UI context)
- **name**: `response.title` or a generated name from participants
- **description**: `response.last_message.text`
- **timestamp**: `response.last_message.timestamp`
- **extraInfo**: List of `response.participants` or additional metadata

```ts
const data: AuthorCardData = {
  type: CardType.Author,
  name: conversation.title,
  description: conversation.last_message.text,
  timestamp: conversation.last_message.timestamp,
  extraInfo: conversation.participants.join(', '),
};
```

---

## 2. Organizations Endpoint

**Source file**: `frontend/src/lib/requests/organizations.ts`

### Endpoint
```
GET {INTERNAL_API_URL}/v2/public/organizations/{orgID}
```

### Request
- Path parameter:
  - `orgID` (string): Unique identifier for an organization

### Response Shape
```json
{
  "id": "string",
  "name": "string",
  "description": "string",
  "created_at": "ISO8601 string",
  "owner": "string",
  // ... other org metadata
}
```

### Consumption & Conversion
The helper `getOrganizationInfo(orgID)` returns the response JSON. To display as a card:
- **CardType**: `Docket` (or `Author` based on UI semantics)
- **name**: `response.name`
- **description**: `response.description`
- **timestamp**: `response.created_at`
- **extraInfo**: `response.owner` or other fields

```ts
const data: DocketCardData = {
  type: CardType.Docket,
  name: org.name,
  description: org.description,
  timestamp: org.created_at,
  extraInfo: `Owner: ${org.owner}`,
};
```

---

## 3. Search File Endpoints

**Source file**: `frontend/src/lib/requests/search.ts`

### Endpoints
- **Search files** (POST):
  ```
  POST {PUBLIC_API_URL}/v2/search/file?{pageQuery}
  Body: {
    query: string,
    filters: BackendFilterObject
  }
  ```

- **Recent filings** (GET):
  ```
  GET {PUBLIC_API_URL}/v2/search/file/recent_updates?{pageQuery}
  ```

### Response Shape
The POST and GET endpoints return an array of _hydrated search results_:
```json
[
  {
    "file": { /* CompleteFileSchema */ },
    "score": number,
    // ... other metadata
  },
  // ... more results
]
```

#### `CompleteFileSchema` (validated by `CompleteFileSchemaValidator`)
Key fields:
- `id`: string
- `name`: string
- `extension`: string (e.g., `.pdf`)
- `lang`: string
- `mdata`: object containing
  - `docID`, `date`, `author`, `item_number`, `file_class`, `docket_id`, `url`
- `authors`: array of detailed author objects

### Conversion Pipeline
1. **Validation**: `CompleteFileSchemaValidator.parse(maybe_file)`
2. **Mapping to `Filing`**:
   ```ts
   interface Filing {
     id: string;
     title: string;
     source: string;
     lang: string;
     date: string;
     author: string;
     authors_information: any[];
     item_number: string;
     file_class: string;
     docket_id: string;
     url: string;
     extension: string;
   }
   
   function generateFilingFromFileSchema(file_schema: CompleteFileSchema): Filing {
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
   }
   ```
3. **Final translation to `CardData`** (`Document` type):
   ```ts
   const data: DocumentCardData = {
     type: CardType.Document,
     name: filing.title,
     description: `Class: ${filing.file_class} • Ext: ${filing.extension}`,
     timestamp: filing.date,
     authors: filing.authors_information.map(a => a.name || a),
     extraInfo: `Docket: ${filing.docket_id}`,
   };
   ```

---

## 4. Lookup Search Endpoints

### 4.1 Organization Lookup

**Source file**: `frontend/src/components/LookupPages/OrgLookup/OrganizationTable.tsx`

#### Endpoint
```
POST {PUBLIC_API_URL}/v2/search/organization?offset={offset}&limit={limit}
Body: { query?: string }
```

#### Response Shape
```json
[
  {
    "name": string,
    "id": string,
    "aliases": string[],
    "files_authored_count": number
  },
  // ... more orgs
]
```

#### Consumption
- Parsed and validated minimally in `organizationsListGet`.
- Rendered in a table with infinite scroll, linking to `/orgs/{id}`.

### 4.2 Conversation Lookup

**Source file**: `frontend/src/components/LookupPages/ConvoLookup/ConversationTable.tsx`

#### Endpoint
```
POST {PUBLIC_API_URL}/v2/search/conversation?offset={offset}&limit={limit}
Body: { query?: string, industry_type?: string, date_from?: string, date_to?: string }
```

#### Response Shape
```json
[
  {
    "docket_gov_id": string,
    "state": string,
    "name": string,
    "description": string,
    "matter_type": string,
    "industry_type": string,
    "metadata": string,        // JSON-stringified metadata
    "extra": string,
    "documents_count": number,
    "date_published": string,
    "id": string
  },
  // ... more conversations
]
```

#### Consumption
- Parsed in `conversationSearchGet`.
- Extracts nested `metadata` JSON fields (`date_filed`, `matter_type`).
- Rendered in table rows with infinite scroll, navigating to `/dockets/{docket_gov_id}`.

---

## 5. Bookmark Endpoints

**Source file**: `frontend/src/lib/bookmark.ts`

### Endpoint (Local Quickwit Service)
```
POST http://localhost:4041/bookmarks/
Body: { id: string }
```

### Note
- Currently throws an error indicating planned refactor to backend Go API.
- On success, sets `this.title` from response.

---

## 6. Chat Endpoints

**Source file**: `frontend/src/lib/chat.ts`

### Endpoints
- **Create new chat** (POST): `/api/chat/new`
- **Load chat history** (GET): `/api/chat/?id={chatId}`
- **Send message & RAG** (POST): `/api/chat/?id={chatId}`
  ```json
  {
    model: string,
    chat_history: Message[],
    filters?: BackendFilterObject
  }
  ```

### Response Shape
```json
{
  message: {
    content: string,
    citations: any[]
  }
}
```

### Consumption
- Methods on `ChatLog` class (`new`, `loadLog`, `sendMessage`).
- `getUpdatedChatHistory` merges RAG filters, fetches, then appends new `Message`.

---

## 7. Document Service Endpoints

**Source file**: `frontend/src/lib/document.ts`

### Endpoints
- **Get metadata** (GET): `/api/v1/files/metadata/{docid}`
- **Get markdown text** (GET): `/api/v1/files/markdown/{docid}`
- **Raw PDF URL**: `/api/v1/files/raw/{docid}` (client-side URL)

### Response Shapes
- **Metadata**: arbitrary JSON assigned to `Document.docMetadata`
- **Markdown**: raw text assigned to `Document.docText`

### Consumption
- `Document` class methods (`getDocumentMetadata`, `getDocumentText`, `getPdfUrl`, `loadDocument`).

---

## 8. Next Steps for Migration
- **Unify schemas**: Define a single TypeScript interface for all legacy endpoints matching `CardData` or service models.
- **Centralize transformation**: Move mapping logic into shared utilities (e.g. `frontend/src/lib/requests/transformers.ts`).
- **Expand coverage**: Ensure all legacy endpoints (search, lookup, chat, bookmarks, document) are included in migration.
- **Update docs**: Once migrated, archive or remove legacy code and update `docs/frontend_search_current_design.md` to reference the new architecture.

---

*This document provides a starting point for consolidating legacy service calls into the new unified API client and UI mapping layers.*
