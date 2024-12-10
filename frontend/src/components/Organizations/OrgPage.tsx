import axios from "axios";
import { useState, useEffect } from "react";
import { getOrganizationInfo } from "@/lib/requests/organizations";
import { BreadcrumbValues } from "../SitemapUtils";
import PageContainer from "../Page/PageContainer";
import OrganizationFileTable from "./OrgFileResults";
import ConversationComponent from "../Conversations/ConversationComponent";
import { FilterField } from "@/lib/filters";

export default async function OrganizationPage({
  breadcrumbs,
}: {
  breadcrumbs: BreadcrumbValues;
}) {
  const orgId =
    breadcrumbs.breadcrumbs[breadcrumbs.breadcrumbs.length - 1].value;

  const orgInfo = await getOrganizationInfo(orgId);
  const actual_breadcrumb_values = [
    ...breadcrumbs.breadcrumbs.slice(0, -1),
    { value: orgId, title: orgInfo.name || "Loading" },
  ];
  const actual_breadcrumbs: BreadcrumbValues = {
    breadcrumbs: actual_breadcrumb_values,
    state: breadcrumbs.state,
  };
  // const [page, setPage] = useState(0);

  return (
    <PageContainer breadcrumbs={actual_breadcrumbs}>
      <h1 className=" text-2xl font-bold">Organization: {orgInfo.name}</h1>
      <p>
        {orgInfo.description ||
          "Automatically generated org descriptions coming soon"}
      </p>
      <h1 className=" text-2xl font-bold">Authored Documents</h1>
      <ConversationComponent
        inheritedFilters={[
          { filter: FilterField.MatchAuthorUUID, value: orgId },
        ]}
      />
    </PageContainer>
  );
}
