import { ConversationView } from "@/components/ConversationView";
import { createClient } from "@/utils/supabase/server";
export default async function Page({
  params,
}: {
  params: Promise<{ conversation_id: string }>;
}) {
  const supabase = createClient();
  const {
    data: { user },
  } = await supabase.auth.getUser();
  const convo_id = (await params).conversation_id;
  return <ConversationView conversation_id={convo_id} user={user} />;
}
