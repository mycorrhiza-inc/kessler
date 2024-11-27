import { ConversationView } from "@/components/Conversations/ConversationView";
import ConversationTable from "@/components/Organizations/ConversationTable";
import PageContainer from "@/components/Page/PageContainer";
import { PageContext } from "@/lib/page_context";
import { createClient } from "@/utils/supabase/server";
import { headers } from "next/headers";
export default async function Page({
  params,
}: {
  params: Promise<{ conversation_id: string }>;
}) {
  const supabase = createClient();
  const headersList = headers();
  const host = headersList.get("host") || "";
  const hostsplits = host.split(".");
  const state = hostsplits.length > 1 ? hostsplits[0] : undefined;
  const pageContext: PageContext = {
    state: state,
    slug: ["proceedings"],
  };
  const breadcrumbs = {
    state: state,
    breadcrumbs: [{ value: "proceedings", title: "Proceedings" }],
  };
  const {
    data: { user },
  } = await supabase.auth.getUser();
  return (
    <PageContainer user={user} breadcrumbs={breadcrumbs}>
      <h1 className="text-3xl font-bold">Proceedings</h1>
      <ConversationTable />
      <h2 className="text-2xl font-bold">All Documents</h2>
      <ConversationView pageContext={pageContext} />
    </PageContainer>
  );
}
