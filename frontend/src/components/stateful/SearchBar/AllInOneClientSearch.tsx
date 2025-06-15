"use client";
import { encodeUrlParams, TypedUrlParams, UrlQueryParams } from "@/lib/types/url_params"
import { PageContextMode } from "@/lib/types/SearchTypes";
import { useState } from "react";
import { usePathname, useRouter } from "next/navigation";

export interface AIOSearchProps {
  urlParams: UrlQueryParams
}

export default function AllInOneClientSearch({ urlParams, pageContext, overrideBaseUrl }: {
  urlParams: UrlQueryParams,
  pageContext: PageContextMode,
  overrideBaseUrl?: string
}) {
  const router = useRouter()
  const impliedUrl: string = usePathname();
  const baseUrl = overrideBaseUrl || impliedUrl || "/search"
  const [query, setQuery] = useState(urlParams.query || "")
  const [filters, setFilters] = useState({})

  const executeSearch = () => {
    var newQueryParams = urlParams;
    newQueryParams.query = query
    const newParams: TypedUrlParams = { queryData: newQueryParams, paginationData: {} }

    const encodedUrlQuery = encodeUrlParams(newParams)
    const url = baseUrl + encodedUrlQuery
    router.push(url)
  }
  return <div>Not Implemented</div>
}
