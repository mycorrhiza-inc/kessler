import DefaultContainer from "@/components/stateful/PageContainer/DefaultContainer";
import RenderedCardObject from "@/components/stateful/RenderedObjectCards/RednderedObjectCard";
import { CardSize } from "@/components/style/cards/GenericResultCard";
import { GenericSearchType } from "@/lib/adapters/genericSearchCallback";

export default async function Page({
  params,
}: {
  params: Promise<{ filling_id: string }>;
}) {

  const filling_id = (await params).filling_id;
  return <DefaultContainer>
    <RenderedCardObject objectType={GenericSearchType.Filling} object_id={filling_id} size={CardSize.Large} />
  </DefaultContainer>
}
