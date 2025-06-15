import React from "react";
import AllInOneClientSearch from "@/componenets/stateful/SearchBar/AllInOneClientSearch";
import DynamicFilters from "@/componenets/stateful/Filters/DynamicFilters";
import { PageContextMode } from "@/lib/types/SearchTypes";
import { generateTypeUrlParams } from "@/lib/types/url_params";

export default async function SearchPage(
  { searchParams }: { searchParams: Promise<{ [key: string]: string | string[] | undefined }> }
) {

  const untypedUrlParams = await searchParams;
  const urlParams = generateTypeUrlParams(untypedUrlParams)


  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold mb-4">Search</h1>
      <AllInOneClientSearch urlParams={urlParams.queryData} pageContext={PageContextMode.Files}
      />
      {/* <DynamicFilters filters={filters} dataset={dataset} /> */}
    </div>
  );
}
