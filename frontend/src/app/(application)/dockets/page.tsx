import React from "react";
import { useUrlParamsClient } from "@/lib/types/url_params";
import AllInOneClientSearch from "@/stateful_components/SearchBar/AllInOneClientSearch";
import DynamicFilters from "@/stateful_components/Filters/DynamicFilters";
import { PageContextMode } from "@/lib/types/SearchTypes";

export default async function DocketSearchPage(
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>
) {
  const urlParams = await searchParams;

  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold mb-4">Docket Search: TODO MAKE IT SEARCH DOCKETS NOT FILINGS</h1>
      <AllInOneClientSearch urlParams={urlParams} pageContext={PageContextMode.Conversations}
      />
      {/* <DynamicFilters filters={filters} dataset={dataset} /> */}
    </div>
  );
}
