import React from "react";
import AllInOneClientSearch from "@/stateful_components/SearchBar/AllInOneClientSearch";
import DynamicFilters from "@/stateful_components/Filters/DynamicFilters";
import { PageContextMode } from "@/lib/types/SearchTypes";

export default async function SearchPage(
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>
) {
  const untypedUrlParams = await searchParams;



  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold mb-4">Search</h1>
      <AllInOneClientSearch urlParams={urlParams} pageContext={PageContextMode.Files}
      />
      {/* <DynamicFilters filters={filters} dataset={dataset} /> */}
    </div>
  );
}
