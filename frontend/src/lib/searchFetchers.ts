import axios from "axios";
import { getSearchResults } from "@/lib/requests/search";
import {
  ConversationTableSchema,
  InstanitateConversationSearchSchema,
  conversationSearchGet,
} from "@/components/LookupPages/ConvoLookup/ConversationTable";
import {
  OrganizationSearchSchema,
  InstanitateOrganizationSearchSchema,
  organizationsListGet,
} from "@/components/LookupPages/OrgLookup/OrganizationTable";
import { Filing } from "@/lib/types/FilingTypes";

// Common search parameter interface
export interface SearchParams {
  q?: string;
  page: number;
  size: number;
  filters?: any;
  sort?: string;
}

// Common response interface
export interface SearchResponse<Hit> {
  hits: Hit[];
  total: number;
  page: number;
  size: number;
}

/**
 * Fetch file search results using existing getSearchResults.
 */
export async function fileSearchFetcher(
  params: SearchParams,
): Promise<SearchResponse<Filing>> {
  const { page, size, q, filters } = params;
  // Assuming getSearchResults returns filings for given page and size
  const hits = await getSearchResults(
    { query: q || "", filters: filters || {} },
    page - 1,
    size,
  );
  // Placeholder total: cannot retrieve total from legacy API
  const total =
    hits.length === size ? page * size + 1 : (page - 1) * size + hits.length;
  return { hits, total, page, size };
}

/**
 * Fetch conversation lookup results.
 */
export async function conversationLookupFetcher(
  params: SearchParams,
): Promise<SearchResponse<ConversationTableSchema>> {
  const { page, size } = params;
  const offset = (page - 1) * size;
  // Use existing schema instantiation and fetch
  const searchData = InstanitateConversationSearchSchema(params.filters);
  const url = `${process.env.NEXT_PUBLIC_API_URL}/v2/search/conversation?offset=${offset}&limit=${size}`;
  const hits = await conversationSearchGet(searchData, url);
  const total =
    hits.length === size ? page * size + 1 : (page - 1) * size + hits.length;
  return { hits, total, page, size };
}

/**
 * Fetch organization lookup results.
 */
export async function organizationLookupFetcher(params: SearchParams): Promise<
  SearchResponse<
    OrganizationSearchSchema & {
      name: string;
      id: string;
      aliases: string[];
      files_authored_count: number;
    }
  >
> {
  const { page, size } = params;
  const offset = (page - 1) * size;
  const searchData = InstanitateOrganizationSearchSchema(params.filters);
  const url = `${process.env.NEXT_PUBLIC_API_URL}/v2/search/organization?offset=${offset}&limit=${size}`;
  const hits = await organizationsListGet(url, searchData);
  const total =
    hits.length === size ? page * size + 1 : (page - 1) * size + hits.length;
  return { hits, total, page, size };
}
