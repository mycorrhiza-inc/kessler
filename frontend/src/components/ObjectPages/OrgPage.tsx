import {
  OrganizationInfo,
  getOrganizationInfo,
} from "@/lib/requests/organizations";
import { BreadcrumbValues } from "../SitemapUtils";

import HeaderCard from "./HeaderCard";
import SearchResultsServerStandalone from "../Search/SearchResultsServer";
import { GenericSearchType } from "@/lib/adapters/genericSearchCallback";

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
      <HeaderCard title={orgInfo.name}>
        <p>
          {orgInfo.description ||
            "Automatically generated org descriptions coming soon"}
        </p>
      </HeaderCard>
      <SearchResultsServerStandalone
        searchInfo={{
          query: "",
          search_type: GenericSearchType.Filling,
        }}
      // inheritedFilters={[
      //   // { filter: FilterField.MatchAuthorUUID, value: orgInfo.id },
      //   { filter: FilterField.MatchAuthor, value: orgInfo.name },
      // ]}
      />
    </>
  );
}
