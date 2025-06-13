"use client";
import { UrlParams } from "@/lib/hooks/useUrlParams"
import { PageContextMode } from "@/lib/types/SearchTypes";

export interface AIOSearchProps {
  urlParams: UrlParams
}

export default function AllInOneClientSearch({ urlParams, pageContext }: { urlParams: UrlParams, pageContext: PageContextMode }) {
  return <div>Not Implemented</div>
}
