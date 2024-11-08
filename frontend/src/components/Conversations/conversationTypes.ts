export type Filing = {
  id: string;
  lang: string;
  title: string;
  date: string;
  author: string;
  source: string;
  language: string;
  extension: string;
  file_class: string;
  item_number: string;
  author_organisation: string;
  url: string;
  uuid: string;
};

export const testFiling: Filing = {
  id: "0",
  url: "https://documents.dps.ny.gov/public/Common/ViewDoc.aspx?DocRefId={7F4AA7FC-CF71-4C2B-8752-A1681D8F9F46}",
  date: "05/12/2022",
  lang: "en",
  title: "Press Release - PSC Announces CLCPA Tracking Initiative",
  author: "Public Service Commission",
  source: "Public Service Commission",
  language: "en",
  extension: "pdf",
  file_class: "Press Releases",
  item_number: "3",
  author_organisation: "Public Service Commission",
  uuid: "3c4ba5f3-febc-41f2-aa86-2820db2b459a",
};
