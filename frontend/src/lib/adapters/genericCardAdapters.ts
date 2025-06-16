import React from "react";

import type { Filing } from "@/lib/types/FilingTypes";
import type { OrganizationInfo } from "@/lib/requests/organizations";
import type { Conversation } from "@/lib/conversations";
import {
  AuthorCardData,
  CardData,
  CardType,
  DocketCardData,
  DocumentCardData,
} from "../types/generic_card_types";
/**
 * Adapter interface for mapping search hits to card data and handling query params.
 */
export interface SearchAdapter<Hit> {
  /** Map a hit to card props */
  mapHitToCard(hit: Hit): CardData;

  /** Optional: Transform raw query params from URL to fetcher params */
  transformQueryParams?(raw: Record<string, any>): Record<string, any>;

  /** Optional: Provide empty state when no hits are returned */
  getEmptyState?(): React.ReactNode;
}


/**
 * Helper: Format ISO date string to human-readable date.
 */
function formatDate(isoString: string | undefined): string {
  if (!isoString) return "";
  try {
    return new Date(isoString).toLocaleDateString();
  } catch {
    return isoString;
  }
}

/**
 * Adapter: Filing → DocumentCardData
 */
export function adaptFilingToCard(filing: Filing): DocumentCardData {
  return {
    type: CardType.Document,
    index: 0,
    name: filing.title,
    description: filing.file_class || "",
    timestamp: formatDate(filing.date),
    authors: filing.authors_information || [],
    extraInfo: filing.docket_id ? `Docket: ${filing.docket_id}` : undefined,
  };
}

/**
 * Adapter: OrganizationInfo → DocketCardData
 */
export function adaptOrganizationToCard(org: OrganizationInfo): DocketCardData {
  return {
    type: CardType.Docket,
    index: 0,
    name: org.name || org.title || "",
    description: org.description || "",
    timestamp: formatDate(org.updated_at || org.created_at || org.createdAt),
    extraInfo: org.location || org.address,
  };
}

/**
 * Adapter: Conversation → AuthorCardData
 */
export function adaptConversationToCard(convo: Conversation): AuthorCardData {
  return {
    type: CardType.Author,
    index: 0,
    name: convo.name || convo.id,
    description: convo.description || "",
    timestamp: formatDate(convo.updated_at || (convo as any).last_active_at),
    // authors: (convo as any).participants || undefined,
    extraInfo: convo.docket_id ? `Docket: ${convo.docket_id}` : undefined,
  };
}

/**
 * Unified adapter for mixed result arrays
 */
export function adaptResult(item: any, type: CardType): CardData {
  switch (type) {
    case CardType.Document:
      return adaptFilingToCard(item as Filing);
    case CardType.Author:
      return adaptOrganizationToCard(item as OrganizationInfo);
    case CardType.Docket:
      return adaptConversationToCard(item as Conversation);
    default:
      throw new Error(`Unsupported result type: ${type}`);
  }
}
