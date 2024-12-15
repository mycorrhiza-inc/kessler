import DocumentPage, {
  generateDocumentPageData,
} from "@/components/Document/DocumentPage";
import { stateFromHeaders } from "@/lib/nextjs_misc";
import { createClient } from "@/utils/supabase/server";
import { Metadata } from "next";
import { headers } from "next/headers";
import { cache } from "react";

const cachedDocPageData = cache(generateDocumentPageData);

export async function generateMetadata({
  params,
}: {
  params: Promise<{ doc_uuid: string }>;
}): Promise<Metadata> {
  // read route params
  const id = (await params).doc_uuid;
  const doc_data = await cachedDocPageData(id);
  return {
    title: doc_data.fileObj.name,
  };
}
export default async function Page({
  params,
}: {
  params: Promise<{ doc_uuid: string }>;
}) {
  const docid = (await params).doc_uuid;
  const headersList = headers();
  const state = stateFromHeaders(headersList);
  const doc_data = await cachedDocPageData(docid);
  doc_data.breadcrumbs.state = state || "";

  return (
    <DocumentPage
      fileObj={doc_data.fileObj}
      breadcrumbs={doc_data.breadcrumbs}
    />
  );
}
