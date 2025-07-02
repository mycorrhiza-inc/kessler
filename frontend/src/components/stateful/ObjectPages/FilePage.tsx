import ErrorMessage from "@/components/style/messages/ErrorMessage"
import { TypedUrlParams } from "@/lib/types/url_params"
import { CardSize } from "@/components/style/cards/SizedCards"
import Card from "../Card/LinkedCard"
import { DocumentMainTabsClient } from "../Document/DocumentBody"
import { fetchDocumentCardData, fetchDocumentPageData } from "../RenderedObjectCards/RednderedObjectCard"
import { DocumentCardDataValidator } from "@/lib/types/generic_card_types"

export default async function FilePage({ file_id, urlParams }: { file_id: string, urlParams: TypedUrlParams }) {
  try {
    const file_page_info = await fetchDocumentPageData(file_id)
    const card_info = DocumentCardDataValidator.parse(file_page_info.card_info)
    return (
      <>
        <Card data={card_info} size={CardSize.Large} disableHref />

        <DocumentMainTabsClient attachmentInfos={file_page_info.attachments} isPage />
      </>
    )
  } catch (err) {
    return <ErrorMessage error={err} message="Could not get filling data from server :(" />
  }
}
