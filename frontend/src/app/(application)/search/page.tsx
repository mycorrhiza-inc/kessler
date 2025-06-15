import React from "react";
import AllInOneClientSearch from "@/components/stateful/SearchBar/AllInOneClientSearch";
import { ObjectQueryType } from "@/lib/types/SearchTypes";
import { generateTypeUrlParams } from "@/lib/types/url_params";

export default async function SearchPage(
  { searchParams }: { searchParams: Promise<{ [key: string]: string | string[] | undefined }> }
) {

  const untypedUrlParams = await searchParams;
  const urlParams = generateTypeUrlParams(untypedUrlParams)


  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold mb-4">Search</h1>
      <AllInOneClientSearch urlParams={urlParams.queryData} queryType={ObjectQueryType.Files}
      />
      {/* <DynamicFilters filters={filters} dataset={dataset} /> */}
    </div>
  );
}
