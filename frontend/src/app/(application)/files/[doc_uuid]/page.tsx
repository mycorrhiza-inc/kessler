import DocumentPage, {
  generateDocumentPageData,
} from "@/components/Document/DocumentPage";
import { stateFromHeaders } from "@/lib/nextjs_misc";
import { createClient } from "@/utils/supabase/server";
import { Metadata } from "next";
import { headers } from "next/headers";

export const metadata: Metadata = {
  title: "ERROR IN SITE NAME",
};

export default async function Page({
  params,
}: {
  params: Promise<{ doc_uuid: string }>;
}) {
  const docid = (await params).doc_uuid;
  const headersList = headers();
  const state = stateFromHeaders(headersList);
  const doc_data = await generateDocumentPageData(docid, state);
  metadata.title = doc_data.fileObj.name;

  return (
    <DocumentPage
      fileObj={doc_data.fileObj}
      breadcrumbs={doc_data.breadcrumbs}
    />
  );
}
