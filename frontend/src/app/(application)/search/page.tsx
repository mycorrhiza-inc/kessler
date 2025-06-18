import React from "react";
import AllInOneServerSearch from "@/components/stateful/SearchBar/AllInOneServerSearch";
import { ObjectQueryType } from "@/lib/types/SearchTypes";
import { generateTypeUrlParams } from "@/lib/types/url_params";
import { LogoHomepage } from "@/components/style/misc/Logo";
import ServerSearchResults from "@/components/stateful/ServerResults/ServerResults";
import { GenericSearchType } from "@/lib/adapters/genericSearchCallback";
import DefaultContainer from "@/components/stateful/PageContainer/DefaultContainer";

export default async function SearchPage(
  { searchParams }: { searchParams: Promise<{ [key: string]: string | string[] | undefined }> }
) {

  const untypedUrlParams = await searchParams;
  const urlParams = generateTypeUrlParams(untypedUrlParams)


  return (
    <DefaultContainer>
      <LogoHomepage bottomText="Search" />
      <AllInOneServerSearch
        urlParams={urlParams}
        queryType={GenericSearchType.Filling}
        baseUrl="/search"
      />
    </DefaultContainer>
  );
}
