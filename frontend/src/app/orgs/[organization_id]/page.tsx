import { ConversationView } from "@/components/Conversations/ConversationView";
import Navbar from "@/components/Navbar";
import OrganizationPage from "@/components/Organizations/OrgPage";
import { BreadcrumbValues } from "@/components/SitemapUtils";
import { PageContext } from "@/lib/page_context";
import { createClient } from "@/utils/supabase/server";
import { headers } from "next/headers";
export default async function Page({
  params,
}: {
  params: Promise<{ organization_id: string }>;
}) {
  const supabase = createClient();
  const slug = (await params).organization_id;
  const headersList = headers();
  const host = headersList.get("host") || "";
  const hostsplits = host.split(".");
  const state = hostsplits.length > 1 ? hostsplits[0] : undefined;
  const pageContext: PageContext = {
    state: state,
    slug: ["proceedings", slug],
    final_identifier: slug,
  };
  const breadcrumbs: BreadcrumbValues = {
    state: state,
    breadcrumbs: [
      { title: "Organizations", value: "orgs" },
      { title: "Test Organization Name", value: slug },
    ],
  };
  const {
    data: { user },
  } = await supabase.auth.getUser();
  return (
    <>
      <Navbar user={user} breadcrumbs={breadcrumbs} />
      <OrganizationPage pageContext={pageContext} />
    </>
  );
}
