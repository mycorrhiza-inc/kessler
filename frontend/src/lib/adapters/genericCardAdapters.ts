import React from "react";
import { Filing } from "@/lib/types/FilingTypes";
import { ConversationTableSchema } from "@/components/LookupPages/ConvoLookup/ConversationTable";
import { OrganizationTableSchema } from "@/components/LookupPages/OrgLookup/OrganizationTable";
import { CardProps } from "@/components/Search/CardGrid";
import { SearchParams } from "@/lib/searchFetchers";

export interface SearchAdapter<Hit> {
  /** Map a single hit into card props for rendering */
  mapHitToCard: (hit: Hit) => CardProps;

  /** Optional: Provide an empty state when no hits are returned */
  getEmptyState?: () => React.ReactNode;

  /** Optional: Transform raw query params to fetcher-compatible shape */
  transformQueryParams?: (rawParams: Record<string, any>) => SearchParams;
}

/** File search adapter */
export const fileSearchAdapter: SearchAdapter<Filing> = {
  mapHitToCard: (hit) => ({
    title: hit.title,
    subtitle: hit.author,
    link: `/files/${hit.id}`,
    details: hit.date,
  }),
  transformQueryParams: (raw) => ({
    q: raw.q as string,
    filters: raw.filters,
    page: Number(raw.page) || 1,
    size: Number(raw.size) || 20,
  }),
  getEmptyState: () => <div>No files found.</div>,
};

/** Conversation lookup adapter */
export const conversationLookupAdapter: SearchAdapter<ConversationTableSchema> = {
  mapHitToCard: (hit) => ({
    title: hit.name,
    subtitle: `ID: ${hit.docket_gov_id}`,
    link: `/dockets/${hit.docket_gov_id}`,
    details: JSON.parse(hit.metadata).date_filed,
  }),
  transformQueryParams: (raw) => ({
    page: Number(raw.page) || 1,
    size: Number(raw.size) || 20,
    filters: raw,
  }),
  getEmptyState: () => <div>No dockets found.</div>,
};

/** Organization lookup adapter */
export const organizationLookupAdapter: SearchAdapter<OrganizationTableSchema> = {
  mapHitToCard: (hit) => ({
    title: hit.name,
    subtitle: `Aliases: ${hit.aliases.join(", ")}`,
    link: `/orgs/${hit.id}`,
    details: `Documents: ${hit.files_authored_count}`,
  }),
  transformQueryParams: (raw) => ({
    q: raw.q as string,
    page: Number(raw.page) || 1,
    size: Number(raw.size) || 20,
    filters: raw,
  }),
  getEmptyState: () => <div>No organizations found.</div>,
};
