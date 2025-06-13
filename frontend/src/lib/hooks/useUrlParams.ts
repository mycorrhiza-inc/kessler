"use client";
import { useSearchParams, usePathname } from "next/navigation";
import { useMemo } from "react";

export interface UrlParams {
  query?: string;
  dataset?: string;
  filters?: Record<string, string>;
}

/**
 * Parses URL parameters: query, dataset, and filters prefixed by f: (e.g., f:status=active)
 */
export function useUrlParams(): UrlParams {
  const searchParams = useSearchParams();
  const pathname = usePathname();

  const params = useMemo<UrlParams>(() => {
    let query = searchParams.get("query") || "";
    let dataset = searchParams.get("dataset") || "";
    const filters: Record<string, string> = {};
    Array.from(searchParams.entries()).forEach(([key, value]) => {
      if (key.startsWith("f:")) {
        const filterKey = key.substring(2);
        filters[filterKey] = value;
      }
    });
    return { query, dataset, filters };
  }, [searchParams, pathname]);

  return params;
}
