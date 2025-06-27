import { generateTypeUrlParams } from "@/lib/types/url_params";
import AllInOneServerSearch from "@/components/stateful/SearchBar/AllInOneServerSearch";
import { ObjectQueryType } from "@/lib/types/SearchTypes";
import { LogoHomepage } from "@/components/style/misc/Logo";
import DefaultContainer from "@/components/stateful/PageContainer/DefaultContainer";
import { GenericSearchType } from "@/lib/adapters/genericSearchCallback";

export default async function Page(
  { searchParams }: { searchParams: Promise<{ [key: string]: string | string[] | undefined }> }
) {
  const untypedUrlParams = await searchParams;
  const urlParams = generateTypeUrlParams(untypedUrlParams)
  const targetSearchUrl = "/search"
  return <DefaultContainer>
    <AllInOneServerSearch
      aboveSearchElement={<LogoHomepage />}
      urlParams={urlParams}
      baseUrl={targetSearchUrl}
      disableFilterSelection
      disableResults
    />
  </DefaultContainer>;
}
