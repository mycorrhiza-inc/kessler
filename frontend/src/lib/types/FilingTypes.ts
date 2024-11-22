import { z } from "zod";

export const FilingValidator = z.object({
  id: z.string(),
  lang: z.string(),
  title: z.string(),
  date: z.string().optional(),
  author: z.string().optional(),
  source: z.string().optional(),
  extension: z.string().optional(),
  file_class: z.string().optional(),
  docket_id: z.string().optional(),
  item_number: z.string().optional(),
  author_organisation: z.string().optional(),
  url: z.string().url().optional(),
});

export type Filing = z.infer<typeof FilingValidator>;

export const testFiling: Filing = {
  id: "3c4ba5f3-febc-41f2-aa86-2820db2b459a",
  url: "https://documents.dps.ny.gov/public/Common/ViewDoc.aspx?DocRefId={7F4AA7FC-CF71-4C2B-8752-A1681D8F9F46}",
  date: "05/12/2022",
  lang: "en",
  title: "Press Release - PSC Announces CLCPA Tracking Initiative",
  author: "Public Service Commission",
  source: "Public Service Commission",
  extension: "pdf",
  file_class: "Press Releases",
  item_number: "3",
  author_organisation: "Public Service Commission",
  // uuid: "",
};
