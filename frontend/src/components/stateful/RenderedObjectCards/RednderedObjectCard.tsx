import { CardData, CardSize } from "@/components/style/cards/SizedCards";
import ErrorMessage from "@/components/style/messages/ErrorMessage";
import { GenericSearchType } from "@/lib/adapters/genericSearchCallback";
import { getContextualAPIUrl } from "@/lib/env_variables";
import { generateFakeResultsRaw } from "@/lib/search/search_utils";
import { AuthorCardDataValidator, DocumentCardDataValidator } from "@/lib/types/generic_card_types";
import { sleep } from "@/lib/utils";
import axios from "axios";
import Card from "../Card/LinkedCard";

async function fetchCardData(object_id: string, objectType: GenericSearchType): Promise<CardData> {
  const api_url = getContextualAPIUrl()


  switch (objectType) {
    case GenericSearchType.Organization:
      const org_endpoint = `${api_url}/card/org/${object_id}`
      console.log("Fetching Organization data from:", org_endpoint)
      const org_response = await axios.get(org_endpoint)
      console.log("Organization response data:", org_response.data)
      const org_raw_data = org_response.data
      const org_data = AuthorCardDataValidator.parse(org_raw_data)
      return org_data

    case GenericSearchType.Docket:
      const docket_endpoint = `${api_url}/card/convo/${object_id}`
      console.log("Fetching Docket data from:", docket_endpoint)
      const docket_response = await axios.get(docket_endpoint)
      console.log("Docket response data:", docket_response.data)
      const docket_raw_data = docket_response.data
      const docket_data = AuthorCardDataValidator.parse(docket_raw_data)
      return docket_data
    case GenericSearchType.Filing:
      const filling_endpoint = `${api_url}/card/file/${object_id}`
      console.log("Fetching Filling data from:", filling_endpoint)
      const filling_response = await axios.get(filling_endpoint)
      console.log("Filling response data:", filling_response.data)
      const filling_raw_data = filling_response.data
      const filling_data = DocumentCardDataValidator.parse(filling_raw_data)
      return filling_data
    case GenericSearchType.Dummy:
      return generateFakeResultsRaw(1)[0]
  }
}


export default async function RenderedCardObject({ object_id, objectType, size }: { object_id: string, objectType: GenericSearchType, size?: CardSize }) {
  try {
    await sleep(500);
    if (!size) {
      size = CardSize.Large as CardSize;
    }
    const cardData = await fetchCardData(object_id, objectType)

    return <Card data={cardData} size={size} disableHref />
  } catch (err) {
    console.log(err)
    return <ErrorMessage message="Could not get object info" error={err} />
  }
}
