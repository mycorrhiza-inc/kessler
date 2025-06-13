import { useUrlParams } from "@/lib/hooks/useUrlParams";
import { PageContextMode } from "@/lib/types/SearchTypes";
import RenderedConvo from "@/stateful_components/RenderedObjectCards/RednderedConvo";
import AllInOneClientSearch from "@/stateful_components/SearchBar/AllInOneClientSearch";
import LoadingSpinner from "@/style_components/misc/LoadingSpinner";
import { Suspense } from "react";

export default async function Page({
  params,
}: {
  params: Promise<{ conversation_id: string }>;
}) {


  const urlParams = useUrlParams();

  const convo_id = (await params).conversation_id;

  return (
    <div className="p-4">
      <Suspense fallback={<LoadingSpinner loadingText="Loading Organization Data" />}>
        <RenderedConvo convo_id={convo_id} />
      </Suspense>
      <h1 className="text-2xl font-bold mb-4">Search [org-name]'s Filings</h1>
      <AllInOneClientSearch urlParams={urlParams} pageContext={PageContextMode.Files}
      />
      {/* <DynamicFilters filters={filters} dataset={dataset} /> */}
    </div>
  );
}
