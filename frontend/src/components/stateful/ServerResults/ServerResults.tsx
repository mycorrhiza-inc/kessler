import Card, { CardSize } from "@/components/style/cards/GenericResultCard";
import { GenericSearchType, searchInvokeFromUrlParams } from "@/lib/adapters/genericSearchCallback";
import { DEFAULT_PAGE_SIZE } from "@/lib/constants";
import { TypedUrlParams, encodeUrlParams } from "@/lib/types/url_params";
import Link from "next/link";

interface ServerSearchResultProps {
  baseUrl: string;
  objectType: GenericSearchType;
  urlParams: TypedUrlParams;
}

export default async function ServerSearchResults({ baseUrl, objectType, urlParams }: ServerSearchResultProps) {
  // Perform search based on URL params
  const cardResults = await searchInvokeFromUrlParams(urlParams, objectType);
  const cardElements = cardResults.map((card_data) => (
    <Card key={card_data.id} size={CardSize.Medium} data={card_data} />
  ));

  // Pagination logic
  const currentPage = urlParams.paginationData.page || 0;
  const limit = urlParams.paginationData.limit || DEFAULT_PAGE_SIZE;
  // If fewer results than limit, we are on the last page
  const isLastPage = cardResults.length < limit;

  // Build URLs for previous/next
  const prevPage = Math.max(currentPage - 1, 0);
  const nextPage = currentPage + 1;

  const prevParams: TypedUrlParams = {
    ...urlParams,
    paginationData: { ...urlParams.paginationData, page: prevPage },
  };
  const nextParams: TypedUrlParams = {
    ...urlParams,
    paginationData: { ...urlParams.paginationData, page: nextPage },
  };

  const prevHref = `${baseUrl}${encodeUrlParams(prevParams)}`;
  const nextHref = `${baseUrl}${encodeUrlParams(nextParams)}`;

  return (
    <div>
      <div className="grid grid-cols-1 gap-4">
        {cardElements}
      </div>

      <div className="flex justify-center my-4 space-x-2">
        <Link
          href={prevHref}
          className={`btn btn-outline${currentPage === 0 ? ' btn-disabled' : ''}`}
        >
          Previous
        </Link>
        <p>Page {urlParams.paginationData.page || 0}</p>
        <Link
          href={nextHref}
          className={`btn btn-outline${isLastPage ? ' btn-disabled' : ''}`}
        >
          Next
        </Link>
      </div>
    </div>
  );
}
