"use client";
import { encodeUrlParams, TypedUrlParams, UrlQueryParams } from "@/lib/types/url_params"
import { ObjectQueryType } from "@/lib/types/SearchTypes";
import { useState } from "react";
import { usePathname, useRouter } from "next/navigation";
import { SearchBar } from "@/components/style/search/SearchBar";

import { KesslerFilterSystem } from '@/components/stateful/Filters/DynamicFilters';
import { makeFilterEndpoints } from '@/lib/filters';
import HardcodedFileFilters from "../Filters/StaticFilingFilters";

export interface AIOSearchProps {
  urlParams: UrlQueryParams
}

export default function AllInOneClientSearch({ urlParams, queryType: pageContext, overrideBaseUrl, disableFilterSelection }: {
  urlParams: UrlQueryParams,
  queryType: ObjectQueryType,
  overrideBaseUrl?: string,
  disableFilterSelection?: boolean,
}) {
  const router = useRouter()
  const impliedUrl: string = usePathname();
  const baseUrl = overrideBaseUrl || impliedUrl || "/search"
  const [query, setQuery] = useState(urlParams.query || "")

  const executeSearch = () => {
    const newQueryParams = { ...urlParams, query };
    const newParams: TypedUrlParams = { queryData: newQueryParams, paginationData: {} }
    const encodedUrlQuery = encodeUrlParams(newParams)
    const url = baseUrl + encodedUrlQuery
    router.push(url)
  }

  if (disableFilterSelection) {
    // if (true) {
    return (
      <div className="flex justify-center">
        <SearchBar placeholder="Search" value={query} setQuery={setQuery} searchExecute={executeSearch} />
      </div>
    )
  }

  // Render search bar with integrated filters
  return (
    <div className="flex flex-col md:flex-row gap-4 items-start">
      {/* Filters Section */}
      <div className="flex-1">
        <HardcodedFileFilters baseUrl={baseUrl} urlParams={urlParams} />
      </div>

      {/* Search Section */}
      <div className="flex-1 flex justify-center">
        <SearchBar placeholder="Search" value={query} setQuery={setQuery} searchExecute={executeSearch} />
      </StaticFilter>
    </div>
  );
}
