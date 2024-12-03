import { ConversationView } from "@/components/Conversations/ConversationView";
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
  const slug = (await params).conversation_id;
  const headersList = headers();
  const host = headersList.get("host") || "";
  const hostsplits = host.split(".");
  const state = hostsplits.length > 1 ? hostsplits[0] : undefined;
  const pageContext: PageContext = {
    state: state,
    slug: ["proceedings", slug],
    final_identifier: slug,
  };

  const breadcrumbs = {
    state: state,
    breadcrumbs: [
      { value: "proceedings", title: "Proceedings" },
      { value: slug, title: slug },
    ],
  };
  const {
    data: { user },
  } = await supabase.auth.getUser();
  return (
    <PageContainer breadcrumbs={breadcrumbs}>
      <ConversationView pageContext={pageContext} />
    </PageContainer>
  );
}
