import Card, { CardSize } from "@/components/style/cards/GenericResultCard";
import LoadingSpinner from "@/components/style/misc/LoadingSpinner";
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

export async function ServerSearchResultsUnsuspended({ baseUrl, objectType, urlParams, inherentRouteFilters }: ServerSearchResultProps) {
  // Perform search based on URL params
  const cardResults = await searchInvokeFromUrlParams(urlParams, objectType, inherentRouteFilters || {});
  const cardElements = cardResults.map((card_data) => (
    <Card key={card_data.id} size={CardSize.Medium} data={card_data} />
  ));

  // Pagination logic
  const currentPage = urlParams.paginationData.page || 0;
  const limit = urlParams.paginationData.limit || DEFAULT_PAGE_SIZE;
  // If fewer results than limit, we are on the last page
  const isLastPage = cardResults.length < limit;

  // Build URLs for previous/next

  const gotoPageHref = (page: number): string => {
    const params: TypedUrlParams = {
      ...urlParams,
      paginationData: { ...urlParams.paginationData, page: page },
    }
    return `${baseUrl}${encodeUrlParams(params)}`
  }

  const prevHref = gotoPageHref(Math.max(currentPage - 1, 0));
  const nextHref = gotoPageHref(currentPage + 1);

  return (
    <div>
      <div className="grid grid-cols-1 gap-4">
        {cardElements}
      </div>
      <div className="flex justify-center my-4 space-x-2">
        <div className="join bg-base-100">
          <Link
            href={prevHref}
            className={
              clsx("join-item btn btn-outline text-2xl",
                currentPage === 0 && 'btn-disabled'
              )}
          >
            «
          </Link>
          <p className="join-item btn btn-outline text-base pointer-events-none select-none">Page {urlParams.paginationData.page || 0}</p>
          <Link
            href={nextHref}
            className={
              clsx("join-item btn btn-outline text-2xl",
                isLastPage && 'btn-disabled'
              )}
          >
            »
          </Link>
        </div>
      </div>
    </div >
  );
}
