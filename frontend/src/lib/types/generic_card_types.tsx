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
    }).optional()
  ),
  conversation: z.object({
    convo_id: z.string().uuid(),
    convo_name: z.string(),
    convo_num: z.string()
  }).optional()
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
