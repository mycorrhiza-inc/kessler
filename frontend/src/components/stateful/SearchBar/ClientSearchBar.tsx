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

export default function ClientSearchBar({ urlParams, baseUrl }: {
  urlParams: UrlQueryParams,
  baseUrl: string,
}) {
  const router = useRouter()
  const [query, setQuery] = useState(urlParams.query || "")

  const executeSearch = () => {
    const newQueryParams = { ...urlParams, query };
    const newParams: TypedUrlParams = { queryData: newQueryParams, paginationData: {} }
    const encodedUrlQuery = encodeUrlParams(newParams)
    const url = baseUrl + encodedUrlQuery
    router.push(url)
  }

  return (
    <div className="flex justify-center">
      <SearchBar placeholder="Search" value={query} setQuery={setQuery} searchExecute={executeSearch} />
    </div>
  )
}
