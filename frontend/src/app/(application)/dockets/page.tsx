import React from "react";
import AllInOneClientSearch from "@/components/stateful/SearchBar/AllInOneClientSearch";
import DynamicFilters from "@/components/stateful/Filters/DynamicFilters";
import { ObjectQueryType } from "@/lib/types/SearchTypes";
import { generateTypeUrlParams } from "@/lib/types/url_params";
import { LogoHomepage } from "@/components/style/misc/Logo";
import DefaultContainer from "@/components/stateful/PageContainer/DefaultContainer";
import { GenericSearchType } from "@/lib/adapters/genericSearchCallback";
import ServerSearchResults from "@/components/stateful/ServerResults/ServerResults";

export default async function DocketSearchPage(
  { searchParams }: { searchParams: Promise<{ [key: string]: string | string[] | undefined }> }
) {
  const untypedUrlParams = await searchParams;
  const urlParams = generateTypeUrlParams(untypedUrlParams)

  return (
    <DefaultContainer>
      <h1 className="text-2xl font-bold mb-4"></h1>

      <LogoHomepage bottomText="Docket Search: TODO MAKE IT SEARCH DOCKETS NOT FILINGS" />
      <AllInOneClientSearch urlParams={urlParams.queryData} queryType={ObjectQueryType.Conversations}
      />
      <ServerSearchResults
        baseUrl={`/dockets`}
        urlParams={urlParams}
        objectType={GenericSearchType.Docket}
      />
    </DefaultContainer>
  );
}
