import { ConversationView } from "@/components/Conversations/ConversationView";
import DocumentPage from "@/components/Document/DocumentPage";
import Navbar from "@/components/Navbar";
import { BreadcrumbValues } from "@/components/SitemapUtils";
import { PageContext } from "@/lib/page_context";
import { createClient } from "@/utils/supabase/server";
import { headers } from "next/headers";
export default async function Page({
  params,
}: {
  params: Promise<{ doc_uuid: string }>;
}) {
  const supabase = createClient();
  const slug = (await params).doc_uuid;
  const headersList = headers();
  const host = headersList.get("host") || "";
  const hostsplits = host.split(".");
  const state = hostsplits.length > 1 ? hostsplits[0] : undefined;
  const pageContext: PageContext = {
    state: state,
    slug: ["proceedings", slug],
    final_identifier: slug,
  };
  const breadcrumbs: BreadcrumbValues = {
    state: state,
    breadcrumbs: [
      { title: "Files", value: "files" },
      { title: "Test Document Name", value: slug },
    ],
  };
  const {
    data: { user },
  } = await supabase.auth.getUser();
  return (
    <>
      <Navbar user={user} breadcrumbs={breadcrumbs} />
      <DocumentPage objectId={slug} />
    </>
  );
}
