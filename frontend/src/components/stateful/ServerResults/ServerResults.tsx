import Card, { CardSize } from "@/components/style/cards/GenericResultCard";
import ErrorMessage from "@/components/style/messages/ErrorMessage";
import LoadingSpinner from "@/components/style/misc/LoadingSpinner";
import ServerSearchResultsRaw from "@/components/style/RawPages/RawServerSearchResults";
import { GenericSearchType, searchInvokeFromUrlParams } from "@/lib/adapters/genericSearchCallback";
import { DEFAULT_PAGE_SIZE } from "@/lib/constants";
import { TypedUrlParams, encodeUrlParams } from "@/lib/types/url_params";
import clsx from "clsx";
import Link from "next/link";
import { Suspense } from "react";

interface ServerSearchResultProps {
  baseUrl: string;
  objectType: GenericSearchType;
  urlParams: TypedUrlParams;
  inherentRouteFilters?: Record<string, string>
}


export default async function ServerSearchResults(params: ServerSearchResultProps) {
  return <Suspense fallback={<LoadingSpinner loadingText="Loading Server Results" />}>
    <ServerSearchResultsUnsuspended {...params} />
  </Suspense>
}

export async function ServerSearchResultsUnsuspended(params: ServerSearchResultProps) {

  try {
    const cardResults = await searchInvokeFromUrlParams(params.urlParams, params.objectType, params.inherentRouteFilters || {});
    return <ServerSearchResultsRaw baseUrl={params.baseUrl} urlParams={params.urlParams} results={cardResults} />
  } catch (err: any) {
    throw err
    return <ErrorMessage error={JSON.stringify(err)} />
  }

}

