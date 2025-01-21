import OrgLookupPage from "@/components/LookupPages/OrgLookup/OrgLookupPage";
import { BreadcrumbValues } from "@/components/SitemapUtils";
import { stateFromHeaders } from "@/lib/nextjs_misc";
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
  const headersList = headers();
  const state = stateFromHeaders(headersList);
  const breadcrumbs: BreadcrumbValues = {
    state: state,
    breadcrumbs: [{ title: "Organizations", value: "orgs" }],
  };
  // TODO: zustand set breadcrumbs

  return <OrgLookupPage />;
}
