import OrganizationPage from "@/components/Organizations/OrgPage";
import { BreadcrumbValues } from "@/components/SitemapUtils";
import { stateFromHeaders } from "@/lib/nextjs_misc";
import { headers } from "next/headers";
export default async function Page({
  params,
}: {
  params: Promise<{ organization_id: string }>;
}) {
  const slug = (await params).organization_id;
  const headersList = headers();
  const state = stateFromHeaders(headersList);
  const breadcrumbs: BreadcrumbValues = {
    state: state,
    breadcrumbs: [
      { title: "Organizations", value: "orgs" },
      { title: "Test Organization Name", value: slug },
    ],
  };
  return <OrganizationPage breadcrumbs={breadcrumbs} />;
}
