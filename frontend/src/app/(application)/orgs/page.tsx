import React from "react";
import { ObjectQueryType } from "@/lib/types/SearchTypes";
import { generateTypeUrlParams } from "@/lib/types/url_params";
import AllInOneClientSearch from "@/components/stateful/SearchBar/AllInOneClientSearch";

export default async function OrgSearchPage(

  { searchParams }: { searchParams: Promise<{ [key: string]: string | string[] | undefined }> }
) {
  const untypedUrlParams = await searchParams;
  const urlParams = generateTypeUrlParams(untypedUrlParams)

  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold mb-4">Organization Search: TODO MAKE SO IT SEARCHES ORGS AND NOT FILINGS</h1>
      <AllInOneClientSearch urlParams={urlParams.queryData} queryType={ObjectQueryType.Organizations}
      />
      {/* <DynamicFilters filters={filters} dataset={dataset} /> */}
    </div>
  );
}
