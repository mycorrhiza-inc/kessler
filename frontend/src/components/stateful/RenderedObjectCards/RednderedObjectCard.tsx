import Card, { CardSize } from "@/components/style/cards/GenericResultCard";
import { GenericSearchType } from "@/lib/adapters/genericSearchCallback";
import { generateFakeResultsRaw } from "@/lib/search/search_utils";
import { sleep } from "@/lib/utils";

export default async function RenderedCardObject({ object_id, objectType, size }: { object_id: string, objectType: GenericSearchType, size?: CardSize }) {
  await sleep(500);
  if (!size) {
    size = CardSize.Large as CardSize;
  }
  switch (objectType) {
    case GenericSearchType.Organization:
      return <Card size={size} data={generateFakeResultsRaw(1)[0]} />
    case GenericSearchType.Docket:
      return <Card size={size} data={generateFakeResultsRaw(1)[0]} />
    case GenericSearchType.Filling:
      return <Card size={size} data={generateFakeResultsRaw(1)[0]} />
    case GenericSearchType.Dummy:
      return <Card size={size} data={generateFakeResultsRaw(1)[0]} />
  }
}
