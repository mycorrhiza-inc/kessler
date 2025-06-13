import React from "react";
import AllInOneClientSearch from "@/stateful_components/SearchBar/AllInOneClientSearch";
import DynamicFilters from "@/stateful_components/Filters/DynamicFilters";
import { PageContextMode } from "@/lib/types/SearchTypes";
import { generateTypeUrlParams } from "@/lib/types/url_params";

export default async function OrgSearchPage(
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>
) {
  const untypedUrlParams = await searchParams;
  const urlParams = generateTypeUrlParams(untypedUrlParams)

  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold mb-4">Organization Search: TODO MAKE SO IT SEARCHES ORGS AND NOT FILINGS</h1>
      <AllInOneClientSearch urlParams={urlParams.queryData} pageContext={PageContextMode.Organizations}
      />
      {/* <DynamicFilters filters={filters} dataset={dataset} /> */}
    </div>
  );
}
