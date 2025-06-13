"use client";
import { useSearchParams, usePathname } from "next/navigation";
import { useMemo } from "react";
import { DEFAULT_PAGE_SIZE } from "../constants";

export interface UrlQueryParams {
  query?: string;
  dataset?: string;
  filters?: Record<string, string>;
}

export interface UrlPaginationParams {
  page: number; // Defaults to zero
  limit: number; // Defaults to DEFAULT_PAGE_SIZE
}

export interface TypedUrlParams {
  queryData: UrlQueryParams;
  paginationData: UrlPaginationParams;
}

/**
 * Parses raw URL search parameters into typed query, dataset, filters, and pagination data.
 * @param untyped_params Raw URL parameters as string or string[]
 * @returns TypedUrlParams object with queryData and paginationData
 */
export function generateTypeUrlParams(
  untyped_params: { [key: string]: string | string[] | undefined }
): TypedUrlParams {
  // Extract known parameters
  const rawQuery = untyped_params["q"];
  const rawDataset = untyped_params["dataset"];
  const rawPage = untyped_params["page"];
  const rawLimit = untyped_params["limit"];

  // Parse single-value params
  const query = Array.isArray(rawQuery) ? rawQuery[0] : rawQuery;
  const dataset = Array.isArray(rawDataset) ? rawDataset[0] : rawDataset;

  // Initialize filters
  const filters: Record<string, string> = {};

  // Extract filter params prefixed with "f:"
  Object.keys(untyped_params).forEach((key) => {
    if (key.startsWith("f:")) {
      const filterKey = key.substring(2);
      const rawVal = untyped_params[key];
      const val = Array.isArray(rawVal) ? rawVal[0] : rawVal;
      if (filterKey && val !== undefined) {
        filters[filterKey] = val;
      }
    }
  });

  // Parse pagination, with defaults
  let page = 0;
  if (rawPage !== undefined) {
    const p = Array.isArray(rawPage) ? rawPage[0] : rawPage;
    const parsed = parseInt(p || "", 10);
    if (!isNaN(parsed) && parsed >= 0) {
      page = parsed;
    }
  }

  let limit = DEFAULT_PAGE_SIZE;
  if (rawLimit !== undefined) {
    const l = Array.isArray(rawLimit) ? rawLimit[0] : rawLimit;
    const parsed = parseInt(l || "", 10);
    if (!isNaN(parsed) && parsed > 0) {
      limit = parsed;
    }
  }

  return {
    queryData: { query, dataset, filters },
    paginationData: { page, limit }
  };
}

/**
 * Encodes typed URL parameters back into a query string.
 * Skips pagination keys if they match default values.
 * @param params Typed URL params
 * @returns Encoded query string starting with `?` or empty string
 */
export function encodeFilterParamaters(params: TypedUrlParams): string {
  const parts: string[] = [];
  const { queryData, paginationData } = params;

  // Query
  if (queryData.query) {
    parts.push(`q=${encodeURIComponent(queryData.query)}`);
  }

  // Dataset
  if (queryData.dataset) {
    parts.push(`dataset=${encodeURIComponent(queryData.dataset)}`);
  }

  // Filters
  if (queryData.filters) {
    Object.entries(queryData.filters).forEach(([key, value]) => {
      if (value !== undefined && value !== "") {
        parts.push(`f:${encodeURIComponent(key)}=${encodeURIComponent(value)}`);
      }
    });
  }

  // Pagination: include only if not default
  if (paginationData.page > 0) {
    parts.push(`page=${paginationData.page}`);
  }
  if (paginationData.limit > DEFAULT_PAGE_SIZE) {
    parts.push(`limit=${paginationData.limit}`);
  }

  if (parts.length === 0) {
    return "";
  }

  return `?${parts.join("&")}`;
}

// /**
//  * React hook to read and parse current URL search parameters.
//  * @returns TypedUrlParams based on the browser's current query string
//  */
// export function useTypedUrlParams(): TypedUrlParams {
//   const searchParams = useSearchParams();
//
//   return useMemo(() => {
//     // Build raw params map
//     const raw: { [key: string]: string | string[] | undefined } = {};
//     // URLSearchParams may have multiple values per key
//     Array.from(searchParams.keys()).forEach((key) => {
//       const values = searchParams.getAll(key);
//       raw[key] = values.length > 1 ? values : values[0];
//     });
//     return generateTypeUrlParams(raw);
//   }, [searchParams]);
// }
//
// /**
//  * React hook to construct a pathname with encoded query params.
//  * @param params TypedUrlParams to encode
//  * @returns Full URL path including current pathname and query string
//  */
// export function useUrlWithParams(params: TypedUrlParams): string {
//   const pathname = usePathname();
//   return useMemo(() => {
//     const queryString = encodeFilterParamaters(params);
//     return `${pathname}${queryString}`;
//   }, [pathname, params]);
// }
