import OrganizationTableInfiniteScroll from "@/components/Organizations/OrganizationTable";
import PageContainer from "@/components/Page/PageContainer";
import { BreadcrumbValues } from "@/components/SitemapUtils";
import { stateFromHeaders } from "@/lib/nextjs_misc";
import { PageContext } from "@/lib/page_context";
import { headers } from "next/headers";

export const metadata = {
  title: "Organizations - Kessler",
  description:
    "Search all availible organizations who write goverment documents.",
};
export default async function Page({
  params,
}: {
  params: Promise<{ organization_id: string }>;
}) {
  const slug = (await params).organization_id;
  const headersList = headers();
  const state = stateFromHeaders(headersList);
  const pageContext: PageContext = {
    state: state,
    slug: ["dockets", slug],
    final_identifier: slug,
  };
  const breadcrumbs: BreadcrumbValues = {
    state: state,
    breadcrumbs: [{ title: "Organizations", value: "orgs" }],
  };
  return (
    <>
      <PageContainer breadcrumbs={breadcrumbs}>
        <OrganizationTableInfiniteScroll />
      </PageContainer>
    </>
  );
}
