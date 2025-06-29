import { generateTypeUrlParams, } from "@/lib/types/url_params";
import LoadingSpinner from "@/components/style/misc/LoadingSpinner";
import { Suspense } from "react";
import DefaultContainer from "@/components/stateful/PageContainer/DefaultContainer";
import ConvoPage from "@/components/stateful/ObjectPages/ConvoPage";

export default async function Page({
  params,
  searchParams
}: {
  params: Promise<{ conversation_id: string }>;
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>
}) {

  const urlParams = generateTypeUrlParams(await searchParams);
  const convo_id = (await params).conversation_id;
  if (convo_id == undefined) {
    throw new Error("Undefined convo id")
  }

  return (
    <DefaultContainer>
      <Suspense fallback={<LoadingSpinner loadingText="Loading Docket Data" />}>
        <ConvoPage convo_id={convo_id} urlParams={urlParams} />
      </Suspense>
    </DefaultContainer>
  );
}
