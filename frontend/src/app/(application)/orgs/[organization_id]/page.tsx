import OrganizationPage, {
  generateOrganizationData,
} from "@/components/Organizations/OrgPage";
import { BreadcrumbValues } from "@/components/SitemapUtils";
import { stateFromHeaders } from "@/lib/nextjs_misc";
import { Metadata } from "next";
import { headers } from "next/headers";

export const metadata: Metadata = {
  title: "ERROR IN SITE NAME",
};

export default async function Page({
  params,
}: {
  params: Promise<{ organization_id: string }>;
}) {
  const slug = (await params).organization_id;
  const headersList = headers();
  const state = stateFromHeaders(headersList);
  const orgData = await generateOrganizationData(slug, state || "");
  metadata.title = orgData.orgInfo.name;

  return (
    <OrganizationPage
      breadcrumbs={orgData.breadcrumbs}
      orgInfo={orgData.orgInfo}
    />
  );
}
