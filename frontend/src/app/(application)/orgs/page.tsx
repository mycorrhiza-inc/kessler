import React from "react";
import { ObjectQueryType } from "@/lib/types/SearchTypes";
import { generateTypeUrlParams } from "@/lib/types/url_params";
import AllInOneClientSearch from "@/components/stateful/SearchBar/AllInOneClientSearch";
import { LogoHomepage } from "@/components/style/misc/Logo";
import DefaultContainer from "@/components/stateful/PageContainer/DefaultContainer";

export default async function OrgSearchPage(

  { searchParams }: { searchParams: Promise<{ [key: string]: string | string[] | undefined }> }
) {
  const untypedUrlParams = await searchParams;
  const urlParams = generateTypeUrlParams(untypedUrlParams)

  return (
    <DefaultContainer >
      <LogoHomepage bottomText="Organization Search: TODO MAKE SO IT SEARCHES ORGS AND NOT FILINGS" />
      <AllInOneClientSearch urlParams={urlParams.queryData} queryType={ObjectQueryType.Organizations}
      />
      {/* <DynamicFilters filters={filters} dataset={dataset} /> */}
    </DefaultContainer>
  );
}
