
import React, { Suspense } from "react";
import { useUrlParams } from "@/lib/hooks/useUrlParams";
import AllInOneClientSearch from "@/stateful_components/SearchBar/AllInOneClientSearch";
import DynamicFilters from "@/stateful_components/Filters/DynamicFilters";
import { PageContextMode } from "@/lib/types/SearchTypes";
import LoadingSpinner from "@/style_components/misc/LoadingSpinner";
import RenderedOrg from "@/stateful_components/RenderedObjectCards/RednderedOrg";

export default async function OrgPage({
  params,
}: {
  params: Promise<{ organization_id: string }>;
}) {

  const urlParams = useUrlParams();

  const org_id = (await params).organization_id;

  return (
    <div className="p-4">
      <Suspense fallback={<LoadingSpinner loadingText="Loading Organization Data" />}>
        <RenderedOrg org_id={org_id} />
      </Suspense>
      <h1 className="text-2xl font-bold mb-4">Search [org-name]'s Filings</h1>
      <AllInOneClientSearch urlParams={urlParams} pageContext={PageContextMode.Files}
      />
      {/* <DynamicFilters filters={filters} dataset={dataset} /> */}
    </div>
  );
}
