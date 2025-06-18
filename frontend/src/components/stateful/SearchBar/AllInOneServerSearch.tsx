import { encodeUrlParams, TypedUrlParams, UrlQueryParams } from "@/lib/types/url_params"
import { ObjectQueryType } from "@/lib/types/SearchTypes";
import { GenericSearchType } from "@/lib/adapters/genericSearchCallback";
import ServerSearchResults from "../ServerResults/ServerResults";
import ClientSearchBar from "./ClientSearchBar";
import HardcodedFileFilters from "../Filters/StaticFilingFilters";

export interface AIOSearchProps {
  urlParams: UrlQueryParams
}

export interface AllInOneServerSearchParams {
  urlParams: TypedUrlParams,
  queryType: GenericSearchType,
  baseUrl: string,
  inherentRouteFilters?: Record<string, string>
  disableFilterSelection?: boolean,
  disableResults?: boolean,
}


export default async function AllInOneServerSearch({ urlParams, queryType, baseUrl, disableFilterSelection, disableResults, inherentRouteFilters }: AllInOneServerSearchParams) {
  return <div>
    <ClientSearchBar urlParams={urlParams.queryData} baseUrl={baseUrl} />
    {!disableFilterSelection && <HardcodedFileFilters urlParams={urlParams.queryData} baseUrl={baseUrl} />}

    {!disableResults && <ServerSearchResults
      baseUrl={baseUrl}
      urlParams={urlParams}
      objectType={GenericSearchType.Filling}
      inherentRouteFilters={inherentRouteFilters}
    />}
  </div>
}
