import React from "react";
import AllInOneServerSearch from "@/components/stateful/SearchBar/AllInOneServerSearch";
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
      <AllInOneServerSearch
        urlParams={urlParams}
        queryType={GenericSearchType.Docket}
        baseUrl={`/dockets`}
      />
    </DefaultContainer>
  );
}
