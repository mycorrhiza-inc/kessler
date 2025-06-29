import DocumentServerPage from "@/components/stateful/Document/DocumentServerPage";
import FilePage from "@/components/stateful/ObjectPages/FilePage";
import DefaultContainer from "@/components/stateful/PageContainer/DefaultContainer";
import LoadingSpinner from "@/components/style/misc/LoadingSpinner";
import { generateTypeUrlParams } from "@/lib/types/url_params";
import { Suspense } from "react";

// In this project when i visit http://localhost/filing/bf4ebf20-78f8-4f1c-82c8-d8e7c70fa44e
// it says the filing id is undefined. How can this be happening?
export default async function Page({
  params,
  searchParams
}: {
  params: Promise<{ file_id: string }>;
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>
}) {

  const urlParams = generateTypeUrlParams(await searchParams);

  const file_id = (await params).file_id;
  if (file_id === undefined) {
    throw new Error("Undefined file id")
  }
  return <DefaultContainer>
    <Suspense fallback={<LoadingSpinner loadingText="Loading Filling Data" />}>
      <FilePage file_id={file_id} urlParams={urlParams} />
    </Suspense>
  </DefaultContainer>
}
