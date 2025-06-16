import { generateTypeUrlParams } from "@/lib/types/url_params";
import AllInOneClientSearch from "@/components/stateful/SearchBar/AllInOneClientSearch";
import { ObjectQueryType } from "@/lib/types/SearchTypes";
import { LogoHomepage } from "@/components/style/misc/Logo";

export default async function Page(
  { searchParams }: { searchParams: Promise<{ [key: string]: string | string[] | undefined }> }
) {
  const untypedUrlParams = await searchParams;
  const urlParams = generateTypeUrlParams(untypedUrlParams)
  const targetSearchUrl = "/search"
  return <div className="flex justify-center items-center flex-col ">
    <LogoHomepage />
    <AllInOneClientSearch urlParams={urlParams.queryData} queryType={ObjectQueryType.Files} overrideBaseUrl={targetSearchUrl} />
  </div>;
}
