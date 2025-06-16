import { generateTypeUrlParams } from "@/lib/types/url_params";
import AllInOneClientSearch from "@/components/stateful/SearchBar/AllInOneClientSearch";
import { ObjectQueryType } from "@/lib/types/SearchTypes";
import { LogoHomepage } from "@/components/style/misc/Logo";
import DefaultContainer from "@/components/stateful/PageContainer/DefaultContainer";

export default async function Page(
  { searchParams }: { searchParams: Promise<{ [key: string]: string | string[] | undefined }> }
) {
  const untypedUrlParams = await searchParams;
  const urlParams = generateTypeUrlParams(untypedUrlParams)
  const targetSearchUrl = "/search"
  return <DefaultContainer>
    <LogoHomepage />
    <AllInOneClientSearch urlParams={urlParams.queryData} queryType={ObjectQueryType.Files} overrideBaseUrl={targetSearchUrl} disableFilterSelection />
  </DefaultContainer>;
}
