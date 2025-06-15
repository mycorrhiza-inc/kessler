import React from "react";
import AllInOneClientSearch from "@/components/stateful/SearchBar/AllInOneClientSearch";
import { ObjectQueryType } from "@/lib/types/SearchTypes";
import { generateTypeUrlParams } from "@/lib/types/url_params";
import { LogoHomepage } from "@/components/style/misc/Logo";

export default async function SearchPage(
  { searchParams }: { searchParams: Promise<{ [key: string]: string | string[] | undefined }> }
) {

  const untypedUrlParams = await searchParams;
  const urlParams = generateTypeUrlParams(untypedUrlParams)


  return (
    <div className="flex justify-center items-center">
      <LogoHomepage bottomText="Search" />
      <AllInOneClientSearch urlParams={urlParams.queryData} queryType={ObjectQueryType.Files}
      />
      <
    </div>
  );
}
