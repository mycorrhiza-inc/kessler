import ErrorMessage from "@/components/style/messages/ErrorMessage"
import { fetchAuthorCardData } from "../RenderedObjectCards/RednderedObjectCard"
import AllInOneServerSearch from "../SearchBar/AllInOneServerSearch"
import { TypedUrlParams } from "@/lib/types/url_params"
import { CardSize } from "@/components/style/cards/SizedCards"
import Card from "../Card/LinkedCard"
import { DocumentMainTabsClient } from "../Document/DocumentBody"

export default async function FilePage({ file_id, urlParams }: { file_id: string, urlParams: TypedUrlParams }) {
  try {
    const card_data = await fetchAuthorCardData(file_id)
    const doc_object = { verifed: true, id: file_id, extension: "pdf" } as any;
    return (
      <>
        <Card data={card_data} size={CardSize.Large} disableHref />

        <DocumentMainTabsClient documentObject={doc_object} isPage />
      </>
    )
  } catch (err) {
    return <ErrorMessage error={err} message="Could not get filling data from server :(" />
  }
}
