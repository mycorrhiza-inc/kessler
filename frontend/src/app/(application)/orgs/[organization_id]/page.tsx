
import React, { Suspense } from "react";
import { generateTypeUrlParams } from "@/lib/types/url_params";
import AllInOneServerSearch from "@/components/stateful/SearchBar/AllInOneServerSearch";
import LoadingSpinner from "@/components/style/misc/LoadingSpinner";
import DefaultContainer from "@/components/stateful/PageContainer/DefaultContainer";
import { GenericSearchType } from "@/lib/adapters/genericSearchCallback";
import { CardSize } from "@/components/style/cards/GenericResultCard";
import RenderedCardObject from "@/components/stateful/RenderedObjectCards/RednderedObjectCard";

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
        <RenderedCardObject objectType={GenericSearchType.Organization} object_id={org_id} size={CardSize.Large} />
      </Suspense>
      <AllInOneServerSearch
        aboveSearchElement={<h1 className="text-2xl font-bold mb-4">Search [org-name]'s Filings</h1>}
        urlParams={urlParams}
        queryType={GenericSearchType.Filling}
        inherentRouteFilters={{ "author_id": org_id }}
        baseUrl={`/orgs/${org_id}`}
      />
    </DefaultContainer>
  );
}
