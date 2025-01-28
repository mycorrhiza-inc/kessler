import OrganizationPage, {
  generateOrganizationData,
} from "@/components/ObjectPages/OrgPage";
import { BreadcrumbValues } from "@/components/SitemapUtils";
import { stateFromHeaders } from "@/lib/nextjs_misc";
import { Metadata } from "next";
import { headers } from "next/headers";
import { cache } from "react";

const cachedOrgData = cache(generateOrganizationData);

export async function generateMetadata({
  params,
}: {
  params: Promise<{ organization_id: string }>;
}): Promise<Metadata> {
  const slug = (await params).organization_id;
  const headersList = headers();
  const state = stateFromHeaders(headersList);
  const orgData = await cachedOrgData(slug);
  return {
    title: orgData.orgInfo.name,
  };
}

export default async function Page({
  params,
}: {
  params: Promise<{ organization_id: string }>;
}) {
  const slug = (await params).organization_id;
  const headersList = headers();
  const state = stateFromHeaders(headersList);
  const orgData = await cachedOrgData(slug);
  orgData.breadcrumbs.state = state || "";

  return <OrganizationPage orgInfo={orgData.orgInfo} />;
}
