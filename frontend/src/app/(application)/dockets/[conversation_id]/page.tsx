import { generateTypeUrlParams, } from "@/lib/types/url_params";
import { PageContextMode } from "@/lib/types/SearchTypes";
import RenderedConvo from "@/componenets/stateful/RenderedObjectCards/RednderedConvo";
import AllInOneClientSearch from "@/componenets/stateful/SearchBar/AllInOneClientSearch";
import LoadingSpinner from "@/componenets/style/misc/LoadingSpinner";
import { Suspense } from "react";

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
    <div className="p-4">
      <Suspense fallback={<LoadingSpinner loadingText="Loading Organization Data" />}>
        <RenderedConvo convo_id={convo_id} />
      </Suspense>
      <h1 className="text-2xl font-bold mb-4">Search [org-name]'s Filings</h1>
      <AllInOneClientSearch urlParams={urlParams.queryData} pageContext={PageContextMode.Files}
      />
      {/* <DynamicFilters filters={filters} dataset={dataset} /> */}
    </div>
  );
}
