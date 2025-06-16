import { z } from "zod";

export enum CardType {
  Author = "author",
  Docket = "docket",
  Document = "document",
}

export const BaseCardDataValidator = z.object({
  name: z.string(),
  object_uuid: z.string().uuid(),
  description: z.string(),
  timestamp: z.string(),
  extraInfo: z.string().optional(),
  index: z.number(),
});

export const AuthorCardDataValidator = BaseCardDataValidator.extend({
  type: z.literal(CardType.Author),
});

export const DocketCardDataValidator = BaseCardDataValidator.extend({
  type: z.literal(CardType.Docket),
});

export const DocumentCardDataValidator = BaseCardDataValidator.extend({
  type: z.literal(CardType.Document),
  authors: z.array(
    z.object({
      author_name: z.string(),
      is_person: z.boolean(),
      is_primary_author: z.boolean(),
      author_id: z.string(),
    })
  ),
});

export const CardDataValidator = z.union([
  AuthorCardDataValidator,
  DocketCardDataValidator,
  DocumentCardDataValidator,
]);

export type AuthorCardData = z.infer<typeof AuthorCardDataValidator>;
export type DocketCardData = z.infer<typeof DocketCardDataValidator>;
export type DocumentCardData = z.infer<typeof DocumentCardDataValidator>;
export type CardData = z.infer<typeof CardDataValidator>;
// In the comments add some examples for the type of json you would want each cardata to be parsed as?

/*
Example JSON for AuthorCardData:
{
  "name": "Jane Doe",
  "description": "Expert in mycology",
  "timestamp": "2025-06-16T12:34:56Z",
  "extraInfo": "PhD in fungal biology",
  "index": 1,
  "type": "author",
  "object_uuid": "550e8400-e29b-41d4-a716-446655440000"
}

Example JSON for DocketCardData:
{
  "name": "Case 1234",
  "description": "Legal docket for case 1234",
  "timestamp": "2025-01-10T09:00:00Z",
  "index": 2,
  "type": "docket",
  "object_uuid": "550e8400-e29b-41d4-a716-446655440001"
}

Example JSON for DocumentCardData:
{
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
}
*/
