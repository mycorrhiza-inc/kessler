import ErrorMessage from "@/components/style/messages/ErrorMessage"
import { TypedUrlParams } from "@/lib/types/url_params"
import { CardSize } from "@/components/style/cards/SizedCards"
import Card from "../Card/LinkedCard"
import { DocumentMainTabsClient } from "../Document/DocumentBody"
import { fetchDocumentCardData, fetchDocumentPageData } from "../RenderedObjectCards/RednderedObjectCard"

export default async function FilePage({ file_id, urlParams }: { file_id: string, urlParams: TypedUrlParams }) {
  try {
    const file_page_info = await fetchDocumentPageData(file_id)
    const doc_object = { verifed: true, id: file_id, extension: "pdf" } as any;
    return (
      <>
        <Card data={file_page_info.card_info} size={CardSize.Large} disableHref />

        <DocumentMainTabsClient documentObject={doc_object} isPage />
      </>
    )
  } catch (err) {
    return <ErrorMessage error={err} message="Could not get filling data from server :(" />
  }
}
