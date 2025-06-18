import { generateTypeUrlParams, } from "@/lib/types/url_params";
import { ObjectQueryType } from "@/lib/types/SearchTypes";
import RenderedConvo from "@/components/stateful/RenderedObjectCards/RednderedConvo";
import AllInOneServerSearch from "@/components/stateful/SearchBar/AllInOneServerSearch";
import LoadingSpinner from "@/components/style/misc/LoadingSpinner";
import { Suspense } from "react";
import DefaultContainer from "@/components/stateful/PageContainer/DefaultContainer";
import ServerSearchResults from "@/components/stateful/ServerResults/ServerResults";
import { GenericSearchType } from "@/lib/adapters/genericSearchCallback";

export default async function Page({
  params,
  searchParams
}: {
  params: Promise<{ conversation_id: string }>;
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>
}) {

  const untypedUrlParams = await searchParams;
  const urlParams = generateTypeUrlParams(untypedUrlParams)


  const convo_id = (await params).conversation_id;

  return (
    <DefaultContainer>
      <Suspense fallback={<LoadingSpinner loadingText="Loading Organization Data" />}>
        <RenderedConvo convo_id={convo_id} />
      </Suspense>
      <AllInOneServerSearch
        aboveSearchElement={<h1 className="text-2xl font-bold mb-4">Search [org-name]'s Filings</h1>}
        urlParams={urlParams}
        queryType={GenericSearchType.Filling}

        inherentRouteFilters={{ "conversation_id": convo_id }}
        baseUrl={`/dockets/${convo_id}`}
      />
    </DefaultContainer>
  );
}
