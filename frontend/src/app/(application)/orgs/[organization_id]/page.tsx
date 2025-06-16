
import React, { Suspense } from "react";
import { generateTypeUrlParams } from "@/lib/types/url_params";
import AllInOneClientSearch from "@/components/stateful/SearchBar/AllInOneClientSearch";
import DynamicFilters from "@/components/stateful/Filters/DynamicFilters";
import { ObjectQueryType } from "@/lib/types/SearchTypes";
import LoadingSpinner from "@/components/style/misc/LoadingSpinner";
import RenderedOrg from "@/components/stateful/RenderedObjectCards/RednderedOrg";
import DefaultContainer from "@/components/stateful/PageContainer/DefaultContainer";
import ServerSearchResults from "@/components/stateful/ServerResults/ServerResults";
import { GenericSearchType } from "@/lib/adapters/genericSearchCallback";

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
    <DefaultContainer>
      <Suspense fallback={<LoadingSpinner loadingText="Loading Organization Data" />}>
        <RenderedOrg org_id={org_id} />
      </Suspense>
      <h1 className="text-2xl font-bold mb-4">Search [org-name]'s Filings</h1>
      <AllInOneClientSearch urlParams={urlParams.queryData} queryType={ObjectQueryType.Files}
      />
      <ServerSearchResults
        baseUrl={`/orgs/${org_id}`}
        urlParams={urlParams}
        objectType={GenericSearchType.Filling}
        inherentRouteFilters={{ "author_id": org_id }}
      />
    </DefaultContainer>
  );
}
