import ConversationComponent from "@/components/Conversations/ConversationComponent";
import PageContainer from "@/components/Page/PageContainer";
import { BreadcrumbValues } from "@/components/SitemapUtils";
import { stateFromHeaders } from "@/lib/nextjs_misc";
import { headers } from "next/headers";

export const metadata = {
  title: "Files - Kessler",
  description: "Search all availible files.",
};
export default function Page() {
  const headersList = headers();
  const state = stateFromHeaders(headersList);
  const breadcrumbs: BreadcrumbValues = {
    state: "ny",
    breadcrumbs: [{ title: "Files", value: "files" }],
  };
  return (
    <PageContainer breadcrumbs={breadcrumbs}>
      <h1 className="text-3xl font-bold">All Files Search</h1>
      <ConversationComponent inheritedFilters={[]} />
    </PageContainer>
  );
}
