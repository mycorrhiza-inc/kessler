import OrganizationTableInfiniteScroll from "@/components/Organizations/OrganizationTable";
import PageContainer from "@/components/Page/PageContainer";
import SearchBox from "@/components/Search/SearchBox";
import { BreadcrumbValues } from "@/components/SitemapUtils";

const OrgLookupPage = ({ breadcrumbs }: { breadcrumbs: BreadcrumbValues }) => {
  return (
    <PageContainer breadcrumbs={breadcrumbs}>
      <h1 className="text-3xl font-bold">Organizations</h1>
      <div className="overflow-x-hidden border-r pr-4">
        <SearchBox />
        <OrganizationTableInfiniteScroll />
      </div>
    </PageContainer>
  );
};

export default OrgLookupPage;
