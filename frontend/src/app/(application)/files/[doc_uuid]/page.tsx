import DocumentPage from "@/components/Document/DocumentPage";
import { stateFromHeaders } from "@/lib/nextjs_misc";
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
  const state = stateFromHeaders(headersList);
  const {
    data: { user },
  } = await supabase.auth.getUser();
  return <DocumentPage objectId={slug} user={user} state={state} />;
}
