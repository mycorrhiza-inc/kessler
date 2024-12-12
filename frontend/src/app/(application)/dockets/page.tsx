import ConversationComponent from "@/components/Conversations/ConversationComponent";
import ConversationTableInfiniteScroll from "@/components/Organizations/ConversationTable";
import PageContainer from "@/components/Page/PageContainer";
import { stateFromHeaders } from "@/lib/nextjs_misc";
import { createClient } from "@/utils/supabase/server";
import { headers } from "next/headers";
export default async function Page({
  params,
}: {
  params: Promise<{ conversation_id: string }>;
}) {
  const headersList = headers();
  const state = stateFromHeaders(headersList);
  const breadcrumbs = {
    state: state,
    breadcrumbs: [{ value: "dockets", title: "Dockets" }],
  };
  return (
    <PageContainer breadcrumbs={breadcrumbs}>
      <h1 className="text-3xl font-bold">Dockets</h1>
      <div className="overflow-x-hidden border-r pr-4">
        <ConversationTableInfiniteScroll />
      </div>
    </PageContainer>
  );
}
