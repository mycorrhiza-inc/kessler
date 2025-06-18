import React from "react";
import { ObjectQueryType } from "@/lib/types/SearchTypes";
import { generateTypeUrlParams } from "@/lib/types/url_params";
import AllInOneServerSearch from "@/components/stateful/SearchBar/AllInOneServerSearch";
import { LogoHomepage } from "@/components/style/misc/Logo";
import DefaultContainer from "@/components/stateful/PageContainer/DefaultContainer";
import { GenericSearchType } from "@/lib/adapters/genericSearchCallback";
import ServerSearchResults from "@/components/stateful/ServerResults/ServerResults";

export default async function OrgSearchPage(

  { searchParams }: { searchParams: Promise<{ [key: string]: string | string[] | undefined }> }
) {
  const untypedUrlParams = await searchParams;
  const urlParams = generateTypeUrlParams(untypedUrlParams)

  return (
    <DefaultContainer >
      <AllInOneServerSearch
        aboveSearchElement={<LogoHomepage bottomText="Organization Search: TODO MAKE SO IT SEARCHES ORGS AND NOT FILINGS" />}
        urlParams={urlParams}
        queryType={GenericSearchType.Organization}
        baseUrl="/orgs"
      />
    </DefaultContainer>
  );
}
