import { generateTypeUrlParams } from "@/lib/types/url_params";
import AllInOneClientSearch from "@/components/stateful/SearchBar/AllInOneClientSearch";
import { ObjectQueryType } from "@/lib/types/SearchTypes";

export default async function Page(
  { searchParams }: { searchParams: Promise<{ [key: string]: string | string[] | undefined }> }
) {
  const untypedUrlParams = await searchParams;
  const urlParams = generateTypeUrlParams(untypedUrlParams)
  return <AllInOneClientSearch urlParams={urlParams.queryData} queryType={ObjectQueryType.Files} />;
}
