import { CardData, CardSize } from "@/components/style/cards/SizedCards";
import ErrorMessage from "@/components/style/messages/ErrorMessage";
import { GenericSearchType } from "@/lib/adapters/genericSearchCallback";
import { getContextualAPIUrl } from "@/lib/env_variables";
import { generateFakeResultsRaw } from "@/lib/search/search_utils";
import { AuthorCardData, AuthorCardDataValidator, DocketCardData, DocketCardDataValidator, DocumentCardData, DocumentCardDataValidator } from "@/lib/types/generic_card_types";
import { sleep } from "@/lib/utils";
import axios from "axios";
import Card from "../Card/LinkedCard";


// Define the API endpoint fetchers
export async function fetchAuthorCardData(object_id: string): Promise<AuthorCardData> {
  const api_url = getContextualAPIUrl();
  const endpoint = `${api_url}/public/organizations/${object_id}/card`;
  console.log("Fetching Organization data from:", endpoint);
  const response = await axios.get(endpoint);
  console.log("Organization response data:", response.data);
  return AuthorCardDataValidator.parse(response.data);
}

export async function fetchDocketCardData(object_id: string): Promise<DocketCardData> {
  const api_url = getContextualAPIUrl();
  const endpoint = `${api_url}/public/conversations/${object_id}/card`;
  console.log("Fetching Docket data from:", endpoint);
  const response = await axios.get(endpoint);
  console.log("Docket response data:", response.data);
  return DocketCardDataValidator.parse(response.data);
}

export async function fetchDocumentCardData(object_id: string): Promise<DocumentCardData> {
  const api_url = getContextualAPIUrl();
  const endpoint = `${api_url}/public/files/${object_id}/card`;
  console.log("Fetching Filling data from:", endpoint);
  const response = await axios.get(endpoint);
  console.log("Filling response data:", response.data);
  return DocumentCardDataValidator.parse(response.data);
}

type FillingAttachmentInfo = {
  atttachment_uuid: string;
  atttachment_hash: string;
  atttachment_name: string;
  atttachment_extension: string;
};

type FilePageInfo = {
  card_info: DocketCardData; // Ideally, replace 'any' with the actual type/interface for search.DocumentCardData
  attachments: FillingAttachmentInfo[];
};
export async function fetchDocumentPageData(object_id: string): Promise<FilePageInfo> {
  const api_url = getContextualAPIUrl();
  const endpoint = `${api_url}/public/files/${object_id}/p`;
  console.log("Fetching Filling data from:", endpoint);
  const response = await axios.get(endpoint);
  console.log("Filling response data:", response.data);
  return response.data;
}

export async function fetchDummyCardData(): Promise<CardData> {
  return generateFakeResultsRaw(1)[0];
}


// Define the main fetchCardData function
export async function fetchCardData(object_id: string, objectType: GenericSearchType): Promise<CardData> {
  switch (objectType) {
    case GenericSearchType.Organization:
      return fetchAuthorCardData(object_id);
    case GenericSearchType.Docket:
      return fetchDocketCardData(object_id);
    case GenericSearchType.Filing:
      return fetchDocumentCardData(object_id);
    case GenericSearchType.Dummy:
      return fetchDummyCardData();
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
