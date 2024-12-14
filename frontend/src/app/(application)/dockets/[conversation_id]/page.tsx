import { ConversationPage } from "@/components/Conversations/ConversationPage";
import { stateFromHeaders } from "@/lib/nextjs_misc";
import { createClient } from "@/utils/supabase/server";
import { headers } from "next/headers";
export default async function Page({
  params,
}: {
  params: Promise<{ conversation_id: string }>;
}) {
  const supabase = createClient();
  const slug = (await params).conversation_id;
  const headersList = headers();
  const state = stateFromHeaders(headersList);

  const breadcrumbs = {
    state: state,
    breadcrumbs: [
      { value: "dockets", title: "Dockets" },
      { value: slug, title: slug },
    ],
  };
  const {
    data: { user },
  } = await supabase.auth.getUser();
  return <ConversationPage breadcrumbs={breadcrumbs} />;
}
