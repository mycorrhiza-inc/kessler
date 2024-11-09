import { ConversationView } from "@/components/ConversationView";
import Navbar from "@/components/Navbar";
import { createClient } from "@/utils/supabase/server";
export default async function Page({
  params,
}: {
  params: Promise<{ conversation_id: string }>;
}) {
  const supabase = createClient();
  const slugs = (await params).conversation_id;
  const slug = slugs?.[0];
  const {
    data: { user },
  } = await supabase.auth.getUser();
  return (
    <>
      <Navbar user={user} />
      <ConversationView conversation_id={slug} />
    </>
  );
}
