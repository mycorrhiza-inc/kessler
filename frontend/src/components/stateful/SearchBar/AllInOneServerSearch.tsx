import { encodeUrlParams, TypedUrlParams, UrlQueryParams } from "@/lib/types/url_params"
import { ObjectQueryType } from "@/lib/types/SearchTypes";
import { GenericSearchType } from "@/lib/adapters/genericSearchCallback";
import ServerSearchResults from "../ServerResults/ServerResults";
import ClientSearchBar from "./ClientSearchBar";
import HardcodedFileFilters from "../Filters/StaticFilingFilters";
import { ReactNode } from "react";

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
  aboveSearchElement?: ReactNode,
}


export default async function AllInOneServerSearch({ urlParams, queryType, baseUrl, disableFilterSelection, disableResults, inherentRouteFilters, aboveSearchElement }: AllInOneServerSearchParams) {
  if (disableFilterSelection) {
    return (
      <div className="bg-base-100">
        <ClientSearchBar urlParams={urlParams.queryData} baseUrl={baseUrl} />
        {!disableResults && (
          <ServerSearchResults
            baseUrl={baseUrl}
            urlParams={urlParams}
            objectType={GenericSearchType.Filling}
            inherentRouteFilters={inherentRouteFilters}
          />
        )}
      </div>
    );
  }

  return (
    <div className="bg-base-100 min-w-screen flex">
      <div className="flex flex-col flex-grow items-center">
        {aboveSearchElement}
        <ClientSearchBar urlParams={urlParams.queryData} baseUrl={baseUrl} />
        {!disableResults && (
          <div className="w-full max-w-4xl">
            <ServerSearchResults
              baseUrl={baseUrl}
              urlParams={urlParams}
              objectType={GenericSearchType.Filling}
              inherentRouteFilters={inherentRouteFilters}
            />
          </div>
        )}
      </div>
      <div className="w-64">
        <HardcodedFileFilters urlParams={urlParams.queryData} baseUrl={baseUrl} />
      </div>
    </div>
  );
}
