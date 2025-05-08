# Legacy API Documentation

This document outlines the existing (legacy) search architecture and API endpoints in the `/frontend/src/lib/requests` directory and describes how their responses can be converted into a unified format consumable by `GenericResultCard.tsx`.

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
- **name**: Use `response.title` or a generated name from participants
- **description**: `response.last_message.text`
- **timestamp**: `response.last_message.timestamp`
- **extraInfo**: List of `response.participants` or additional metadata

```ts
// Example mapping
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
// Example mapping
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
    // ... potentially other metadata
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

## 4. Next Steps for Migration
- **Unify schemas**: Define a single TypeScript interface for the legacy endpoints that matches `CardData`.
- **Centralize transformation**: Move the mapping logic into a shared utility (`frontend/src/lib/requests/transformers.ts`).
- **Expand coverage**: Identify additional legacy endpoints (e.g. users, tags) and document mappings.
- **Update docs**: Once migrated, archive or remove legacy code and update `docs/frontend_search_current_design.md` to reference the new architecture.

---

*This document is a starting point for consolidating legacy search logic into the new search architecture.*