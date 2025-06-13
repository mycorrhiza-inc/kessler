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
  queryData: UrlQueryParams,
  paginationData: UrlPaginationParams

}

// For these two functions I want to encode a 
//
// ?q= must populate the search box
//     a new search should update this for the next page
// ?f:{filterkey}= must populate its respective filter
//     these should be validated, when submitted, by the backend
// ?dataset= must have an indicator of which dataset it is in
//     a new search should update this for the next page
//     the filter struct already has an enabled field and can be greyed out if the previously picked filters do no apply to the new dataset
//     pagination
export function generateTypeUrlParams(untyped_params: { [key: string]: string | string[] | undefined }): TypedUrlParams {

}

export function encodeFilterParamaters(params: TypedUrlParams): string {
  // Skip the pagination encoding if the values equal the default ones.
  return ""
}
