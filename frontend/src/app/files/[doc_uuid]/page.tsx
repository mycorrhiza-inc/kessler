import { ConversationView } from "@/components/Conversations/ConversationView";
import DocumentPage from "@/components/Document/DocumentPage";
import Navbar from "@/components/Navbar";
import { PageContext } from "@/lib/page_context";
import { createClient } from "@/utils/supabase/server";
import { headers } from "next/headers";
export default async function Page({
  params,
}: {
  params: Promise<{ doc_uuid: string }>;
}) {
  const slug = (await params).doc_uuid;
  const supabase = createClient();
  const {
    data: { user },
  } = await supabase.auth.getUser();
  return (
    <>
      <Navbar user={user} />
      <DocumentPage objectId={slug} />
    </>
  );
}
