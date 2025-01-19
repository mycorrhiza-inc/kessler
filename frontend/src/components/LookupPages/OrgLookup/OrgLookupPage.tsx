import OrganizationTableInfiniteScroll from "@/components/Organizations/OrganizationTable";
import PageContainer from "@/components/Page/PageContainer";
import SearchBox from "@/components/Search/SearchBox";
import { BreadcrumbValues } from "@/components/SitemapUtils";

const OrgLookupPage = ({ breadcrumbs }: { breadcrumbs: BreadcrumbValues }) => {
  return (
    <PageContainer breadcrumbs={breadcrumbs}>
      <SearchBox />
      <OrganizationTableInfiniteScroll />
    </PageContainer>
  );
};

export default OrgLookupPage;
