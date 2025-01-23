import {
  OrganizationInfo,
  getOrganizationInfo,
} from "@/lib/requests/organizations";
import { BreadcrumbValues } from "../SitemapUtils";

import FileSearchView from "@/components/Search/FileSearch/FileSearchView";
import { FilterField } from "@/lib/filters";

export const generateOrganizationData = async (orgId: string) => {
  const orgInfo = await getOrganizationInfo(orgId);

  const breadcrumbs: BreadcrumbValues = {
    breadcrumbs: [
      { title: "Organizations", value: "orgs" },
      { title: orgInfo.name, value: orgId },
    ],
  };
  return { breadcrumbs: breadcrumbs, orgInfo: orgInfo };
};

export default function OrganizationPage({
  orgInfo,
}: {
  orgInfo: OrganizationInfo;
}) {
  return (
    <>
      <h1 className=" text-2xl font-bold">Organization: {orgInfo.name}</h1>
      <p>
        {orgInfo.description ||
          "Automatically generated org descriptions coming soon"}
      </p>
      <h1 className=" text-2xl font-bold">Authored Documents</h1>
      <FileSearchView
        inheritedFilters={[
          { filter: FilterField.MatchAuthorUUID, value: orgInfo.id },
          { filter: FilterField.MatchAuthor, value: orgInfo.name },
        ]}
      />
    </>
  );
}
