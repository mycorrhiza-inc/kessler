import axios from 'axios';
import { z } from "zod";
import { CardType } from "@/lib/types/generic_card_types";


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

// API response validator
const SearchResponseValidator = z.object({
  results: z.array(CardDataValidator),
  total: z.number(),
  page: z.number(),
});

type SearchResponse = z.infer<typeof SearchResponseValidator>;

// API function with validation
export const searchCards = async (
  query: string,
  page: number = 1,
  limit: number = 10
): Promise<SearchResponse> => {
  try {
    const response = await axios.get('/api/search', {
      params: {
        q: query,
        page,
        limit
      }
    });

    // Validate the response data
    const validatedData = SearchResponseValidator.parse(response.data);
    return validatedData;
  } catch (error) {
    if (error instanceof z.ZodError) {
      console.error('Invalid response format:', error.errors);
      throw new Error('Invalid response format from server');
    }
    console.error('Search request failed:', error);
    throw error;
  }
};

// Utility function to format timestamp
const formatTimestamp = (timestamp: string): string => {
  return new Date(timestamp).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  });
};
