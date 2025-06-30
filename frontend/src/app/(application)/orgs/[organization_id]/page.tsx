
import React, { Suspense } from "react";
import { generateTypeUrlParams } from "@/lib/types/url_params";
import AllInOneServerSearch from "@/components/stateful/SearchBar/AllInOneServerSearch";
import LoadingSpinner from "@/components/style/misc/LoadingSpinner";
import DefaultContainer from "@/components/stateful/PageContainer/DefaultContainer";
import { GenericSearchType } from "@/lib/adapters/genericSearchCallback";
import { CardSize } from "@/components/style/cards/SizedCards";
import RenderedCardObject from "@/components/stateful/RenderedObjectCards/RednderedObjectCard";
import OrgPage from "@/components/stateful/ObjectPages/OrgPage";

export default async function Page({
  params,
  searchParams
}: {
  params: Promise<{ organization_id: string }>;
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>
}) {

  const untypedUrlParams = await searchParams;
  const urlParams = generateTypeUrlParams(untypedUrlParams)

  const org_id = (await params).organization_id;

  if (org_id == undefined) {
    throw new Error("Undefined org id")
  }

  return (
    <DefaultContainer>
      <Suspense fallback={<LoadingSpinner loadingText="Loading Organization Data" />}>
        <OrgPage org_id={org_id} urlParams={urlParams} />
      </Suspense>
    </DefaultContainer>
  );
}

