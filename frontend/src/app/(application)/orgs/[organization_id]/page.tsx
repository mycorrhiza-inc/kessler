
import React, { Suspense } from "react";
import { generateTypeUrlParams } from "@/lib/types/url_params";
import AllInOneClientSearch from "@/componenets/stateful/SearchBar/AllInOneClientSearch";
import DynamicFilters from "@/componenets/stateful/Filters/DynamicFilters";
import { PageContextMode } from "@/lib/types/SearchTypes";
import LoadingSpinner from "@/componenets/style/misc/LoadingSpinner";
import RenderedOrg from "@/componenets/stateful/RenderedObjectCards/RednderedOrg";

export default async function OrgPage({
  params,
  searchParams
}: {
  params: Promise<{ organization_id: string }>;
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>
}) {

  const untypedUrlParams = await searchParams;
  const urlParams = generateTypeUrlParams(untypedUrlParams)

  const org_id = (await params).organization_id;

  return (
    <div className="p-4">
      <Suspense fallback={<LoadingSpinner loadingText="Loading Organization Data" />}>
        <RenderedOrg org_id={org_id} />
      </Suspense>
      <h1 className="text-2xl font-bold mb-4">Search [org-name]'s Filings</h1>
      <AllInOneClientSearch urlParams={urlParams.queryData} pageContext={PageContextMode.Files}
      />
      {/* <DynamicFilters filters={filters} dataset={dataset} /> */}
    </div>
  );
}
