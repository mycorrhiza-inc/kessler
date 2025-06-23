import { DocumentMainTabsClient } from "@/components/stateful/Document/DocumentBody";
import DocumentServerPage from "@/components/stateful/Document/DocumentServerPage";
import DefaultContainer from "@/components/stateful/PageContainer/DefaultContainer";

export default async function Page({
  params,
}: {
  params: Promise<{ file_id: string }>;
}) {

  const file_id = (await params).file_id;
  if (file_id === undefined) {
    throw new Error("Undefined file id")

  }
  return <DefaultContainer>
    <DocumentServerPage filling_id={file_id} />
  </DefaultContainer>
}
