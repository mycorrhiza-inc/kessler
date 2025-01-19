import OrgLookupPage from "@/components/LookupPages/OrgLookup/OrgLookupPage";
import ConversationTableInfiniteScroll from "@/components/Organizations/ConversationTable";
import PageContainer from "@/components/Page/PageContainer";
import { stateFromHeaders } from "@/lib/nextjs_misc";
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
  return <OrgLookupPage breadcrumbs={breadcrumbs} />;
}
