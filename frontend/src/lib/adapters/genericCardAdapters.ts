import {
  CardData,
  AuthorCardData,
  DocketCardData,
  DocumentCardData,
} from "@/components/NewSearch/GenericResultCard";
import type { Filing } from "@/lib/types/FilingTypes";
import type { OrganizationInfo } from "@/lib/requests/organizations";
import type { Conversation } from "@/lib/conversations";

/**
 * Helper: Format ISO date string to human-readable date.
 */
function formatDate(isoString: string): string {
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
    type: "document",
    name: filing.title,
    description: filing.file_class || "",
    timestamp: formatDate(filing.date),
    authors: filing.authors_information
      ? [
          filing.author,
          ...filing.authors_information.map((a) =>
            typeof a === "string" ? a : a.username || String(a),
          ),
        ]
      : [filing.author],
    extraInfo: filing.docket_id ? `Docket: ${filing.docket_id}` : undefined,
  };
}

/**
 * Adapter: OrganizationInfo → DocketCardData
 */
export function adaptOrganizationToCard(org: OrganizationInfo): DocketCardData {
  return {
    type: "docket",
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
    type: "author",
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
export function adaptResult(
  item: any,
  type: "filing" | "organization" | "conversation",
): CardData {
  switch (type) {
    case "filing":
      return adaptFilingToCard(item as Filing);
    case "organization":
      return adaptOrganizationToCard(item as OrganizationInfo);
    case "conversation":
      return adaptConversationToCard(item as Conversation);
    default:
      throw new Error(`Unsupported result type: ${type}`);
  }
}
