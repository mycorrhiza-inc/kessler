import RenderedCardObject from "@/components/stateful/RenderedObjectCards/RednderedObjectCard";
import { CardSize } from "@/components/style/cards/SizedCards";
import { GenericSearchType } from "@/lib/adapters/genericSearchCallback";
import { DocumentMainTabsClient } from "./DocumentBody";

export default async function DocumentServerPage({
  filing_id,
}: {
  filing_id: string
}) {
  const doc_object = { verifed: true, id: filing_id, extension: "pdf" } as any;

  return <>
    <RenderedCardObject objectType={GenericSearchType.Filing} object_id={filing_id} size={CardSize.Large} />
    <DocumentMainTabsClient documentObject={doc_object} isPage />
  </>

}
