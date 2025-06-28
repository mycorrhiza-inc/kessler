
import { z } from "zod";
import {
  AuthorCardDataValidator,
  DocketCardDataValidator,
  DocumentCardDataValidator,
} from "@/lib/types/generic_card_types";
import { describe, it, expect } from "vitest";

describe("CardDataValidators", () => {
  const authorJson = `{
    "name": "Jane Doe",
    "description": "Expert in mycology",
    "timestamp": "2025-06-16T12:34:56Z",
    "extraInfo": "PhD in fungal biology",
    "index": 1,
    "type": "author",
    "object_uuid": "550e8400-e29b-41d4-a716-446655440000"
  }`;

  const docketJson = `{
    "name": "Case 1234",
    "description": "Legal docket for case 1234",
    "timestamp": "2025-01-10T09:00:00Z",
    "index": 2,
    "type": "docket",
    "object_uuid": "550e8400-e29b-41d4-a716-446655440001"
  }`;

  const documentJson = `{
    "name": "Research Paper on Mycorrhizae",
    "description": "Detailed study on mycorrhizal networks",
    "timestamp": "2024-11-20T15:20:30Z",
    "extraInfo": "Published in Nature",
    "index": 3,
    "type": "document",
    "object_uuid": "550e8400-e29b-41d4-a716-446655440002",
    "authors": [
      {
        "author_name": "Jane Doe",
        "is_person": true,
        "is_primary_author": true,
        "author_id": "author-123"
      },
      {
        "author_name": "Research Institute",
        "is_person": false,
        "is_primary_author": false,
        "author_id": "org-456"
      }
    ]
  }`;

  it("parses valid AuthorCardData JSON correctly", () => {
    const parsed = AuthorCardDataValidator.parse(JSON.parse(authorJson));
    expect(parsed.name).toBe("Jane Doe");
    expect(parsed.type).toBe("author");
  });

  it("parses valid DocketCardData JSON correctly", () => {
    const parsed = DocketCardDataValidator.parse(JSON.parse(docketJson));
    expect(parsed.name).toBe("Case 1234");
    expect(parsed.type).toBe("docket");
  });

  it("parses valid DocumentCardData JSON correctly", () => {
    const parsed = DocumentCardDataValidator.parse(JSON.parse(documentJson));
    expect(parsed.name).toBe("Research Paper on Mycorrhizae");
    expect(parsed.type).toBe("document");
    expect(parsed.authors.length).toBe(2);
    expect(parsed.authors[0]?.author_name).toBe("Jane Doe");
  });

  it("fails parsing invalid JSON for AuthorCardData", () => {
    const invalidJson = `{
      "name": "Jane Doe",
      "description": "Expert in mycology",
      "timestamp": "invalid-date",
      "index": 1,
      "type": "author",
      "object_uuid": "not-a-uuid"
    }`;
    expect(() => AuthorCardDataValidator.parse(JSON.parse(invalidJson))).toThrow();
  });
});
