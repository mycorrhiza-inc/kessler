import DocumentServerPage from "@/components/stateful/Document/DocumentServerPage";
import DefaultContainer from "@/components/stateful/PageContainer/DefaultContainer";

// In this project when i visit http://localhost/filing/bf4ebf20-78f8-4f1c-82c8-d8e7c70fa44e
// it says the filing id is undefined. How can this be happening?
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
    <DocumentServerPage filing_id={file_id} />
  </DefaultContainer>
}
