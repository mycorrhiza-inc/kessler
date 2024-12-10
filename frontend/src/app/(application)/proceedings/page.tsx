import ConversationComponent from "@/components/Conversations/ConversationComponent";
import ConversationTableSimple from "@/components/Organizations/ConversationTable";
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
    breadcrumbs: [{ value: "proceedings", title: "Proceedings" }],
  };
  return (
    <PageContainer breadcrumbs={breadcrumbs}>
      <h1 className="text-3xl font-bold">Proceedings</h1>
      <div className="max-h-[600px] overflow-x-hidden border-r pr-4">
        <ConversationTableSimple />
      </div>
      <h2 className="text-2xl font-bold">All Documents</h2>
      <ConversationComponent inheritedFilters={[]} />
    </PageContainer>
  );
}
