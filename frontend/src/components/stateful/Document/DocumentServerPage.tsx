import RenderedCardObject from "@/components/stateful/RenderedObjectCards/RednderedObjectCard";
import { CardSize } from "@/components/style/cards/GenericResultCard";
import { GenericSearchType } from "@/lib/adapters/genericSearchCallback";
import { DocumentMainTabsClient } from "./DocumentBody";

export default async function DocumentServerPage({
  filling_id,
}: {
  filling_id: string
}) {
  const doc_object = { verifed: true, id: filling_id } as any;

  return <>
    <RenderedCardObject objectType={GenericSearchType.Filling} object_id={filling_id} size={CardSize.Large} />
    <DocumentMainTabsClient documentObject={doc_object} isPage />
  </>

}
