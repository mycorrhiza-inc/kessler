import { encodeUrlParams, TypedUrlParams } from "@/lib/types/url_params";
import Card, { CardSize } from "../cards/GenericResultCard";
import { DEFAULT_PAGE_SIZE } from "@/lib/constants";
import clsx from "clsx";

interface ServerSearchResultsRawParams {
  baseUrl: string;
  urlParams: TypedUrlParams;
  results: any[]
}
export default function ServerSearchResultsRaw({ baseUrl, urlParams, results }: ServerSearchResultsRawParams) {
  // Perform search based on URL params
  const cardElements = results.map((card_data, index) => (
    <Card key={`server-results-${card_data.id}-${index}`} size={CardSize.Medium} data={card_data} />
  ));

  // Pagination logic
  const currentPage = urlParams.paginationData.page || 0;
  const limit = urlParams.paginationData.limit || DEFAULT_PAGE_SIZE;
  // If fewer results than limit, we are on the last page
  const isLastPage = results.length < limit;

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
          {/* Using an a tag instead of a Link tag, because Link tags wait for the entire page to load before doing a transition even with stuff thats behind a suspense. */}
          {/* It does seem like there is a solution to this problem described here using some magic routing flags, I should implement this, but not now: https://medium.com/@amirilovic/how-to-use-react-suspense-for-data-loading-4b68f9200c19 */}
          <a
            href={prevHref}
            className={
              clsx("join-item btn btn-outline text-2xl",
                currentPage === 0 && 'btn-disabled'
              )}
          >
            «
          </a>
          <p className="join-item btn btn-outline text-base pointer-events-none select-none">Page {(urlParams.paginationData.page || 0) + 1}</p>
          <a
            href={nextHref}
            className={
              clsx("join-item btn btn-outline text-2xl",
                isLastPage && 'btn-disabled'
              )}
          >
            »
          </a>
        </div>
      </div>
    </div >
  );
}
