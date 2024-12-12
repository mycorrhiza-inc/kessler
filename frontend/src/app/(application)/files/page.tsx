import ConversationComponent from "@/components/Conversations/ConversationComponent";
import PageContainer from "@/components/Page/PageContainer";
import { BreadcrumbValues } from "@/components/SitemapUtils";
import { stateFromHeaders } from "@/lib/nextjs_misc";
import { headers } from "next/headers";

export default function Page() {
  const headersList = headers();
  const state = stateFromHeaders(headersList);
  const breadcrumbs: BreadcrumbValues = {
    state: "ny",
    breadcrumbs: [{ title: "Files", value: "files" }],
  };
  return (
    <PageContainer breadcrumbs={breadcrumbs}>
      <h1 className="text-3xl font-bold">
        We dont really know what to put on this page, a view of all files is
        probably better gotten from the filter view on dockets, but we can
        put a table here if you really want to browse. If you have any better
        ideas about what to put here, please let me know - nicole
      </h1>
      <ConversationComponent inheritedFilters={[]} />
    </PageContainer>
  );
}
