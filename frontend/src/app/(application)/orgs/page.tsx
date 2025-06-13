import React from "react";
import { useUrlParams } from "@/lib/hooks/useUrlParams";
import AllInOneClientSearch from "@/stateful_components/SearchBar/AllInOneClientSearch";
import DynamicFilters from "@/stateful_components/Filters/DynamicFilters";
import { PageContextMode } from "@/lib/types/SearchTypes";

export default function OrgSearchPage() {
  const urlParams = useUrlParams();

  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold mb-4">Organization Search: TODO MAKE SO IT SEARCHES ORGS AND NOT FILINGS</h1>
      <AllInOneClientSearch urlParams={urlParams} pageContext={PageContextMode.Organizations}
      />
      {/* <DynamicFilters filters={filters} dataset={dataset} /> */}
    </div>
  );
}
