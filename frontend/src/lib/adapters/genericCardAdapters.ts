import { CardData } from "@/components/NewSearch/GenericResultCard";
import type { Filing } from "@/lib/types/FilingTypes";
import type { OrganizationInfo } from "@/lib/requests/organizations";
import type { Conversation } from "@/lib/conversations";
import {
  AuthorCardData,
  CardType,
  DocketCardData,
  DocumentCardData,
} from "../types/generic_card_types";

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
  const author_strs =
    filing.authors_information?.map((a) =>
      typeof a === "string" ? a : a.author_name || String(a),
    ) || [];
  return {
    type: CardType.Document,
    name: filing.title,
    description: filing.file_class || "",
    timestamp: formatDate(filing.date),
    authors: author_strs,
    extraInfo: filing.docket_id ? `Docket: ${filing.docket_id}` : undefined,
  };
}

/**
 * Adapter: OrganizationInfo → DocketCardData
 */
export function adaptOrganizationToCard(org: OrganizationInfo): DocketCardData {
  return {
    type: CardType.Docket,
    name: org.name || org.title || "",
    description: org.description || "",
    timestamp: formatDate(org.updated_at || org.created_at || org.createdAt),
    authors: Array.isArray(org.admins)
      ? org.admins.map((u: any) => u.username || u.id || String(u))
      : undefined,
    extraInfo: org.location || org.address,
  };
}

/**
 * Adapter: Conversation → AuthorCardData
 */
export function adaptConversationToCard(convo: Conversation): AuthorCardData {
  return {
    type: CardType.Author,
    name: convo.name || convo.id,
    description: convo.description || "",
    timestamp: formatDate(convo.updated_at || (convo as any).last_active_at),
    authors: (convo as any).participants || undefined,
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
